package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func parse(s string) (time.Time, error) {
	if s == "" || s == "now" {
		return time.Now(), nil
	}

	if strings.HasPrefix(s, "now-") {
		d, err := time.ParseDuration(s[4:])
		if err != nil {
			return time.Time{}, err
		}
		return time.Now().Add(-d), nil

	}

	if d, err := time.ParseDuration(s); err == nil {
		return time.Now().Add(d), nil
	}

	return dateparse.ParseLocal(s)
}

var stepFlag = flag.Duration("s", 10*time.Minute, "step")

var formatFlag = flag.String("f", "s", "format")

var delimiterFlag = flag.String("d", " ", "delimiter")

func format(t time.Time) string {
	switch *formatFlag {
	case "s":
		return strconv.FormatInt(t.Unix(), 10)
	case "ms":
		return strconv.FormatInt(t.UnixMilli(), 10)
	case "rfc":
		return t.Format(time.RFC3339)
	case "dt":
		return t.Format(time.DateTime)
	default:
		panic("unknown format")
	}
}

func main() {
	flag.Parse()
	if !flag.Parsed() {
		panic("arg error")
	}
	start, err := parse(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	end, err := parse(flag.Arg(1))
	if err != nil {
		panic(err)
	}

	step := *stepFlag

	delimiter := *delimiterFlag
	if delimiter == "tab" {
		delimiter = "\t"
	}

	for ts := start.Round(step); ts.Before(end); ts = ts.Add(step) {
		t1 := ts
		t2 := ts.Add(step)
		if t2.After(end) {
			t2 = end
		}

		fmt.Printf("%s%s%s\n", format(t1), delimiter, format(t2))

	}

}
