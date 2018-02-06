// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import "io"

type MsgPto struct {
	Name 		string		`json:"name"`
	Operation 	string		`json:"operation"`
	Data 		interface{}	`json:"data"`
}

// 路由数据解析协议
type Protocal interface {
	// 解析通信内容,同时负责添加主动连接过来的peer
	Parse(r io.Reader, msg []byte) (MsgPto, error)
	// 处理请求
	Handle(r MsgPto) ([]byte, error)
	// 获取协议底层路由
	GetRouter() Router
}
