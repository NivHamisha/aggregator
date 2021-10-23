package utils

type GroupByQuery struct {
	Params   []string  `json:"groupBy"`
}
func (q *GroupByQuery)GetQueryResultFileName()string{
	resultFileName := "result("
	for _,param := range q.Params{
		resultFileName = resultFileName + param + "-"
	}
	return resultFileName[:len(resultFileName)-1] + ").json"
}