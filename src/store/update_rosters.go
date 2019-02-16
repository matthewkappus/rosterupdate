package store

import (
	"time"

	"github.com/matthewkappus/rosterUpdate/src/synergy"
)

// DownloadRosters prompts user for Synergy Credentials, downloads Stu415s and emails, and
// returns error if the db can't be updated or provided wait time exceeded
func (r Roster) DownloadRosters(wait time.Duration, synergyUser, synergyPassword string) error {

	// u, pw, err := synergy.CredentialPrompt()
	// if err != nil {
	// 	return err
	// }
	ac, err := synergy.NewClient(synergyUser, synergyPassword, wait)
	if err != nil {
		return err
	}

	emails, err := ac.DownloadEmails()
	if err != nil {
		return err
	}

	s415s, err := ac.DownloadCurrentStu415s()
	if err != nil {
		return err
	}

	if err := s415s.SetSyncIDs(); err != nil {
		return err
	}

	if err := r.CreateNewStu415AndStaffEmails(); err != nil {
		return err
	}
	s415s.TeacherNameToEmail(emails)

	if err := r.InsertStu415s(s415s); err != nil {
		return err
	}

	return r.CreateMatthewADV()
}
