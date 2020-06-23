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

type PerformanceDate struct {
	Value string
	data  string
}

type PerformanceDateSlice []PerformanceDate
func (a PerformanceDateSlice) Len() int  {
	return len(a)
}
func (a PerformanceDateSlice) Swap(i,j int)  {
	a[i],a[j] = a[j],a[i]
}
func (a PerformanceDateSlice) Less(i,j int) bool  {
	return a[j].data < a[i].data
}

type LongTrxPath struct {

	Regionid        string  //可用区ID
	AccessKeyid     string  //AccessKey ID
	AccessKeysecret string  //AccessKey Secret
}
func InitLongTrxPath() *LongTrxPath  {

	return &LongTrxPath {
		Regionid:        "",
		AccessKeyid:     "",
		AccessKeysecret: "",
	}
}
func GetPerformanceRequest(regionId string,accessKeyId string,accessKeySecret string,DBInstance string) (result []PerformanceDate) {
	var StartTime string
	var EndTime string

	EndTime=time.Now().UTC().Format("2006-01-02T15:04Z")
	StartTime=time.Now().UTC().Add(-time.Minute*10).Format("2006-01-02T15:04Z")
	client, err := rds.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)

	request := rds.CreateDescribeDBInstancePerformanceRequest()
	request.Scheme = "https"

	request.DBInstanceId = mail.Format_Dbinstance(DBInstance)
	request.Key = "MySQL_MemCpuUsage"
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

	performancedate := make([]PerformanceDate,1000)

	var inct int
	for _, d := range values{

		if performancedate[inct].Value == " " {
			continue
		}

		performancedate[inct].Value =d.Value
		performancedate[inct].data=d.Date
		inct++
	}
	sort.Sort(PerformanceDateSlice(performancedate))


	return performancedate[:1]


}
