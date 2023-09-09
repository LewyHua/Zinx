package znet

import "zinx/ziface"

// BaseRouter 实现router时，先嵌入BaseRouter基类，后续重写方法即可
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (b *BaseRouter) Handle(request ziface.IRequest)     {}
func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
