package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func CalcLongEpochFromEndDate(year, month, day int) int64 {

	future := time.Date(year, time.Month(month), day, 23, 59, 59, 0, time.Local)
	return future.Unix()

}

func CalcSearchYYYYMMFromDate(month string) string {

	monthIncrement, err := strconv.Atoi(month)
	if err != nil {
		fmt.Println("calcSearchYYYYMMFromDate: Could not convert month to YYYYMM date")
		monthIncrement = 1
	}
	dateToFetch := time.Now().AddDate(0, monthIncrement-1, 0).String()
	dateValues := strings.SplitN(dateToFetch, "-", 3)
	return dateValues[0] + dateValues[1]

}

func GetFormattedDateToday() string {

	dateToday := time.Now().String()
	dateToday = string(dateToday[0:4]) + string(dateToday[5:7]) + string(dateToday[8:10])

	return dateToday

}

func convertShortDate(inDate string) string {

	//Dates coming in take the form Jan 20, 2020 00:00:00
	//so just need to convert the initial month characters

	inDate = strings.Replace(inDate, "Jan ", "January ", 1)
	inDate = strings.Replace(inDate, "Feb ", "February ", 1)
	inDate = strings.Replace(inDate, "Mar ", "March ", 1)
	inDate = strings.Replace(inDate, "Apr ", "April ", 1)
	//May doesnt need converted
	inDate = strings.Replace(inDate, "Jun ", "June ", 1)
	inDate = strings.Replace(inDate, "Jul ", "July ", 1)
	inDate = strings.Replace(inDate, "Aug ", "August ", 1)
	inDate = strings.Replace(inDate, "Sep ", "September ", 1)
	inDate = strings.Replace(inDate, "Oct ", "October ", 1)
	inDate = strings.Replace(inDate, "Nov ", "November ", 1)
	inDate = strings.Replace(inDate, "Dec ", "December ", 1)

	return inDate
}

func FormatEventDates(dateString string) (startDate, endDate string, err error) {

	var splitDates []string
	splitDates = strings.Split(dateString, " - ")
	startDate = convertShortDate(splitDates[0])
	t, _ := time.Parse(layoutUS, startDate)
	startDate = t.String()
	startDate = string(startDate[0:4]) + string(startDate[5:7]) + string(startDate[8:10])

	if len(splitDates) > 1 {
		endDate = convertShortDate(splitDates[1])
		t, _ = time.Parse(layoutUS, endDate)
		endDate = t.String()
		endDate = string(endDate[0:4]) + string(endDate[5:7]) + string(endDate[8:10])
	} else {
		endDate = startDate
	}
	return startDate, endDate, nil

}
