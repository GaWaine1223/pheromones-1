package protocal

import (
	"io"
	"net"
	"fmt"
	"time"
	"encoding/json"

	p2p "github.com/GaWaine1223/Lothar/pheromone"
)

const (
	// 连接请求
	ConnectReq = "connectreq"
	// 获取一个
	GetReq = "getreq"
	// 批量获取
	FetchReq = "fetchreq"
	// 同步更新
	NoticeReq = "noticereq"

	// 连接请求返回
	ConnectResp = "connectresp"
	// 获取一个返回
	GetResp = "getresp"
	// 批量获取返回
	FetchResp = "fetchresp"
	// 同步更新返回
	NoticeResp = "noticeresp"

	// 未知操作
	UnknownOp = "unknownop"

	defultByte = 10240

)

type MsgGreetingReq struct {
	Addr 	string		`json:"add"`
	Account int		`json:"account"`
}

type Protocal struct {
	HostName string
	Router 	p2p.Router
	to 	time.Duration
}

func NewProtocal(name string, r p2p.Router, to time.Duration) *Protocal {
	return &Protocal{name, r, to}
}

func (p *Protocal) Handle(msg []byte) ([]byte, error) {
	req := &p2p.MsgPto{}
	resp := &p2p.MsgPto{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		resp.Name = p.HostName
		resp.Operation = UnknownOp
		ret, _ := json.Marshal(resp)
		return ret, p2p.Error(p2p.ErrMismatchProtocalReq)
	}
	resp.Name = p.HostName
	switch req.Operation {
	case ConnectReq:
		subReq := &MsgGreetingReq{}
		err := json.Unmarshal(req.Data, subReq)
		if err != nil {
			return nil, p2p.Error(p2p.ErrMismatchProtocalResp)
		}
		err = p.Router.AddRoute(req.Name, subReq.Addr)
		if err != nil {
			fmt.Printf("@%s@report: %s operation from @%s@ failed, err=%s\n", p.HostName, req.Operation, req.Name, err)
		}
		resp.Operation = ConnectResp
	case GetReq:
		resp.Operation = GetResp
	case FetchReq:
		resp.Operation =FetchResp
	case NoticeReq:
		resp.Operation = NoticeResp
	case ConnectResp:
		resp.Operation = GetReq
	case GetResp:
		resp.Operation = FetchReq
	case FetchResp:
		resp.Operation = NoticeReq
	case NoticeResp:
		fmt.Printf("@%s@report: %s operation from @%s@ finished\n", p.HostName, req.Operation, req.Name)
		return nil, nil
	default:
		resp.Operation = UnknownOp
	}
	ret, err := json.Marshal(resp)
	fmt.Printf("@%s@report: %s operation from @%s@ succeed\n", p.HostName, req.Operation, req.Name)
	return ret, nil
}

func (p *Protocal) GetRouter() p2p.Router {
	return p.Router
}

func (p *Protocal) Add(name string, addr string) error{
	if p.Router.GetConnType() == p2p.ShortConnection {
		return p.Router.AddRoute(name, addr)
	}
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	err = p.Router.AddRoute(name, c)
	go p.IOLoop(c)
	return err
}

// 长连接的话，需要在加入路由的时刻起携程 循环监控
func (p *Protocal) IOLoop(c net.Conn) {
	for {
		msg, err := p.read(c)
		if err != nil {
			return
		}
		resp, err := p.Handle(msg)
		if err != nil && resp != nil {
			continue
		}
		c.SetWriteDeadline(time.Now().Add(p.to))
		_, err = c.Write(resp)
		if err != nil {
			return
		}
	}
}

func (p *Protocal) read(r io.Reader) ([]byte, error) {
	buf := make([]byte, defultByte)
	_, err := r.Read(buf[:])
	if err != nil {
		return nil, err
	}
	return	buf, nil
}

func (p *Protocal) DispatchAll(msg []byte) map[string][]byte {
	return p.Router.DispatchAll(msg)
}

func (p *Protocal) Dispatch(name string, msg []byte) ([]byte, error) {
	return p.Router.Dispatch(name, msg)
}