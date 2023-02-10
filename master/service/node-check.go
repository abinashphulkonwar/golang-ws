package service

import (
	"github.com/jasonlvhit/gocron"
)

func CheckNode() {

	job := gocron.NewScheduler()

	job.Every(30).Seconds().Do(func() {
		println("Hello")
	})

}
