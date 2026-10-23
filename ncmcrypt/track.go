package ncmcrypt

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

var (
	albumAPIBaseURL = "https://music.163.com/api/v1/album"
	albumHTTPClient = &http.Client{Timeout: 10 * time.Second}
	albumTrackCache sync.Map
)

func fetchTrackNumber(albumID, musicID int64) (int64, error) {
	if albumID <= 0 || musicID <= 0 {
		return 0, fmt.Errorf("album id or music id is missing")
	}

	if cached, ok := albumTrackCache.Load(albumID); ok {
		return cached.(map[int64]int64)[musicID], nil
	}

	req, err := http.NewRequest(http.MethodGet, albumAPIBaseURL+"/"+strconv.FormatInt(albumID, 10), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "ncmdump-go")

	resp, err := albumHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("album request returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return 0, err
	}

	tracks := make(map[int64]int64)
	for _, song := range gjson.GetBytes(body, "songs").Array() {
		id := song.Get("id").Int()
		trackNumber := song.Get("no").Int()
		if id > 0 && trackNumber > 0 {
			tracks[id] = trackNumber
		}
	}
	if len(tracks) == 0 {
		return 0, fmt.Errorf("album response contains no tracks")
	}

	albumTrackCache.Store(albumID, tracks)
	return tracks[musicID], nil
}
