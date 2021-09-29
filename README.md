### 概述

- Component : 一般表示外部资源，或者需要被托管的持久化对象，如cache对象~
- Service :表示可能会长期运行的服务，这里主要包含 grpc server ,web server 或 daemon job,crontab task等
- Client : 表示调用其他使用的client
- sd 服务发现，app会托管服务发现的注册和取消注册

### TODO

- app支持服务发现的生命周期管理
- 考虑job运行
- 长期运行的Job实现 , crontab ? demain
- 考虑client管理
- 公用 middlewares 支持
- grpc middlewares 支持