package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/runstp/bilibili-dmj/types"
)

func SendAuthMessage(ws *websocket.Conn, authMsg types.AuthMessage) error {
	raw, err := json.Marshal(authMsg)
	if err != nil {
		return err
	}

	return ws.WriteMessage(websocket.TextMessage, PackMessage(OperationUserAuthentication, raw).ToBytes())
}

func SendHeartbeat(ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.BinaryMessage, PackMessage(OperationHeartbeat, []byte(`[object Object]`)).ToBytes())
}
