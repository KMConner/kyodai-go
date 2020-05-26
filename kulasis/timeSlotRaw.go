package kulasis

import (
	"time"
)

type lectureRaw struct {
	DepartmentName   string `json:"departmentName"`
	DepartmentNo     int    `json:"departmentNo"`
	IsNew            bool   `json:"isNew"`
	IsShownOnKouki   bool   `json:"isShownOnKouki"`
	IsShownOnZenki   bool   `json:"isShownOnZenki"`
	IsSyutyuSemester bool   `json:"isSyutyuSemester"`
	LectureName      string `json:"lectureName"`
	LectureNo        int    `json:"lectureNo"`
	PeriodNo         int    `json:"periodNo"`
	RoomName         string `json:"roomName"`
	TeacherName      string `json:"teacherName"`
	WeekdayNo        int    `json:"weekdayNo"`
}

type timeSlotRaw struct {
	Timetables []lectureRaw `json:"timetables"`
}

func (lr *lectureRaw) extractSemester() Semester {
	if lr.IsSyutyuSemester {
		if lr.IsShownOnZenki && lr.IsShownOnKouki {
			return FullYearIntensive
		}
		if lr.IsShownOnZenki {
			return FirstIntensive
		}
		return SecondIntensive
	}
	if lr.IsShownOnZenki && lr.IsShownOnKouki {
		return FullYear
	}
	if lr.IsShownOnZenki {
		return First
	}
	return Second
}

func (lr *lectureRaw) extractDay() time.Weekday {
	switch lr.WeekdayNo {
	case 1:
		return time.Monday
	case 2:
		return time.Tuesday
	case 3:
		return time.Wednesday
	case 4:
		return time.Thursday
	default:
		return time.Friday
	}
}

func (lr *lectureRaw) extractPeriod() DayPeriod {
	return DayPeriod{
		Semester: lr.extractSemester(),
		Period:   lr.PeriodNo,
		Day:      lr.extractDay(),
	}
}

func (lr *lectureRaw) extractId() lectureId {
	return lectureId{
		DepartmentNo: lr.DepartmentNo,
		LectureNo:    lr.LectureNo,
	}
}

func (lr *lectureRaw) toLecture() *Lecture {
	return &Lecture{
		DepartmentName: lr.DepartmentName,
		DepartmentNo:   lr.DepartmentNo,
		IsNew:          lr.IsNew,
		LectureName:    lr.LectureName,
		LectureNo:      lr.LectureNo,
		RoomName:       lr.RoomName,
		TeacherName:    lr.TeacherName,
	}
}

func (tsr *timeSlotRaw) toTimeSlot() *TimeSlot {
	ret1 := make(map[DayPeriod]lectureId)
	ret2 := make(map[lectureId]*Lecture)
	for _, lec := range tsr.Timetables {
		lecId := lec.extractId()
		if _, ok := ret2[lecId]; !ok {
			ret2[lecId] = lec.toLecture()
		}
		ret1[lec.extractPeriod()] = lecId
	}
	return &TimeSlot{
		times:    ret1,
		lectures: ret2,
	}
}
