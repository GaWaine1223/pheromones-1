package pheromone

type MsgPto struct {
	Name 		string		`json:"name"`
	Operation 	string		`json:"operation"`
	Data 		interface{}	`json:"data"`
}

// 路由数据解析协议
type Protocal interface {
	// 解析通信内容
	Parse(msg []byte) (MsgPto, error)
	// 处理请求
	Handle(r MsgPto) (MsgPto, error)
}
