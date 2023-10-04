package subscribers

import (
	"github.com/icloudeng/platform-installer/internal/pubsub"
	proxyhost "github.com/icloudeng/platform-installer/internal/resources/proxy_host"
	"github.com/icloudeng/platform-installer/internal/resources/utilities"
)

func EventSubscribers() {
	// Notifier
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_NOTIFIER_EVENT, utilities.SendNotification)

	// Remove Proxy Host
	pubsub.BusEvent.Subscribe(pubsub.RESOURCES_CLEANUP_EVENT, proxyhost.DeleteProxyHost)
}
