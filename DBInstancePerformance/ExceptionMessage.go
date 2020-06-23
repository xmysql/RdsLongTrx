package DBInstancePerformance

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"moniter/DBConnect"
	"moniter/mail"
	"strconv"
	"strings"
)
type LongTrxPerformance struct {
	Id   int
	User string
	Host string
	Db   string
	Time string
	Info string
	Ip   string
	Cpu  string
	Mem  string
	Iops string
}

type PerformancePath struct {
	AegionId string
	AccessKeyId string
	AccessKeySecret string
	DBInstance string
	QueryTimes string
	FromEmail string
	FromPassword string
	EmailServerHost string
	EmailServerPort int
	EmailRecivers string
}

func InitPerformancePath() *PerformancePath{
	return &PerformancePath{
		AegionId:"",
		AccessKeyId: "",
		AccessKeySecret: "",
		DBInstance: "",
		QueryTimes: "",
		FromEmail:"",
		FromPassword:"",
		EmailServerHost:"",
		EmailServerPort:0,
		EmailRecivers:"",
	}
}
func SendLongTrxMail(path *PerformancePath)  {
	u := mail.InitNewUser()
	if u.Toers == "" {
		u.Toers = path.EmailRecivers
	}
	if u.FromEmail == "" {
		u.FromEmail = path.FromEmail
	}
	if u.FromPassword == "" {
		u.FromPassword = path.FromPassword
	}
	if u.ServerHost == "" {
		u.ServerHost = path.EmailServerHost
	}
	if u.ServerPort == 0 {
		u.ServerPort = path.EmailServerPort
	}
	mail.InitEmail(u)

	rend_result := new(bytes.Buffer)
	sqlrequestcount := len(DBConnect.GetLongTrx(path.QueryTimes,path.DBInstance))
	sqlrequest := make([]LongTrxPerformance,sqlrequestcount)
	t := make([]PerformanceDate,sqlrequestcount)
	var incr int
	var incv int
	for _, v := range DBConnect.GetLongTrx(path.QueryTimes,path.DBInstance){
		sqlrequest[incv].Ip= path.DBInstance
		sqlrequest[incv].Id=v.Id
		sqlrequest[incv].User =v.User
		sqlrequest[incv].Host = v.Host
		sqlrequest[incv].Db =v.Db
		sqlrequest[incv].Time =v.Time
		sqlrequest[incv].Info =v.Info
		for _ , d := range GetPerformanceRequest(path.AegionId,path.AccessKeyId,path.AccessKeySecret,path.DBInstance){
			t[incr].Value=d.Value
			for _, t := range mail.Format_date(t[incr].Value){
				sqlrequest[incr].Cpu=t.Cpu
				sqlrequest[incr].Mem=t.Mem
			}

		}
		for _,s := range GetPerformanceIops(path.AegionId,path.AccessKeyId,path.AccessKeySecret,path.DBInstance) {
			sqlrequest[incv].Iops=s.Value
		}
		incv++
	}

	LongTrx := template.New("RdsLongTrx.html")
	LongTrx = LongTrx.Funcs(template.FuncMap{"format_id": mail.Format_id})
	LongTrx, err := LongTrx.Parse(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<style type="text/css">
		#hor-minimalist-b
		{
			font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
			font-size: 14px;
			background: #fff;
			margin: 10px;
			width: auto;
			border-collapse: collapse;
			text-align: center;
		}
		#hor-minimalist-b th
		{
			font-size: 14px;
			font-weight: normal;
			color: #039;
			padding: 10px 8px;
			border-bottom: 2px solid #6678b1;
		}
		#hor-minimalist-b td
		{
			border-bottom: 1px solid #ccc;
			color: #669;
			padding: 6px 8px;
		}
		#hor-minimalist-b tbody tr:hover td
		{
			color: #009;
		}
	</style>
</head>
<body>
<table id="hor-minimalist-b" style="table-layout:fixed;word-break:break-all;">
	<thead>
	<tr>
		<th width="3%">序号</th>
        <th width="5.5%">数据库地址</th>
		<th width="5.5%">线程号</th>
		<th width="5.5%">查询用户</th>
		<th>语句详情</th>
		<th width="4.5%">来访地址</th>
		<th width="4.5%">数据库</th>
		<th width="4.5%">执行时间</th>
		<th width="5.5%">当前MEM(%)</th>
		<th width="5.5%">当前CPU(%)</th>
        <th width="5.5%">IOPS</th>
	</tr>
	</thead>
	<tbody>

	{{ range $i,$v := . }}
		<tr>
			<td title="序号" width="3%">{{ $i | format_id }}</td>
            <td title="数据库地址" width="5.5%">{{ .Ip }}</td>
			<td title="线程号" width="5.5%">{{ .Id }}</td>
			<td title="查询用户" width="5.5%">{{ .User }}</td>
			<td title="语句详情" style="text-align: left">{{ .Info }}</td>
			<td title="来访地址" width="9.5%">{{ .Host }}</td>
			<td title="数据库" width="4.5%">{{ .Db }}</td>
			<td title="执行时间" width="4.5%">{{ .Time }}</td>
			{{/*保留一位精度，四舍五入*/}}
			<td title="当前MEM(%)" width="5.5%">{{ .Mem }}</td>
			<td title="当前CPU(%)" width="5.5%">{{ .Cpu }}</td>
            <td title="IOPS" width="5.5%">{{ .Iops }}</td>

		</tr>
	{{ end }}
	</tbody>
</table>
</body>
</html>`)
	if err != nil {
		fmt.Print(err)
	}
	if err := LongTrx.Execute(rend_result, sqlrequest); err !=nil {
		log.Panicf("render to html have errors->%s", err)
	}


	iops,errio:=strconv.ParseFloat(strings.TrimSpace(sqlrequest[0].Iops),64)
	cpu,errcpu:=strconv.ParseFloat(strings.TrimSpace(sqlrequest[0].Cpu),64)
	mem,errmem:=strconv.ParseFloat(strings.TrimSpace(sqlrequest[0].Mem),64)
	if errio != nil {
		fmt.Println(errio)
	}
	if errcpu != nil {
		fmt.Println(errcpu)
	}
	if errmem != nil {
		fmt.Println(errmem)
	}
	if iops >= 12000.00 || cpu >= 90.00 || mem > 90.00 {
		err := mail.SendEmail("LongTrx异常", rend_result.String(),"")
		if err != nil {
			fmt.Println(err)
		}

	}else {
		fmt.Println("数据库资源正常")
	}


}