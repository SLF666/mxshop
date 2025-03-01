基于gin，grpc的微服务商城系统

service层向内提供grpc通讯，web层向外提供http通讯。

业务逻辑在web层开发，底层交互在service层开发。

所有服务集成nacos配置中心、consul注册服务发现、健康检查、负载均衡。

redis,rocketmq中间件，mysql存储，elasticsearch全文搜索。
