package initialization

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	path := "../config.yaml"
	conf := GetTestConfig(path)
	fmt.Println(conf)
}