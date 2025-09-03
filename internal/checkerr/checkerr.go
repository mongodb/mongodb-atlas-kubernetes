package checkerr

import "log"

type funcErrs func() error

func CheckErr(msg string, f funcErrs) {
	if err := f(); err != nil {
		log.Printf("%s failed: %v", msg, err)
	}
}