package adapters

import (
	"context"
)

type CommandRepository interface {
	AddAuditTrail(ctx context.Context, in *AuditTrailModel) error
}

type QueryRepository interface {
}
