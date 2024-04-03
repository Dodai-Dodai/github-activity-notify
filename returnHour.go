package main

import "time"

func ReturnHour() int {
	return (time.Now().Hour() + 9)
}
