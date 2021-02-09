package common

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Eventdata struct {
	EventID          string `json:"EventID"`
	StartDate        string `json:"StartDate"`
	EndDate          string `json:"EndDate"`
	EventName        string `json:"EventName"`
	EventContact     string `json:"EventContact"`
	EventLocation    string `json:"EventLocation"`
	ImgURL           string `json:"ImgURL"`
	EventURL         string `json:"EventURL"`
	EventDescription string `json:"EventDescription"`
	EventExpiry      int64  `json:"EventExpiry"`
}

//These constants are for HTML parsing
const LISTING_BLOCK = ".listing-block.listing-calendar"

const LISTING_IMG = ".swipebox.expand-img"
const LISTING_DESCRIPTION = ".listing-desc"
const LISTING_LOCATION = ".listing-location"
const LISTING_DATE = ".listing-date"
const LISTING_NAME = ".listing-name"
const LISTING_PHONE = ".listing-phone"

func makeSSMLCompatible(descriptionText string) (string) {

	returnText := descriptionText
	returnText = strings.ReplaceAll(returnText, "& ", "and ")
	returnText = strings.ReplaceAll(returnText, "&nbsp;", " ")
	returnText = strings.ReplaceAll(returnText, "<a href=\"", " ")
	returnText = strings.ReplaceAll(returnText, "\">", " ")
	returnText = strings.ReplaceAll(returnText, "here</a>", " ")
	returnText = strings.ReplaceAll(returnText, "</a>", " ")

	return returnText;

}


func (ed *Eventdata) ExtractEventData(s *goquery.Selection) (err error) {
	// Load the HTML document

	var iQuery *goquery.Selection

	fmt.Println("")
	//We're sitting on the base Node for all other nodes for this Event. The Attr ID is what we can use for the unique
	//Event ID
	ed.EventID = s.Nodes[0].Attr[1].Val

	iQuery = s.Find(LISTING_DATE)
	if iQuery.Nodes != nil {
		startDate, endDate, err := FormatEventDates(iQuery.Nodes[0].LastChild.Data)
		if err != nil {
			return errors.New("could not format Event dates")
		} else {
			if startDate == endDate {
				ed.StartDate = startDate
				ed.EndDate = startDate
			} else {
				ed.StartDate = startDate
				ed.EndDate = endDate
			}
		}
	} else {
		return errors.New("could not format Event dates")
	}

	iQuery = s.Find(LISTING_NAME)
	if iQuery.Nodes != nil {

		if iQuery.Nodes[0].FirstChild.Data == "a"  {    //This entry has an EventURL
			ed.EventURL = iQuery.Nodes[0].FirstChild.Attr[1].Val
			ed.EventName = makeSSMLCompatible(iQuery.Nodes[0].FirstChild.FirstChild.Data)
		} else {   //This entry does NOT have an EventURL
			ed.EventURL = " "  //Place a space in EventURL so DynamoDB doesn't null the value
			ed.EventName = makeSSMLCompatible(iQuery.Nodes[0].FirstChild.Data)
		}

	}

	iQuery = s.Find(LISTING_DESCRIPTION)
	if iQuery.Nodes != nil {
		ed.EventDescription = makeSSMLCompatible(iQuery.Nodes[0].FirstChild.Data)
	}

	iQuery = s.Find(LISTING_PHONE)
	if iQuery.Nodes != nil {
		ed.EventContact = iQuery.Nodes[0].LastChild.Data
	} else {
		ed.EventContact = "No contact phone provided"
	}

	iQuery = s.Find(LISTING_IMG)
	if iQuery.Nodes != nil {
		ed.ImgURL = iQuery.Nodes[0].Attr[1].Val
	}

	iQuery = s.Find(LISTING_LOCATION)
	if iQuery.Nodes != nil {
		var locationData []string
		locationData = strings.Split(iQuery.Nodes[0].LastChild.LastChild.Data, ": ")
		if (len(locationData) > 1) {
			ed.EventLocation = strings.ToLower(strings.TrimLeft(locationData[1], " "))
		} else {
			ed.EventLocation = strings.ToLower(strings.TrimLeft(locationData[0], " "))
		}
		ed.EventLocation = strings.Replace(ed.EventLocation, "the ", "", 1)
		ed.EventLocation = strings.Replace(ed.EventLocation, " ", "-", 1)


	}

	return nil

}
