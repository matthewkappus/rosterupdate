package service

import (
	"fmt"
	"time"

	"github.com/matthewkappus/rosterUpdate/src/synergy"
)

// UpdateRosters prompts user for Synergy Credentials, downloads Stu415s and emails, and
// returns error if the db can't be updated or provided wait time exceeded
func (c *Classroom) UpdateRosters(wait time.Duration) error {

	u, pw, err := synergy.CredentialPrompt()
	if err != nil {
		return err
	}
	ac, err := synergy.NewClient(u, pw, wait)
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

	if err := c.DB.CreateNewStu415AndStaffEmails(); err != nil {
		return err
	}
	s415s.TeacherNameToEmail(emails)

	if err := c.DB.InsertStu415s(s415s); err != nil {
		return err
	}

	fmt.Println("Download complete")
	return c.DB.CreateMatthewADV()
}
