package common

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func TypeAssertion(tp any) {
	b := (*string)(unsafe.Pointer(uintptr(tp)))
	switch b {
	case string:
		if *b == "" {
			log.Panicf("配置文件中字段：%v 不能为空", b)
		}
	case reflect.Int.String():
		if b == nil {
			log.Panicf("配置文件中字段：%v 不能为空", b)
		}
	case reflect.Bool.String():
		if b == nil {
			log.Panicf("配置文件中字段：%v 不能为空", b)
		}
	default:
		fmt.Print("aa")
	}
}
