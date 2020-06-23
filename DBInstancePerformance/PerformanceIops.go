package DBInstancePerformance

import (
"encoding/json"
_ "encoding/json"
"fmt"
"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"moniter/RdsMoniter"
	"moniter/mail"
	"sort"
"time"
)

type PerformanceIops struct {
	Value string
	Data  string
}

type PerformanceIopsSlice []PerformanceIops
func (a PerformanceIopsSlice) Len() int  {
	return len(a)
}
func (a PerformanceIopsSlice) Swap(i,j int)  {
	a[i],a[j] = a[j],a[i]
}
func (a PerformanceIopsSlice) Less(i,j int) bool  {
	return a[j].Data < a[i].Data
}


func GetPerformanceIops(regionId string,accessKeyId string,accessKeySecret string,DBInstance string) (result []PerformanceIops) {
	var StartTime string
	var EndTime string

	EndTime=time.Now().UTC().Format("2006-01-02T15:04Z")
	StartTime=time.Now().UTC().Add(-time.Minute*10).Format("2006-01-02T15:04Z")
	client, err := rds.NewClientWithAccessKey(regionId,accessKeyId ,accessKeySecret )

	request := rds.CreateDescribeDBInstancePerformanceRequest()
	request.Scheme = "https"

	request.DBInstanceId = mail.Format_Dbinstance(DBInstance)
	request.Key = "MySQL_IOPS"
	request.StartTime = StartTime
	request.EndTime = EndTime

	response, err := client.DescribeDBInstancePerformance(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	Performance := RdsMoniter.Performances{}
	PerformanceValue := json.Unmarshal([]byte(response.GetHttpContentString()),&Performance)

	if PerformanceValue != nil {
		panic(PerformanceValue)
	}


	values := make([]RdsMoniter.Value,1000)




	for _, v := range Performance.PerformanceKeys.PerformanceKey{

		values = v.Values.PerformanceValue


	}

	performanceiops := make([]PerformanceIops,1000)

	var inct int
	for _, d := range values{

		if performanceiops[inct].Value == " " {
			continue
		}

		performanceiops[inct].Value =d.Value
		performanceiops[inct].Data =d.Date
		inct++
	}
	sort.Sort(PerformanceIopsSlice(performanceiops))


	return performanceiops[:1]


}

