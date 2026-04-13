package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Format controls the output format style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter converts alert events into human-readable or structured strings.
type Formatter struct {
	format    Format
	timestamp bool
}

// NewFormatter creates a Formatter with the given format and timestamp option.
func NewFormatter(format Format, timestamp bool) *Formatter {
	return &Formatter{format: format, timestamp: timestamp}
}

// FormatEvent returns a formatted string representation of a single event.
func (f *Formatter) FormatEvent(e alert.Event) string {
	if f.format == FormatJSON {
		return f.formatJSON(e)
	}
	return f.formatText(e)
}

func (f *Formatter) formatText(e alert.Event) string {
	var sb strings.Builder
	if f.timestamp {
		sb.WriteString(time.Now().UTC().Format(time.RFC3339))
		sb.WriteString(" ")
	}
	sb.WriteString(fmt.Sprintf("[%s] %s/%d", strings.ToUpper(e.Kind), e.Port.Protocol, e.Port.Number))
	if e.Port.Process != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", e.Port.Process))
	}
	return sb.String()
}

func (f *Formatter) formatJSON(e alert.Event) string {
	ts := ""
	if f.timestamp {
		ts = fmt.Sprintf(`"timestamp":%q,`, time.Now().UTC().Format(time.RFC3339))
	}
	return fmt.Sprintf(
		`{%s"kind":%q,"protocol":%q,"port":%d,"process":%q}`,
		ts, e.Kind, e.Port.Protocol, e.Port.Number, e.Port.Process,
	)
}
