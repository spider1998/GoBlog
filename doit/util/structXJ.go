package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//结构体转XML
func StructToMapXML(s interface{}) map[string]string {
	return StructToMap(s, "xml")
}

//结构体转JSON
func StructToMapJSON(s interface{}) map[string]string {
	return StructToMap(s, "json")
}

func StructToMap(s interface{}, tag string) map[string]string {
	params := make(map[string]string)
	v := reflect.Indirect(reflect.ValueOf(s))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		isUnexported := ft.PkgPath != ""
		if isUnexported {
			continue
		}
		fv := v.Field(i)
		switch ft.Type.Kind() {
		case reflect.Struct:
			for k, v := range StructToMapXML(fv.Interface()) {
				params[k] = v
			}
		case reflect.String:
			tags := strings.Split(ft.Tag.Get(tag), ",")
			if len(tags) > 0 && tags[0] != "-" {
				params[tags[0]] = v.Field(i).String()
			}
		case reflect.Int:
			tags := strings.Split(ft.Tag.Get(tag), ",")
			if len(tags) > 0 && tags[0] != "-" {
				params[tags[0]] = strconv.Itoa(int(v.Field(i).Int()))
			}
		default:
			panic(fmt.Sprintf("invalid type \"%s\" of field \"%s\" in struct \"%s\".", ft.Type.Kind(), ft.Name, t.Kind()))
		}
	}
	return params
}
