# IM-Backend

## 一、如何运行项目

### 1.服务器

#### 1.1 首先你需要获取镜像
有两种方式可供参考：
+ 先在本地构建镜像，然后再上传至公共镜像仓库（如阿里个人镜像仓库），再从服务器端拉取运行
+ 当然你也可以先从**GitHub**上先拉取代码，然后执行命令`docker build -t im_backend .`构建镜像
#### 1.2 启动nacos、pg、redis
**nacos**是用来管理配置的，尤其重要
**pg**管理数据
**redis**负责管理缓存以及生成ID

当然你也可以按照[docker-compose.yaml](./docker-compose.yaml)（你也可以按你自己的需求修改）里的运行这三个，注意最好不要把本服务先启动了

之后你需要修改nacos的配置（修改前建议了解下nacos的配置管理），在自己的电脑上登录**nacos管理配置网站**(***服务器公网IP:8848/nacos***，如果你是直接运行docker-compose.yaml里的，如果需要登录，用户名和密码均为nacos)，进去后上传配置，即除配置中心外的配置
[参考配置](./configs/app_conf.yaml)在此

#### 1.3 启动
在创建获运行容器前，请修改并确认好连接nacos的配置，[参考配置](./configs/config-example.yaml)在此，确认好docker-compose的挂载目录后正常启动即可

### 2.本地

步骤可以参考[服务器](#12-启动nacospgredis)，只不过最后的运行指令变为了`go run main.go -conf configs/config.yaml`

## 二、关于appKey

对于任何请求都需要appKey，获得它你需要让管理员帮你注册服务（服务名），之后你会得到**secret**（应该为**16个字符的字符串**），然后用其使用AES加密UTC的秒级时间戳，即可得到appKey

参考JavaScript加密代码如下
```
const crypto = require('crypto-js');

// 加密函数：使用 AES CFB 模式
function encryptAES(plaintext, key) {
    // 生成 IV
    const iv = crypto.lib.WordArray.random(16); // AES 块大小是 16 字节
    const keyWordArray = crypto.enc.Utf8.parse(key); // 将密钥转换为 WordArray

    // 使用 CFB 模式加密（不使用填充）
    const encrypted = crypto.AES.encrypt(plaintext, keyWordArray, {
        iv: iv,
        mode: crypto.mode.CFB,
        padding: crypto.pad.NoPadding // 禁用填充
    });

    // 将 IV 和加密数据拼接在一起，并转换为 Hex 编码
    const result = iv.concat(encrypted.ciphertext);
    return result.toString(crypto.enc.Hex);
}


// 获取当前的秒级时间戳
const plaintext = Math.floor(Date.now() / 1000).toString();  // 获取当前秒级时间戳并转换为字符串
console.log("timestamp:",plaintext);


const key = "W7K8pJ3aQv2LcXgH"; // 16 字节的密钥

// 调用加密函数
const encryptedText = encryptAES(plaintext, key);
console.log("Encrypted Hex:", encryptedText);

```

## 三、返回的错误
返回的响应格式(**json**)为
```
{
    "code":200,
    "msg":success,
    "data":null
}
```

默认的成功响应的code为200，msg为"success"，data部分根据请求而变


**错误码**
| 错误码 (Code) | 错误信息 (Message)                          |
|---------------|---------------------------------------------|
| 10001         | the table of the service has not been created |
| 10002         | can not create table of the service          |
| 10003         | create data failed                           |
| 10004         | update data failed                           |
| 10005         | delete data failed                           |
| 10006         | find data failed                             |
| 10007         | count data failed                            |
| 10008         | find query is empty                          |
| 10009         | set k-v in cache failed                      |
| 10010         | sadd val in cache failed                     |
| 10011         | get k-v from cache failed                    |
| 10012         | convert json failed                          |
| 10013         | get set from cache failed                    |
| 10014         | no right record found                        |
| 10015         | cache miss                                   |
| 10016         | update query is empty                        |
| 10017         | generate id failed                           |
| 10018         | del k-v in cache failed                      |
| 504           | Request timeout                              |
| 400           | bind param failed                            |
| 400           | something of params is wrong                 |
| 400           | bind authentication                          |
| 200           | success                                      |


## 四、接口文档

你可将[接口文档](./Muxi-IM.openapi.json)导入apifox中详细查看