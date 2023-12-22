package event_dispatcher

import "github.com/freddyouellette/ai-dashboard/internal/models"

type Logger interface {
	Info(msg string, fields map[string]interface{})
}

type EventDispatcher struct {
	logger    Logger
	listeners map[models.EventType][]func(payload interface{})
}

func NewEventDispatcher(logger Logger) *EventDispatcher {
	d := &EventDispatcher{
		logger: logger,
	}
	d.listeners = make(map[models.EventType][]func(payload interface{}))
	return d
}

func (d *EventDispatcher) Register(event models.EventType, listener func(payload interface{})) {
	d.listeners[event] = append(d.listeners[event], listener)
}

func (d *EventDispatcher) Dispatch(event models.EventType, payload interface{}) {
	d.logger.Info("Event dispatched", map[string]interface{}{
		"event":   event,
		"payload": payload,
	})
	for _, listener := range d.listeners[event] {
		listener(payload)
	}
}
