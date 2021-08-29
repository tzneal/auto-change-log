package changelog

import "fmt"

type EntryType byte

const (
	AddedEntry EntryType = iota
	ChangedEntry
	DeprecatedEntry
	RemovedEntry
	FixedEntry
	SecurityEntry
)

func (e *EntryType) UnmarshalYAML(b []byte) error {
	s := string(b)
	switch s {
	case "Added":
		*e = AddedEntry
	case "Changed":
		*e = ChangedEntry
	case "Deprecated":
		*e = DeprecatedEntry
	case "Removed":
		*e = RemovedEntry
	case "Fixed":
		*e = FixedEntry
	case "Security":
		*e = SecurityEntry
	default:
		return fmt.Errorf("unknown value '%s'", s)
	}
	return nil
}

func (e EntryType) MarshalYAML() (interface{}, error) {
	return e.String(), nil
}
func (e EntryType) String() string {
	switch e {
	case AddedEntry:
		return "Added"
	case ChangedEntry:
		return "Changed"
	case DeprecatedEntry:
		return "Deprecated"
	case RemovedEntry:
		return "Removed"
	case FixedEntry:
		return "Fixed"
	case SecurityEntry:
		return "Security"
	}
	return "Invalid"
}

type Entry struct {
	Type    EntryType
	Summary string
}
