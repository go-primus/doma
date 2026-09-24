package command

type CommandType string

func (ct CommandType) String() string {
	return string(ct)
}

type Command struct {
	CommandId   string `json:"command_id,omitempty"`
	CommandType string `json:"command_type,omitempty"`
	Stream      string `json:"stream,omitempty"`
	Payload     any    `json:"payload,omitempty"`
}
