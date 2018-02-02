// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import (
	"net"
	"time"
	"io"
	"strings"
	"runtime"
	"fmt"
)

const defultByte = 1024

type Host struct {
	HostName 	string
	Router 		Router
	Type 		ConnType
	temp		int32
}

func NewHost(name string, r Router, t ConnType) *Host {
	return &Host{name, r, t, 0}
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

func (h *Host) handler(c net.Conn) {
	if h.Type == ShortConnection {
		msg, err := h.read(c, h.Router.(*SRouter).to)
		println(string(msg))
		if err != nil {
			return
		}
		req, err := h.Router.GetProtocal().Parse(msg)
		if err != nil {
			return
		}
		// TODO 处理完要返回
		err = h.Router.GetProtocal().Handle(req)
		if err != nil {
			return
		}
		h.Router.AddRoute(req.Name, c.RemoteAddr().String())
		fmt.Printf("%s:::将%s添加入路由，地址为%s", h.HostName, req.Name, c.RemoteAddr().String())
	} else {
		msg, err := h.read(c, h.Router.(*PRouter).to)
		if err != nil {
			return
		}
		req, err := h.Router.GetProtocal().Parse(msg)
		if err != nil {
			return
		}
		err = h.Router.GetProtocal().Handle(req)
		if err != nil {
			return
		}
		// TODO 没名字或者已添加，都返回失败
		err = h.Router.AddRoute(req.Name, c)
		// 如果没添加对方路由成功，则长连接关闭
		if err != nil {
			return
		}
		for {
			if h.Router.(*PRouter).pool[req.Name].status == 1 {
				h.Router.(*PRouter).Lock()
				h.Router.(*PRouter).pool[req.Name] = endPointP{h.Router.(*PRouter).pool[req.Name].c, 2}
				h.Router.(*PRouter).Unlock()
				return
			}
			msg, err := h.read(c, h.Router.(*PRouter).to)
			if err != nil {
				return
			}
			req, _ := h.Router.GetProtocal().Parse(msg)
			err = h.Router.GetProtocal().Handle(req)
		}
	}
}

// 主动连接peer
func (h *Host) Connect(name, addr string) error {
	if h.Type == ShortConnection {
		return h.Router.AddRoute(name, addr)
	} else {
		c, err := net.DialTimeout("tcp", addr, h.Router.(*PRouter).to)
		if err != nil {
			return err
		}
		err = h.Router.AddRoute(name, c)
		if err != nil {
			return err
		}
		go h.listen(name, c)
	}
	return nil
}

// 长连接的话，需要在加入路由的时刻起携程 循环监控
func (h *Host) listen(s string, c net.Conn) {
	for {
		if h.Router.(*PRouter).pool[s].status == 1 {
			h.Router.(*PRouter).Lock()
			h.Router.(*PRouter).pool[s] = endPointP{h.Router.(*PRouter).pool[s].c, 2}
			h.Router.(*PRouter).Unlock()
			break
		}
		msg, err := h.read(c, h.Router.(*PRouter).to)
		if err != nil {
			return
		}
		req, _ := h.Router.GetProtocal().Parse(msg)
		err = h.Router.GetProtocal().Handle(req)
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