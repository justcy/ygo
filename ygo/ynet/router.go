package ynet

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
)

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request yiface.IRequest) {
}

func (br *BaseRouter) Handle(request yiface.IRequest) {
}

func (br *BaseRouter) AfterHandle(request yiface.IRequest) {
}
