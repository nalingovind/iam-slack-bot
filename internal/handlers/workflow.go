package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nalingovind/iam-slack-bot/internal/dynamo"
	"github.com/nalingovind/iam-slack-bot/internal/models"
	slackclient "github.com/nalingovind/iam-slack-bot/internal/slack"
)

// HandleWorkflow processes incoming Slack Workflow webhooks
func HandleWorkflow(w http.ResponseWriter, r *http.Request) {
	var req models.WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Lookup project owner
	ownerID, err := dynamo.LookupProjectOwner(r.Context(), req.Project)
	if err != nil {
		http.Error(w, "project lookup failed", http.StatusInternalServerError)
		return
	}

	// Create request in DynamoDB
	requestID, err := dynamo.CreateRequest(r.Context(), req, ownerID)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	// Acknowledge Slack Workflow step
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"outputs":{}}`))

	// Send approval prompt asynchronously
	go func() {
		_ = slackclient.SendApprovalMessage(ownerID, req, requestID)
	}()
}
