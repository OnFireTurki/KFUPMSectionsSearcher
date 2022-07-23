package LFS

import "fmt"

// This file was created in order to simplify the json response handling
// Note : we have to create struct to store the data but not all the data, only what we need
// I hate writing comments anyway

// facultyInfo : struct storing the instructor inforamtion, sadly everything we need is just the name, maybe the email also
type facultyInfoS struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// meetingTime : to get metting time
type meetingTimeS struct {
	BeginTime           string `json:"beginTime"`
	EndTime             string `json:"endTime"`
	Sunday              bool   `json:"sunday"`
	Monday              bool   `json:"monday"`
	Tuesday             bool   `json:"tuesday"`
	Wednesday           bool   `json:"wednesday"`
	Thursday            bool   `json:"thursday"`
	MeetingScheduleType string `json:"meetingScheduleType"`
}

// GetDays : Return days letters,
// Any other ideas ?
func (c *meetingTimeS) GetDays() string {
	Days := ""
	if c.Sunday {
		Days += "U"
	}
	if c.Monday {
		Days += "M"
	}
	if c.Tuesday {
		Days += "T"
	}
	if c.Wednesday {
		Days += "W"
	}
	if c.Thursday {
		Days += "R"
	}
	return Days
}

// meetingsFaculty : I am tired of typing comments
type meetingsFacultyS struct {
	MeetingsFaculty meetingTimeS `json:"meetingTime"`
}

// CourseDetails : to store the section details
type CourseDetails struct {
	CourseReferenceNumber   string             `json:"courseReferenceNumber"`
	SequenceNumber          string             `json:"sequenceNumber"`
	ScheduleTypeDescription string             `json:"scheduleTypeDescription"`
	MaximumEnrollment       int                `json:"maximumEnrollment"`
	Enrollment              int                `json:"enrollment"`
	SeatsAvailable          int                `json:"seatsAvailable"`
	WaitCapacity            int                `json:"waitCapacity"`
	WaitCount               int                `json:"waitCount"`
	WaitAvailable           int                `json:"waitAvailable"`
	OpenSection             bool               `json:"openSection"`
	MeetingsFaculty         []meetingsFacultyS `json:"meetingsFaculty"`
	Faculty                 []facultyInfoS     `json:"faculty"`
}

// CheckWaitList : Can I get this section?
func (Course *CourseDetails) CheckWaitList() bool {
	return Course.OpenSection && Course.WaitAvailable > 0
}

// PrintIT : to print out the details
func (Course *CourseDetails) PrintIT() {
	fmt.Printf("CRN: %v | Section : %v | Available Seats : %v | Instructor : %v | Include : \n", Course.CourseReferenceNumber, Course.SequenceNumber,
		Course.SeatsAvailable, Course.Faculty[0].DisplayName)
	for i := 0; i < len(Course.MeetingsFaculty); i++ {
		fmt.Printf("- %v on %v at %v\n", Course.MeetingsFaculty[i].MeetingsFaculty.MeetingScheduleType, Course.MeetingsFaculty[i].MeetingsFaculty.GetDays(), fmt.Sprintf("%v-%v", Course.MeetingsFaculty[i].MeetingsFaculty.BeginTime, Course.MeetingsFaculty[i].MeetingsFaculty.EndTime))
	}
	fmt.Println("=======================================================================")
}

// JsonRes : To Get the Json Res
type JsonRes struct {
	// Total count of all available seciton
	TotalCount int `json:"totalCount"`
	// For CourseDetails
	Data []CourseDetails `json:"data"`
}
