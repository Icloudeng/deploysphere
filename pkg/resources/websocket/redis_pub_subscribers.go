package websocket

import (
	"encoding/base64"
	"encoding/json"
	"smatflow/platform-installer/pkg/events/redis_events"
	"smatflow/platform-installer/pkg/http/ws"
)

func EmitRedisEvent(playload redis_events.ResourceRedisEventPayload) {
	data := playload
	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)
	if err != nil {
		return
	}

	data.Payload = string(decodedBytes)

	if json_data, err := json.Marshal(data); err == nil {
		ws.Broadcast(json_data)
	}
}
