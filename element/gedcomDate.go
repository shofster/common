package element

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"strings"
)

/*

  File:    gedcomDate.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Model and Validation of GEDCOM Date formats.
*/

// find first 4 digit number
var yearOnly = regexp.MustCompile("\\d{4}")

func GetYear(datePhrase string) string {
	year := yearOnly.FindString(datePhrase)
	return year
}

// 3 forms of a GEDCOM date:
var dateMod = regexp.MustCompile("(?i)(ABT|BEF|AFT|EST)\\s+(.*)")
var dateRange = regexp.MustCompile("(?i)BET\\s*(.*)AND\\s*(.*)")
var dateExact = regexp.MustCompile("(\\w+)\\s*(\\w+)?\\s*(\\w+)?")

//goland:noinspection GoUnusedExportedFunction
func NewGedcomDateEntry(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	//	entry.PlaceHolder = "Enter Date in form [[dd ]Mon ]yyyy"
	entry.PlaceHolder = placeholder
	entry.Text = ""
	entry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		if GetYear(s) == "" {
			return errors.New("missing year")
		}
		c, err := GetGedcomString(s)
		if err != nil {
			return err
		}
		if c != s { // replace with pretty date
			entry.Text = c
		}
		return nil
	}
	return entry
}

func GetGedcomString(dateString string) (string, error) {
	t := strings.TrimSpace(strings.ToLower(dateString))
	// at least a year is present
	p := t[0:3]
	switch p {
	case "abt", "bef", "aft", "est":
		// see if date Modifier
		modDates := dateMod.FindStringSubmatch(dateString)
		if modDates != nil {
			// a match will always have 3: original, mod, something?
			// mod is one of ABT|BEF|AFT|EST
			if g, err := getExactGedcomString(modDates[2]); err == nil {
				return fmt.Sprintf("%s %s", strings.ToUpper(modDates[1]), g), nil
			}
		}
		return dateString, errors.New("incomplete mod type")
	case "bet":
		// see if a BET / AND phrase
		fromTo := dateRange.FindStringSubmatch(t)
		if fromTo != nil {
			// a match will always have 3: original, from, to
			// BET and AND are used and ignored
			if from, ef := getExactGedcomString(fromTo[1]); ef == nil {
				if to, et := getExactGedcomString(fromTo[2]); et == nil {
					return fmt.Sprintf("BET %s AND %s", from, to), nil
				}
			}
			return dateString, errors.New("incomplete range type")
		}
		return dateString, errors.New("incomplete range type")
	}
	// a 1 to 3  (dd mmm yyyy) part GEDCOM date
	return getExactGedcomString(t)
}
func getExactGedcomString(dateString string) (string, error) {
	dateString = strings.TrimSpace(dateString)
	if dateString == "" {
		return "", errors.New("empty")
	}
	if GetYear(dateString) == "" {
		return "", errors.New("missing year")
	}
	dates := dateExact.FindStringSubmatch(dateString)
	if dates == nil {
		return "", errors.New("invalid")
	}
	n := 0
	for _, v := range dates {
		if v != "" {
			n++
		}
	}
	switch n {
	case 2:
		y := strings.TrimSpace(dates[1])
		return fmt.Sprintf("%4s", y), nil
	case 3:
		m := strings.TrimSpace(dates[1])
		y := strings.TrimSpace(dates[2])
		return fmt.Sprintf("%3s %4s", month(m), y), nil
	case 4:
		d := strings.TrimSpace(dates[1])
		m := strings.TrimSpace(dates[2])
		y := strings.TrimSpace(dates[3])
		return fmt.Sprintf("%s %3s %4s", d, month(m), y), nil
	}
	return dateString, nil
}
func month(s string) string {
	m := strings.TrimSpace(strings.ToUpper(s))
	switch {
	case strings.HasPrefix(m, "JA"):
		return "Jan"
	case strings.HasPrefix(m, "F"):
		return "Feb"
	case strings.HasPrefix(m, "MAR"):
		return "Mar"
	case strings.HasPrefix(m, "AP"):
		return "Apr"
	case strings.HasPrefix(m, "MAY"):
		return "May"
	case strings.HasPrefix(m, "JUN"):
		return "Jun"
	case strings.HasPrefix(m, "JUL"):
		return "Jul"
	case strings.HasPrefix(m, "AU"):
		return "Aug"
	case strings.HasPrefix(m, "S"):
		return "Sep"
	case strings.HasPrefix(m, "O"):
		return "Oct"
	case strings.HasPrefix(m, "N"):
		return "Nov"
	case strings.HasPrefix(m, "D"):
		return "Dec"
	}
	return m
}
