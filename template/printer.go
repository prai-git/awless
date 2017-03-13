package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/oklog/ulid"
)

type Printer struct {
	Sep         string
	IncludeErrs bool
	OnOK        func(string) string
	OnKO        func(string) string
}

func NewPrinter() *Printer {
	noop := func(s string) string { return s }

	return &Printer{
		Sep:  "\t",
		OnOK: noop,
		OnKO: noop,
	}
}

func (p *Printer) PrintReport(t *Template) string {
	var buff bytes.Buffer

	buff.WriteString(fmt.Sprintf("Date: %s", parseULIDDate(t.ID)))

	if IsRevertible(t) {
		buff.WriteString(fmt.Sprintf(", RevertID: %s", t.ID))
	} else {
		buff.WriteString(", RevertID: <not revertible>")
	}
	buff.WriteString("\n")

	w := tabwriter.NewWriter(&buff, 0, 8, 0, '\t', 0)
	for _, cmd := range t.CommandNodesIterator() {
		var result, status string

		exec := fmt.Sprint("%s", cmd.String())

		if cmd.CmdErr != nil {
			status = "KO"
			if p.IncludeErrs {
				exec = fmt.Sprint("%s\n%s", exec, formatMultiLineErrMsg(cmd.CmdErr.Error()))
			}
		} else {
			status = "OK"
		}

		if v, ok := cmd.CmdResult.(string); ok && v != "" {
			result = v
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t\n", status, result, cmd.String())
	}

	w.Flush()

	return buff.String()
}

func formatMultiLineErrMsg(msg string) string {
	notabs := strings.Replace(msg, "\t", "", -1)
	var indented []string
	for _, line := range strings.Split(notabs, "\n") {
		indented = append(indented, fmt.Sprintf("\t\t%s", line))
	}
	return strings.Join(indented, "\n")
}

func parseULIDDate(uid string) string {
	parsed, err := ulid.Parse(uid)
	if err != nil {
		panic(err)
	}

	date := time.Unix(int64(parsed.Time())/int64(1000), time.Nanosecond.Nanoseconds())

	return date.Format(time.Stamp)
}
