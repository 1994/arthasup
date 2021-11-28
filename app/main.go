package main

import (
	"fmt"
	"time"

	"1994.arthasup/common"
)

func main() {
	fmt.Printf("start execute at: %v\n", time.Now())
	execute()
}

func execute() {
	// start := time.Now()
	name, err := common.Download()

	// cost := time.Since(start)

	// fmt.Printf("download cost: %v,end at: %v\n", cost, time.Now())
	if err != nil && name == "" {
		panic(err)
	}
	version, _ := common.Unzip(name)
	common.Alias(version)
}
