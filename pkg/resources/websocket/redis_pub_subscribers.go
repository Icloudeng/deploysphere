package websocket

import (
	"encoding/base64"
	"encoding/json"
	"smatflow/platform-installer/pkg/http/ws"
	"smatflow/platform-installer/pkg/pubsub"
)

func EmitDecodedEvent(playload pubsub.NetworkEventPayload) {
	if json_data, err := json.Marshal(playload); err == nil {
		ws.Broadcast(json_data)
	}
}

func EmitEncodedEvent(playload pubsub.NetworkEventPayload) {
	data := playload
	decodedBytes, err := base64.StdEncoding.DecodeString(data.Payload)
	if err != nil {
		return
	}

	data.Payload = string(decodedBytes)

	if json_data, err := json.Marshal(data); err == nil {
		ws.Broadcast(json_data)
	}
}
