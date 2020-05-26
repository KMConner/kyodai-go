package main

import (
	"fmt"
	"github.com/KMConner/kyodai-go/internal/auth"
	"github.com/KMConner/kyodai-go/internal/kulasis"
	"github.com/jessevdk/go-flags"
	"time"
)

type timeslotOptions struct {
	defaultOptions
	Term string `short:"s" long:"semester" choice:"first" choice:"second" required:"true"`
}

func (opt *timeslotOptions) Execute(_ []string) error {
	authInfo := auth.Info{
		AccessToken: opt.Token,
		Account:     opt.AccountId,
	}
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

type getMailOptions struct {
	defaultOptions
	GetNew bool `short:"n" long:"long"`
}

type defaultOptions struct {
	AccountId string `short:"a" long:"account" required:"true" env:"ACCOUNT_ID"`
	Token     string `short:"t" long:"token" required:"true" env:"ACCESS_TOKEN"`
}

func main() {
	defaults := defaultOptions{}
	parser := flags.NewParser(&defaults, flags.Default)
	timeslot := timeslotOptions{}
	mail := getMailOptions{}
	_, e := parser.AddCommand("timeslot", "Show timeslot",
		"Print time slot to console", &timeslot)
	if e != nil {
		println(e.Error())
		return
	}

	_, e = parser.AddCommand("mail", "Get mails",
		"Get mail", &mail)
	if e != nil {
		println(e.Error())
		return
	}

	_, e = parser.Parse()
	if e != nil {
		return
	}
}
