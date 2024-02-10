package utils

import "log"

func Check(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			log.Println("Error: " + msg[0] + " -- " + e.Error())
		} else {
			log.Println("Error: " + e.Error())
		}
		panic(e)
	}
}
