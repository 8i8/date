// Copyright 2012 Sonia Keys
// License: MIT

package date

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func ExampleCalendarGregorianToJD_sputnik() {
	// Example 7.a, p. 61.
	jd := CalendarGregorianToJD(1957, 10, 4.81)
	fmt.Printf("%.2f\n", jd)
	// Output:
	// 2436116.31
}

func ExampleCalendarGregorianToJD_halley() {
	// Example 7.c, p. 64.
	jd1 := CalendarGregorianToJD(1910, 4, 20)
	jd2 := CalendarGregorianToJD(1986, 2, 9)
	fmt.Printf("%.0f days\n", jd2-jd1)
	// Output:
	// 27689 days
}

func TestMyFunc(t *testing.T) {
	t1 := time.Date(1000, 10, 4, 22, 30, 0, 0, time.UTC)
	jd1 := TimeToJD(t1)
	jd2 := TimeToJD2(t1)

	if jd1 == jd2 {
		t.Errorf("TestMyFunc: recieved %v expected %v", jd2, jd1)
	}
}

func TestGreg(t *testing.T) {
	for _, tp := range []struct {
		y, m  int
		d, jd float64
	}{
		{2000, 1, 1.5, 2451545}, // more examples, p. 62
		{1999, 1, 1, 2451179.5},
		{1987, 1, 27, 2446822.5},
		{1987, 6, 19.5, 2446966},
		{1988, 1, 27, 2447187.5},
		{1988, 6, 19.5, 2447332},
		{1900, 1, 1, 2415020.5},
		{1600, 1, 1, 2305447.5},
		{1600, 12, 31, 2305812.5},
	} {
		dt := CalendarGregorianToJD(tp.y, tp.m, tp.d) - tp.jd
		if math.Abs(dt) > .1 {
			t.Logf("%#v", tp)
			t.Fatal("dt:", time.Duration(dt*24*float64(time.Hour)))
		}
	}
}

func ExampleCalendarJulianToJD() {
	// Example 7.b, p. 61.
	jd := CalendarJulianToJD(333, 1, 27.5)
	fmt.Printf("%.1f\n", jd)
	// Output:
	// 1842713.0
}

func TestJuli(t *testing.T) {
	for _, tp := range []struct {
		y, m  int
		d, jd float64
	}{
		{837, 4, 10.3, 2026871.8}, // more examples, p. 62
		{-123, 12, 31, 1676496.5},
		{-122, 1, 1, 1676497.5},
		{-1000, 7, 12.5, 1356001},
		{-1000, 2, 29, 1355866.5},
		{-1001, 8, 17.9, 1355671.4},
		{-4712, 1, 1.5, 0},
	} {
		dt := CalendarJulianToJD(tp.y, tp.m, tp.d) - tp.jd
		if math.Abs(dt) > .1 {
			t.Logf("%#v", tp)
			t.Fatal("dt:", time.Duration(dt*24*float64(time.Hour)))
		}
	}
}

func TestJuliLeap(t *testing.T) {
	for _, tp := range []struct {
		year int
		leap bool
	}{
		{900, true},
		{1236, true},
		{750, false},
		{1429, false},
	} {
		if LeapYearJulian(tp.year) != tp.leap {
			t.Logf("%#v", tp)
			t.Fatal("JuliLeapYear")
		}
	}
}

func TestGregLeap(t *testing.T) {
	for _, tp := range []struct {
		year int
		leap bool
	}{
		{1700, false},
		{1800, false},
		{1900, false},
		{2100, false},
		{1600, true},
		{2000, true},
		{2400, true},
	} {
		if LeapYearGregorian(tp.year) != tp.leap {
			t.Logf("%#v", tp)
			t.Fatal("JuliLeapYear")
		}
	}
}

func ExampleJDToCalendar() {
	// Example 7.c, p. 64.
	y, m, d := JDToCalendar(2436116.31)
	fmt.Printf("%d %s %.2f\n", y, time.Month(m), d)
	// Output:
	// 1957 October 4.81
}

func TestYMD(t *testing.T) {
	for _, tp := range []struct {
		jd   float64
		y, m int
		d    float64
	}{
		{1842713, 333, 1, 27.5},
		{1507900.13, -584, 5, 28.63},
	} {
		y, m, d := JDToCalendar(tp.jd)
		if y != tp.y || m != tp.m || math.Abs(d-tp.d) > .01 {
			t.Logf("%#v", tp)
			t.Fatal("JDToYMD", y, m, d)
		}
	}
}

func ExampleDayOfWeek() {
	// Example 7.e, p. 65.
	fmt.Println(time.Weekday(DayOfWeek(2434923.5)))
	// Output:
	// Wednesday
}

func ExampleDayOfYear_f() {
	// Example 7.f, p. 65.
	fmt.Println(DayOfYear(1978, 11, 14, false))
	// Output:
	// 318
}

func ExampleDayOfYear_g() {
	// Example 7.g, p. 65.
	fmt.Println(DayOfYear(1988, 4, 22, true))
	// Output:
	// 113
}

func TestDOYToCal(t *testing.T) {
	for _, tp := range []struct {
		y, m, d int
		leap    bool
		doy     int
	}{
		// same data as examples above
		{1978, 11, 14, false, 318},
		{1988, 4, 22, true, 113},
	} {
		m, d := DayOfYearToCalendar(tp.doy, tp.leap)
		if m != tp.m || d != tp.d {
			t.Logf("%#v", tp)
			t.Fatal("DayOfYearToCalendar", m, d)
		}
	}
}

func BenchmarkTimeToJD(b *testing.B) {
	t1 := time.Date(1000, 10, 4, 22, 30, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		TimeToJD(t1)
	}
}

func BenchmarkTimeToJD2(b *testing.B) {
	t1 := time.Date(1000, 10, 4, 22, 30, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		TimeToJD2(t1)
	}
}
