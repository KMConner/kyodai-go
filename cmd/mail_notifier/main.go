package main

import (
	"github.com/KMConner/kyodai-go/kulasis"
	"os"
)

func main() {
	account := os.Getenv("KULASIS_ID")
	token := os.Getenv("KULASIS_TOKEN")
	info := kulasis.Info{
		AccessToken: token,
		Account:     account,
	}
	slot, err := kulasis.RetrieveTimeSlot(info)

	if err != nil {
		println(err.Error())
		return
	}

	news := slot.GetNewLecture()
	for _, v := range news {
		println(v.LectureName)
		m, _ := v.GetCourseMailTitles()
		println(len(*m))
	}
}
