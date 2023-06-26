package events

import "github.com/asaskevich/EventBus"

var Bus EventBus.Bus

// Events
const RESOURCES_CLEANUP_EVENT = "resources:cleanup"

func init() {
	Bus = EventBus.New()

	// Remove Proxy Host
	Bus.Subscribe(RESOURCES_CLEANUP_EVENT, deleteProxyHost)
}
