package handle

import (
	"fmt"
	"github.com/jiangzhx/receive/util/date"
	"testing"
)

func Test_heartbeatHandle(t *testing.T) {
	when := "2006-01-02 15:04:05"
	cm := date.Parse(when).Minute()
	cm = cm / 5 * 5
	minute := "00"
	if cm < 10 {
		minute = fmt.Sprintf("0%d", cm)
	} else {
		minute = fmt.Sprintf("%d", cm)
	}
	minute = fmt.Sprintf("%s:%sZ", date.Format(when, "2006-01-02T15"), minute)

	println(minute)

}
