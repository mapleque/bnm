Database Design
====

所有的数据库表设计都可以在[sql](/sql)文件夹中找到。

这里主要解释数据之间的关系和设计初衷。

```
business_user        business_profile         customer_user
  id <-----------+     id <---------+    +----> id
  open_id(unq)   +---- buid(unq)    |    |      open_id(unq)
  union_id       |     name         |    |      union_id
  token(unq)     |     avatar       |    |      token(unq)
  expired_at     |     desc         |    |      expired_at
                 |     qrcode       |    |
                 |     wxid         |    |
                 |     status       |    |    customer_address
                 |     create_at    |    |      id
                 +---------------+  |    +----- cuid(idx)
         +--------------------------+    |      label
item     |         order         |  |    |      reciever
  id <---|----+      id <-----+  |  |    |      address
  bid ---+    |      cuid ----------|----+      phone
  name <----+ |      bid -----------+    |      update_at
  price <-+ | +----- itid     |  |       |      create_at
  pic     | +------> name     |  |       |
  desc    +--------> price    |  |       |
  status             counts   |  |       |  order_log
  update_at          reciever |  |       |    id
  create_at          address  |  +----------- buid
  -idx(bid,status)   phone    |          +--- cuid
                     stage    +-------------- oid(idx)
                     status                   op
                     additional               new
                     exp_no                   create_at
                     update_at
                     create_at
                     -idx(cuid,status)
                     -idx(bid,status)
```
