package ylog_test

import (
	"github.com/justcy/ygo/ygo/ylog"
	"testing"
)
func TestStdYLog(t *testing.T) {
	//测试 默认debug输出
	ylog.Debug("zinx debug content1")
	ylog.Debug("zinx debug content2")

	ylog.Debugf(" zinx debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	ylog.ResetFlags(ylog.BitDate | ylog.BitLongFile | ylog.BitLevel)
	ylog.Info("zinx info content")

	//设置日志前缀，主要标记当前日志模块
	ylog.SetPrefix("MODULE")
	ylog.Error("zinx error content")

	//添加标记位
	ylog.AddFlag(ylog.BitShortFile | ylog.BitTime)
	ylog.Stack(" Zinx Stack! ")

	//设置日志写入文件
	//ylog.SetLogFile("./log", "testfile.log")
	ylog.Debug("===> zinx debug content ~~666")
	ylog.Debug("===> zinx debug content ~~888")
	ylog.Error("===> zinx Error!!!! ~~~555~~~")

	//关闭debug调试
	ylog.CloseDebug()
	ylog.Debug("===> 我不应该出现~！")
	ylog.Debug("===> 我不应该出现~！")
	ylog.Error("===> zinx Error  after debug close !!!!")
}