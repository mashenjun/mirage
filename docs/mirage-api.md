# API文档说明

## 调用接口说明
- request请求时必须添加header头：Content-Type:application/json
- 所有的接口的返回形式都是统一为：
    - 正常返回
	```
	{
  	  "code": 0,
 	  "message": "ok",
  	  "data": {}
	}
	```
   - 错误返回

	```
	{
  	  "code": 错误码,
  	  "message": "错误信息"
	}
	```


## 公共错误码定义
| 错误码 | 说明  |
| ------ | ----- |
| 0      | 正常  |
| xxx    | ooooo |


## 广告配置 API

根据ad_code获取对应广告配置内容

### 请求

```
GET /api/v1/config/advertise?type=<xxx> HTTP/1.1
Content-Type: application/json

```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| type | `string` | 是   | 广告位code |
| extra | `string` | 否   | 额外信息字段 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "candidates":[
            {
                "ad_code": "<ad_code>",
                "ad_id": "<ad_id>",
                "ad_channel": "<ad_channel>", 
                "width": 9,
                "height": 16,
                "cool_down": 100, /*广告展示冷却时间，单位为秒*/
                "count_down": 100, /*广告倒计时，单位为秒*/
                "action":1, /* 0 广告 1 原生 2 webView*/
                "image_url":"http://via.placeholder.com/300",
                "location": "http://www.baidu.com",
            }]
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.candidates[i].ad_code   | `string` | 广告位code           |
| data.candidates[i].ad_id | `string` | 广告ID |
| data.candidates[i].ad_channel   | `string` | 广告渠道           |
| data.candidates[i].width | `int` | 宽长比(宽) |
| data.candidates[i].height   | `int` |  宽长比(长)           |
| data.candidates[i].cool_down | `int` | 广告展示冷却时间 |
| data.candidates[i].count_down | `int` | 广告倒计时 |
| data.candidates[i].image_url | `string` | 封面图地址 |
| data.candidates[i].location | `string` | 跳转链接 |
| data.candidates[i].action | `int` | 广告类型 0 广告 1 原生页面 2 webView|

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 检测人脸年龄 API

检测人脸年龄

### 请求

```
POST /api/v1/face/detect HTTP/1.1
Content-Type: application/json

{
    "image": "<base64 string>"
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| image | `string` | 是   | 图片oss地址 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "age": 10.5
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.age   | `float` | 年龄           |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 修改人脸 API

修改人脸

### 请求

```
POST /api/v1/face/edit_attr HTTP/1.1
Content-Type: application/json

{
    "image": "<base64 string>",
    "action_type": "TO_KID"
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| image | `string` | 是   | 图片oss地址 |
| action_type | `string` | 是   | 可选值 TO_KID; TO_OLD; TO_FEMALE; TO_MALE |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "image": "<base64 string>"
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的base64字符串          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 修改图片风格 API

修改图片风格

### 请求

```
POST /api/v1/image_process/style_trans HTTP/1.1
Content-Type: application/json

{
    "image": "<base64 string>",
    "option": "cartoon"
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| image | `string` | 是   | 图片oss地址 |
| option | `string` | 是   |cartoon：卡通画风格; pencil：铅笔风格; color_pencil：彩色铅笔画风格; warm：彩色糖块油画风格; wave：神奈川冲浪里油画风格; lavender：薰衣草油画风格; mononoke：奇异油画风格; scream：呐喊油画风格; gothic：哥特油画风格 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "image": "<base64 string>"
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的base64字符串          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |



## 人脸动漫化 API

人脸动漫化

### 请求

```
POST /api/v1/image_process/selie_anime HTTP/1.1
Content-Type: application/json

{
    "image": "<base64 string>",
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| image | `string` | 是   | 图片oss地址 |

### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "image": "<base64 string>"
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的base64字符串          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## OSS上传签名 API

OSS上传签名

### 请求

```
POST /api/v1/upload_signature HTTP/1.1
Content-Type: application/json

{}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |

### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "end_point":"<endpoint>",
        "access_key_id": "ak",	
        "access_key_secret": "sk",	
        "bucket_name": "bn",	
        "expiration": 3600,	
        "security_token": "stk",	
        "path": "/",
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的base64字符串          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 人脸融合 API

人脸融合

### 请求

```
POST /api/v1/face/merge HTTP/1.1
Content-Type: application/json

{
    "template_image": "<url>",
    "target_image": "<url>",
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| template_image | `string` | 是   | 模板图片oss地址 |
| target_image | `string` | 是   | 目标图片oss地址 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "image": "<url>"
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的oss地址          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 人像分割 API

人像分割

### 请求

```
POST /api/v1/image_process/body_seg HTTP/1.1
Content-Type: application/json

{
    "image": "<url>",
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| image | `string` | 是   | 图片oss地址 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "image": "<url>"
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.image   | `string` |  处理后图片的oss地址          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |


## 获取模板图 API

获取模板图

### 请求

```
GET /api/v1/config/template?type=<xxx> HTTP/1.1
Content-Type: application/json

```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| type | `string` | 是   | 场景类型 |


### 返回

```
200 OK / 202 Accepted
Content-Type: application/json

{
    "code": "<error code>",
    "message": "<error message>",
    "data": {
        "templates": [
        {
            "name":<xxx>,
            "image":<url>
        },
        ...]
    }
}
```

#### 结果说明

| 名称    | 类型     | 描述                    |
| :------ | :------- | :---------------------- |
| data.templates   | `array` | 模板信息          |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |