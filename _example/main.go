package main

import (
	"time"

	p2p "github.com/GaWaine1223/Lothar/pheromone"
	pto "github.com/GaWaine1223/Lothar/pheromone/_example/protocal"
	"encoding/json"
)

var (
	hello = p2p.ReqMsg{
		Name:"luda",
		Operation:"greeting",
	}
	hellomsg = pto.MsgGreeting{
		Addr:"北京",
		Account:11900,
	}
)

func main() {
	j, _ := json.Marshal(hellomsg)
	hello.Data = string(j)
	p := pto.NewProtocal()
	router := p2p.NewSRouter(time.Millisecond * 100, p)
	h1 := p2p.NewHost("luda", router, p2p.ShortConnection)
	h2 := p2p.NewHost("yoghurt", router, p2p.ShortConnection)
	h3 := p2p.NewHost("diudiu", router, p2p.ShortConnection)
	println("h1 监听 12345")
	go h1.Listen("127.0.0.1:12345")
	println("h2 监听 12345")
	go h2.Listen("127.0.0.1:12346")
	println("h3 监听 12345")
	go h3.Listen("127.0.0.1:12347")
	h1.Connect("yoghurt", "127.0.0.1:12346")
	helloMsg, _ := json.Marshal(hello)
	h1.Router.DispatchAll(helloMsg)
	//h2.Connect("diudiu", "127.0.0.1:12347")
	println("done")
	for {
		time.Sleep(time.Second)
	}
}

