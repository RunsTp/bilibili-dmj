package types

type AuthMessage struct {
	Uid      int    `json:"uid"`
	Roomid   int    `json:"roomid"`
	Protover int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

func NewDefaultAuthMessage(uid, roomId int, key string) AuthMessage {
	return AuthMessage{
		Uid:      uid,
		Roomid:   roomId,
		Protover: 2,
		Platform: "web",
		Type:     2,
		Key:      key,
	}
}
