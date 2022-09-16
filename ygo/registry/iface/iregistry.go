package iface

type Service struct {
	Id       string
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
}
type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type IRegistry interface {
	Register(service Service)
	UnRegister(service Service)
	UnRegisterById(id Service)
	GetService(name string)
}
