package iface

type ITask interface {
	String() interface{}
	Run()
	Reset()
}
