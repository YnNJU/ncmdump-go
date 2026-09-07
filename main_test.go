package main

import "testing"

func TestIsNCMFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "lowercase", path: "/Users/test/Music/song.ncm", want: true},
		{name: "uppercase", path: "/Users/test/Music/song.NCM", want: true},
		{name: "unicode path", path: "/Users/test/音乐/歌曲.ncm", want: true},
		{name: "macOS AppleDouble file", path: "/Users/test/Music/._song.ncm", want: false},
		{name: "short name", path: "a", want: false},
		{name: "extension only", path: ".ncm", want: true},
		{name: "different extension", path: "song.flac", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNCMFile(tt.path); got != tt.want {
				t.Fatalf("isNCMFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
