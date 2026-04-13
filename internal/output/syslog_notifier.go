package output

import (
	"fmt"
	"log/syslog"

	"github.com/user/portwatch/internal/alert"
)

// SyslogNotifier sends port change events to the system syslog daemon.
type SyslogNotifier struct {
	writer   *syslog.Writer
	formatter *Formatter
}

// NewSyslogNotifier creates a SyslogNotifier that writes to syslog with the
// given priority and tag. It returns an error if the syslog connection fails.
func NewSyslogNotifier(priority syslog.Priority, tag string) (*SyslogNotifier, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog_notifier: connect: %w", err)
	}
	return &SyslogNotifier{
		writer:    w,
		formatter: NewFormatter(FormatText, false),
	}, nil
}

// NewSyslogNotifierWithWriter creates a SyslogNotifier using an existing
// syslog.Writer. Useful for testing or pre-configured writers.
func NewSyslogNotifierWithWriter(w *syslog.Writer) *SyslogNotifier {
	return &SyslogNotifier{
		writer:    w,
		formatter: NewFormatter(FormatText, false),
	}
}

// Notify sends each event as a separate syslog message. Events with action
// "opened" are logged at warning level; closed ports at info level.
func (s *SyslogNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}
	for _, ev := range events {
		line := s.formatter.FormatOne(ev)
		var err error
		if ev.Action == alert.ActionOpened {
			err = s.writer.Warning(line)
		} else {
			err = s.writer.Info(line)
		}
		if err != nil {
			return fmt.Errorf("syslog_notifier: write: %w", err)
		}
	}
	return nil
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}
