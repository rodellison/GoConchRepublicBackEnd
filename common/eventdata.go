package common

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Eventdata struct {
	EventID          string `json:"eventid"`
	EventStartDate   string `json:"eventstartdate"`
	EventEndDate     string `json:"eventenddate"`
	EventName        string `json:"eventname"`
	EventContact     string `json:"eventcontact"`
	EventLocation    string `json:"eventlocation"`
	EventImgURL      string `json:"eventimgurl"`
	EventURL         string `json:"eventurl"`
	EventDescription string `json:"eventdescription"`
}

//ToDo: set URLBASE as a Env value
const URLBASE = "https://fla-keys.com"
const URLBASE2 = "/calendar/all/florida-keys/"

//These constants are for HTML parsing
const LISTING_BLOCK = ".listing-block.listing-calendar"

const LISTING_IMG = ".swipebox.expand-img"
const LISTING_DESCRIPTION = ".listing-desc"
const LISTING_LOCATION = ".listing-location"
const LISTING_DATE = ".listing-date"
const LISTING_NAME = ".listing-name"
const LISTING_PHONE = ".listing-phone"


func (ed *Eventdata) ExtractEventData(i int, s *goquery.Selection) (err error) {
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
				ed.EventStartDate = startDate
				ed.EventEndDate = startDate
			} else {
				ed.EventStartDate = startDate
				ed.EventEndDate = endDate
			}
		}
	} else {
		return errors.New("could not format Event dates")
	}

	iQuery = s.Find(LISTING_NAME)
	if iQuery.Nodes != nil {
		ed.EventURL = iQuery.Nodes[0].FirstChild.Attr[1].Val
		ed.EventName = iQuery.Nodes[0].FirstChild.FirstChild.Data
	}

	iQuery = s.Find(LISTING_DESCRIPTION)
	if iQuery.Nodes != nil {
		ed.EventDescription = iQuery.Nodes[0].FirstChild.Data
	}

	iQuery = s.Find(LISTING_PHONE)
	if iQuery.Nodes != nil {
		ed.EventContact = iQuery.Nodes[0].LastChild.Data
	}

	iQuery = s.Find(LISTING_IMG)
	if iQuery.Nodes != nil {
		ed.EventImgURL = URLBASE +iQuery.Nodes[0].Attr[1].Val
	}

	iQuery = s.Find(LISTING_LOCATION)
	if iQuery.Nodes != nil {
		var locationData []string
		locationData = strings.Split(iQuery.Nodes[0].LastChild.LastChild.Data, ": ")
		ed.EventLocation = locationData[1]
	}

	return nil

}
