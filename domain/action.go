package domain

type Action struct {
	ID    string
	Type  string
	Value map[string]interface{}
	// User   User
	ChatID string
}
