package mail

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type  MemCpu struct {
	Cpu string
	Mem string
}
// 序号从1开始
func Format_id(args ...interface{}) int {
	id, _ := args[0].(int)
	return id+1
}

func Dedumplicate(data interface{}) interface{} {
	inArr := reflect.ValueOf(data)
	if inArr.Kind() != reflect.Slice && inArr.Kind() != reflect.Array {
		return data
	}

	existMap := make(map[interface{}]bool)
	outArr := reflect.MakeSlice(inArr.Type(), 0, inArr.Len())

	for i := 0; i < inArr.Len(); i++ {
		iVal := inArr.Index(i)

		if _, ok := existMap[iVal.Interface()]; !ok {
			outArr = reflect.Append(outArr, inArr.Index(i))
			existMap[iVal.Interface()] = true
		}
	}

	return outArr.Interface()
}

func GetName(params ...interface{}) {
	var paramSlice []string
	for _, param := range params {
		switch v := param.(type) {
		case string:
			fmt.Printf("#%d",v)
		case float64:
			fmt.Printf("#%d",v)
		default:
			panic("params type not supported")
		}

	}

	res := strings.Join(paramSlice, "_")
	fmt.Println(res)
}

// 切分api返回的内存和cpu数据
func Format_date(args ...interface{}) (result []MemCpu) {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}

	substrs := strings.Split(s, "&")
	if len(substrs) != 2 {
		fmt.Println("split err %s",s)
	}

	cpumem := make([]MemCpu,1)
	cpumem[0].Cpu =substrs[0]
	cpumem[0].Mem =substrs[1]
	return cpumem
}


func Get_Dbinstance(dbs string) (result []string){
	Dbinstances := []string{}
	for _,d := range strings.Split(dbs,","){
		Dbinstances=append(Dbinstances,strings.TrimSpace(d))
	}
	return Dbinstances
}
// 返回地址
func Format_Dbinstance(args ...interface{}) (result string) {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}

	substrs := strings.Split(s, ".")
	dbinstance := substrs[0]
	return dbinstance
}




// 判断文件或目录是否存在
func FileOrDirIfExists(binfile string) bool {
	_, err := os.Stat(binfile)
	if err != nil {
		return false
	}
	return true
}

// 执行shell命令
func Exec_shell_cmd(shellcmd string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", shellcmd)
	out,err := cmd.CombinedOutput()
	return string(out), err
}
