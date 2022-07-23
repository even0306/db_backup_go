package common

import (
	"fmt"
)

func TypeAssertion(tp any) (string, error) {
	switch tp.(type) {
	case int:
		return "int", nil
	case string:
		return "string", nil
	case bool:
		return "bool", nil
	default:
		return "", fmt.Errorf("无法识别的类型")
	}
}
