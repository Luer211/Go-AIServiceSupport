package request

type CreateTaskRequest struct {
	Prompt string `json:"prompt" binding:"required,min=1,max=5000"`
}

type GetTaskStatusRequest struct {
	TaskID string `uri:"task_id" binding:"required"`
}
