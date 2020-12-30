// Copyright 2012 Sonia Keys
// License: MIT

// Package julian Chapter 7, Julian day.
//
// Under "General remarks", Meeus describes the INT function as used in the
// book.  In some contexts, math.Floor might be suitable, but I think more
// often, the functions FloorDiv or FloorDiv64 will be more
// appropriate.  See documentation in package
//
// On p. 63, Modified Julian Day is defined.  See constant JMod in package
//
//
// See also related functions JulianToGregorian and GregorianToJulian in
// package jm.
package julian

import (
	"math"
	"time"
)

// ConvertDate returns the Julian Date for a given Gregorian or Julian
// calendar date and time, the time is given as a decimal of the day.
func ConvertDate(year, month int, day float64) float64 {

	// year, month, day := timeToJulianTime(t)

	if month < 3 {
		year--
		month = month + 12
	}

	// Convert the date into YYYY.MMDDdd format
	date := float64(year) + ((float64(month) + day) / 100)

	// Calculate the Julian date.
	a := int64(365.25 * float64(year))
	b := int64(30.6001 * float64(month+1))
	c := day + 1720994.5
	JD := float64(a) + float64(b) + c

	// If the number is larger than or equal to 1582.1015 that is within
	// the Gregorian calendar.
	if date >= 1582.1015 {
		A := int64(year / 100)
		B := 2 - A + int64(A/4)
		JD = JD + float64(B)
	}

	// If YYYY.MMDDdd is smaller than 1582.1015 then it is not necessary to
	// calculate A and B.

	return JD
}

// CalendarGregorianToJD converts a Gregorian year, month, and day of month
// to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
func CalendarGregorianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	a := FloorDiv(y, 100)
	b := 2 - a + FloorDiv(a, 4)
	// (7.1) p. 61
	return float64(FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(FloorDiv(306*(m+1), 10)+b) + d - 1524.5
}

// CalendarJulianToJD converts a Julian year, month, and day of month to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
func CalendarJulianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	return float64(FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(FloorDiv(306*(m+1), 10)) + d - 1524.5
}

// LeapYearJulian returns true if year y in the Julian calendar is a leap year.
func LeapYearJulian(y int) bool {
	return y%4 == 0
}

// LeapYearGregorian returns true if year y in the Gregorian calendar is a leap year.
func LeapYearGregorian(y int) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

// JDToCalendar returns the calendar date for the given jd.
//
// Note that this function returns a date in either the Julian or Gregorian
// Calendar, as appropriate.
func JDToCalendar(jd float64) (year, month int, day float64) {
	zf, f := math.Modf(jd + .5)
	z := int64(zf)
	a := z
	if z >= 2299151 {
		α := FloorDiv64(z*100-186721625, 3652425)
		a = z + 1 + α - FloorDiv64(α, 4)
	}
	b := a + 1524
	c := FloorDiv64(b*100-12210, 36525)
	d := FloorDiv64(36525*c, 100)
	e := int(FloorDiv64((b-d)*1e4, 306001))
	// compute return values
	day = float64(int(b-d)-FloorDiv(306001*e, 1e4)) + f
	switch e {
	default:
		month = e - 1
	case 14, 15:
		month = e - 13
	}
	switch month {
	default:
		year = int(c) - 4716
	case 1, 2:
		year = int(c) - 4715
	}
	return
}

// jdToCalendarGregorian returns the Gregorian calendar date for the given jd.
//
// Note that it returns a Gregorian date even for dates before the start of
// the Gregorian calendar.  The function is useful when working with Go
// time.Time values because they are always based on the Gregorian calendar.
func jdToCalendarGregorian(jd float64) (year, month int, day float64) {
	zf, f := math.Modf(jd + .5)
	z := int64(zf)
	α := FloorDiv64(z*100-186721625, 3652425)
	a := z + 1 + α - FloorDiv64(α, 4)
	b := a + 1524
	c := FloorDiv64(b*100-12210, 36525)
	d := FloorDiv64(36525*c, 100)
	e := int(FloorDiv64((b-d)*1e4, 306001))
	// compute return values
	day = float64(int(b-d)-FloorDiv(306001*e, 1e4)) + f
	switch e {
	default:
		month = e - 1
	case 14, 15:
		month = e - 13
	}
	switch month {
	default:
		year = int(c) - 4716
	case 1, 2:
		year = int(c) - 4715
	}
	return
}

