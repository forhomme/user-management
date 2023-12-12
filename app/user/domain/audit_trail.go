package domain

type AuditTrail struct {
	UserId   string
	Menu     string
	Method   string
	Request  string
	Response string
}
