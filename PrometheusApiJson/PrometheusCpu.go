package PrometheusCpuJson

type CpuInfo struct {
	Status string
	Data CpuResult

}
type CpuResult struct {
	ResultType string
	Result []CpuMetric
}
type CpuMetric struct {
	Metric MetricInfo
	Value interface{}
}
type MetricInfo struct {
	__Name__ string
	Cpu string
	Instance string
	Job string
	Mode string
}




