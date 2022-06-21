# 任务管理

+ secret
+ provider
+ host service


/task/
```
# 这个是资源同步的任务
type: "sync"
# secret, 比如腾讯云secret
secret_id: "xxx"
# operater 按照资源划分, 不如操作主机
resource_type: "host"
# 操作那个地域的资源
region: "shanghai"
```

任务的状态：
+ 状态: running
+ 开启时间
+ 介绍时间
+ 执行日志