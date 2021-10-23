package main

import (
	"aggregator/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	dataDir = "input"
	queryFile = "query.json"
	resultsDir = "results"
)

func main() {
	expenses, err := utils.GetExpensesFromDir(dataDir)
	if err != nil {
		panic(err)
	}

	query, err := utils.ReadQuery(queryFile)
	if err != nil {
		panic(err)
	}

	v := reflect.ValueOf(utils.Expense{})
	typeOfS := v.Type()
	compareFunctions := utils.CompereFunctions{}

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		compareFunctions[fieldName] = func(c1, c2 *utils.Expense)bool{
			value1, valueType:= c1.GetField(fieldName)
			value2, _:= c2.GetField(fieldName)
			switch valueType.Name(){
			case "float64":
				return value1.(float64)< value2.(float64)
			case "int":
				return value1.(int)< value2.(int)
			default:
				return value1.(string)< value2.(string)
			}
		}
	}

	requiredOrderFunctions := func (params []string)([]utils.LessFunc, error){
		var paramFunctions []utils.LessFunc
		for _, param := range params{
			if function, ok := compareFunctions[utils.ModifyParam(param)]; !ok{
				return nil, fmt.Errorf("param is not exists: %s", param)
			}else {
				paramFunctions = append(paramFunctions, function)
			}
		}
		return paramFunctions, nil
	}

	paramOrderFunctions, err := requiredOrderFunctions(query.Params)
	if err != nil {
		panic(err)
	}

	utils.OrderedBy(paramOrderFunctions).Sort(expenses)
	groupedData, err := utils.GroupBy(expenses, query.Params)
	if err != nil {
		panic(err)
	}

	fmt.Println(utils.GetGroupedDataString(groupedData, &query))

	jsonString, _ := json.Marshal(utils.GetJsonStringData(groupedData, &query))
	if err = ioutil.WriteFile(resultsDir + "/" +query.GetQueryResultFileName(), jsonString, os.ModePerm); err != nil{
		panic(err)
	}
}
