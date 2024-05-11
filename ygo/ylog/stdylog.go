package ylog

/*
	A global Log handle is provided by default for external use, which can be called directly through the API series.
	The global log object is StdYLog.
	Note: The methods in this file do not support customization and cannot replace the log recording mode.

	If you need a custom logger, please use the following methods:
	ylog.SetLogger(yourLogger)
	ylog.Ins().InfoF() and other methods.

   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdYLog
   注意：本文件方法不支持自定义，无法替换日志记录模式，如果需要自定义Logger:

   请使用如下方法:
   ylog.SetLogger(yourLogger)
   ylog.Ins().InfoF()等方法
*/

// StdYLog creates a global log
var StdYLog = NewYLog("", BitDefault)

// Flags gets the flags of StdYLog
func Flags() int {
	return StdYLog.Flags()
}

// ResetFlags sets the flags of StdYLog
func ResetFlags(flag int) {
	StdYLog.ResetFlags(flag)
}

// AddFlag adds a flag to StdYLog
func AddFlag(flag int) {
	StdYLog.AddFlag(flag)
}

// SetPrefix sets the log prefix of StdYLog
func SetPrefix(prefix string) {
	StdYLog.SetPrefix(prefix)
}

// SetLogFile sets the log file of StdYLog
func SetLogFile(fileDir string, fileName string) {
	StdYLog.SetLogFile(fileDir, fileName)
}

// SetMaxAge 最大保留天数
func SetMaxAge(ma int) {
	StdYLog.SetMaxAge(ma)
}

// SetMaxSize 单个日志最大容量 单位：字节
func SetMaxSize(ms int64) {
	StdYLog.SetMaxSize(ms)
}

// SetCons 同时输出控制台
func SetCons(b bool) {
	StdYLog.SetCons(b)
}

// SetLogLevel sets the log level of StdYLog
func SetLogLevel(logLevel int) {
	StdYLog.SetLogLevel(logLevel)
}

func Debugf(format string, v ...interface{}) {
	StdYLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdYLog.Debug(v...)
}

func Infof(format string, v ...interface{}) {
	StdYLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdYLog.Info(v...)
}

func Warnf(format string, v ...interface{}) {
	StdYLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdYLog.Warn(v...)
}

func Errorf(format string, v ...interface{}) {
	StdYLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdYLog.Error(v...)
}

func Fatalf(format string, v ...interface{}) {
	StdYLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdYLog.Fatal(v...)
}

func Panicf(format string, v ...interface{}) {
	StdYLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdYLog.Panic(v...)
}

func Stack(v ...interface{}) {
	StdYLog.Stack(v...)
}

func init() {
	// Since the StdYLog object wraps all output methods with an extra layer, the call depth is one more than a normal logger object
	// The call depth of a regular zinxLogger object is 2, and the call depth of StdYLog is 3
	// (因为StdZinxLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	// 一般的zinxLogger对象 calldDepth=2, StdZinxLog的calldDepth=3)
	StdYLog.calldDepth = 3
}
