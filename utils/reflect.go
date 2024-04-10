package utils

import (
	"errors"
	"reflect"
)

/*
* @Author: zouyx

* @Date:   2024/4/8 15:11
* @Package: 反射相关
 */

// GetFieldValue 获取结构体或结构体指针的fieldName 属性值
func GetFieldValue(subject any, fieldName string) (interface{}, error) {
	valueOfSubject := reflect.ValueOf(subject)
	// 还是先判断一下传入参数的数据类型，如果是指针则进行取值处理
	if valueOfSubject.Kind() == reflect.Ptr {
		for valueOfSubject.Kind() == reflect.Ptr {
			valueOfSubject = valueOfSubject.Elem()
		}
	}
	if valueOfSubject.Kind() != reflect.Struct {
		return nil, errors.New("subject is not a pointer of struct or struct")
	}

	// 如果该属性存在的话，field不是零值
	if field := valueOfSubject.FieldByName(fieldName); !field.IsZero() {
		return field.Interface(), nil
	} else {
		// 如果属性不存在，则直接返回错误
		return nil, errors.New("field: " + fieldName + " not exist in subject")
	}
}
