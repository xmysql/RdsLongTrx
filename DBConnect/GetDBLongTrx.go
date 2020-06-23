package DBConnect

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type LongTrx struct {
	Id   int
	User string
	Host string
	Db   string
	Time string
	Info string
}

func GetLongTrx(sqltime string,dbinstance string) (result []LongTrx){
	dbpath:= fmt.Sprintf("exporter:123456@tcp(%s"+")/mysql?charset=utf8",dbinstance)
	//打开数据库
	//DSN数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
	db, err := sql.Open("mysql", dbpath)
	db.SetConnMaxLifetime(30)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	if err != nil {
		fmt.Println(err)
	}
	//关闭数据库，db会被多个goroutine共享，可以不调用
	defer db.Close()
	var command string
	command="Query"
	var querytime int
	querytime, err = strconv.Atoi(sqltime)
	if err != nil {
		fmt.Println(err)
	}
	sql:=fmt.Sprintf("SELECT ID,user,host,db,time,info from INFORMATION_SCHEMA.processlist where command='%s'"+
		" and time>=%d",command,querytime)
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, err := db.Query(sql)
	var ids []LongTrx
	if err == nil {
		for rows.Next() {
			var id LongTrx
			err := rows.Scan(&id.Id,&id.User,&id.Host,&id.Db,&id.Time,&id.Info)
			if err != nil {
				fmt.Println("Rows fail: ", err)
			}
			ids=append(ids,id)
		}

	}else {
		fmt.Println(err)
	}

	return  ids

}
