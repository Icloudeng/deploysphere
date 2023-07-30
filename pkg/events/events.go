package events

import (
	// "smatflow/platform-installer/pkg/resources"
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/resources/utilities"

	"github.com/asaskevich/EventBus"
)

var BusEvent EventBus.Bus

// Events
const RESOURCES_CLEANUP_EVENT = "resources:cleanup"
const RESOURCES_NOTIFIER_EVENT = "resources:utilities:notifier"

const RESOURCES_DB_STORE_UPDATE = "resources:state:db:store"

func init() {
	BusEvent = EventBus.New()

	// Notifier
	BusEvent.Subscribe(RESOURCES_NOTIFIER_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	BusEvent.Subscribe(RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)

	// BusEvent.Subscribe(RESOURCES_DB_STORE_UPDATE, resources.StoreOrUpdateResourceState)
}
