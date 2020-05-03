package kulasis

import (
	"github.com/KMConner/kyodai-go/internal/auth"
	"github.com/KMConner/kyodai-go/internal/network"
	"net/url"
	"time"
)

type Semester int

const (
	First Semester = iota
	FirstIntensive
	Second
	SecondIntensive
	FullYear
	FullYearIntensive
)

type LectureId struct {
	DepartmentNo int
	LectureNo    int
}

type Lecture struct {
	DepartmentName string
	DepartmentNo   int
	IsNew          bool
	LectureName    string
	LectureNo      int
	RoomName       string
	TeacherName    string
	info           *auth.Info
}

type DayPeriod struct {
	Semester Semester
	Day      time.Weekday
	Period   int
}

type TimeSlot struct {
	times    map[DayPeriod]LectureId
	lectures map[LectureId]*Lecture
}

func RetrieveTimeSlot(info auth.Info) (*TimeSlot, error) {
	var timeSlotRaw timeSlotRaw
	timeslotUrl, err := url.Parse("https://www.k.kyoto-u.ac.jp/api/app/v1/timetable/get_table")
	if err != nil {
		return nil, err
	}

	err = network.AccessWithToken(*timeslotUrl, &info, &timeSlotRaw)
	if err != nil {
		return nil, err
	}
	timeSlot := timeSlotRaw.toTimeSlot()
	for _, v := range timeSlot.lectures {
		v.info = &info
	}
	return timeSlot, nil
}

func (slot *TimeSlot) GetLecture(dp DayPeriod) *Lecture {
	if v, ok := slot.times[dp]; ok {
		if l, ok := slot.lectures[v]; ok {
			return l
		}
	}
	return nil
}

func (slot *TimeSlot) GetNewLecture() []*Lecture {
	var ret []*Lecture
	for _, v := range slot.lectures {
		if v.IsNew {
			ret = append(ret, v)
		}
	}
	return ret
}
