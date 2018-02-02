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
	p  Protocal
	to time.Duration
	// 长链接池
	pool map[string]endPointP
}

func NewPRouter(to time.Duration, p Protocal) *PRouter {
	var r PRouter
	r.to = to
	r.pool = make(map[string]endPointP, 0)
	r.p = p
	return &r
}

func (r *PRouter) AddRoute(s string, addr interface{}) error {
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
	return errors.New("shutdown success")
}

func (r *PRouter) DispatchAll(msg []byte) {
	r.RLock()
	defer r.RUnlock()
	for _, v := range r.pool {
		if v.status != 0 {
			continue
		}
		go func() {
			r.Add(1)
			defer r.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("panic: %v", err)
				}
			}()
			v.c.SetWriteDeadline(time.Now().Add(r.to))
			v.c.Write(msg)
		}()
	}
	r.Wait()
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

func (r *PRouter) Dispatch(s string, msg []byte) error {
	r.RLock()
	defer r.RUnlock()
	r.pool[s].c.SetWriteDeadline(time.Now().Add(r.to))
	_, err := r.pool[s].c.Write(msg)
	return err
}

func (r *PRouter) GetProtocal() Protocal {
	return r.p
}
