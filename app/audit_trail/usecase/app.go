package usecase

import (
	"user-management/app/audit_trail/usecase/command"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddAuditTrail command.AddAuditTrailHandler
}

type Queries struct {
}
