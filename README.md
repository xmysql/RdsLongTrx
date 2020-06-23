## MySQL慢查询邮件报表
```bash
一个基于Golang html/template开发的
监测阿里云mysql数据库进程信息中
执行时间过长的语句
并关联阿里云监控中的cpu/mem/iops等参数
通过邮件报表的方式发送给相关人员的程序
```

## 特性
```bash
- 检测Rds Mysql进程中执行大于指定时间的语句(通过参数指定)
- 通过调用阿里云API实现和监控中CPU/MEM/IOPS值进行绑定
- 检测时间频率自控(通过crontab)
- 可以通过IOPS/CPU/MEM过滤结果，如果达到某个设置的阀值进行邮件告警(项目设置的阀值:IOPS>12000,CPU>90,MEM>90)
```
### 编译
```bash
- 程序开发的golang版本:go1.14.2,请自行安装
- 程序监控需要创建mysql用户(程序写死)
grant select, replication slave , replication client on *.* to 'export'@'%' identified '123456'
- 编译项目
#go build RdsLongTrxSql.go
```
### 运行
示例
```bash
usage: RdsLongTrxSql [<flags>]

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
      --accessKeyId=""       阿里云控制台获取accessKeyId
      --accessKeySecret=""   阿里云控制台获取accessKeySecret
      --regionId=""          样例:cn-beijing
      --DBInstance=""        链接数据库的IP:端口，多个用逗号隔开，样例: 192.168.10.1:3306,192.168.10.2:3306
      --QueryTimes=""        要过滤show full processlist里边中time大于该时间的进程
      --emailRecivers=""     邮件接收方,可多个用逗号隔开
      --email.serverHost=""  邮箱服务器地址
      --email.serverPort=0   邮箱服务器端口
      --email.from=""        邮件发送方
      --email.password=""    邮件发送方密码
      --version              Show application version.
```
