# pheromones
Peer to Peer network.

## 网络模型
    在P2P网络环境中，彼此连接的多台计算机之间都处于对等的地位，各台计算机有相同的功能，无主从之分，一台计算机既可作为服务器，设定共享资源供网络中其他计算机所使用，又可以作为工作站，整个网络一般来说不依赖专用的集中服务器，也没有专用的工作站。网络中的每一台计算机既能充当网络服务的请求者，又对其它计算机的请求做出响应，提供资源、服务和内容。通常这些资源和服务包括：信息的共享和交换、计算资源（如CPU计算能力共享）、存储共享（如缓存和磁盘空间的使用）、网络共享、打印机共享等。
p2p网络的实现要基于传输层协议(TCP/UDP)，而使用TCP协议时又分为短连接和长连接：
![image](https://github.com/GaWaine1223/pheromones/raw/master/readme_image/short.png)

短连接实现中，每次由router层重新创建连接，进行通信，并等待回复。
![image](https://github.com/GaWaine1223/pheromones/raw/master/readme_image/perminent.png)

长连接实现中，应该由protocal层、更高层创建连接，并开启线程保持监听(因为要对监听结果进行处理，也就是调用protocal层实现的解析协议。因此监听线程必须由protocal层或更高层进行维护)
由于在广播数据的时候，需要进行传递式的广播，会形成广播风暴
![image](https://github.com/GaWaine1223/pheromones/raw/master/readme_image/broadcast.png)

因此需要在protocal层实现的传输协议中保证，传输的内容幂等(同一个消息请求n次的结果一致)且转发的信息不会再二次转发，这部分具体在协议层的实现中来规定。

## 使用

## 实现
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

## TODO 需要解决的问题
1. 如何保存路由表？
允许重名，使用ip+server端口。
2. 对路由表更安全的操作？
目前用尽可能小的锁。
3. 在会话尚未结束的时候，修改conn，重新连接，之前未传完的协议状态机会失效。
允许失效，同时应该保证有有效conn存在的时候，且正在状态流转的时候不允许修改conn。
4. 同时添加对方路由，并同时向对方发送hello的时候，互换conn地址，无法统一。
没解决，打电话和视频聊天都会出现这种情况，他们也都没解决，让用户自己重试。