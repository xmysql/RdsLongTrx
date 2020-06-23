package main

import (
	_ "encoding/json"
	"gopkg.in/alecthomas/kingpin.v2"
	"moniter/DBInstancePerformance"
	"moniter/mail"
	"time"
)



var (
	accessKeyId=kingpin.Flag("accessKeyId","阿里云控制台获取accessKeyId").Default("").String()
	accessKeySecret=kingpin.Flag("accessKeySecret","阿里云控制台获取accessKeySecret").Default("").String()
	regionId=kingpin.Flag("regionId","样例:cn-beijing").Default("").String()
	DBInstance = kingpin.Flag("DBInstance", "链接数据库的IP:端口，多个用逗号隔开，样例: 192.168.10.1:3306,192.168.10.2:3306").Default("").String()
	QueryTimes=kingpin.Flag("QueryTimes","要过滤show full processlist里边中time大于该时间的进程").Default("").String()
	emailRecivers = kingpin.Flag("emailRecivers", "邮件接收方,可多个用逗号隔开").Default("").String()
	emailServerHost = kingpin.Flag("email.serverHost", "邮箱服务器地址").Default("").String()
	emailServerPort = kingpin.Flag("email.serverPort", "邮箱服务器端口").Default("0").Int()
	fromEmail = kingpin.Flag("email.from", "邮件发送方").Default("").String()
	fromPassword = kingpin.Flag("email.password", "邮件发送方密码").Default("").String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Version("1.0")
	kingpin.Parse()
	var i int
	db:= *DBInstance
	DbinstanceCount := len(mail.Get_Dbinstance(db))
	dbs := mail.Get_Dbinstance(db)
	if DbinstanceCount >= 1 {
		for  i = 0; i < DbinstanceCount ;i++{
			p := DBInstancePerformance.InitPerformancePath()
			if p.AegionId == "" {
				p.AegionId = *regionId
			}
			if p.AccessKeyId == "" {
				p.AccessKeyId = *accessKeyId
			}
			if p.AccessKeySecret == "" {
				p.AccessKeySecret = *accessKeySecret
			}
			if p.DBInstance == "" {
				p.DBInstance = dbs[i]
			}
			if p.QueryTimes == "" {
				p.QueryTimes = *QueryTimes
			}
			if p.EmailRecivers == "" {
				p.EmailRecivers = *emailRecivers
			}
			if p.FromEmail == "" {
				p.FromEmail = *fromEmail
			}
			if p.FromPassword == "" {
				p.FromPassword = *fromPassword
			}
			if p.EmailServerPort == 0 {
				p.EmailServerPort = *emailServerPort
			}
			if p.EmailServerHost == "" {
				p.EmailServerHost = *emailServerHost
			}
			DBInstancePerformance.SendLongTrxMail(p)
			time.Sleep(5 * time.Millisecond)
		}
	}

}

