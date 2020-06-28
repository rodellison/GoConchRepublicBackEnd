package common

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
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
}

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
		ed.ImgURL = os.Getenv("URLBASE") + iQuery.Nodes[0].Attr[1].Val
	}

	iQuery = s.Find(LISTING_LOCATION)
	if iQuery.Nodes != nil {
		var locationData []string
		locationData = strings.Split(iQuery.Nodes[0].LastChild.LastChild.Data, ": ")
		ed.EventLocation = locationData[1]
	}

	return nil

}
