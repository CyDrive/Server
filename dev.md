# 项目结构
CyDrive 主要分为 master（服务层）和 storage nodes（存储层）两个部分，分别位于 master 和 node 目录。他们有一些公共的部分，置于项目的根目录：
- config：配置相关的目录
- consts：常量，枚举
- docs：文档
- model：定义了一些 struct，例如各种请求的 request 和 response
- rpc：定义了 master 与 node 之间的 rpc 相关的结构
- utils：封装了一些常用的方法和类型

## 结构定义
CyDrive 中的很多 枚举量，struct 都是通过 protobuf 来定义的，例如状态码枚举，各种接口的 request 和 response，以及一些简单的结构。 **当你需要更新它们时，应当修改对应的 proto 文件，并在项目根目录运行 uproto.sh，这个脚本会生成对应的源代码文件。** 你会注意到这个脚本还会生成 .cs 文件，可以将它们用于开发 C# SDK，这个 repo 中 ignore 了这些 .cs 文件。

# 编译
使用项目根目录下的的 build.sh 脚本来进行编译，该脚本会在 master/bin 下生成一个名为 master 的程序