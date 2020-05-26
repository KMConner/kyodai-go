package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/KMConner/kyodai-go/internal/kulasis"
	"os"
	"strconv"
	"strings"
)

type getMailOptions struct {
	defaultOptions
	GetNew bool `short:"n" long:"long"`
}

func (opt *getMailOptions) Execute(_ []string) error {
	authInfo := opt.GetInfo()
	timeSlot, err := kulasis.RetrieveTimeSlot(authInfo)
	if err != nil {
		return err
	}
	var lectures []*kulasis.Lecture
	if opt.GetNew {
		lectures = timeSlot.GetNewLecture()
	} else {
		lectures = timeSlot.GetAllLectures()
	}

	for i, l := range lectures {
		fmt.Printf("%d: %s\n", i+1, l.LectureName)
	}
	println("Select lectures to read course mail.")

	reader := bufio.NewReader(os.Stdin)
	numStr, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return err
	}
	if num < 1 || num > len(lectures) {
		return errors.New("INVALID SELECTION")
	}
	lecture := lectures[num-1]

	titles, err := lecture.GetCourseMailTitles()
	if err != nil {
		return err
	}

	for _, t := range *titles {
		mail, err := t.GetContent()
		if err != nil {
			return err
		}

		fmt.Printf("[%s] - %s\n%s\n##########\n\n", mail.Title, mail.Date, mail.TextBody)
	}
	return nil
}
