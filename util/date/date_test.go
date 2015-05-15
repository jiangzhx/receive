package date

import (
	"strings"
	"testing"
)

func Test_Format(t *testing.T) {
	d := Format("2015-11-22 09:23:14", "2006-01-02T15Z")
	if strings.EqualFold("2015-11-22T09Z", d) {
		t.Log(d)
	} else {
		t.Error("format error!!")
	}
}

func Test_PreDay(t *testing.T) {

	d := PreDay("2015-11-22 09:23:14", -1)

	if strings.EqualFold(d.Format("2006-01-02 15:04:05"), "2015-11-21 09:23:14") {
		t.Log(d)
	} else {
		t.Error("format error!!" + d.Format("2006-01-02 15:04:05"))
	}
}
