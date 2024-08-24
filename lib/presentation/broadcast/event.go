package broadcast

type EventId string

type Event[P interface{}] struct {
	ID      EventId `json:"id"`
	Payload P       `json:"payload"`
}

func NewEvent[P interface{}](id EventId, payload P) *Event[P] {
	return &Event[P]{
		ID:      id,
		Payload: payload,
	}
}
