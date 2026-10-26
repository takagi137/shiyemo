# Feature: 上传头像成功

## 场景：上传 JPG 头像

### 前置条件
- 用户已登录
- 准备一个 JPG 图片文件

### 测试步骤

1. 发送 POST 请求到 `/api/v1/users/avatar`
2. 请求头携带 `Authorization: Bearer <token>`
3. 请求体为 multipart/form-data，包含 avatar 字段

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "avatar_url": "/uploads/avatars/1.jpg"
  }
}
```

### 验证点

- [x] 上传成功返回头像 URL
- [x] 头像文件保存到 /uploads/avatars/ 目录
- [x] 文件名使用用户 ID

## 场景：上传 PNG 头像

### 测试步骤

1. 上传 PNG 格式图片

### 预期结果

- 上传成功
- 返回 .png 后缀的 URL

## 场景：上传超大文件失败

### 测试步骤

1. 上传超过 2MB 的图片

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40001,
  "message": "文件大小超过限制",
  "data": null
}
```

## 场景：上传非图片文件失败

### 测试步骤

1. 上传 .txt 文件

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40001,
  "message": "只支持图片格式",
  "data": null
}
```
