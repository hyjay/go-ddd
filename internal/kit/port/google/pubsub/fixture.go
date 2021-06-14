package pubsub

type greetedEvent struct {
	Greeting string `json:"greeting"`
}

func (e *greetedEvent) Name() string {
	return "greetedEvent"
}
