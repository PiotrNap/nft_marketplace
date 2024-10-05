package services

import (
	"fmt"
)

func HelloService(s interface{}) (string, error) {
    if s != s.(string) {
        return "", fmt.Errorf("Invalid type for HelloService : %s", s)
    }

    return "Work done!", nil
}
