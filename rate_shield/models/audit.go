package models

// AuditLog represents an audit trail entry for rule modifications
type AuditLog struct {
	ID        string `json:"id"`         // Unique identifier (UUID)
	Timestamp int64  `json:"timestamp"`  // Unix timestamp
	Actor     string `json:"actor"`      // User ID/email who performed the action
	Action    string `json:"action"`     // CREATE, UPDATE, or DELETE
	Endpoint  string `json:"endpoint"`   // The API endpoint affected by the rule
	OldRule   *Rule  `json:"old_rule"`   // State before change (null for CREATE)
	NewRule   *Rule  `json:"new_rule"`   // State after change (null for DELETE)
	IPAddress string `json:"ip_address"` // IP address of the requester
	UserAgent string `json:"user_agent"` // User agent of the requester
}

// AuditAction constants for different types of actions
const (
	AuditActionCreate = "CREATE"
	AuditActionUpdate = "UPDATE"
	AuditActionDelete = "DELETE"
)

// PaginatedAuditLogs represents a paginated response of audit logs
type PaginatedAuditLogs struct {
	PageNumber  int         `json:"page_number"`
	TotalItems  int         `json:"total_items"`
	HasNextPage bool        `json:"has_next_page"`
	Logs        []AuditLog  `json:"logs"`
}
