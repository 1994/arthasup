package main

import (
	"1994.arthasup/common"
)

func main() {
	common.Init()
	execute()
}

func execute() {
	// start := time.Now()
	common.Pre()
	name, err := common.Download()

	// cost := time.Since(start)

	// fmt.Printf("download cost: %v,end at: %v\n", cost, time.Now())
	if err != nil && name == "" {
		panic(err)
	}
	version, _ := common.Unzip(name)
	common.Alias(version)
}
