package service

import (
	"fmt"

	"github.com/matthewkappus/syncup/src/types"
	gc "google.golang.org/api/classroom/v1"
)

// NewCourse takes a teacher email, Google Classroom name, section and description and returns an  API error
func (c *Classroom) NewCourse(name, section, description string) (*gc.Course, error) {
	return gc.NewCoursesService(c.api).Create(&gc.Course{
		Name:        name,
		Section:     section,
		Description: description,
		OwnerId:     "me",
	}).Do()

}

// RemoveStudents deltes the provided students from gcid returning possible error
func (c *Classroom) RemoveStudents(gcid string, toRemove []*gc.Student) error {
	cssvc := gc.NewCoursesStudentsService(c.api)
	for _, s := range toRemove {
		// todo handle error
		cssvc.Delete(gcid, s.UserId).Do()
	}
	return nil
}

// InviteStudents takes  slice of class Students and Google Classroom Course Id
// It returns an API error on serice invite call
func (c *Classroom) InviteStudents(gcid string, s415s types.Stu415s) error {
	// TODO: Batch invites for 1 api call
	var inviteError string
	isvc := gc.NewInvitationsService(c.api)
	for _, s := range s415s {
		_, err := isvc.Create(&gc.Invitation{
			CourseId: gcid,
			Role:     "STUDENT",
			UserId:   s.PermID,
		}).Do()
		if err != nil {
			inviteError += "\n" + err.Error()
		}
	}
	if inviteError != "" {
		return fmt.Errorf("InviteStudents error: %s", inviteError)
	}
	return nil
}

// GetRoster returns a slice of Students or an error from provided Course ID
func (c *Classroom) GetRoster(gcid string) ([]*gc.Student, error) {
	cssvc := gc.NewCoursesStudentsService(c.api)
	cl, err := cssvc.List(gcid).Do()
	if err != nil {
		return nil, err
	}
	return cl.Students, nil
}

// GetInvites returns a slice of Students for provided gcid returning err
func (c *Classroom) GetInvites(gcid string) ([]*gc.Invitation, error) {
	isvc := gc.NewInvitationsService(c.api)
	li, err := isvc.List().Do()
	if err != nil {
		return nil, err
	}
	return li.Invitations, nil
}
