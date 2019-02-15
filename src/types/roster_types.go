package types

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strings"
)

type (
	// GRID is a unique id associating stu415s
	GRID string

	// stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, prescheduled, sync_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`


	// Stu415 holds Synergery report course data
	Stu415 struct {
		OrganizationName string `json:"organization_name,omitempty"`
		SchoolYear       string `json:"school_year,omitempty"`
		StudentName      string `json:"student_name,omitempty"`
		PermID           string `json:"perm_id,omitempty"`
		Gender           string `json:"gender,omitempty"`
		Grade            string `json:"grade,omitempty"`
		TermName         string `json:"term_name,omitempty"`
		Per              string `json:"per,omitempty"`
		Term             string `json:"term,omitempty"`
		SectionID        string `json:"section_id,omitempty"`
		CourseID         string `json:"course_id,omitempty"`
		CourseIDAndTitle string `json:"course_id_and_title,omitempty"`
		MeetDays         string `json:"meet_days,omitempty"`
		Teacher          string `json:"teacher,omitempty"`
		Room             string `json:"room,omitempty"`
		Prescheduled     string `json:"prescheduled,omitempty"`
		SyncID           string `json:"sync_id,omitempty"`
	}

	// Stu415s is a list of Stu415
	Stu415s []*Stu415

	// Roster groups Stu415s by period
	Roster struct {
		ID       string  `json:"id,omitempty"`
		Title    string  `json:"title,omitempty"`
		Per      string  `json:"per,omitempty"`
		Teacher  string  `json:"teacher,omitempty"`
		Students Stu415s `json:"students,omitempty"`
	}

	// Rosters list
	Rosters []*Roster

	// SyncClass associates Stu415s with Classroom Courses
	// gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active
	SyncClass struct {
		GCID             string  `json:"gcid,omitempty"`
		SyncID           string  `json:"sync_id,omitempty"`
		Per              string  `json:"per,omitempty"`
		Term             string  `json:"term,omitempty"`
		Name             string  `json:"name,omitempty"`
		CourseIDAndTitle string  `json:"course_id_and_title,omitempty"`
		Description      string  `json:"description,omitempty"`
		Teacher          string  `json:"teacher,omitempty"`
		IsActive         bool    `json:"is_active,omitempty"`
		Roster           Stu415s `json:"roster,omitempty"`
		Invited          Stu415s `json:"invited,omitempty"`
		Enrolled         Stu415s `json:"enrolled,omitempty"`
	}
)

// SetSyncIDs hashes the stu415 props to create new sync id
func (s415s Stu415s) SetSyncIDs() error {
	if len(s415s) == 0 {
		return fmt.Errorf("no students to create syncid")
	}
	enc := fnv.New32()

	for _, s := range s415s {
		enc.Write([]byte(s.Per + s.CourseIDAndTitle + s.SectionID + s.Term))
		s.SyncID = hex.EncodeToString(enc.Sum(nil))
		enc.Reset()
	}
	return nil
}

// Stu415sToRoster takes students and groups them by periods
func Stu415sToRoster(s415s Stu415s) []*Roster {
	// Group students by period
	periods := make(map[string]Stu415s)
	for _, s := range s415s {
		if r, ok := periods[s.Per]; !ok {
			periods[s.Per] = Stu415s{s}
		} else {
			periods[s.Per] = append(r, s)
		}
	}

	var rosters = make([]*Roster, 0)
	// For each period, create a roster
	for _, students := range periods {
		r := &Roster{
			ID:       students[0].SyncID,
			Title:    students[0].CourseIDAndTitle,
			Per:      students[0].Per,
			Teacher:  students[0].Teacher,
			Students: students,
		}
		rosters = append(rosters, r)
	}
	return rosters
}

// TeacherNameToEmail takes a [email,name] list and changes s.Teacher to their lowercase email
func (s415s Stu415s) TeacherNameToEmail(emails [][]string) {
	nameEmail := make(map[string]string, len(emails))
	for _, r := range emails {
		nameEmail[r[1]] = strings.ToLower(r[0])
	}
	for _, s := range s415s {
		if email, ok := nameEmail[s.Teacher]; ok {
			s.Teacher = email
		}
	}
}

// ToList returns stu415s indexed by CourseIDAndTitle
func ToList(s415s Stu415s) map[string]Stu415s {
	m := make(map[string]Stu415s)
	for _, s := range s415s {
		if list, ok := m[s.CourseIDAndTitle]; !ok {
			m[s.CourseIDAndTitle] = Stu415s{s}
		} else {
			m[s.CourseIDAndTitle] = append(list, s)
		}
	}
	return m
}
