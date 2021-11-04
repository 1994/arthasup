package main

import "1994.arthasup/common"

func main() {
	name, err := common.Download()
	if err != nil && name == "" {
		panic(err)
	}
	common.Unzip(name)
}
