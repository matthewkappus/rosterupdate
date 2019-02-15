package synergy

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"time"
)

var reEmailGUID = regexp.MustCompile(`<REV_ELEMENT>(.{36})</REV_ELEMENT>`)

// DownloadEmails returns a csv slice or an error if Synergy does not return csv
func (ac *AuthClient) DownloadEmails() (emails [][]string, err error) {
	guid, err := ac.requestEmailGUID()
	if err != nil {
		return nil, err
	}
	res, err := ac.c.Get(fmt.Sprintf("https://synergy.aps.edu/ReportOutput/%s.TXT", guid))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	r := csv.NewReader(res.Body)
	r.LazyQuotes = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		emails = append(emails, record)
	}

	return emails, nil
}

// requestEmailGUID returns a guid or an error if failure in 5 second
// It uploads postEmailParams (http) then an xml request emailGetProperties to activate email report
func (ac *AuthClient) requestEmailGUID() (guid string, err error) {
	res, err := ac.c.PostForm("https://synergy.aps.edu/ST_UploadFile.aspx", url.Values{"data": []string{setFocusKey(postEmailParams, ac.focusKey)}})
	if err != nil {
		return "", err
	}

	body, err := readClose(res)
	if err != nil {
		return "", err
	}

	guid = parseSubmatch(reEmailGUID, body)
	if len(guid) != 36 {
		return "", fmt.Errorf("Did not get 36-char guid")
	}

	res, err = ac.c.PostForm("https://synergy.aps.edu/Service/RTCommunication.asmx/XMLDoRequest", url.Values{"xml": []string{setFocusKey(emailGetProperties, ac.focusKey)}})
	if err != nil {
		return "", err
	}


	time.Sleep(time.Second)
	getEmailResults = setJobGUID(getEmailResults, guid)
	getEmailResults = setFocusKey(getEmailResults, ac.focusKey)
	res, err = ac.c.PostForm("https://synergy.aps.edu/Service/RTCommunication.asmx/XMLDoRequest", url.Values{"xml": []string{getEmailResults}})
	if err != nil {
		return "", err
	}

	return guid, nil
}

