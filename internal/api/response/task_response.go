package response

type CreateTaskResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

type GetTaskStatusResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}
