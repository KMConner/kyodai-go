package main

import (
	"fmt"
	"github.com/KMConner/kyodai-go/kulasis"
	"time"
)

type timeslotOptions struct {
	defaultOptions
	Term string `short:"s" long:"semester" choice:"first" choice:"second" required:"true"`
}

func (opt *timeslotOptions) Execute(_ []string) error {
	authInfo := opt.GetInfo()
	timeSlot, err := kulasis.RetrieveTimeSlot(authInfo)
	if err != nil {
		return err
	}
	var semester kulasis.Semester
	if opt.Term == "first" {
		semester = kulasis.First
	} else {
		semester = kulasis.Second
	}
	for d := 1; d <= 5; d++ {
		for p := 1; p <= 5; p++ {
			dp := kulasis.DayPeriod{
				Semester: semester,
				Day:      time.Weekday(d),
				Period:   p,
			}
			lecture := timeSlot.GetLecture(dp)
			if lecture != nil {
				fmt.Printf("[%s %d] %s\n", dp.Day.String(), dp.Period, lecture.LectureName)
			}
		}
	}
	return nil
}
