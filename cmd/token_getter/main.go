package main

import (
	"bufio"
	"github.com/KMConner/kyodai-go/internal"
	"github.com/KMConner/kyodai-go/kulasis"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func main() {
	printToConsole := len(os.Args) == 1 && os.Args[0] == "-o"
	reader := bufio.NewReader(os.Stdin)
	print("Enter ECS ID:")
	id, err := reader.ReadString('\n')
	if err != nil {
		println(err.Error())
		return
	}

	print("Enter Password:")
	bpass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		println(err.Error())
		return
	}

	pass := string(bpass)
	println()

	info, err := kulasis.SignIn(id, pass)
	if err != nil {
		println(err.Error())
		return
	}

	if printToConsole {
		println("AccountNo: " + info.Account)
		println("Access Token: " + info.AccessToken)
	} else {
		err = internal.Store(*info)
		if err != nil {
			println(err.Error())
		}
	}
}
