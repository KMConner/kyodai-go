package main

import (
	"bufio"
	"github.com/KMConner/kyodai-go/internal"
	"github.com/KMConner/kyodai-go/kulasis"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

type loginOptions struct {
	defaultOptions
	PrintToConsole bool `hide:"true" short:"o"`
}

func (opts *loginOptions) Execute(_ []string) error {
	reader := bufio.NewReader(os.Stdin)
	print("Enter ECS ID:")
	id, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	print("Enter Password:")
	bpass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	pass := string(bpass)
	println()

	info, err := kulasis.SignIn(id, pass)
	if err != nil {
		return err
	}

	if opts.PrintToConsole {
		println("AccountNo: " + info.Account)
		println("Access Token: " + info.AccessToken)
	} else {
		err = internal.Store(*info)
		return err
	}
	return nil
}
