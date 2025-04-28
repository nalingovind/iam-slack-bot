# AWS Slack IAM Provisioning Bot

A Slack-driven workflow and approval system that grants time-limited AWS IAM access, emails credentials via SES, and automatically revokes them upon expiration.

## Summary

This repository implements a fully automated, approval-based AWS credential provisioning system using Slack Workflow Builder, a Go backend service, Python Lambdas (Boto3), Amazon SES for emailing credentials, and EventBridge Scheduler for auto-revocation. The system routes requests per project ownership, enforces least-privilege, tracks state in DynamoDB, and is deployed via Infrastructure as Code.

## Prerequisites

- AWS account with permissions to create IAM roles, Lambda functions, EventBridge rules, SES identities, and DynamoDB tables
- Verified SES "From" email address or domain in your AWS region
- Slack workspace with admin access to install and configure Slack apps and Workflows
- Go 1.18+ installed locally
- Python 3.9+ for AWS Lambda development
- Terraform (or AWS CDK/SAM) for IaC

## Repository Structure

```text
├── infrastructure/         # Terraform/SAM/CDK definitions
├── slack-workflow/         # Documentation on Workflow Builder setup
├── backend-go/             # Go service handling Slack Webhooks & Interactivity
├── lambda-provisioning/    # Python (Boto3) Lambda for creating credentials & SES emailing
├── lambda-revoke/          # Python Lambda for periodic deprovisioning
└── README.md               # Project overview and setup instructions
```

## 1. Slack Workflow Setup

1. Open **Slack > Workflow Builder** and create a new workflow named **Request AWS Access**.
2. **Trigger:** Global Shortcut (e.g., “Request AWS Access”).
3. **Form Step:** Collect the following fields:
   - **Project** (dropdown or short text)
   - **Role** (dropdown: e.g., ReadOnly, PowerUser)
   - **Duration** (short text, e.g., `2h`, `1d`)
   - **Justification** (optional long text)
4. **Send Web Request:** Configure to POST to your Go service endpoint (`/workflows/aws-access`) with JSON body mapping form inputs:
   ```json
   {
     "user_id": "{{triggered_by.id}}",
     "user_name": "{{triggered_by.name}}",
     "project": "{{form.project}}",
     "role": "{{form.role}}",
     "duration": "{{form.duration}}",
     "justification": "{{form.justification}}"
   }
   ```
5. Publish the workflow.

## 2. Backend Service (Go)

### 2.1 Initialize Module

```bash
go mod init github.com/yourorg/aws-slackbot
go get github.com/slack-go/slack
```

### 2.2 Environment Variables

- `SLACK_SIGNING_SECRET`
- `SLACK_BOT_TOKEN`
- `AWS_REGION`
- `DYNAMODB_TABLE_PROJECTS`
- `DYNAMODB_TABLE_REQUESTS`

### 2.3 Workflow Webhook Handler

Implement an HTTP POST handler at `/workflows/aws-access`:

```go
func handleWorkflow(w http.ResponseWriter, r *http.Request) {
  // 1. Parse JSON payload
  // 2. Validate project & role
  // 3. Lookup project owner Slack ID from DynamoDB
  // 4. Write PENDING request to DynamoDB
  // 5. Respond 200 OK with {"outputs":{}}
  // 6. Asynchronously post approval message
}
```

### 2.4 Approval Message

Use Slack Block Kit to send to the project owner’s channel:

```go
msg := slack.MsgOptionBlocks(
  slack.NewSectionBlock(...),
  slack.NewActionBlock("aws_actions",
    slack.NewButtonBlockElement("approve", requestID, ...).WithStyle(slack.StylePrimary),
    slack.NewButtonBlockElement("deny", requestID, ...).WithStyle(slack.StyleDanger),
  ),
)
api.PostMessage(ownerChannelID, msg)
```

### 2.5 Interaction Endpoint

Handle button clicks at `/slack/interactions`:

```go
func handleInteraction(w http.ResponseWriter, r *http.Request) {
  // 1. Parse slack.InteractionCallback
  // 2. Switch on ActionID: "approve" or "deny"
  // 3. Update DynamoDB Status
  // 4. If approved: invoke provisioning Lambda
  // 5. Notify requester and approver
}
```

## 3. Provisioning Lambda (Python + Boto3)

### 3.1 Setup

- Use AWS SAM or Terraform to define the Lambda.
- Include dependencies: `boto3`, `slack-sdk` (optional for Slack callbacks).

### 3.2 Handler Logic

```python
def handler(event, context):
    # 1. Read request from DynamoDB
    # 2. AssumeRole or create IAM user + access key
    # 3. Store credentials & expiration back to DynamoDB
    # 4. Send credentials via SES
    # 5. Optionally post to Slack
```

### 3.3 SES Email

Use Boto3 SES:

```python
ses = boto3.client('ses')
ses.send_email(
    Source='no-reply@yourdomain.com',
    Destination={'ToAddresses':[requester_email]},
    Message={...}
)
```

Email should include:

- Access Key ID
- Secret Access Key
- Role session ARN or console link
- Expiration timestamp

## 4. Auto-Revocation Lambda

1. **EventBridge Rule:** Schedule every 5 minutes.
2. **Lambda Handler:** Scan DynamoDB for expired APPROVED items:
   ```python
   def handler(event, context):
       # Query items where now >= ExpirationTimestamp
       # Call delete_access_key or RemoveRolePolicy
       # Update Status=EXPIRED
       # Notify via Slack/email (optional)
   ```

## 5. Infrastructure as Code

Define all resources in Terraform (recommended) or AWS SAM/CDK:

- DynamoDB tables: `projects`, `requests`
- IAM roles and policies for Go service and Lambdas
- Lambda functions and triggers
- EventBridge Scheduler rule
- SES identity and settings
- API Gateway (if using for webhooks)

## 6. Testing & Monitoring

- Unit tests for Go handlers and Python Lambdas
- Integration tests in a sandbox account
- CloudWatch Alarms for:
  - Lambda errors
  - SES bounces
  - DynamoDB throttling
- Slack/email alerts for critical failures

## 7. Contributing

1. Fork the repo and create a feature branch.
2. Run tests locally:
   ```bash
   go test ./backend-go/...
   pytest lambda-provisioning/
   ```

```
3. Submit a Pull Request.

## 8. License
This project is licensed under the MIT License. Please see [LICENSE](LICENSE) for details.

```
