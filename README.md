# Github Uploader

上传资源到Github仓库，返回 `jsdelivr` 访问地址

## 部署
### 修改配置
`yaml`中各种配置项都很简单，看一眼就能明白

### 上传页面
`/upload`

## Api

### Upload文件上传接口

- /upload
- POST
- multipart/form-data

| 参数名称 | 参数类型 | 必选 |说明 |
| :-----| :---- | :---- |:---- |
| file | 二进制文件 | 是 |可以有多个文件 |

- response 

```json5
{
  "success": true,              // 整个请求是否成功
  "data": [{
    "success": true,              // 这个文件上传是否成功
    "data": "https://xxxx,jpg",   // 资源路径，如果上传失败，这个字段则为github完整的响应体
    "code": "ok",                 
    "message": "ok"
  }],
  "code": "ok",
  "message": "ok"
}
```