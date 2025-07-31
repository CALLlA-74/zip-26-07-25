package domain

type CreateTaskResponse struct {
	TaskUuid string `json:"taskUuid" validate:"required"`
}

type AddLinksRequest struct {
	Links []string `json:"links" validate:"required"`
}

type AddLinksResponse struct {
	AddedLinks      []string `json:"addedLinks" validate:"required"`
	HasReachedLimit bool     `json:"hasReachedLimit" validate:"required"`
}

type TaskStatusResponse struct {
	Status      TStatuses `json:"status" validate:"required"`
	FailedLinks []string  `json:"failedLinks" validate:"required"`
	ArchiveLink string    `json:"archiveLink"`
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
}
