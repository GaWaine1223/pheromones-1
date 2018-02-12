// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import (
	"fmt"
	"net"
	"sync"
	"time"
	"io"
)

// 长链接对象
type endPointP struct {
	c net.Conn
}

// 长链接路由
type PRouter struct {
	sync.RWMutex
	sync.WaitGroup
	to time.Duration
	// 长链接池
	Pool map[string]endPointP
}

// 长链接路由
func NewPRouter(to time.Duration) *PRouter {
	var r PRouter
	r.to = to
	r.Pool = make(map[string]endPointP, 0)
	return &r
}

// 添加路由时，已添加或者地址为空是都返回有错误，防止收到请求和主动连接重复建立
// 如果名字相同地址不同，则将原来的地址删除
func (r *PRouter) AddRoute(name string, addr interface{}) error {
	if _, ok := addr.(net.Conn); !ok {
		return Error(ErrRemoteSocketMisType)
	}
	if addr.(net.Conn) == nil {
		return Error(ErrRemoteSocketEmpty)
	}
	if _, ok := r.Pool[name]; ok {
		if addr.(net.Conn) == r.Pool[name].c {
			return Error(ErrRemoteSocketExist)
		}
		r.Delete(name)
	}
	r.Lock()
	r.Pool[name] = endPointP{addr.(net.Conn)}
	r.Unlock()
	fmt.Printf("添加路由, peername=@%s@||peeraddress=%s\n", name, addr.(net.Conn).RemoteAddr())

	return nil
}

func (r *PRouter) Delete(name string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Pool[name] ; !ok {
		return Error(ErrRemoteSocketEmpty)
	}
	r.Pool[name].c.Close()
	delete(r.Pool, name)
	return nil
}

func (r *PRouter) GetConnType() ConnType {
	return PersistentConnection
}

func (r *PRouter) DispatchAll(msg []byte) map[string][]byte {
	for k, v := range r.Pool {
		r.Add(1)
		go func(name string, c net.Conn) {
			defer r.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("panic: %v", err)
				}
			}()
			fmt.Printf("dispatchall||发送请求, peername=%s||peeraddr=%s||msg=%s\n", name, c.RemoteAddr(), string(msg))
			r.RLock()
			c.SetWriteDeadline(time.Now().Add(r.to))
			_, err := c.Write(msg)
			r.RUnlock()
			if err != nil {
				r.Delete(name)
			}
		}(k, v.c)
	}
	r.Wait()
	return nil
}

func (r *PRouter) FetchPeers() map[string]interface{} {
	p2 := make(map[string]interface{})
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.Pool {
		p2[k] = v
	}
	return p2
}

func (r *PRouter) Dispatch(name string, msg []byte) ([]byte, error) {
	r.RLock()
	if _, ok := r.Pool[name] ; !ok {
		return nil, Error(ErrUnknuowPeer)
	}
	fmt.Printf("发送请求, peername=%s||msg=%s\n", name, string(msg))
	r.Pool[name].c.SetWriteDeadline(time.Now().Add(r.to))
	_, err := r.Pool[name].c.Write(msg)
	r.RUnlock()
	if err != nil {
		r.Delete(name)
	}
	return nil, err
}

func (r *PRouter) read(io io.Reader, to time.Duration) ([]byte, error) {
	buf := make([]byte, defultByte)
	messnager := make(chan int)
	go func() {
		n, _ := io.Read(buf[:])
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