package store

import (
	"database/sql"
	"fmt"
	"strings"

	// kosher
	"github.com/matthewkappus/syncup/src/types"
	// Required by law
	_ "github.com/mattn/go-sqlite3"
)

//  head -n1 stu415.csv | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | tr ',' ' NOT NULL,'

const (
	createStu415Table            = `CREATE TABLE IF NOT EXISTS stu415(student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id TEXT)`
	dropStu415Table              = `DROP TABLE IF EXISTS stu415`
	createStaffEmailsTable       = `CREATE TABLE IF NOT EXISTS staff_emails(email PRIMARY KEY NOT NULL, teacher)`
	dropStaffEmailsTable         = `DROP TABLE IF EXISTS staff_emails`
	insertStaffEmails            = `INSERT INTO staff_emails(email, teacher) VALUES(?,?)`
	insertStu415                 = `INSERT INTO stu415(student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`
	selectStu415sByTeacherPeriod = `SELECT student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id FROM stu415 WHERE teacher=? AND per=?`
	selectStu415sByTeacher       = `SELECT student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id FROM stu415 WHERE teacher=?`
	selectStu415BySection        = `SELECT student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id FROM stu415 WHERE section_id=?`
	selectStu415BySID            = `SELECT student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, teacher, room, sync_id FROM stu415 WHERE sync_id=?`
	selectEmailByName            = `SELECT email FROM staff_emails WHERE teacher=?`
	selectNameFromEmail          = `SELECT teacher from staff_emails where email=?`

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
func (rs *Rosters) CreateNewStu415AndStaffEmails() error {
	var err error
	if _, err = rs.Exec(dropStaffEmailsTable); err != nil {
		return err
	}
	if _, err = rs.Exec(createStaffEmailsTable); err != nil {
		return err
	}

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
func (rs *Rosters) UpdateRosters(s415s types.Stu415s, emails [][]string) error {
	if len(s415s) == 0 || len(emails) == 0 {
		return fmt.Errorf("can't update empty s415s len(%d) or emails len(%d)", len(s415s), len(emails))
	}

	var err error
	if rs.CreateNewStu415AndStaffEmails(); err != nil {
		return err
	}

	if err = rs.InsertStu415s(s415s); err != nil {
		return err
	}

	return rs.InsertStaffEmails(emails)
}

// CreateMatthewADV adds matthew.kappus@aps.edu as teacher to advisories
func (rs *Rosters) CreateMatthewADV() error {
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

func (rs *Rosters) initTables() error {
	if _, err := rs.Exec(createStaffEmailsTable); err != nil {
		return err
	}
	if _, err := rs.Exec(createStu415Table); err != nil {
		return err
	}

	_, err := rs.Exec(createSyncClasses)
	return err
}

// InsertStaffEmails inserts list if [name, email] and returns an error
func (rs *Rosters) InsertStaffEmails(emails [][]string) error {
	tx, err := rs.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStaffEmails)
	if err != nil {
		return err
	}
	for _, e := range emails {
		if len(e) != 2 {
			continue
		}
		_, err = stmt.Exec(e[0], e[1])

	}

	return tx.Commit()
}

// InsertStu415s deletes (old) and inserts provided types.Stu415 slice into table and returns an error
func (rs *Rosters) InsertStu415s(s415s types.Stu415s) error {
	tx, err := rs.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStu415)
	if err != nil {
		return err
	}

	for _, s := range s415s {
		// if s.Teacher, err = rs.TeacherEmailFromName(s.Teacher); err != nil {
		// 	continue
		// }

		stmt.Exec(
			s.StudentName,
			s.PermID,
			s.Gender,
			s.Grade,
			s.TermName,
			s.Per,
			s.Term,
			s.SectionID,
			s.CourseIDAndTitle,
			s.Teacher,
			s.Room,
			s.SyncID,
		)

	}

	return tx.Commit()
}

// IsTeacher returns true if teacher_emails contains provided email
func (rs *Rosters) IsTeacher(email string) bool {
	_, err := rs.SelectTeacherNameFromEmail(email)
	return err == nil
}

// TeacherEmailFromName takes a name and returns their email
func (rs *Rosters) TeacherEmailFromName(name string) (email string, err error) {
	err = rs.QueryRow(selectEmailByName, name).Scan(&email)
	return email, err
}

// SelectTeacherNameFromEmail takes a name and returns their email
func (rs *Rosters) SelectTeacherNameFromEmail(email string) (name string, err error) {
	err = rs.QueryRow(selectNameFromEmail, cleanEmail(email)).Scan(&name)
	return name, err
}

func scan415(rows *sql.Rows) (*types.Stu415, error) {
	s := new(types.Stu415)

	err := rows.Scan(
		&s.StudentName,
		&s.PermID,
		&s.Gender,
		&s.Grade,
		&s.TermName,
		&s.Per,
		&s.Term,
		&s.SectionID,
		&s.CourseIDAndTitle,
		&s.Teacher,
		&s.Room,
		&s.SyncID,
	)
	return s, err
}

// SelectStu415sBySyncID returns students by their sync id
func (rs *Rosters) SelectStu415sBySyncID(sid string) (types.Stu415s, error) {
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
func (rs *Rosters) SelectStu415sByTeacherPeriod(email, per string) (s415s types.Stu415s, err error) {
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
func (rs *Rosters) SelectStu415sByTeacher(email string) (s415s types.Stu415s, err error) {
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

// GetRosters returns rosters by provided email
func (rs *Rosters) GetRosters(email string) (rosters types.Rosters, err error) {
	stu415s, err := rs.SelectStu415sByTeacher(email)
	if err != nil {
		return nil, err
	}
	return types.Stu415sToRoster(stu415s), nil
}
