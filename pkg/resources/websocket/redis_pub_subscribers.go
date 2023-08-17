package websocket

import (
	"encoding/base64"
	"encoding/json"
	"smatflow/platform-installer/pkg/http/ws"
	"smatflow/platform-installer/pkg/pubsub"
)

func EmitRedisEvent(playload pubsub.NetworkEventPayload) {
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
