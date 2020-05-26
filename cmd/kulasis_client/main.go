package main

import (
	"github.com/KMConner/kyodai-go/internal/auth"
	"github.com/jessevdk/go-flags"
)

type defaultOptions struct {
	AccountId string `short:"a" long:"account" required:"true" env:"ACCOUNT_ID"`
	Token     string `short:"t" long:"token" required:"true" env:"ACCESS_TOKEN"`
}

func (opt *defaultOptions) GetInfo() auth.Info {
	authInfo := auth.Info{
		AccessToken: opt.Token,
		Account:     opt.AccountId,
	}
	return authInfo
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
