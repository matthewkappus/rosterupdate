// Package synergy makes authenticated job requests for retrieving csv files
// package synergy
package synergy

import (
	"github.com/matthewkappus/syncup/src/types"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

// AuthClient singleton for retrieving different job requests
var ac *AuthClient

// Parse values for Synergy authenticated Report requests
var (
	reViewState          = regexp.MustCompile(`id="__VIEWSTATE" value="(.+)?"`)
	reViewStateGenerator = regexp.MustCompile(`id="__VIEWSTATEGENERATOR" value="(.{8})"`)
	reFocusKey           = regexp.MustCompile(`ST.RevFocusKey = '(.*)?'`)

	reStu415GUID = regexp.MustCompile(`<ROW GUID="(.{36})".+STU415`)
)

// Synergy site endpoints
const (
	loginURL        = "https://synergy.aps.edu/Login.aspx"
	logoutURL       = "https://synergy.aps.edu/ST_Content.aspx?logout=true"
	xmlDoRequestURL = "https://synergy.aps.edu/Service/RTCommunication.asmx/XMLDoRequest"
	downloadURL     = "https://synergy.aps.edu/Download.aspx"
)

// AuthClient encapsulates values for synergy authentication and methods for
// authenticated http requests
type AuthClient struct {
	wait     time.Duration
	focusKey string
	c        *http.Client
}

// CredentialPrompt asks for a Synergy username and password
// It returns error if invalid or a username and password
func CredentialPrompt() (userName, password string, err error) {
	fmt.Println("Synergy username (e#####) and password requried for update")

	fmt.Print("Username: ")
	if _, err := fmt.Scanln(&userName); err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	if _, err := fmt.Scanln(&password); err != nil {
		return "", "", err
	}

	// TODO: Check login err
	return userName, password, nil
}

// Logout ends authenticated session
func (ac *AuthClient) Logout() error {
	_, err := ac.c.Get(logoutURL)
	return err
}

// NewClient takes synergy credentials to create an auth session with cookiejar and a
// duration in which to Timeout on http requests
// Returns encapsulated *http.Client if login successful, else error
func NewClient(synergyUser, synergyPassword string, wait time.Duration) (*AuthClient, error) {
	res, err := http.Get(loginURL)
	if err != nil {
		return nil, err
	}
	body, err := readClose(res)
	if err != nil {
		return nil, err
	}

	viewstate := parseSubmatch(reViewState, body)
	viewstateGenerator := parseSubmatch(reViewStateGenerator, body)

	jar, _ := cookiejar.New(&cookiejar.Options{})
	c := &http.Client{Jar: jar}

	loginResponse, err := c.PostForm(loginURL, url.Values{
		"__VIEWSTATE":          []string{viewstate},
		"__VIEWSTATEGENERATOR": []string{viewstateGenerator},
		"login_name":           []string{synergyUser},
		"password":             []string{synergyPassword},
	})
	if err != nil {
		return nil, err
	}

	if !isLoginSuccess(loginResponse) {
		return nil, fmt.Errorf("Login Unsuccessfull")
	}
	loginBody, err := readClose(loginResponse)
	if err != nil {
		return nil, err
	}
	// FocusKey set package-level in case of future uses outside of next 3:
	focusKey := parseSubmatch(reFocusKey, loginBody)
	if focusKey == "" {
		return nil, fmt.Errorf("Could not create a focus key")
	}
	return &AuthClient{c: c, focusKey: focusKey, wait: wait}, nil
}

// notifyFinish takes a job guid and sends it back thru the guid chan when the job is ready to download
func (ac *AuthClient) notifyFinish(guid string, guidChan chan string) error {

	if guid == "" {
		return fmt.Errorf("notifyFinished called without guid")
	}
	revJobQueueGetStatus = setFocusKey(revJobQueueGetStatus, ac.focusKey)
	revJobQueueGetStatus = setJobGUID(revJobQueueGetStatus, guid)
	formValues := url.Values{"xml": []string{revJobQueueGetStatus}}
	// if jobStatusFinished break the loop to get result
	// if chan is closed, break the loop to return an error
	for {

		res, err := ac.c.PostForm(xmlDoRequestURL, formValues)
		if err != nil {
			return err
		}
		body, err := readClose(res)
		if err != nil {
			return err
		}

		if jobStatusFinished(guid, body) {
			guidChan <- guid
			break
		}
		time.Sleep(time.Second)
	}
	return nil

}

// getJob takes a Rev_Queue_ReportJob xml string (with a set focus key)
// and returns a guid which is empty if an error
func (ac *AuthClient) requestJobGUID(xmlRequest string) (res []byte, err error) {
	xmlRequest = setFocusKey(xmlRequest, ac.focusKey)
	// get jobguid
	r, err := ac.c.PostForm(xmlDoRequestURL, url.Values{"xml": []string{xmlRequest}})
	if err != nil {
		log.Printf("requestJobGUID: PostForm error %v", err)
		return nil, err
	}
	return readClose(r)
}

// DownloadCurrentStu415s returns a parsed stu415 report downloaded from Synergy or an error if actions fail
func (ac *AuthClient) DownloadCurrentStu415s() (stu415s types.Stu415s, err error) {
	res, err := ac.requestJobGUID(setFocusKey(requestStu415Today, ac.focusKey))

	if err != nil {
		return stu415s, err
	}

	guid := parseSubmatch(reStu415GUID, res)
	if guid == "" {
		return stu415s, fmt.Errorf("DownloadS1Stu415s: requestJobGUID got empty return")
	}

	guidChan := make(chan string, 1)
	go ac.notifyFinish(guid, guidChan)
	err = ac.downloadWhenFinished(guidChan)
	if err != nil {
		log.Println("DownloadS1Stu415s: downloadWhenFinishedErr")
	}

	b, err := ac.getReport(guid)
	if err != nil {
		log.Println("DownloadStu415s: getReport err")
		return stu415s, err
	}
	return types.Stu415sFromCSV(bytes.NewBuffer(b))
}

func (ac *AuthClient) downloadWhenFinished(guidChan chan string) error {
	killCh := make(chan bool, 1)
	time.AfterFunc(ac.wait, func() { killCh <- true })
	for {
		select {
		case guid := <-guidChan:
			return ac.requestFinishedJob(guid)

		case <-killCh:
			close(guidChan)
			return fmt.Errorf("Timeout: Could not download roster in %v", ac.wait)

		default:
			time.Sleep(time.Second)
		}
	}
}
func readClose(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func (ac *AuthClient) requestFinishedJob(jobGUID string) (err error) {
	emailGetResults = setFocusKey(emailGetResults, ac.focusKey)
	emailGetResults = setJobGUID(emailGetResults, jobGUID)

	_, err = ac.c.PostForm(xmlDoRequestURL, url.Values{"xml": []string{emailGetResults}})
	return
}

func (ac *AuthClient) getReport(jobGUID string) (csv []byte, err error) {

	// Prepare the file location

	_, err = ac.c.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	res, err := ac.c.Get(fmt.Sprintf("https://synergy.aps.edu/ReportOutput/%s.CSV", jobGUID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func jobStatusFinished(guid string, body []byte) bool {
	pattern := fmt.Sprintf(`<ROW GUID="%s" State="4"`, guid)
	foundState4, _ := regexp.Match(pattern, body)
	return foundState4

}

func parseMatch(re *regexp.Regexp, body []byte) (string, error) {
	matches := re.Find(body)
	if len(matches) == 0 {
		return "", fmt.Errorf("parseMatch did not find any matches")
	}
	return string(matches), nil
}

func parseSubmatch(re *regexp.Regexp, body []byte) string {
	sm := re.FindSubmatch(body)
	if len(sm) > 1 {
		return string(sm[1])
	}
	return ""

}

func isLoginSuccess(loginResponse *http.Response) bool {
	return loginResponse.Request.URL.String() == "https://synergy.aps.edu/ST_Content.aspx"
}
