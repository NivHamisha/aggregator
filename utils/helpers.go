package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
)

type jsonStringData map[string]interface{}

func GetExpensesFromDir(dirName string) ([]Expense, error) {
	var allExpenses []Expense
	files, err := ioutil.ReadDir(dirName)
	if err != nil{
		return nil, err
	}

	for _, file := range files{
		fullFilePath := dirName + "/" + file.Name()
		fileData, err := ioutil.ReadFile(fullFilePath)
		if err != nil {
			return nil, err
		}

		var fileExpense Expense
		if err = json.Unmarshal(fileData, &fileExpense); err != nil {
			return nil, err
		}
		allExpenses = append(allExpenses, fileExpense)
	}
	return allExpenses, nil
}

func ReadQuery(filename string) (GroupByQuery, error) {
	queryData, err := ioutil.ReadFile(filename)
	if err != nil {
		return GroupByQuery{}, err
	}

	var query GroupByQuery
	if err = json.Unmarshal(queryData, &query); err != nil {
		return GroupByQuery{}, err
	}
	return query, nil
}

func IsExpensesEqual(expense1, expense2 *Expense, params []string) bool {
	for _,param := range params{
		value1, _ := expense1.GetField(param)
		value2, _:= expense2.GetField(param)
		if value1 != value2{
			return false
		}
	}
	return true
}

func ModifyParam(field string) string {
	return strings.ToUpper(field[:1]) + field[1:]
}

func GroupBy(expenses []Expense, params []string)([]Expense, error){
	if expenses == nil{
		return nil, fmt.Errorf("expenses is empty")
	}

	var groupedExpenses []Expense
	lastExpense := expenses[0]

	if exists, param := lastExpense.CheckParamsExists(params); !exists {
		return nil, fmt.Errorf("param is not exists: %s", param)
	}

	for _, expense := range expenses[1:]{
		if IsExpensesEqual(&lastExpense, &expense, params){
			lastExpense.MoneySpent = lastExpense.MoneySpent + expense.MoneySpent
		} else {
			groupedExpenses = append(groupedExpenses, lastExpense)
			lastExpense = expense
		}
	}
	groupedExpenses = append(groupedExpenses, lastExpense)
	return groupedExpenses, nil
}

func GetGroupedDataString(groupedData []Expense, query *GroupByQuery)string{
	dataJsonString := "["
	for _, result:= range groupedData{
		resultString := result.GetPrintParamsString(append(query.Params,"moneySpent"))
		dataJsonString = dataJsonString + resultString + fmt.Sprintf(",\n")
	}
	return dataJsonString[:len(dataJsonString)-2] + fmt.Sprintf("]")
}


func GetJsonStringData(groupedData []Expense, query *GroupByQuery)[]jsonStringData{
	var jsonData []jsonStringData
	newParams := append(query.Params,"moneySpent")
	jsonString := map[string]interface{}{}
	for _, result:= range groupedData{
		for _, param := range newParams {
			value,valueType := result.GetField(param)
			var floatTest float64
			if valueType == reflect.TypeOf(floatTest){
				value = math.Round(value.(float64)*100)/100
			}
			jsonString[param] = value
		}
		jsonData = append(jsonData, jsonString)
		jsonString = map[string]interface{}{}
	}
	return jsonData
}
