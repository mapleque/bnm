Server Api
====

本服务接口除登陆外全部遵循RESTful风格设计。

例如：
- `GET /资源名` 表示获取资源列表
- `POST /资源名` 表示新建资源
- `GET /资源名/id` 表示获取指定资源
- `POST /资源名/id` 表示修改指定资源
- `DELETE /资源名/id` 表示删除指定资源

注意：这里资源支持嵌套，其含义参考如下示例：
- `GET /资源a/aid/资源b/bid` 表示获取id为aid的资源a所属的id为bid的资源b

分页参数
----

所有获取列表接口都需要携带分页参数。分页参数通过URL Query形式传递，例如：
- `GET /资源名?s=10&l=24` 表示从标识24往后（不包含）获取10条数据

其中：
- `s` 表示获取数据条数，必选参数
- `l` 表示上次获取数据的最后标识，可选参数

接口返回数据顶层结构和状态码说明
----

返回数据顶层结构：

```json
{
    "status": <int>, // 状态码
    "data": <json>, // 返回数据
    "message": <json>, // 附加信息
}
```

其中，状态码表如下：

| 状态码 | 说明 |
| --: | :-- |
|200|请求成功|
|400|请求参数有问题|
|401|无访问权限|
|404|请求资源不存在|
|500|服务端问题|

- 当状态码为`200`时，通常有返回数据`data`
- 当状态码不是`200`是，通常有附加信息`message`

认证
----

本项目主要用于为微信小程序提供服务，除登陆接口外都需要携带认证标识：

```
// HTTP Header:
SessionKey: <token>
```

其中，`token`由登陆接口([/bussiness/login]()和[/customer/login]())返回。

注意：`token`有过期时间，当`token`过期后，所有接口都会返回状态码`status:401`，请实现跳转登陆逻辑。

特别的：由于本服务登陆接口都是通过微信授权，不排除会出现服务不稳定情况，此时接口也可能返回状态码`status:401`，这种情况是不可预知且不可控的，因此自动登录逻辑需要注意避免死循环，设置自动重试次数。

`/business/login`
----

### POST

1\. 请求示例
```sh
curl
```

1\. 请求参数
```json
{
}
```

2\. 返回数据
3\. 可能出现的状态码

`/business/profile`
----

### GET
### POST

`/business/items`
----

`/business/items/<itid>`
----

`/business/orders`
----

`/business/orders/<oid>`
----

`/customer/login`
----

`/customer/business`
----

`/customer/business/<bid>`
----

`/customer/business/<bid>/items`
----

`/customer/business/<bid>/items/<itid>`
----

`/customer/orders`
----

`/customer/orders/<oid>`
----

`/customer/address`
----

`/customer/address/<aid>`
----
