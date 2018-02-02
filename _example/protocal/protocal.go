package protocal

import (
	p2p "github.com/GaWaine1223/Lothar/pheromone"
	"encoding/json"
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

	// 连接请求
	ConnectResp = "connectresp"
	// 获取一个
	GetResp = "getresp"
	// 批量获取
	FetchResp = "fetchresp"
	// 同步更新
	NoticeResp = "noticeresp"
)

type MsgGreeting struct {
	Addr 	string		`json:"add"`
	Account int		`json:"account"`
}

type protocal struct {
}

func NewProtocal() *protocal {
	return &protocal{}
}
// 解析通信内容
func (p *protocal) Parse(data []byte) (p2p.ReqMsg, error) {
	msg := &p2p.ReqMsg{}
	err := json.Unmarshal(data, msg)
	return msg, err
}

// 处理收到的请求
func (p *protocal) Handle(r p2p.ReqMsg) error {
	switch r.Operation {
	case ConnectReq:
	case GetReq:
	case FetchReq:
	case NoticeReq:
	case ConnectResp:
	case GetResp:
	case FetchResp:
	case NoticeResp:
	}
	return "", nil
}

