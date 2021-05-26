## 雷神加速器自动暂停工具

### 编译
- 下载并安装Go https://golang.google.cn/dl/go1.16.4.windows-amd64.msi
- 在powershell中 cd到根目录执行下面两条命令
- $env:GOPROXY="https://goproxy.cn,direct"
- go build -ldflags "-H windowsgui"


### 申明

- 本工具仅适用于交流学习，不得用作商业用途
- 如果本项目可能侵犯您的权益，请及时联系我删除，邮箱: yixinin@outlook.com
- 本项目只有源码，没有可执行文件，任何形式的文件传播均与本项目及账号无关

### 特性

- 退出游戏时自动暂停加速器
- 电脑关机时自动暂停加速器
- 长时间没有运行游戏时暂停加速器(每1小时检测一次)


### 使用说明

#### 填写config.toml
- 将雷神账号密码填写在username,password
- 将需要监听的游戏进程名(不要.exe)填写在games, 以英文逗号分隔(默认监听pubg，盗贼之海。)
- 将需要启动的应用程序填写在[start_with], (默认启动雷神加速器和steam)

#### 启动
- 启动leigodhelper.exe 
- leigodhelper.exe 会自动启动雷神加速器，steam等应用程序
- 可以将leigodhelper.exe 调整为以管理员权限运行

