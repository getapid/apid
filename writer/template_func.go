package writer

import (
	"strings"
	"time"

	"github.com/fatih/color"
)

func increment(a int) int {
	return a + 1
}

func indent(indent int, txt string) string {
	prefix := strings.Repeat(" ", indent)
	return strings.Replace(txt, "\n", "\n"+prefix, -1)
}

func renderTime(t time.Duration) string {
	t = t.Round(time.Millisecond)
	if t <= 0 {
		t = 0
	}
	return t.String()
}

func add(nums ...int) (sum int) {
	for _, n := range nums {
		sum += n
	}
	return
}

func red(txt string) string {
	return color.New(color.FgHiRed, color.Bold).Sprint(txt)
}

func green(txt string) string {
	return color.New(color.FgHiGreen, color.Bold).Sprint(txt)
}

func bold(txt string) string {
	return color.New(color.FgHiCyan, color.Bold).Sprint(txt)
}
