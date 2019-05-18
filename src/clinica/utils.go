package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func fmtDateTime(dt, layout string) (dts string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	r := regexp.MustCompile(`\d+`)
	m := r.FindAllString(dt, -1)
	if len(m) < 3 {
		panic(fmt.Errorf("date/time string too short"))
	}
	p := func(idx int) string {
		if idx >= len(m) {
			return "00"
		}
		s := m[idx]
		v, _ := strconv.Atoi(s)
		switch idx {
		case 0:
			if v < 1900 || v > 2100 {
				panic(fmt.Errorf("year out-of-bound"))
			}
			return s
		case 1:
			if v < 1 || v > 12 {
				panic(fmt.Errorf("month out-of-bound"))
			}
			return fmt.Sprintf("%02d", v)
		case 2:
			if v < 1 || v > 31 {
				panic(fmt.Errorf("day out-of-bound"))
			}
			return fmt.Sprintf("%02d", v)
		case 3:
			if v < 0 || v > 23 {
				panic(fmt.Errorf("hour out-of-bound"))
			}
			return fmt.Sprintf("%02d", v)
		case 4:
			if v < 0 || v > 59 {
				panic(fmt.Errorf("minute out-of-bound"))
			}
			return fmt.Sprintf("%02d", v)
		case 5:
			if v < 0 || v > 59 {
				panic(fmt.Errorf("second out-of-bound"))
			}
			return fmt.Sprintf("%02d", v)
		}
		return ""
	}
	dts = strings.Replace(layout, "Y", p(0), 1)
	dts = strings.Replace(dts, "m", p(1), 1)
	dts = strings.Replace(dts, "d", p(2), 1)
	dts = strings.Replace(dts, "H", p(3), 1)
	dts = strings.Replace(dts, "i", p(4), 1)
	dts = strings.Replace(dts, "s", p(5), 1)
	return
}
