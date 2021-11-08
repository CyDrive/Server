# Share
分享功能本质是授权其它人访问自己的文件，CyDrive 提供下面三种分享模式：
1. 拥有分享链接的人都可以访问
2. 输入密码的人可以访问
3. 文件可以被特定的人访问

同时 CyDrive 支持对分享的文件进行过期设置，分享过期后自动失效。用户删除分享的文件也会导致分享失效。

## Schema
我们需要一张分享表：
- uri
- file_path
- from
- to
- password
- left_access_count
- expire
- created_at

