package lib

import "log"

var (
	LogWarning *log.Logger
	LogInfo    *log.Logger
	LogError   *log.Logger
)

func IfErrInfo(s string, e error) {
	if e != nil {
		LogInfo.Println(s, e)
	}
}

func IfErrWarn(s string, e error) {
	if e != nil {
		LogWarning.Println(s, e)
	}
}

func IfErrError(s string, e error) {
	if e != nil {
		LogError.Println(s, e)
	}
}

func IfErrFatal(s string, e error) {
	if e != nil {
		LogError.Fatalln(s, e)
	}
}
