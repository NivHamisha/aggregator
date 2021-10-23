package utils

import (
	"fmt"
	"reflect"
)

type Expense struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	MoneySpent float64 `json:"moneySpent"`
}

func (e *Expense) GetField(field string) (interface{}, reflect.Type){
	field = ModifyParam(field)
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field).String()
	if f == "<float64 Value>"{
		var floatTest float64
		return reflect.Indirect(r).FieldByName(field).Float(), reflect.TypeOf(floatTest)
	}
	var stringTest string
	return f, reflect.TypeOf(stringTest)
}

func (e *Expense) IsFieldExists(field string) bool{
	v := reflect.ValueOf(*e)
	typeOfS := v.Type()
	field = ModifyParam(field)

	for i := 0; i< v.NumField(); i++ {
		if field == typeOfS.Field(i).Name{
			return true
		}
	}
	return false
}

func (e *Expense)CheckParamsExists(params []string) (bool, string){
	for _, param := range params{
		if !e.IsFieldExists(param){
			return false, param
		}
	}
	return true, ""
}

func (e *Expense) GetPrintParamsString(params []string)string{
	toPrint := fmt.Sprintf("{")

	for _, param := range params{
		value,valueType := e.GetField(param)
		var floatTest float64
		if valueType == reflect.TypeOf(floatTest){
			toPrint = toPrint + fmt.Sprintf("\"%s\":%.2f, ", param, value)
		}else {
			toPrint = toPrint + fmt.Sprintf("\"%s\":\"%s\", ", param, value)
		}
	}
	return toPrint[:len(toPrint)-2] + fmt.Sprintf("}")
}
