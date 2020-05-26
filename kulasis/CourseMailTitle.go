package kulasis

import (
	"fmt"
	"net/url"
)

type CourseMailTitle struct {
	MailNo       int    `json:"courseMailNo"`
	Date         string `json:"date"`
	DepartmentNo int    `json:"departmentNo"`
	IsNew        bool   `json:"isNew"`
	Title        string `json:"title"`
	info         *Info
}

func (title *CourseMailTitle) GetContent() (*CourseMail, error) {
	mailUrl, err := url.Parse(fmt.Sprintf(
		"https://www.k.kyoto-u.ac.jp/api/app/v1/support/course_mail?departmentNo=%d&courseMailNo=%d",
		title.DepartmentNo, title.MailNo))
	if err != nil {
		return nil, err
	}

	mail := CourseMail{}
	err = accessWithToken(*mailUrl, title.info, &mail)
	if err != nil {
		return nil, err
	}

	return &mail, nil
}
