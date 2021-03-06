package store

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strings"
)

type (

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
		CourseIDAndTitle string `json:"course_id_and_title,omitempty"`
		MeetDays         string `json:"meet_days,omitempty"`
		Teacher          string `json:"teacher,omitempty"`
		Room             string `json:"room,omitempty"`
		Prescheduled     string `json:"prescheduled,omitempty"`
		SyncID           string `json:"sync_id,omitempty"`
	}

	// Stu415s is a list of Stu415
	Stu415s []*Stu415
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
