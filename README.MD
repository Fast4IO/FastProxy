# FastProxy



## 安装

​	根据您的系统去[发布页]( https://github.com/Fast4IO/fastProxy/releases)选择对应的二进制文件

2. **此文件需要使用7999的端口为出口代理端口，若有开启防火墙，请添加出口**

3. 运行文件

​		windows或linux使用 命令行启动，命令为

```
windows server 
./windows-release.exe -c /path/conf.yaml

linux server 
./linux-release	-c /path/conf.yaml

```



## 使用

### 修改配置文件

在conf.yml 中可修改程序**监听端口**，**认证账密**

```yaml
server:
  host: 0.0.0.0
  port: 7999

proxy:
  username: admin
  password: password123
  
```







## 声明

本项目遵循 [AGPL 3.0](https://github.com/Fast4IO/fastProxy/master/LICENSE)开源，在此基础上，所有使用本项目提供服务者都必须在网站首页保留指向本项目的链接

禁止使用本项目进行营利和做其他违法事情，产生的一切后果本项目概不负责