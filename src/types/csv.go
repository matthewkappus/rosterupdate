package types

import (
	"encoding/csv"
	"io"
)

// Parse returns an error if provided row can't be assigned to s415s
// adds @aps.edu to perm
func (s *Stu415) Parse(r []string) error {
	s.OrganizationName = r[0]
	s.SchoolYear = r[1]
	s.StudentName = r[2]
	s.PermID = r[3] + "@aps.edu"
	s.Gender = r[4]
	s.Grade = r[5]
	s.TermName = r[6]
	s.Per = r[7]
	s.Term = r[8]
	s.SectionID = r[9]
	// no CourseID
	s.CourseIDAndTitle = r[10]
	s.MeetDays = r[11]
	s.Teacher = r[12]
	s.Room = r[13]
	s.Prescheduled = r[14]
	s.SyncID = r[15]

	return nil
}

// func (s *Stu415) Parse(r []string) error {
// 	s.StudentName = r[0]
// 	s.PermID = r[1] + "@aps.edu"
// 	s.Gender = r[2]
// 	s.Grade = r[3]
// 	s.TermName = r[4]
// 	s.Per = r[5]
// 	s.Term = r[6]
// 	s.SectionID = r[7]
// 	s.CourseIDAndTitle = r[8]
// 	s.Teacher = r[9]
// 	s.Room = r[10]

// 	return nil
// }

// Use s415s.Map to return a syncable list of rosters

// Stu415sFromCSV takes a csv string and returns
// Stu408s from parsing its rows
func Stu415sFromCSV(r io.Reader) (stus Stu415s, err error) {
	csvR := csv.NewReader(r)
	csvR.LazyQuotes = true
	records, err := csvR.ReadAll()

	if err != nil {
		return stus, err
	}

	for _, record := range records {
		s := new(Stu415)
		s.Parse(record)
		stus = append(stus, s)

	}
	return stus, nil
}
