package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal server error")
	ErrBusyServer          = errors.New("Server is busy at the moment")
	ErrAddingImpossible    = errors.New("Archiver is proccesing -- adding new files is impossible")
	ErrTaskNotFound        = errors.New("Task is not found")

	ErrVersionConflict = errors.New("Version conflict")

	ErrLoadFile   = errors.New("Load file error")
	ErrUnsuppType = errors.New("File has unsupportable type")
)
