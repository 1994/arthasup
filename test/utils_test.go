package test

import (
	"fmt"
	"testing"

	"1994.arthasup/common"
)

func TestProfile(t *testing.T) {
	s := common.Profile()
	fmt.Printf("s: %v\n", s)
}

func TestDownload(t *testing.T) {
	a := make([]string, 1)
	b := append(a, "1")
	fmt.Printf("b: %v\n", b)
}
func Test(t *testing.T) {
	name, err := common.Download()
	if err != nil && name == "" {
		panic(err)
	}
	version, _ := common.Unzip(name)
	common.Alias(version)
}
