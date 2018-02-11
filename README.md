# pheromones
Peer to Peer network.

###Host层 Server功能
ListenAndServe() 开启接口监听，将读到的数据传输给prtocal层解析

###Protocal层 handler功能
Handle() 处理接收到的数据
IOLoop() 维护长连接协程
...

###Router层 Client功能：发送数据
Dispatch*() []byte 发送数据，并返回结果  
AddRoute() 添加路由表
...