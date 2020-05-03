package kulasis

import "github.com/KMConner/kyodai-go/internal/auth"

type courseMailTitle struct {
	MailNo       int    `json:"courseMailNo"`
	Date         string `json:"date"`
	DepartmentNo int    `json:"departmentNo"`
	IsNew        bool   `json:"isNew"`
	Title        string `json:"title"`
	info         *auth.Info
}
