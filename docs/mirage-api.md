# API文档说明

## 调用接口说明
- 统一采用POST + JSON 方式请求
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
POST /api/v1/config/advertise HTTP/1.1
Content-Type: application/json

{
    "ad_code": "xxxx",
    "extra":"xxxx"
}
```

#### 参数说明

| 名称 | 类型     | 必选 | 描述                 |
| :--- | :------- | :--- | :------------------- |
| ad_code | `string` | 是   | 广告位code |
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
| data.candidates[i].ad_id | `string` | 微信用户唯一标识 OpenID |

#### 错误码

| 错误码 | 说明  |
| ------ | ----- |
| xxx    | ooooo |

