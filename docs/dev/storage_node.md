# Storage Node
在集群部署的情况下，数据存储在 storage nodes 。根据网络情况，我们把 storage nodes 分为两类：
1. Private：无公网 IP
2. Public：有公网 IP

对于前者，可以尝试 NAT 穿透，但如果失败，那么流量还是需要经过 Master。

在这个设计中，我们假设 Master 是可靠的。

## Private Storage Node
这里主要描述流量经过 Master 的设计，我们用一个 server-side streaming RPC 来实现通知和对 Node 的主动管理。例如修改 Node 的状态，通知 Node 发送/接收文件等。Node 收到通知后，建立相应的连接来进行数据传输。没有采用双向流是因为这样我们不需要一直维护大量的长连接，而且可以每个传输任务使用单独的连接，让它们之间不相互影响。另一个原因就是这样的代码可维护性会更高。

### Notifier

### File
对于读写文件的情况，Master 作为主动方，向 storag node 发出读写请求，我们定义下面的 message 作为 Master 发出的请求：
```protobuf
message MasterFileChunk {
    int64 task_id = 1;          // r&w, first for read, continue for w
    string file_path = 2;       // r&w, first
    bool should_truncate = 3;   // w, first
    bool should_append = 4;     // w, first
    
    bytes data = 5;             // w, continue
}
```

#### Read
读流程分为下面的几步：
1. Master 生成 task_id，并发出带有 `task_id` 和 `file_path` 的请求
2. Node 返回一个 stream，stream 中的每个 message 是带有 `task_id` 和 `data` 的响应，如果响应中 `data` 字段为空，则表示已经读到 EOF

如果这个过程中 Node 发生了 crash 或网络故障的情况，那么 Master 会检测到读超时，从而切换到其它 storage node 去读取数据。

#### Write
写流程分为下面的几步：
1. Master 生成 `task_id`，并发出带有 `task_id`，`file_path`，以及 `data` 的请求
2. Node 返回一个带有 `task_id` 与 `offset` 的相应，后者指示当前写到的位置
3. Master 接收到请求后写后续数据，请求中只带 `task_id` 和 `data`
4. 如果发生错误，则 Node 会进行一次响应，从而让 Master 停止写数据

#### Optimizations
这里介绍一个优化机制：group read。

对于读请求，如果多个请求都是在读一个文件，则 Master 应该能够检测到，只会去 storage node 读一次。

#### Conflicts
对于写的情况，如果多个请求都在写同一个文件，则应当检测到冲突。

