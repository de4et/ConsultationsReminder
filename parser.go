package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/takuoki/clmconv"
	"github.com/xuri/excelize/v2"
)

type Day struct {
	name string
	num  int
}

type Time struct {
	Hour   int
	Minute int
}

type Lesson struct {
	Name string

	Day   Day
	Start Time
	End   Time
}

type Group struct {
	Number  string
	Lessons []Lesson

	colNumber int
}

var dayStN = map[string]int{ // String to number converter
	"понедельник": 1,
	"вторник":     2,
	"среда":       3,
	"четверг":     4,
	"пятница":     5,
	"суббота":     6,
	"воскресенье": 0,
}

func initGroups(rows [][]string) map[string]*Group {
	const groupNumbersRow = 1 // #2

	groups := make(map[string]*Group)
	for i, col := range rows[groupNumbersRow] {
		if col != "" {
			groups[col] = &Group{Number: col, colNumber: i}
		}
	}
	return groups
}

func replaceBlankSymbols(text string) string {
	return strings.Replace(strings.Replace(text, " ", "", -1), "\n", "", -1)

}

func readCell(file *excelize.File, rowN int, colN int) string {
	value, err := file.GetCellValue(studyType, clmconv.Itoa(colN)+strconv.Itoa(rowN+1))
	if err != nil {
		panic(err)
	}
	return value
}

func parseTimePeriod(text string) (start Time, end Time) {
	tmp := strings.Split(text, "-")
	startT := strings.Split(tmp[0], ".")
	endT := strings.Split(tmp[1], ".")

	startHT, _ := strconv.Atoi(startT[0])
	startMT, _ := strconv.Atoi(startT[1])
	endHT, _ := strconv.Atoi(endT[0])
	endMT, _ := strconv.Atoi(endT[1])
	return Time{startHT, startMT}, Time{endHT, endMT}
}

func ParseFile(name string) map[string]*Group {
	// opening schedule file (download from google.com? itf)
	file, err := excelize.OpenFile(name)
	if err != nil {
		log.Fatal(err)
	}

	// reading all rows from file
	rows, err := file.GetRows(studyType)
	if err != nil {
		log.Fatal(err)
	}
	rowsLen := len(rows)

	// init groups, filling Number and colNumber
	groups := initGroups(rows)

	// for rectal use only)
	const (
		scheduleStartRow = 2 // #3
		scheduleStartCol = 2 // #3
		dayStrsCol       = 0 // #1
		timePeriodsCol   = 1 // #2
	)

	for rowN := scheduleStartRow; rowN < rowsLen; rowN++ {
		dayStr := replaceBlankSymbols(readCell(file, rowN, dayStrsCol)) // convert to Day
		timePeriodStr := readCell(file, rowN, timePeriodsCol)           // parse
		_ = timePeriodStr

		for k, g := range groups {
			t := readCell(file, rowN, g.colNumber)
			timeStart, timeEnd := parseTimePeriod(timePeriodStr)
			if t != "" {
				groups[k].Lessons = append(groups[k].Lessons, Lesson{
					Name: t,
					Day: Day{
						name: dayStr,
						num:  dayStN[dayStr],
					},
					Start: timeStart,
					End:   timeEnd,
				})
			}
		}
	}
	// for _, v := range groups {
	// 	fmt.Println(v)
	// 	fmt.Println("---")
	// }
	// fmt.Printf("\n%+v", groups)
	fmt.Println(groups["11-402"])
	return groups

}
