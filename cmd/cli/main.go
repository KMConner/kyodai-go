package main

import "github.com/KMConner/kyodai-go/internal/auth"

func main() {
	info, err := auth.SignIn("a0000000", "xxxxxxxxxx")
	if err != nil {
		println(err.Error())
		return
	}

	println(info.AccessToken)
	println(info.Account)
}
