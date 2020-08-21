开发文档
====

- [服务端接口文档](./server_api.md)

订单状态变化流程
```
 order stage change process

  customer order       customer cancel
  ---------------> new ---------------> c_cancel
                    |
                    |  business cancel
                    + ----------------> b_cancel
                    |
       customer pay |
                    |
                    v   customer want cancel
                   paid --------------------> c_w_cancel
                    |                            |
                    |                            |
                    +----------------------------+
                    |                            |
     business trans |            business repaid |
                    |                            |
                    v                            v
                 transport                     repaid
                    |
                    |
    customer commit |
                    |
                    v
                  finish
```

