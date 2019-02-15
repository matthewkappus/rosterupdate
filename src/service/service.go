package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/matthewkappus/syncup/src/handlers"
	"github.com/matthewkappus/syncup/src/store"
	"github.com/matthewkappus/syncup/src/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gc "google.golang.org/api/classroom/v1"
)

// Classroom provides data services to handlers
type Classroom struct {
	DB      *store.Rosters
	Session *handlers.Session
	api     *gc.Service
}

// RostersHandler Returns json handler rendering rosters of logged in user
func (c *Classroom) RostersHandler(w http.ResponseWriter, r *http.Request) {
	u, err := c.Session.User(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	rosters, err := c.DB.GetRosters(u)
	if len(rosters) == 0 {
		fmt.Fprint(w, "no rosters found for", u)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ToJSON(w, rosters)
}

// SampRosters Returns json handler rendering rosters of logged in user
func (c *Classroom) SampRosters(w http.ResponseWriter, r *http.Request) {

	rosters, err := c.DB.GetRosters("danza@aps.edu")
	if len(rosters) == 0 {
		fmt.Fprint(w, "no rosters found for")
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ToJSON(w, rosters)
}

// AddSyncClassHandler adds to Google Classroom and Database the POSTed SynClass
func (c *Classroom) AddSyncClassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	sc := new(types.SyncClass)
	_, err := c.Session.User(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewDecoder(r.Body).Decode(sc)

	if err := c.createSyncClass(sc); err != nil {
		fmt.Println("createSyncClass made err", err.Error())
	}

	json.NewEncoder(w).Encode(sc)
}

// createSyncClass takes a posted SyncClass (partial value) and
// validates and updates the fields after a google API call
func (c *Classroom) createSyncClass(sc *types.SyncClass) (err error) {
	if sc.SyncID == "" || sc.Name == "" {
		return fmt.Errorf("Course not created: Need SyncID and Name")
	}
	sc.Roster, err = c.DB.SelectStu415sBySyncID(sc.SyncID)
	if err != nil {
		return err
	}

	course, err := c.NewCourse(sc.Name, sc.Per, sc.Description)
	if err != nil {
		return err
	}

	sc.GCID = course.Id
	if err = c.InviteStudents(sc.GCID, sc.Roster); err != nil {
		return err
	}
	return c.DB.InsertSyncClass(sc)
}

// GetUser returns a UserProfile if the API is not nil
func (c *Classroom) getUser() (*gc.UserProfile, error) {
	if c.api == nil {
		return nil, fmt.Errorf("Can't get user: Classroom service is not initiated")
	}

	return gc.NewUserProfilesService(c.api).Get("me").Do()
}

var config = &oauth2.Config{
	ClientID:     "419661903175-ibl4qvtb8rtkhjs30vugj8s70vgiifg6.apps.googleusercontent.com",
	ClientSecret: "oUM46Ty8v4sfESkPheFU8RR8",
	RedirectURL:  "http://localhost:8080/oauth2callbackhandler",
	Scopes:       []string{gc.ClassroomProfileEmailsScope, gc.ClassroomRostersScope, gc.ClassroomCoursesScope},
	Endpoint:     google.Endpoint,
}
