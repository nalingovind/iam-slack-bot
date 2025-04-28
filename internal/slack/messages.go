package slack

import (
	"fmt"

	"github.com/nalingovind/iam-slack-bot/internal/models"
	"github.com/slack-go/slack"
)

// SendApprovalMessage sends an approval prompt to the owner
func SendApprovalMessage(ownerID string, req models.WorkflowRequest, requestID string) error {
	text := fmt.Sprintf(`*AWS Access Request* by <@%s>
• *Project*: %s
• *Role*: %s
• *Duration*: %s
• *Justification*: %s`,
		req.UserID, req.Project, req.Role, req.Duration, req.Justification)

	blocks := []slack.Block{
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", text, false, false), nil, nil),
		slack.NewActionBlock("aws_request_actions",
			slack.NewButtonBlockElement("approve_btn", requestID, slack.NewTextBlockObject("plain_text", "Approve", false, false)).WithStyle(slack.StylePrimary),
			slack.NewButtonBlockElement("deny_btn", requestID, slack.NewTextBlockObject("plain_text", "Deny", false, false)).WithStyle(slack.StyleDanger),
		),
	}
	_, _, err := Client.PostMessage(ownerID, slack.MsgOptionBlocks(blocks...))
	return err
}
