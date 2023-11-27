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
#### 安装
- 将config.toml 和 leigodhelper.exe 放置在目录leigodhelper
- 编辑config.toml, 填入 username,password,games
- 用管理员身份打开powershell/cmd 执行以下命令
```
cd leigodhelper
./leigodhelper.exe install
```
- leigodhelper service会自动启动
- 编辑config.toml会自动生效
- 安装路径为C:\Program Files\LeigodHelper

### 卸载
- 用管理员身份打开powershell/cmd
- ./leigodhelper.exe remove