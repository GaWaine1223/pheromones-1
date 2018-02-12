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

## 需要解决的问题
1. 如何保存路由表, 允许重名，使用ip+server端口？
2. 对路由表更安全的操作，目前用尽可能小的锁。