package domain

type TStatuses string

const (
	WAITING_LINKS TStatuses = "WAITING_LINKS"
	PROCESSING    TStatuses = "PROCESSING"
	FINISHED      TStatuses = "FINISHED"
)

type Task struct {
	TaskUuid    string
	TaskStatus  TStatuses
	Files       []*File
	ArchivePath string
	ArchiveLink string
}

type File struct {
	Link         string
	Path         string
	ErrorMessage string
}
