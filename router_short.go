// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"io"
)

// 短链接对象
type endPointS struct {
	addr string
}

type SRouter struct {
	sync.RWMutex
	sync.WaitGroup
	to time.Duration
	// 短链接池
	pool map[string]endPointS
}

// 短链接路由
func NewSRouter(to time.Duration) *SRouter {
	var r SRouter
	r.to = to
	r.pool = make(map[string]endPointS, 0)
	return &r
}

func (r *SRouter) AddRoute(s string, addr interface{}) error {
	r.Lock()
	defer r.Unlock()
	r.pool[s] = endPointS{addr.(string)}
	return nil
}

func (r *SRouter) Delete(s string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.pool, s)
	return errors.New("shutdown success")
}

func (r *SRouter) DispatchAll(msg []byte) map[string][]byte {
	resp := make(map[string][]byte)
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.pool {
		go func(name string) {
			r.Add(1)
			defer r.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("panic: %v", err)
				}
			}()
			c, err := net.DialTimeout("tcp", v.addr, r.to)
			if err != nil {
				return
			}
			defer c.Close()
			for i := 0; i < 3; i++ {
				_, err = c.Write(msg)
				if err != nil {
					continue
				}
				msg, err := r.read(c, r.to)
				if err != nil {
					break
				}
				resp[k] = msg
			}
		}(k)
	}
	r.Wait()
	return resp
}

func (r *SRouter) GetConnType() ConnType {
	return ShortConnection
}

func (r *SRouter) FetchPeers() map[string]interface{} {
	p2 := make(map[string]interface{})
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.pool {
		p2[k] = v
	}
	return p2
}

func (r *SRouter) Dispatch(s string, msg []byte) ([]byte, error) {
	var resp []byte
	r.RLock()
	defer r.RUnlock()
	c, err := net.DialTimeout("tcp", r.pool[s].addr, r.to)
	if err != nil {
		return err
	}
	defer c.Close()
	for i := 0; i < 3; i++ {
		_, err = c.Write(msg)
		if err != nil {
			continue
		}
		resp, err = r.read(c, r.to)
		if err != nil {
			break
		}
	}
	return resp, err
}

func (h *SRouter) read(r io.Reader, to time.Duration) ([]byte, error) {
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