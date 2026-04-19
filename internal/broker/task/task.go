package task

type Created struct {
	TaskID  string         `json:"task_id"`
	Type    string         `json:"type"`
	Payload CreatedPayload `json:"payload"`
}

type CreatedPayload struct {
	Msg string `json:"msg"`
}

type Progress struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

type Log struct {
	TaskID  string `json:"task_id"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Result struct {
	TaskID string `json:"task_id"`
	Result string `json:"result"`
}

type Error struct {
	TaskID string `json:"task_id"`
	Error  string `json:"error"`
}