// JDToTime takes a JD and returns a Go time.Time value.
func JDToTime(jd float64) time.Time {
	// time.Time is always Gregorian
	y, m, d := jdToCalendarGregorian(jd)
	t := time.Date(y, time.Month(m), 0, 0, 0, 0, 0, time.UTC)
	return t.Add(time.Duration(d * 24 * float64(time.Hour)))
}

// TimeToJD takes a Go time.Time and returns a JD as float64.
//
// Any time zone offset in the time.Time is ignored and the time is
// treated as UTC.
func TimeToJD(t time.Time) float64 {
	ut := t.UTC()
	y, m, _ := ut.Date()
	d := ut.Sub(time.Date(y, m, 0, 0, 0, 0, 0, time.UTC))
	// time.Time is always Gregorian
	return CalendarGregorianToJD(y, int(m), float64(d)/float64(24*time.Hour))
}

// TimeToJD2 test
func TimeToJD2(t time.Time) float64 {
	ut := t.UTC()
	y, m, _ := ut.Date()
	d := ut.Sub(time.Date(y, m, 0, 0, 0, 0, 0, time.UTC))
	// time.Time is always Gregorian
	return ConvertDate(y, int(m), float64(d)/float64(24*time.Hour))
}

// DateNow returns the current Julian date.
func DateNow() float64 {
	return TimeToJD(time.Now())
}

// DayOfWeek determines the day of the week for a given JD.
//
// The value returned is an integer in the range 0 to 6, where 0 represents
// Sunday.  This is the same convention followed in the time package of the
// Go standard library.
func DayOfWeek(jd float64) int {
	return int(jd+1.5) % 7
}

// DayOfYearGregorian computes the day number within the year of the Gregorian
// calendar.
func DayOfYearGregorian(y, m, d int) int {
	return DayOfYear(y, m, d, LeapYearGregorian(y))
}

// DayOfYearJulian computes the day number within the year of the Julian
// calendar.
func DayOfYearJulian(y, m, d int) int {
	return DayOfYear(y, m, d, LeapYearJulian(y))
}

// DayOfYear computes the day number within the year.
//
// This form of the function is not specific to the Julian or Gregorian
// calendar, but you must tell it whether the year is a leap year.
func DayOfYear(y, m, d int, leap bool) int {
	k := 2
	if leap {
		k--
	}
	return wholeMonths(m, k) + d
}

// DayOfYearToCalendar returns the calendar month and day for a given
// day of year and leap year status.
func DayOfYearToCalendar(n int, leap bool) (m, d int) {
	k := 2
	if leap {
		k--
	}
	if n < 32 {
		m = 1
	} else {
		m = (900*(k+n) + 98*275) / 27500
	}
	return m, n - wholeMonths(m, k)
}

func wholeMonths(m, k int) int {
	return 275*m/9 - k*((m+9)/12) - 30
}

/* ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *  Math
 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ */

// FloorDiv returns the integer floor of the fractional value (x / y).
//
// It uses integer math only, so is more efficient than using floating point
// intermediate values.  This function can be used in many places where INT()
// appears in AA.  As with built in integer division, it panics with y == 0.
func FloorDiv(x, y int) (q int) {
	q = x / y
	if (x < 0) != (y < 0) && x%y != 0 {
		q--
	}
	return
}

// FloorDiv64 returns the integer floor of the fractional value (x / y).
//
// It uses integer math only, so is more efficient than using floating point
// intermediate values.  This function can be used in many places where INT()
// appears in AA.  As with built in integer division, it panics with y == 0.
func FloorDiv64(x, y int64) (q int64) {
	q = x / y
	if (x < 0) != (y < 0) && x%y != 0 {
		q--
	}
	return
}
