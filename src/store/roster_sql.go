package store

// import (
// 	"github.com/matthewkappus/rosterUpdate/src/types"
// )

// // create table queries
// const (
// 	createStaffEmails = `CREATE TABLE IF NOT EXISTS staff_emails(email, name text)`
// )

// // select queries
// const (
// 	// id, group_id, sync_id, title, email, is_active
// 	selectRoster         = `SELECT group_id, sync_id, title, email, is_active FROM rosters WHERE id=?`
// 	selectRostersByEmail = `SELECT id, group_id, sync_id, title, is_active FROM rosters where email=?`
// 	selectTeacherByGroup = `SELECT email, name FROM teachers WHERE group_id=?`

// 	selectNameByEmail = `SELECT name FROM staff_emails where email=?`
// )

// // insert queries
// const (
// 	insertSyncer  = `INSERT INTO rosters(group_id, sync_id, title, email) values(?,?,?,?)`
// 	insertStudent = `INSERT INTO students(perm_id, name, group_id) values(?,?,?)`
// 	insertTeacher = `INSERT INTO teachers(email, name, group_id) values(?,?,?)`
// )

// // update queries
// const (
// 	// When a Stu415 set changes studens, it creates new GRID which matches the SID of a course
// 	syncRoster = `UPDATE rosters SET group_id=? WHERE sync_id=?`
// )

// // IsStaff returns true if provided email found in stasff_emails
// func (rs *Roster) IsStaff(email string) bool {
// 	var name string
// 	err := rs.QueryRow(selectNameByEmail, email).Scan(&name)
// 	return err != nil
// }

