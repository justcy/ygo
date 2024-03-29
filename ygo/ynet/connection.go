package ynet

import (
	"context"
	"errors"
	"fmt"
	"github.com/justcy/ygo/ygo/utils"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	sync.RWMutex
	//当前Conn属于哪个Server
	TcpServer yiface.IServer //当前conn属于哪个server，在conn初始化的时候添加即可
	//当前连接的套接字
	Conn *net.TCPConn
	//当前连接的ID，也可以作为sessionID，全局唯一
	ConnId uint32
	//当前链接的关闭状态
	isClosed bool
	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler yiface.IMsgHandle
	//告知该链接已经退出
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc

	//ExitBuffChan chan bool
	//无缓冲通道，用于读写两个goroutine之间的通信
	msgChan chan []byte
	//有关冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := c.TcpServer.Packet()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		ylog.Infof("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := c.TcpServer.Packet()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgChan <- msg //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取
	return nil
}

func NewConnection(server yiface.IServer, conn *net.TCPConn, connId uint32, handle yiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnId:     connId,
		isClosed:   false,
		MsgHandler: handle,
		//ExitBuffChan: make(chan bool, 1),
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]interface{}), //对链接属性map初始化
	}
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c) //将当前新创建的连接添加到ConnManager中
	return c
}

/* 处理conn读数据的Goroutine */
func (c *Connection) StartReader() {
	ylog.Info("Reader Goroutine is  running")
	defer ylog.Info(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		select {
		case <- c.TcpServer.GetCtx().Done():
			ylog.Info("connection  recive server stop 2")
			return
		case <- c.ctx.Done():
			return
		default:
			//读取客户端Msg head
			headData := make([]byte, c.TcpServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
				ylog.Errorf("read msg head error %s", err)
				return
			}
			//拆包，得到msgid 和 datalen 放在msg中
			msg, err := c.TcpServer.Packet().UnPack(headData)
			if err != nil {
				ylog.Errorf("unpack error %s", err)
				return
			}
			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
					ylog.Errorf("read msg data error %s", err)
					return
				}
			}
			msg.SetData(data)
			//得到当前客户端请求的Request数据
			req := Request{
				Conn: c,
				Msg:  msg,
			}
			c.SetProperty(LAST_ACTIVE,time.Now().Unix())
			if utils.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(&req)
			}
			
		}


	}
}

func (c *Connection) StartWriter() {

	ylog.Info("[Writer Goroutine is running]")
	defer ylog.Info(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case <- c.TcpServer.GetCtx().Done():
			ylog.Info("connection  recive server stop 1")
			return
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				ylog.Errorf("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					ylog.Errorf("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				break
				ylog.Info("msgBuffChan is Closed")
			}
		case <-c.ctx.Done():
			//conn已经关闭
			return
		}
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	//关闭Writer
	c.cancel()

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}
