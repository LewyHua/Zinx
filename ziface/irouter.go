package ziface

// IRouter 路由抽象接口，路由里的数据都是IRequest
type IRouter interface {
	PreHandle(request IRequest)  // 处理conn业务之前的hook
	Handle(request IRequest)     // 处理conn业务的主hook
	PostHandle(request IRequest) // 处理conn业务之后的hook
}
