package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredService  []ServiceName
	ServiceUpdateURL string
	HeartBeatURL     string
}
type ServiceName string

const (
	LogService   = ServiceName("LogService")
	GradeService = ServiceName("GradeService")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
