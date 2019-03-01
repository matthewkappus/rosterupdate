package store

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// kosher
	// Required by law
	"github.com/matthewkappus/rosterUpdate/src/types"
	//kosher
	_ "github.com/mattn/go-sqlite3"
)

//  head -n1 stu415.csv | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | tr ',' ' NOT NULL,'
// Organization Name,School Year,Student Name,Perm ID,Gender,Grade,Term Name,Per,Term,Section ID,Course ID And Title,Meet Days,Teacher,Room,PreScheduled
const (
	createStu415Table = `CREATE TABLE IF NOT EXISTS stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, prescheduled, sync_id TEXT)`
	// createStu415Table            = `CREATE TABLE IF NOT EXISTS stu415(student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id TEXT)`
	dropStu415Table              = `DROP TABLE IF EXISTS stu415`
	insertStu415                 = `INSERT INTO stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, prescheduled, sync_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	selectStu415sByTeacherPeriod = `SELECT * FROM stu415 WHERE teacher=? AND per=?`
	selectStu415sByTeacher       = `SELECT * FROM stu415 WHERE teacher=?`
	selectStu415BySection        = `SELECT * FROM stu415 WHERE section_id=?`
	selectStu415BySID            = `SELECT * FROM stu415 WHERE sync_id=?`

	// Executed by ac.CreateMatthewADV
	createTmp         = `CREATE TABLE tmp AS SELECT * FROM stu415 WHERE course_id_and_title LIKE "08%00%"`
	updateTmp         = `UPDATE tmp SET teacher="matthew.kappus@aps.edu"`
	insertTmpToStu415 = `INSERT INTO stu415 SELECT * FROM tmp`
	dropTmp           = `DROP TABLE IF EXISTS tmp`
)

func cleanEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// CreateNewStu415AndStaffEmails drops old roster tables
func (rs *Roster) CreateNewStu415AndStaffEmails() error {
	var err error
	if _, err = rs.Exec(dropStu415Table); err != nil {
		return err
	}

	if _, err = rs.Exec(createStu415Table); err != nil {
		return err
	}

	return nil
}

// UpdateRosters creates new tables and inserts arguments Stu415 and staff emails
// TODO: Use this in types.Update
func (rs *Roster) UpdateRosters(s415s types.Stu415s, emails [][]string) error {
	if len(s415s) == 0 || len(emails) == 0 {
		return fmt.Errorf("can't update empty s415s len(%d) or emails len(%d)", len(s415s), len(emails))
	}

	var err error
	if rs.CreateNewStu415AndStaffEmails(); err != nil {
		return err
	}

	return rs.InsertStu415s(s415s)

}

// CreateMatthewADV adds matthew.kappus@aps.edu as teacher to advisories
func (rs *Roster) CreateMatthewADV() error {
	if _, err := rs.Exec(dropTmp); err != nil {
		println("dropTmp err")
		return err
	}
	if _, err := rs.Exec(createTmp); err != nil {
		println("createTmp err")
		return err
	}
	if _, err := rs.Exec(updateTmp); err != nil {
		println("updateTmp err")
		return err
	}
	if _, err := rs.Exec(insertTmpToStu415); err != nil {
		println("insertTmpToStu415 err")
		return err
	}

	return nil
}

func (rs *Roster) initTables() error {

	if _, err := rs.Exec(dropTmp); err != nil {
		return err
	}
	if _, err := rs.Exec(createStu415Table); err != nil {
		return err
	}

	return nil
}

// InsertStu415s deletes (old) and inserts provided types.Stu415 slice into table and returns an error
func (rs *Roster) InsertStu415s(s415s types.Stu415s) error {
	tx, err := rs.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStu415)
	if err != nil {
		return err
	}

	for _, s := range s415s {

		// organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, prescheduled, sync_id
		if _, err := stmt.Exec(
			s.OrganizationName,
			s.SchoolYear,
			s.StudentName,
			s.PermID,
			s.Gender,
			s.Grade,
			s.TermName,
			s.Per,
			s.Term,
			s.SectionID,
			s.CourseIDAndTitle,
			s.MeetDays,
			s.Teacher,
			s.Room,
			s.Prescheduled,
			s.SyncID,
		); err != nil {
			log.Fatal(err)
		}

	}

	return tx.Commit()
}

func scan415(rows *sql.Rows) (*types.Stu415, error) {
	s := new(types.Stu415)

	err := rows.Scan(
		&s.OrganizationName,
		&s.SchoolYear,
		&s.StudentName,
		&s.PermID,
		&s.Gender,
		&s.Grade,
		&s.TermName,
		&s.Per,
		&s.Term,
		&s.SectionID,
		&s.CourseIDAndTitle,
		&s.MeetDays,
		&s.Teacher,
		&s.Room,
		&s.Prescheduled,
		&s.SyncID,
	)
	return s, err
}

// SelectClassesByTeacher gets stu415s by provided teacher email and returns a Roster list or error
func (rs *Roster) SelectClassesByTeacher(email string) ([]*types.Class, error) {
	s415s, err := rs.SelectStu415sByTeacher(email)
	if err != nil {
		return nil, err
	}
	return types.Stu415sToClasses(s415s), nil
}

// SelectClassByID takes a syncid and returns a class or error if not found
func (rs *Roster) SelectClassByID(sid string) (*types.Class, error) {
	s415s, err := rs.SelectStu415sBySyncID(sid)
	if err != nil {
		return nil, err
	}
	return types.Stu415sToClass(s415s), nil
}

// SelectStu415sBySyncID returns students by their sync id
func (rs *Roster) SelectStu415sBySyncID(sid string) (types.Stu415s, error) {
	stmt, err := rs.Prepare(selectStu415BySID)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(sid)
	if err != nil {
		return nil, err
	}
	s415s := make(types.Stu415s, 0)
	var scanErr string
	var scanCount int
	for rows.Next() {
		s, err := scan415(rows)
		if err != nil {
			scanCount++
			scanErr += fmt.Sprintf("Scan Error: %v\n", err)
		}
		s415s = append(s415s, s)
	}
	if scanCount > 0 {
		return nil, fmt.Errorf("SelectStu415sBySyncID encountered %d errors:\n%s", scanCount, scanErr)
	}
	return s415s, nil
}

// SelectStu415sByTeacherPeriod returns stu415s where provided email is the Teacher,
func (rs *Roster) SelectStu415sByTeacherPeriod(email, per string) (s415s types.Stu415s, err error) {
	stmt, err := rs.Prepare(selectStu415sByTeacherPeriod)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(email, per)
	if err != nil {
		return nil, err
	}

	var scanErr string
	var scanCount int
	for rows.Next() {
		s, err := scan415(rows)
		if err != nil {
			scanCount++
			scanErr += fmt.Sprintf("Scan Error: %v\n", err)
		}
		s415s = append(s415s, s)
	}
	if scanCount > 0 {
		return nil, fmt.Errorf("SelectStu415sByTeacher encountered %d errors:\n%s", scanCount, scanErr)
	}
	return s415s, nil
}

// SelectStu415sByTeacher returns stu415s where provided email is the Teacher,
func (rs *Roster) SelectStu415sByTeacher(email string) (s415s types.Stu415s, err error) {
	stmt, err := rs.Prepare(selectStu415sByTeacher)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}

	var scanErr string
	var scanCount int
	for rows.Next() {
		s, err := scan415(rows)
		if err != nil {
			scanCount++
			scanErr += fmt.Sprintf("Scan Error: %v\n", err)
		}
		s415s = append(s415s, s)
	}
	if scanCount > 0 {
		return nil, fmt.Errorf("SelectStu415sByTeacher encountered %d errors:\n%s", scanCount, scanErr)
	}
	return s415s, nil
}
