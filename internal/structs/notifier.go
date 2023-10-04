package structs

type Notifier struct {
	Status   string // info | succeeded | failed
	Logs     string
	Details  string
	Metadata string
}
