package domain

type MessageCard interface {
	ToJson() (string, error)
}

type Button struct {
	Text      string                 `json:"text"`
	Value     map[string]interface{} `json:"value"`
	Behaviors map[string]interface{} `json:"behaviors"`
}
