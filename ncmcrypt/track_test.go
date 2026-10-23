package ncmcrypt

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMetadataIDs(t *testing.T) {
	metadata := NewNeteaseCloudMusicMetadata(`{"musicId":123,"musicName":"Song","albumId":456,"album":"Album"}`)
	if metadata.mMusicID != 123 || metadata.mAlbumID != 456 {
		t.Fatalf("unexpected ids: music=%d album=%d", metadata.mMusicID, metadata.mAlbumID)
	}
}

func TestFetchTrackNumberCachesAlbum(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"songs":[{"id":11,"no":1},{"id":22,"no":2}]}`))
	}))
	defer server.Close()

	oldURL := albumAPIBaseURL
	oldClient := albumHTTPClient
	albumAPIBaseURL = server.URL
	albumHTTPClient = server.Client()
	albumTrackCache = sync.Map{}
	t.Cleanup(func() {
		albumAPIBaseURL = oldURL
		albumHTTPClient = oldClient
		albumTrackCache = sync.Map{}
	})

	for _, musicID := range []int64{11, 22} {
		trackNumber, err := fetchTrackNumber(99, musicID)
		if err != nil {
			t.Fatal(err)
		}
		if trackNumber != musicID/11 {
			t.Fatalf("music %d: got track %d", musicID, trackNumber)
		}
	}
	if requests != 1 {
		t.Fatalf("got %d requests, want 1", requests)
	}
}
