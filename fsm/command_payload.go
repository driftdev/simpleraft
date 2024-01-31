package fsm

type CommandPayload struct {
	Operation string
	Key       string
	Value     interface{}
}
