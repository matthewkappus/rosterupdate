package store

import (
	"database/sql"
	"log"

	"github.com/matthewkappus/syncup/src/types"
)

const (
	createSyncClasses          = `CREATE TABLE IF NOT EXISTS sync_classes(gcid, sync_id, per, term, name, course_id_and_title, description, teacher TEXT, is_active BOOLEAN)`
	insertSyncClasses          = `INSERT INTO sync_classes(gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active) VALUES(?,?,?,?,?,?,?,?,?)`
	inactivateSyncClass        = `UPDATE sync_classes SET is_active=0 WHERE gcid=?`
	selectSyncClassesByTeacher = `SELECT gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active FROM sync_classes where teacher=?`
	selectSyncClassesBySID     = `SELECT gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active FROM sync_classes where sync_id=?`
	selectByGCID               = `SELECT gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active FROM sync_classes where gcid=?`
)

func scanClass(r *sql.Rows) (*types.SyncClass, error) {
	// SELECT gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active FROM sync_classes where teacher=?
	s := new(types.SyncClass)
	err := r.Scan(
		&s.GCID,
		&s.SyncID,
		&s.Per,
		&s.Term,
		&s.Name,
		&s.CourseIDAndTitle,
		&s.Description,
		&s.Teacher,
		&s.IsActive,
	)

	return s, err
}

func scanRoster(rows *sql.Rows) (*types.Roster, error) {
	r := new(types.Roster)
	err := rows.Scan(
		&r.ID,
		&r.Title,
		&r.Per,
		&r.Teacher,
	)
	return r, err
}

// SelectClassesByTeacher returns a syncclass slice by belonging to provided teacher
func (rs *Rosters) SelectClassesByTeacher(email string) ([]*types.SyncClass, error) {
	scs := make([]*types.SyncClass, 0)
	rows, err := rs.Query(selectSyncClassesByTeacher, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := new(types.SyncClass)
		err := rows.Scan(
			&s.GCID,
			&s.SyncID,
			&s.Per,
			&s.Term,
			&s.Name,
			&s.CourseIDAndTitle,
			&s.Description,
			&s.Teacher,
			&s.IsActive,
		)
		if err != nil {
			log.Println("scan err:", err.Error())
			continue
		}
		scs = append(scs, s)

	}
	return scs, nil

}

// SelectSyncClassByID returns a syncclass slice by belonging to provided teacher
func (rs *Rosters) SelectSyncClassByID(sid string) ([]*types.SyncClass, error) {
	scs := make([]*types.SyncClass, 0)
	rows, err := rs.Query(selectSyncClassesBySID, sid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := new(types.SyncClass)
		err := rows.Scan(
			&s.GCID,
			&s.SyncID,
			&s.Per,
			&s.Term,
			&s.Name,
			&s.CourseIDAndTitle,
			&s.Description,
			&s.Teacher,
			&s.IsActive,
		)
		if err != nil {
			log.Println("scan err:", err.Error())
			continue
		}
		scs = append(scs, s)

	}
	return scs, nil

}

func scanStu415(rows *sql.Rows) (*types.Stu415, error) {
	s := new(types.Stu415)
	err := rows.Scan(
		&s.StudentName,
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
		&s.PermID,
	)
	return s, err
}

// SelectSyncClass returns a syncclass by provided GCID
func (rs *Rosters) SelectSyncClass(gcid string) (*types.SyncClass, error) {
	r := rs.QueryRow(selectByGCID)

	s := new(types.SyncClass)
	// gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active
	err := r.Scan(
		&s.GCID,
		&s.SyncID,
		&s.Per,
		&s.Term,
		&s.Name,
		&s.CourseIDAndTitle,
		&s.Description,
		&s.Teacher,
		&s.IsActive,
	)
	s.Roster, err = rs.SelectStu415sBySyncID(s.SyncID)
	if err != nil {
		return nil, err
	}

	return s, err
}

// InsertSyncClass takes a syncclass and returns an error if it can't be inserted into store
func (rs *Rosters) InsertSyncClass(sc *types.SyncClass) error {
	// gcid, sync_id, per, term, name, course_id_and_title, description, teacher, is_active
	_, err := rs.Exec(insertSyncClasses, sc.GCID, sc.SyncID, sc.Per, sc.Term, sc.Name, sc.CourseIDAndTitle, sc.Description, sc.Teacher, sc.IsActive)
	return err
}

// InactivateSyncClass sets isactive to 0 or returns error
func (rs *Rosters) InactivateSyncClass(sc *types.SyncClass) error {
	if _, err := rs.Exec(inactivateSyncClass, sc.GCID); err != nil {
		return err
	}
	sc.IsActive = false
	return nil
}
