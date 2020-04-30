package main

import (
	"github.com/KMConner/kyodai-go/internal/auth"
	"github.com/KMConner/kyodai-go/internal/kulasis"
	"os"
	"time"
)

func main() {
	account := os.Getenv("KULASIS_ID")
	token := os.Getenv("KULASIS_TOKEN")
	info := auth.Info{
		AccessToken: token,
		Account:     account,
	}
	slot, _ := kulasis.RetrieveTimeSlot(info)

	lec :=slot.GetLecture(kulasis.DayPeriod{
		Semester: kulasis.First,
		Day:      time.Monday,
		Period:   5,
	})
	print(lec.LectureName)
}
