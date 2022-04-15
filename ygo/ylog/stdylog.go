package ylog

import (
	"os"
	"time"
)

//StdYLog 创建全局log
var StdYLog = NewYLog(os.Stderr, "", BitDefault)
//Flags 获取StdYLog 标记位
func Flags() int {
	return StdYLog.Flags()
}

//ResetFlags 设置StdYLog标记位
func ResetFlags(flag int) {
	StdYLog.ResetFlags(flag)
}

//AddFlag 添加flag标记
func AddFlag(flag int) {
	StdYLog.AddFlag(flag)
}

//SetPrefix 设置StdYLog 日志头前缀
func SetPrefix(prefix string) {
	StdYLog.SetPrefix(prefix)
}

//SetLogFile 设置StdYLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdYLog.SetLogFile(fileDir, fileName)
}

//CloseDebug 设置关闭debug
func CloseDebug() {
	StdYLog.CloseDebug()
}
//CloseDebug 设置关闭debug
func TestReset(t time.Time) {
	StdYLog.TestReset(t)
}
//OpenDebug 设置打开debug
func OpenDebug() {
	StdYLog.OpenDebug()
}

//Debugf ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdYLog.Debugf(format, v...)
}

//Debug Debug
func Debug(v ...interface{}) {
	StdYLog.Debug(v...)
}

//Infof ====> Info <====
func Infof(format string, v ...interface{}) {
	StdYLog.Infof(format, v...)
}

//Info -
func Info(v ...interface{}) {
	StdYLog.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdYLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdYLog.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdYLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdYLog.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdYLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdYLog.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdYLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdYLog.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdYLog.Stack(v...)
}

func init() {
	//因为StdYLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的zinxLogger对象 calldDepth=2, StdYLog的calldDepth=3
	StdYLog.calldDepth = 3
}