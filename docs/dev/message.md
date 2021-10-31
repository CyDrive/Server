# Message Module

消息模块负责在同一账号的不同设备之间共享临时数据，支持的数据类型有：
- 文本
- 图片
- 文件

## 整体实现
CyDrive 只支持同账号的不同设备之间发送消息，我们为每一个账号维护一个 Hub，每个设备对应一个 Websocket 连接。同一个账号的设备都连接到同一个 Hub。Hub 负责设备的上线、下线处理以及消息分发。当一个设备发送消息，Websocket 连接会把这个消息转发到 Hub，由 Hub 来进行消息分发。

## MessageManager
MessageManager 管理所有的 Hub，并负责将 message 存储到 MessageStore 中，用于客户端主动拉取消息记录。

## Hub
一个账号对应一个 Hub，Hub 负责处理设备的上线，下线以及消息分发。CyDrive 支持广播消息，这些都由 Hub 来实现。Hub 只负责将消息推给具体的连接去处理，而不实际发送消息

## MessageConn
表示与一个设备的连接，其负责读取客户端发送的消息，也负责将其它的消息推送给客户端。内部有一个 channel 来存储要推送的消息队列。

