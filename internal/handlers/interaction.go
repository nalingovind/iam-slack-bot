package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nalingovind/iam-slack-bot/internal/dynamo"
	"github.com/nalingovind/iam-slack-bot/internal/slack"
	slackapi "github.com/slack-go/slack"
)

// HandleInteraction processes Slack button clicks (approve/deny)
func HandleInteraction(w http.ResponseWriter, r *http.Request) {
	payloadStr := r.FormValue("payload")
	var payload slackapi.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		http.Error(w, "invalid interaction payload", http.StatusBadRequest)
		return
	}

	action := payload.ActionCallback.BlockActions[0]
	requestID := action.Value
	userID := payload.User.ID

	var newStatus, replyText string
	if action.ActionID == "approve_btn" {
		newStatus = "APPROVED"
		replyText = ":white_check_mark: Approved by <@" + userID + ">"
		// TODO: invoke provisioning Lambda here
	} else {
		newStatus = "DENIED"
		replyText = ":x: Denied by <@" + userID + ">"
	}

	// Update DynamoDB
	_ = dynamo.UpdateRequestStatus(r.Context(), requestID, newStatus)

	// Post feedback in channel
	slack.Client.PostMessage(payload.Channel.ID, slackapi.MsgOptionText(replyText, false))

	w.WriteHeader(http.StatusOK)
}
