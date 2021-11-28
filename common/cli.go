package common

import "flag"

var (
	Force bool
	Java  string
)

func Init() {
	flag.BoolVar(&Force, "f", false, "reinstall latest arthas")
	flag.StringVar(&Java, "j", "java", "specify java home")
	flag.Parse()
}
