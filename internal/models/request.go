package models

// WorkflowRequest represents incoming form data
// from Slack Workflow Builder
type WorkflowRequest struct {
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	Project       string `json:"project"`
	Role          string `json:"role"`
	Duration      string `json:"duration"`
	Justification string `json:"justification"`
}
