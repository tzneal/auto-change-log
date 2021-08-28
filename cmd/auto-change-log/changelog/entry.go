package changelog

type EntryType byte

const (
	AddedEntry EntryType = iota
	ChangedEntry
	DeprecatedEntry
	RemovedEntry
	FixedEntry
	SecurityEntry
)

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
