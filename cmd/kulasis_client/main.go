package main

import "github.com/jessevdk/go-flags"

type timeslotOptions struct {
	Term string `short:"t" long:"term" choice:"first" choice:"second" required:"true"`
}

type getMailOptions struct {
	GetNew bool `short:"n" long:"long"`
}

type defaultOptions struct {
}

func main() {
	defaults := defaultOptions{}
	parser := flags.NewParser(&defaults, flags.Default)
	timeslot := timeslotOptions{}
	mail:=getMailOptions{}
	c, e := parser.AddCommand("timeslot", "Show timeslot",
		"Print time slot to console", &timeslot)
	if e != nil {
		println(e.Error())
		return
	}

	c, e = parser.AddCommand("mail", "Get mails",
		"Get mail", &mail)
	if e != nil {
		println(e.Error())
		return
	}

	ss, e := parser.Parse()
	if e != nil {
		return
	}
	print(c.Hidden)
	print(ss)
}
