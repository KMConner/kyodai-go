package main

import (
	"bufio"
	"github.com/KMConner/kyodai-go/internal/auth"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	print("Enter ECS ID:")
	id, err := reader.ReadString('\n')
	if err != nil {
		println(err.Error())
		return
	}

	print("Enter Password")
	bpass, err := terminal.ReadPassword(0)
	if err != nil {
		println(err.Error())
		return
	}

	pass := string(bpass)
	println()

	info, err := auth.SignIn(id, pass)
	if err != nil {
		println(err.Error())
		return
	}

	println("AccountNo: " + info.Account)
	println("Access Token: " + info.AccessToken)
}
