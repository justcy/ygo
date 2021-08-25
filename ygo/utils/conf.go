package utils

import (
	"encoding/json"
	"github.com/justcy/ygo/ygo/yiface"
	"io/ioutil"
)

type Config struct {
	TCPServer        yiface.IServer //Ygo全局Server对象
	Host             string         //IP
	TcpPort          int            //当前服务器监听的端口
	Name             string         //当前服务器名称
	Version          string         //当前Ygo版本号
	MaxPacketSize    uint32         //数据包的最大值
	MaxConn          int            //当前服务器允许的最大连接个数
	WorkerPoolSize   uint32         //业务工作Worker池的数量
	MaxWorkerTaskLen uint32         //业务工作Worker对应负责的任务队列最大任务存储数量
}

func (g Config) Reload() {
	data, err := ioutil.ReadFile("conf/ygo.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

var GlobalObject *Config

func init() {
	GlobalObject = &Config{
		Name:          "YgoServerApp",
		Version:       "0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
