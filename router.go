// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

// 链接类型
type ConnType int

const (
	// 长链接方式
	PersistentConnection = iota
	// 短链接方式
	ShortConnection
)

// Router 路由接口
// 提供了长链接／短链接两种通信方式
type Router interface {
	// 短链接传的是地址；长链接传的是net.Conn
	AddRoute(s string, addr interface{}) error
	Delete(s string) error
	DispatchAll(msg []byte)
	GetConnType() ConnType
	FetchPeers() map[string]interface{}
	Dispatch(s string, msg []byte) error
	GetProtocal() Protocal
}
