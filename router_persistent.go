// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// 长链接对象
type endPointP struct {
	c net.Conn
	// 0 : ON
	// 1 : Closing
	// 2 : OFF
	status int32
}

// 长链接路由
type PRouter struct {
	sync.RWMutex
	sync.WaitGroup
	to time.Duration
	// 长链接池
	pool map[string]endPointP
}

func NewPRouter(to time.Duration) *PRouter {
	var r PRouter
	r.to = to
	r.pool = make(map[string]endPointP, 0)
	return &r
}

// 添加路由时，已添加或者地址为空是都返回有错误，防止收到请求和主动连接重复建立
// 如果名字相同地址不同，则将原来的地址删除
func (r *PRouter) AddRoute(s string, addr interface{}) error {
	if addr == nil {
		return Error(ErrRemoteSocketEmpty)
	}
	r.RLock()
	a, b := r.pool[s]
	if b && a == addr.(string) {
		return Error(ErrRemoteSocketExist)
	}
	r.RUnlock()
	if b {
		r.Delete(s)
	}
	r.Lock()
	r.pool[s] = endPointP{addr.(net.Conn), 0}
	r.Unlock()
	return nil
}

func (r *PRouter) Delete(s string) error {
	r.Lock()
	defer r.Unlock()
	r.pool[s].c.Close()
	if r.pool[s].status != 0 {
		return errors.New("shutdown fail")
	}
	r.pool[s] = endPointP{r.pool[s].c, 1}
	for {
		if r.pool[s].status == 2 {
			break
		}
	}
	delete(r.pool, s)
	return nil
}

func (r *PRouter) DispatchAll(msg []byte) map[string][]byte {
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.pool {
		if v.status != 0 {
			continue
		}
		go func(name string) {
			r.Add(1)
			defer r.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("panic: %v", err)
				}
			}()
			v.c.SetWriteDeadline(time.Now().Add(r.to))
			_, err := v.c.Write(msg)
			if err != nil {
				r.Delete(name)
			}
		}(k)
	}
	r.Wait()
	return nil
}

func (r *PRouter) GetConnType() ConnType {
	return PersistentConnection
}

func (r *PRouter) FetchPeers() map[string]interface{} {
	p2 := make(map[string]interface{})
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.pool {
		p2[k] = v
	}
	return p2
}

func (r *PRouter) Dispatch(name string, msg []byte) ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	r.pool[name].c.SetWriteDeadline(time.Now().Add(r.to))
	_, err := r.pool[name].c.Write(msg)
	if err != nil {
		r.Delete(name)
	}
	return "", err
}
