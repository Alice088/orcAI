package agent

type Stop struct {
	AgentID string `json:"agent_id"`
}

type Stopped struct {
	AgentID string `json:"agent_id"`
	Reason  string `json:"reason"`
}

type Health struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}