var (
	postEmailParams = `<REV_REQUEST><EVENT NAME="Rev_Do_Command"><REQUEST><REV_DATA_ROOT VIEW_GUID="E51430E2-DFD1-4348-9266-BDCFB437820D" ACTION="COMMAND" PRIMARY_OBJECT="7F746134-5F24-4958-BA04-1EB42C44632E" VIEW_TYPE="BOUND" REV_VIEW_TYPE="REV_QUERY" CUR_TAB_GUID="52037271-8916-4C2E-B208-C9CF0B74412C" BUTTON_ID="EXECUTE_BUTTON" BUTTON_OBJ="" BUTTON_TEXT="Execute" BUTTON_URL="WebData.aspx" VIEW_ID="E51430E2-DFD1-4348-9266-BDCFB437820D" BUTTON_OPEN_TYPE="0" FOCUS_KEY="{{.FocusKey}}" FRAME="0"><REV_DATA_ROOT FOCUS_KEY="{{.FocusKey}}" VIEW_TYPE="BOUND" VIEW_GUID="E51430E2-DFD1-4348-9266-BDCFB437820D" ORIGINAL_VIEW_GUID="E51430E2-DFD1-4348-9266-BDCFB437820D" ACTION="SAVE" PRIMARY_OBJECT="7F746134-5F24-4958-BA04-1EB42C44632E" REV_VIEW_TYPE="REV_QUERY" CUR_TAB_GUID="52037271-8916-4C2E-B208-C9CF0B74412C" ORDER="1"><REV_DATA_REQUEST><REV_VIEW GUID="E51430E2-DFD1-4348-9266-BDCFB437820D"><REV_TAB GUID="52037271-8916-4C2E-B208-C9CF0B74412C"/></REV_VIEW><REV_ELEMENT ALIAS="Name" SRC_NAME="Revelation-Query-RevQuery-Name" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Name">staff_emails</REV_ELEMENT><REV_ELEMENT ALIAS="Group" SRC_NAME="Revelation-Query-RevQuery-Group" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Group">STAFF</REV_ELEMENT><REV_ELEMENT ALIAS="Type" SRC_NAME="Revelation-Query-RevQuery-Type" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Type">Select</REV_ELEMENT><REV_ELEMENT ALIAS="OutputType" SRC_NAME="Revelation-Query-RevQuery-OutputType" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="OutputType">CSV</REV_ELEMENT><REV_ELEMENT ALIAS="Orientation" SRC_NAME="Revelation-Query-RevQuery-Orientation" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Orientation">Portrait</REV_ELEMENT><REV_ELEMENT ALIAS="QueryType" SRC_NAME="Revelation-Query-RevQuery-QueryType" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="QueryType">User</REV_ELEMENT><REV_ELEMENT ALIAS="Template" SRC_NAME="Revelation-Query-RevQuery-Template" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Template"></REV_ELEMENT><REV_ELEMENT ALIAS="DelimeterDD" SRC_NAME="Revelation-Query-RevQuery-DelimeterDD" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="DelimeterDD">Comma</REV_ELEMENT><REV_ELEMENT ALIAS="DelimeterOther" SRC_NAME="Revelation-Query-RevQuery-DelimeterOther" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="DelimeterOther"></REV_ELEMENT><REV_ELEMENT ALIAS="SuppressHeader" SRC_NAME="Revelation-Query-RevQuery-SuppressHeader" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="SuppressHeader">N</REV_ELEMENT><REV_ELEMENT ALIAS="FixedLength" SRC_NAME="Revelation-Query-RevQuery-FixedLength" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="FixedLength">N</REV_ELEMENT><REV_ELEMENT ALIAS="Description" SRC_NAME="Revelation-Query-RevQuery-Description" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Description">name, email</REV_ELEMENT><REV_ELEMENT ALIAS="MyRatingValue" SRC_NAME="Revelation-Query-RevQuery-MyRatingValue" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="MyRatingValue"></REV_ELEMENT><REV_ELEMENT ALIAS="QueryText" SRC_NAME="Revelation-Query-RevQuery-QueryText" SRC_OBJECT="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="QueryText"></REV_ELEMENT><REV_ELEMENT ALIAS="LabelSelect" SRC_NAME="Revelation-Reports-ReportUI-LabelSelect" SRC_OBJECT="867F81B0-C2F8-44C9-9BFE-5F8B38315DF5" SRC_ELEMENT="LabelSelect"></REV_ELEMENT><REV_ELEMENT ALIAS="TopMargin" SRC_NAME="Revelation-Reports-Label-TopMargin" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="TopMargin"></REV_ELEMENT><REV_ELEMENT ALIAS="SideMargin" SRC_NAME="Revelation-Reports-Label-SideMargin" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="SideMargin"></REV_ELEMENT><REV_ELEMENT ALIAS="VerticalPitch" SRC_NAME="Revelation-Reports-Label-VerticalPitch" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="VerticalPitch"></REV_ELEMENT><REV_ELEMENT ALIAS="HorizontalPitch" SRC_NAME="Revelation-Reports-Label-HorizontalPitch" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="HorizontalPitch"></REV_ELEMENT><REV_ELEMENT ALIAS="LabelHeight" SRC_NAME="Revelation-Reports-Label-LabelHeight" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="LabelHeight"></REV_ELEMENT><REV_ELEMENT ALIAS="LabelWidth" SRC_NAME="Revelation-Reports-Label-LabelWidth" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="LabelWidth"></REV_ELEMENT><REV_ELEMENT ALIAS="NumberAcross" SRC_NAME="Revelation-Reports-Label-NumberAcross" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="NumberAcross"></REV_ELEMENT><REV_ELEMENT ALIAS="NumberDown" SRC_NAME="Revelation-Reports-Label-NumberDown" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="NumberDown"></REV_ELEMENT><REV_ELEMENT ALIAS="PageSizeGU" SRC_NAME="Revelation-Reports-Label-PageSizeGU" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="PageSizeGU"></REV_ELEMENT><REV_ELEMENT ALIAS="PageOrientation" SRC_NAME="Revelation-Reports-Label-PageOrientation" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="PageOrientation"></REV_ELEMENT><REV_ELEMENT ALIAS="RowHeight" SRC_NAME="Revelation-Reports-Label-RowHeight" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="RowHeight"></REV_ELEMENT><REV_ELEMENT ALIAS="RowSpace" SRC_NAME="Revelation-Reports-Label-RowSpace" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="RowSpace"></REV_ELEMENT><REV_ELEMENT ALIAS="ScaleFields" SRC_NAME="Revelation-Reports-Label-ScaleFields" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="ScaleFields"></REV_ELEMENT><REV_ELEMENT ALIAS="FontSize" SRC_NAME="Revelation-Reports-Label-FontSize" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="FontSize"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_RecurType" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_RecurType" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_RecurType"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_StartTime" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_StartTime" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_StartTime"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_StartDate" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_StartDate" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_StartDate"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_StopDate" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_StopDate" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_StopDate"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_DayCount" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_DayCount" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_DayCount"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_WeekCount" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_WeekCount" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_WeekCount"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Monday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Monday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Monday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Tuesday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Tuesday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Tuesday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Wednesday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Wednesday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Wednesday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Thursday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Thursday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Thursday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Friday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Friday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Friday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Saturday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Saturday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Saturday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_Sunday" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_Sunday" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_Sunday"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_MonthType" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_MonthType" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_MonthType"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_MonthDayofMonth" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_MonthDayofMonth" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_MonthDayofMonth"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_MonthDayType" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_MonthDayType" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_MonthDayType"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_MonthDayOfWeek" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_MonthDayOfWeek" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_MonthDayOfWeek"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_January" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_January" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_January"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_February" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_February" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_February"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_March" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_March" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_March"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_April" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_April" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_April"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_May" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_May" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_May"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_June" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_June" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_June"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_July" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_July" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_July"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_August" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_August" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_August"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_September" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_September" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_September"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_October" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_October" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_October"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_November" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_November" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_November"></REV_ELEMENT><REV_ELEMENT ALIAS="Z_December" SRC_NAME="Revelation-JobQueueInfo-JobQueueRecur-Z_December" SRC_OBJECT="56D84CE5-2FDB-4F23-953B-E7C6C0E96050" SRC_ELEMENT="Z_December"></REV_ELEMENT><REV_ELEMENT SRC_OBJECT="RevQuery" SRC_NAME="Revelation-Query-RevQuery-GUID" SRC_ELEMENT="GUID" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E">BF397580-721B-44AF-8A94-686EBBF1491B</REV_ELEMENT><REV_ELEMENT ALIAS="ShowAllBO" SRC_NAME="Revelation-Query-RevQuery-ShowAllBO" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="ShowAllBO"></REV_ELEMENT><REV_ELEMENT ALIAS="ShowAllProperties" SRC_NAME="Revelation-Query-RevQuery-ShowAllProperties" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="ShowAllProperties"></REV_ELEMENT><REV_ELEMENT ALIAS="QueryXML" SRC_NAME="Revelation-Query-RevQuery-QueryXML" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="QueryXML"></REV_ELEMENT><REV_ELEMENT ALIAS="EditableResults" SRC_NAME="Revelation-Query-RevQuery-EditableResults" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="EditableResults"></REV_ELEMENT></REV_DATA_REQUEST><IDENTITY>
	<REV_ELEMENT SRC_OBJECT="RevQuery" SRC_NAME="Revelation-Query-RevQuery-GUID" SRC_ELEMENT="GUID" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E">BF397580-721B-44AF-8A94-686EBBF1491B</REV_ELEMENT>
  </IDENTITY></REV_DATA_ROOT><GROUP_FIELDS></GROUP_FIELDS><IDENTITY>
	<REV_ELEMENT SRC_OBJECT="RevQuery" SRC_NAME="Revelation-Query-RevQuery-GUID" SRC_ELEMENT="GUID" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E">BF397580-721B-44AF-8A94-686EBBF1491B</REV_ELEMENT>
  </IDENTITY><REV_DATA_GROUP><REV_ELEMENT ALIAS="ShowAllBO" SRC_NAME="Revelation-Query-RevQuery-ShowAllBO" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="ShowAllBO"></REV_ELEMENT><REV_ELEMENT ALIAS="ShowAllProperties" SRC_NAME="Revelation-Query-RevQuery-ShowAllProperties" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="ShowAllProperties"></REV_ELEMENT><REV_ELEMENT ALIAS="QueryXML" SRC_NAME="Revelation-Query-RevQuery-QueryXML" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="QueryXML"></REV_ELEMENT><REV_ELEMENT ALIAS="OutputType" SRC_NAME="Revelation-Query-RevQuery-OutputType" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="OutputType">CSV</REV_ELEMENT><REV_ELEMENT ALIAS="Orientation" SRC_NAME="Revelation-Query-RevQuery-Orientation" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="Orientation">Portrait</REV_ELEMENT><REV_ELEMENT ALIAS="LabelSelect" SRC_NAME="Revelation-Reports-ReportUI-LabelSelect" SRC_OBJECT="867F81B0-C2F8-44C9-9BFE-5F8B38315DF5" SRC_ELEMENT="LabelSelect"></REV_ELEMENT><REV_ELEMENT ALIAS="LabelHeight" SRC_NAME="Revelation-Reports-Label-LabelHeight" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="LabelHeight"></REV_ELEMENT><REV_ELEMENT ALIAS="LabelWidth" SRC_NAME="Revelation-Reports-Label-LabelWidth" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="LabelWidth"></REV_ELEMENT><REV_ELEMENT ALIAS="NumberAcross" SRC_NAME="Revelation-Reports-Label-NumberAcross" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="NumberAcross"></REV_ELEMENT><REV_ELEMENT ALIAS="NumberDown" SRC_NAME="Revelation-Reports-Label-NumberDown" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="NumberDown"></REV_ELEMENT><REV_ELEMENT ALIAS="PageOrientation" SRC_NAME="Revelation-Reports-Label-PageOrientation" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="PageOrientation"></REV_ELEMENT><REV_ELEMENT ALIAS="PageSizeGU" SRC_NAME="Revelation-Reports-Label-PageSizeGU" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="PageSizeGU"></REV_ELEMENT><REV_ELEMENT ALIAS="SideMargin" SRC_NAME="Revelation-Reports-Label-SideMargin" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="SideMargin"></REV_ELEMENT><REV_ELEMENT ALIAS="TopMargin" SRC_NAME="Revelation-Reports-Label-TopMargin" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="TopMargin"></REV_ELEMENT><REV_ELEMENT ALIAS="VerticalPitch" SRC_NAME="Revelation-Reports-Label-VerticalPitch" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="VerticalPitch"></REV_ELEMENT><REV_ELEMENT ALIAS="HorizontalPitch" SRC_NAME="Revelation-Reports-Label-HorizontalPitch" SRC_OBJECT="823DC988-BA33-4022-B15A-D350B79228B2" SRC_ELEMENT="HorizontalPitch"></REV_ELEMENT><REV_ELEMENT ALIAS="MyRatingValue" SRC_NAME="Revelation-Query-RevQuery-MyRatingValue" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="MyRatingValue"></REV_ELEMENT><REV_ELEMENT ALIAS="EditableResults" SRC_NAME="Revelation-Query-RevQuery-EditableResults" SRC_OBJECT="RevQuery" SRC_OBJECT_GUID="7F746134-5F24-4958-BA04-1EB42C44632E" SRC_ELEMENT="EditableResults"></REV_ELEMENT></REV_DATA_GROUP><CLIENT_STATE><CLIENT_ACTION TYPE="SAVE_PARENT_COMMAND" BUTTON_ID="EXECUTE_BUTTON" ELEMENT_ID="REV_BUTTON"></CLIENT_ACTION></CLIENT_STATE></REV_DATA_ROOT></REQUEST></EVENT><QUERY COMMUNITY="N" FIXEDLENGTH="N" GROUP="STAFF" SUPPRESSHEADER="N" ORIENTATION="Portrait" SAVEDQUERYTYPE="User" NAME="staff_emails" EXPORTFILTERTYPE="CSV" DELIMETERDD="Comma" QUERYTYPE="Select" GUID="F3A25400-A67E-43F8-A68D-8E9FC42E4337">
	<DESCRIPTION>name, email</DESCRIPTION>
	<BO ID="E0E8FF04-A965-4365-A5DF-238CD2A76FF5" BOID="52D78195-371A-4A32-ADF4-4C19AA7CED7B" NAMEORIGINAL="Staff" NAME="Staff" ALIAS="R0" NAMESPACE="K12"><PROPERTY SRCELEMENT="Email" ID="E6A988D1-1E4C-4BA5-9956-F14024ABB03D" ALIAS="R0" ORDER="1"/><PROPERTY SRCELEMENT="FormattedName" ID="F53A7946-42D5-4C8D-BC06-F4C460009636" ALIAS="R0" ORDER="2"/></BO>
	<LABELDEF/>
  </QUERY></REV_REQUEST>
	`

	emailGetProperties = `<?xml version="1.0" encoding="utf-8"?>
	<REV_REQUEST><EVENT NAME="Query_Get_BOProperties"><REQUEST FOCUS_KEY="{{.FocusKey}}" BOID="52D78195-371A-4A32-ADF4-4C19AA7CED7B" WINDOW_ID="d90cb432-6dce-48e1-9ad1-7fa3d823b48d"></REQUEST></EVENT></REV_REQUEST>`

	getEmailResults = `<?xml version="1.0" encoding="utf-8"?><REV_REQUEST><EVENT NAME="JobQueue_Get_Results"><REQUEST FOCUS_KEY="{{.FocusKey}}" JOB_GUID="{{.JobGUID}}" FILE_GUID="" WINDOW_ID="d90cb432-6dce-48e1-9ad1-7fa3d823b48d"><SERVER_STATE><D K="DebugGroupGU" V=""/></SERVER_STATE></REQUEST></EVENT></REV_REQUEST>`
)