package ylog_test

import (
	"github.com/justcy/ygo/ygo/ylog"
	"testing"
)

func TestStdylog(t *testing.T) {

	//测试 默认debug输出
	ylog.Debug("ylog debug content1")
	ylog.Debug("ylog debug content2")

	ylog.Debugf(" ylog debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	ylog.ResetFlags(ylog.BitDate | ylog.BitLongFile | ylog.BitLevel)
	ylog.Info("ylog info content")

	//设置日志前缀，主要标记当前日志模块
	ylog.SetPrefix("MODULE")
	ylog.Error("ylog error content")

	//添加标记位
	ylog.AddFlag(ylog.BitShortFile | ylog.BitTime)
	ylog.Stack(" ylog Stack! ")

	//设置日志写入文件
	ylog.SetLogFile("./log", "testfile.log")
	ylog.Debug("===> ylog debug content ~~666")
	ylog.Debug("===> ylog debug content ~~888")
	ylog.Error("===> ylog Error!!!! ~~~555~~~")

	//调试隔离级别
	ylog.Debug("=================================>")
	//1.debug
	ylog.SetLogLevel(ylog.LogInfo)
	ylog.Debug("===> 调试Debug：debug不应该出现")
	ylog.Info("===> 调试Debug：info应该出现")
	ylog.Warn("===> 调试Debug：warn应该出现")
	ylog.Error("===> 调试Debug：error应该出现")
	//2.info
	ylog.SetLogLevel(ylog.LogWarn)
	ylog.Debug("===> 调试Info：debug不应该出现")
	ylog.Info("===> 调试Info：info不应该出现")
	ylog.Warn("===> 调试Info：warn应该出现")
	ylog.Error("===> 调试Info：error应该出现")
	//3.warn
	ylog.SetLogLevel(ylog.LogError)
	ylog.Debug("===> 调试Warn：debug不应该出现")
	ylog.Info("===> 调试Warn：info不应该出现")
	ylog.Warn("===> 调试Warn：warn不应该出现")
	ylog.Error("===> 调试Warn：error应该出现")
	//4.error
	ylog.SetLogLevel(ylog.LogPanic)
	ylog.Debug("===> 调试Error：debug不应该出现")
	ylog.Info("===> 调试Error：info不应该出现")
	ylog.Warn("===> 调试Error：warn不应该出现")
	ylog.Error("===> 调试Error：error不应该出现")
}

func Testylogger(t *testing.T) {
}
