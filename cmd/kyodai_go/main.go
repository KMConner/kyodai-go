package main

import (
	"github.com/KMConner/kyodai-go/internal"
	"github.com/KMConner/kyodai-go/kulasis"
	"github.com/jessevdk/go-flags"
)

type defaultOptions struct {
	authInfo *kulasis.Info
}

func (opt *defaultOptions) loadCredential() error {
	info, err := internal.Load()
	if err != nil {
		return err
	}

	opt.authInfo = info
	return nil
}

func main() {
	defaults := defaultOptions{}
	parser := flags.NewParser(&defaults, flags.Default)
	timeslot := timeslotOptions{}
	mail := getMailOptions{}
	login := loginOptions{}
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

	_, e = parser.AddCommand("login", "Login",
		"Login to KULASIS", &login)
	if e != nil {
		println(e.Error())
		return
	}

	_, e = parser.Parse()
	if e != nil {
		return
	}
}
