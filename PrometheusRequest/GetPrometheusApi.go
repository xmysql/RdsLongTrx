package PrometheusRequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	PrometheusCpuJson "moniter/PrometheusApiJson"
	"net/http"
)

type memcpuinfo struct {
	cpu interface{}
	mem interface{}
}

func PromethuesHttpGet() (result []memcpuinfo) {

	cpumemrequest := make([]memcpuinfo,1)
	cpuurl := "http://10.162.175.1:9090/api/v1/query?query=100%20-%20(avg(irate(node_cpu_seconds_total%7Binstance%3D%22"+"10.162.175.1"+"%3A9100%22%2Cmode%3D%22idle%22%7D%5B1m%5D))%20*%20100)"
	cpu, err := http.Get(cpuurl)
	if err != nil{
		fmt.Println(err)
	}
	defer cpu.Body.Close()
	cpurequest, _ := ioutil.ReadAll(cpu.Body)
	cpujson := PrometheusCpuJson.CpuInfo{}
	jsonRecord := json.Unmarshal([]byte(cpurequest),&cpujson)
	if jsonRecord != nil {
		panic(jsonRecord)
	}
	var cpudate  interface{}
	for _, v := range cpujson.Data.Result{
		cpudate= v.Value
	}

	memurl :="http://10.162.175.1:9090/api/v1/query?query=(1%20-%20(node_memory_MemAvailable_bytes%7Binstance%3D%22"+"10.162.175.1"+"%3A9100%22%7D%20%2F%20(node_memory_MemTotal_bytes%7Binstance%3D%22"+"10.162.175.1"+"%3A9100%22%7D)))*%20100"
	mem, err := http.Get(memurl)
	if err != nil{
		fmt.Println(err)
	}
	defer mem.Body.Close()
	memrequest, _ := ioutil.ReadAll(mem.Body)
	memjson := PrometheusCpuJson.CpuInfo{}
	memjsonRecord := json.Unmarshal([]byte(memrequest),&memjson)
	if memjsonRecord != nil {
		panic(memjsonRecord)
	}
	var memdate  interface{}
	for _, v := range memjson.Data.Result{
		memdate= v.Value
	}
	cpumemrequest[0].cpu=cpudate
	cpumemrequest[0].mem=memdate

	return  cpumemrequest

}