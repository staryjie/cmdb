## Provider 资源提供方

资源都在第三方服务上，比如：
 - 阿里云
 - 腾讯云
 - AWS
 - ...

云商资源 <-- Provider --> CMDB


1. 已经能够从云商（第三方接口）查询数据
2. 数据需要分页（控制每一次查询数据的量，从而保证接口性能）
    - 设计一个分页查询器（pagger），基础条件：query page_size, page_number，一直查询直到下一页没有数据，就停止查询
        - Next() bool：控制是否有下一页数据
        - Scan()：获取当前页的数据