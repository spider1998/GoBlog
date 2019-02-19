package util

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Pair struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

type Pairs []Pair

func (p Pairs) Map() map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, len(p))
	for _, item := range p {
		m[item.Key] = item.Value
	}
	return m
}

func (p Pairs) Keys() []interface{} {
	m := make([]interface{}, len(p))
	for i, item := range p {
		m[i] = item.Key
	}
	return m
}

func (p Pairs) Values() []interface{} {
	m := make([]interface{}, len(p))
	for i, item := range p {
		m[i] = item.Value
	}
	return m
}

type sliceType int8

const (
	sliceInt sliceType = iota
	sliceUInt
	sliceString
)

// Map 可以提取一个结构体切片某个 int 或 string 字段。
// a := []struct{Foo int}{{1}, {2}}
// Map(a, "Foo") -> []int{1, 2}
func Map(a interface{}, field string) interface{} {
	T := reflect.TypeOf(a)
	if T.Kind() != reflect.Slice {
		panic("only slice allowed in map")
	}
	V := reflect.ValueOf(a)
	elemT := T.Elem()
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}
	elem := reflect.Zero(elemT)
	if elem.Kind() != reflect.Struct {
		panic("only slice of struct type allowed")
	}
	kind := elem.FieldByName(field).Kind()

	var (
		ret interface{}
		t   sliceType
	)
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret = make([]int, V.Len())
		t = sliceInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret = make([]int, V.Len())
		t = sliceUInt
	case reflect.String:
		ret = make([]string, V.Len())
		t = sliceString
	default:
		panic("invalid type " + kind.String())
	}

	for i := 0; i < V.Len(); i++ {
		elem := reflect.Indirect(V.Index(i)).FieldByName(field)
		if elem.Kind() != kind {
			panic(kind.String() + " and " + elem.Kind().String() + " is not same")
		}
		switch t {
		case sliceInt:
			ret.([]int)[i] = int(elem.Int())
		case sliceUInt:
			ret.([]int)[i] = int(elem.Uint())
		case sliceString:
			ret.([]string)[i] = elem.String()
		}
	}
	return ret
}

//按特定字段映射成map接口类型
func MapByKey(a interface{}, field string) interface{} {
	T := reflect.TypeOf(a)         //a的类型
	if T.Kind() != reflect.Slice { //判断是否为切片
		panic("only slice allowed in map")
	}
	V := reflect.ValueOf(a)
	elemT := T.Elem() //a中元素的种类
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}
	elem := reflect.Zero(elemT) //初始化结构体
	if elem.Kind() != reflect.Struct {
		panic("only slice of struct type allowed")
	}
	kind := elem.FieldByName(field).Kind()

	var (
		ret interface{}
		t   sliceType
	)
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ret = make(map[int]interface{}, V.Len())
		t = sliceInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ret = make(map[int]interface{}, V.Len())
		t = sliceUInt
	case reflect.String:
		ret = make(map[string]interface{}, V.Len())
		t = sliceString
	default:
		panic("invalid type " + kind.String())
	}

	for i := 0; i < V.Len(); i++ {
		elem := reflect.Indirect(V.Index(i)).FieldByName(field)
		if elem.Kind() != kind {
			panic(kind.String() + " and " + elem.Kind().String() + " is not same")
		}
		switch t {
		case sliceInt:
			ret.(map[int]interface{})[int(elem.Int())] = V.Index(i).Interface()
		case sliceUInt:
			ret.(map[int]interface{})[int(elem.Uint())] = V.Index(i).Interface()
		case sliceString:
			ret.(map[string]interface{})[elem.String()] = V.Index(i).Interface()
		}
	}
	return ret
}

//转换为切片接口类型
// SliceAnyToSliceInterface convert slice []T to []interface{}
func SliceAnyToSliceInterface(from interface{}) (to []interface{}) {
	switch reflect.TypeOf(from).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(from)
		to = make([]interface{}, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			to = append(to, s.Index(i).Interface())
		}
	default:
		panic("only slice allowed to be converted")
	}
	return
}

// SliceIntToSliceString convert slice []int to []string
func SliceIntToSliceString(from []int) (to []string) {
	for _, i := range from {
		to = append(to, strconv.Itoa(i))
	}
	return
}

// SliceIntToSliceString convert slice []int to []string
func SliceStringToSliceInt(from []string) (to []int) {
	for _, i := range from {
		str, err := strconv.Atoi(i)
		if err == nil {
			to = append(to, str)
		}
	}
	return
}

//字段映射
func MapFields(a interface{}, mappers ...func(string) string) map[string]interface{} {
	V := reflect.ValueOf(a)
	V = reflect.Indirect(V)
	T := reflect.TypeOf(a)
	if V.Kind() != reflect.Struct {
		panic("only struct type supported")
	}

	m := make(map[string]interface{})
	for i := 0; i < V.NumField(); i++ {
		if len(mappers) > 0 {
			m[mappers[0](T.Field(i).Name)] = V.Field(i).Interface()
		} else {
			m[T.Field(i).Name] = V.Field(i).Interface()
		}
	}
	return m
}

//将空字段容量收回
func FillNullSlices(a interface{}) {
	V := reflect.ValueOf(a)
	if V.Type().Kind() != reflect.Ptr {
		panic("only pointer supported")
	}
	V = reflect.Indirect(V)
	T := V.Type()
	if V.Kind() != reflect.Struct {
		panic("only struct type supported")
	}

	//遍历字段
	for i := 0; i < V.NumField(); i++ {
		if T.Field(i).Type.Kind() == reflect.Slice {
			VF := V.Field(i)
			if VF.IsNil() {
				VF.Set(reflect.MakeSlice(VF.Type(), 0, 0))
			}
		}
	}
}

//身份证正则匹配
var identityRegexp = regexp.MustCompile(`^\d{6}(\d{4})(\d{2})(\d{2})`)

//解析身份证年龄
func ParseAgeFromIdentity(identity string) (age int, err error) {
	if len(identity) == 0 {
		err = errors.New("invalid identity")
		return
	}
	matches := identityRegexp.FindStringSubmatch(identity)
	if len(matches) != 4 {
		err = fmt.Errorf("invalid identity: %s", identity)
		return
	}

	t, err := time.ParseInLocation("20060102", strings.Join(matches[1:], ""), time.Local)
	if err != nil {
		return
	}

	age = int(time.Now().Sub(t) / (time.Hour * 24 * 365))
	return
}

type M map[string]interface{}
