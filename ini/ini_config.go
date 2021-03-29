package ini

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	MysqlConfig  `ini:"mysql"`
	ServerConfig `ini:"server"`
}

type MysqlConfig struct {
	Host     string  `ini:"host"`
	Username string  `ini:"username"`
	Password string  `ini:"password"`
	Port     int     `ini:"port"`
	Timeout  float64 `ini:"timeout"`
}

type ServerConfig struct {
	Ip   string `ini:"ip"`
	Port int    `ini:"port"`
}

func UnMarshal(data []byte, res interface{}) (err error) {
	// res为映射结构体指针
	resTypeInfo := reflect.TypeOf(res)
	if resTypeInfo.Kind() != reflect.Ptr {
		err = errors.New("res is not a ptr")
		return
	}

	resStruct := resTypeInfo.Elem()
	if resStruct.Kind() != reflect.Struct {
		err = errors.New("res is not a struct")
		return
	}

	configArr := strings.Split(string(data), "\n")
	var sectionName string

	for k, v := range configArr {
		// 对字符串进行前后去除空格
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if v[0] == '#' {
			continue
		}

		// section
		if v[0] == '[' {
			sectionName, err = parseSection(v, k, resStruct)
			if err != nil {
				return
			}
			continue
		}

		// item
		if strings.Contains(v, "=") {
			err = parseItem(v, k, sectionName, res)
			if err != nil {
				return
			}
		}
	}

	return
}

// parse section in ini
func parseSection(v string, k int, resStruct reflect.Type) (fieldName string, err error) {
	if !strings.Contains(v, "]") || len(v) <= 2 {
		err = fmt.Errorf("error section in line %d", k+1)
		return
	}
	sectionName := strings.TrimSpace(v[1 : len(v)-1])
	if len(sectionName) == 0 {
		err = fmt.Errorf("error section in line %d", k+1)
		return
	}

	for i := 0; i < resStruct.NumField(); i++ {
		field := resStruct.Field(i)
		tagName := field.Tag.Get("ini")
		if tagName == sectionName {
			fieldName = field.Name
			break
		}
	}

	return
}

// parse item in ini
func parseItem(v string, k int, sectionName string, res interface{}) (err error) {
	v = strings.Replace(v, " ", "", -1)
	itemArr := strings.Split(v, "=")
	// 获取到ini配置项的key和value
	key := itemArr[0]
	value := itemArr[1]
	if len(key) == 0 {
		err = fmt.Errorf("error config in line %d", k+1)
		return
	}

	resValue := reflect.ValueOf(res)
	sectionValue := resValue.Elem().FieldByName(sectionName)
	sectionType := sectionValue.Type()
	if sectionType.Kind() != reflect.Struct {
		err = fmt.Errorf("field: %s must be a struct", sectionName)
		return
	}

	for i := 0; i < sectionType.NumField(); i++ {
		structField := sectionType.Field(i)
		fieldValue := sectionValue.Field(i)
		tagName := structField.Tag.Get("ini")
		// 如果找到了结构体中对应的字段 反射给其赋值
		if tagName == key {
			switch structField.Type.Kind() {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error config item in line %d (need int)", k+1)
				}
				fieldValue.SetInt(int64(intValue))
			case reflect.Float32, reflect.Float64:
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("error config item in line %d", k+1)
				}
				fieldValue.SetFloat(floatValue)
			default:
				err = fmt.Errorf("error type: %v", structField.Type.Kind())
			}

			break
		}
	}

	return
}
