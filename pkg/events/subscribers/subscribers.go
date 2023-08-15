package subscribers

import (
	"smatflow/platform-installer/pkg/events"
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/resources/utilities"
)

func EventSubscribers() {
	// Notifier
	events.BusEvent.Subscribe(events.RESOURCES_NOTIFIER_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	events.BusEvent.Subscribe(events.RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
