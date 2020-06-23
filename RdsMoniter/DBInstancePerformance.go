package RdsMoniter


type Performances struct {
	PerformanceKeys Performance
	RequestId string
	EndTime   string
	DBInstanceId string
	StartTime    string
	Engine    string
}
type Performance struct {
	PerformanceKey []ValueFormat

}
type ValueFormat struct {
	ValueFormat string
	Values PerformanceValue
	Unit string
	Key string
}
type PerformanceValue struct {

	PerformanceValue []Value

}
type Value struct {
	Value string
	Date string
}

