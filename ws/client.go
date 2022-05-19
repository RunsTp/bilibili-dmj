package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	bilibilidmj "github.com/runstp/bilibili-dmj"
	"github.com/runstp/bilibili-dmj/types"
)

const (
	liveBaseURL = "https://live.bilibili.com"

	liveId2RoomIdAPI      = "https://api.live.bilibili.com/room/v1/Room/room_init?id=%d"
	getWsDanmuHostListAPI = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?type=0&id=%d"
)

func CreateLiveConn(wsHostServer types.WsHostServer) (*websocket.Conn, error) {
	url := fmt.Sprintf("wss://%s:%d/sub", wsHostServer.Host, wsHostServer.WssPort)
	ws, _, err := websocket.DefaultDialer.Dial(url, http.Header{"origin": []string{"https://live.bilibili.com"}})
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func GetDanmuHostList(roomId int, cookie string) (*types.GetDanmuHostListResp, error) {
	url := fmt.Sprintf(getWsDanmuHostListAPI, roomId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", fmt.Sprintf("%s/%d", liveBaseURL, roomId))
	req.Header.Set("Cookie", strings.TrimRight(cookie, "\r\n"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, bilibilidmj.NewHttpRequestError(fmt.Sprintf("request URL: %s fail", url), resp)
	}

	defer resp.Body.Close()

	var bResp *types.GetDanmuHostListResp
	if err := json.NewDecoder(resp.Body).Decode(&bResp); err != nil {
		return nil, err
	}

	if bResp.Code != 0 {
		return nil, bilibilidmj.NewHttpRequestError(fmt.Sprintf("get ws danmu host list fail, msg:%s", bResp.Message), resp)
	}

	return bResp, nil
}

func LiveId2RoomId(liveId int) (int, error) {

	url := fmt.Sprintf(liveId2RoomIdAPI, liveId)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return 0, err
	} else if resp.StatusCode != http.StatusOK {
		return 0, bilibilidmj.NewHttpRequestError(fmt.Sprintf("request URL: %s fail", url), resp)
	}

	defer resp.Body.Close()

	var bResp *types.LiveId2RoomIdResp
	if err := json.NewDecoder(resp.Body).Decode(&bResp); err != nil {
		return 0, err
	}

	if bResp.Code != 0 {
		return 0, bilibilidmj.NewHttpRequestError(fmt.Sprintf("live id to room id fail, msg:%s", bResp.Message), resp)
	}

	return bResp.Data.RoomId, nil
}
