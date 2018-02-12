# pheromones
Peer to Peer network.

### Host层 Server功能
ListenAndServe() 开启接口监听，将读到的数据传输给prtocal层解析

### Protocal层 handler功能
Handle() 处理接收到的数据
IOLoop() 维护长连接协程
...

### Router层 Client功能：发送数据
Dispatch*() []byte 发送数据，并返回结果  
AddRoute() 添加路由表
...

## 需要解决的问题
1. 如何保存路由表, 允许重名，使用ip+server端口？
2. 对路由表更安全的操作，目前用尽可能小的锁。
3. 在会话尚未结束的时候，修改conn，重新连接，之前未传完的协议状态机会失效。允许失效。
4. 同时添加对方路由，并同时向对方发送hello的时候，互换conn地址，无法统一。没解决