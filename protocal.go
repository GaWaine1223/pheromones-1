// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import "net"

type MsgPto struct {
	Name 		string		`json:"name"`
	Operation 	string		`json:"operation"`
	Data 		[]byte		`json:"data"`
}

// 路由数据解析协议
type Protocal interface {
	// 解析请求通信内容,并返回数据,双工协议
	Handle(c net.Conn, msg []byte) ([]byte, error)
}
