// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import (
	"net"
	"time"
	"io"
	"strings"
	"runtime"
)

const defultByte = 1024

type Host struct {
	HostName 	string
	Type 		ConnType
	temp		int32
	Proto 		Protocal
}

func NewHost(name string, t ConnType, p Protocal) *Host {
	return &Host{name, t, 0, p}
}

// 监听peer的链接请求
func (h *Host) Listen(localAddr string) error{
	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
			}
			break
		}
		go h.handler(c)
	}
	return nil
}

// 被动连接peer  类比connect函数，
// 一个是主动连接：直接就有姓名，可以直接添pool
// 一个是被动连接：需要进行通信，才能知道对方姓名，并添加pool
func (h *Host) handler(c net.Conn) {
	if h.Type == ShortConnection {
		msg, err := h.read(c, h.Proto.GetRouter().(*SRouter).to)
		println(string(msg))
		if err != nil {
			return
		}
		req, err := h.Proto.Parse(c, msg)
		if err != nil {
			return
		}
		resp, err := h.Proto.Handle(req)
		if err != nil {
			return
		}
		c.Write(resp)
	} else {
		// server接受长连接时，需要一个通信来确认对方姓名。
		msg, err := h.read(c, h.Proto.GetRouter().(*PRouter).to)
		if err != nil {
			return
		}
		req, err := h.Proto.Parse(c, msg)
		if err != nil {
			return
		}
		resp, err := h.Proto.Handle(req)
		if err != nil {
			return
		}
		_, err = c.Write(resp)
		if err != nil {
			return
		}
		err = h.Proto.GetRouter().AddRoute(req.Name, c)
		// 如果没成功添加对方路由成功，则长连接关闭
		if err != nil {
			return
		}
		// 下面这段和listen函数可以合并
		h.listen(req.Name, c)
	}
}

// 主动连接peer  类比handler函数，
// 一个是主动连接：直接就有姓名，可以直接添pool
// 一个是被动连接：需要进行通信，才能知道对方姓名，并添加pool
func (h *Host) Connect(name, addr string) error {
	routerI := h.Proto.GetRouter()
	if h.Type == ShortConnection {
		return routerI.AddRoute(name, addr)
	} else {
		c, err := net.DialTimeout("tcp", addr, routerI.(*PRouter).to)
		if err != nil {
			return err
		}
		err = routerI.AddRoute(name, c)
		if err != nil {
			return err
		}
		go h.listen(name, c)
	}
	return nil
}

// 长连接的话，需要在加入路由的时刻起携程 循环监控
func (h *Host) listen(s string, c net.Conn) {
	routerI := h.Proto.GetRouter()
	for {
		if routerI.(*PRouter).pool[s].status == 1 {
			routerI.(*PRouter).Lock()
			routerI.(*PRouter).pool[s] = endPointP{routerI.(*PRouter).pool[s].c, 2}
			routerI.(*PRouter).Unlock()
			break
		}
		msg, err := h.read(c, routerI.(*PRouter).to)
		if err != nil {
			return
		}
		req, err := h.Proto.Parse(c, msg)
		if err != nil {
			return
		}
		resp, err := h.Proto.Handle(req)
		if err != nil {
			return
		}
		_, err = c.Write(resp)
		if err != nil {
			return
		}
	}
}

func (h *Host) read(r io.Reader, to time.Duration) ([]byte, error) {
	buf := make([]byte, defultByte)
	messnager := make(chan int)
	go func() {
		n, _ := r.Read(buf[:])
		messnager <- n
		close(messnager)
	}()
	select {
	case n := <-messnager:
		return buf[:n], nil
	case <-time.After(to):
		return nil, Error(ErrLocalSocketTimeout)
	}
	return	buf, nil
}