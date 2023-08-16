package events

import (
	"github.com/asaskevich/EventBus"
)

var BusEvent EventBus.Bus

type NetworkEventPayload struct {
	Type      string
	Reference string
	Channel   string
	Payload   string
}

// Events
const RESOURCES_CLEANUP_EVENT = "resources:cleanup"
const RESOURCES_NOTIFIER_EVENT = "resources:utilities:notifier"

const RESOURCES_DB_STORE_UPDATE = "resources:state:db:store"

func init() {
	BusEvent = EventBus.New()
}
