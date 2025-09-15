package utils

import (
	"reflect"
	"strconv"
)

// String2Int 将字符串数组转换为整数数组
func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))
	for index, val := range strArr {
		res[index], _ = strconv.Atoi(val)
	}
	return res
}

// Contain 判断 obj 是否在 target 中
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// FilterRules removes (zeroes) rule items whose Level=="deep" when reportType is "basic".
// It walks through first-level nested structs via reflection.
// For simplicity we only support structs that embed a field named "Level string".
// If struct has Level=="deep", it will be zeroed so later analysis treats it as absent.
func FilterRules(rule interface{}, reportType string) {
	if reportType != "basic" {
		return
	}
	rv := reflect.ValueOf(rule).Elem()
	filterStruct(rv)
}

func filterStruct(rv reflect.Value) {
	if rv.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		if !fv.CanSet() {
			continue
		}
		switch fv.Kind() {
		case reflect.Struct:
			// check Level field inside
			lf := fv.FieldByName("Level")
			if lf.IsValid() && lf.Kind() == reflect.String {
				if lf.String() == "deep" {
					fv.Set(reflect.Zero(fv.Type()))
					continue
				}
			}
			filterStruct(fv)
		}
	}
}

// GetNmLevelMap builds a map[nm]level from the loaded rules for quick filtering at presentation.
func GetNmLevelMap() (map[string]string, error) {
	rules, err := GetRule()
	if err != nil {
		return nil, err
	}
	nmToLevel := make(map[string]string)
	var collectNmLevel func(v reflect.Value)
	collectNmLevel = func(v reflect.Value) {
		if v.Kind() != reflect.Struct {
			return
		}
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Struct {
				// try to read Nm and Level fields
				nmField := fv.FieldByName("Nm")
				lvlField := fv.FieldByName("Level")
				if nmField.IsValid() && nmField.Kind() == reflect.String && lvlField.IsValid() && lvlField.Kind() == reflect.String {
					nm := nmField.String()
					if nm != "" {
						nmToLevel[nm] = lvlField.String()
					}
				}
				// dive deeper
				collectNmLevel(fv)
			}
		}
	}
	rv := reflect.ValueOf(rules).Elem()
	collectNmLevel(rv)
	return nmToLevel, nil
}
