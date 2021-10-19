package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
)
func main() {
	err := ffmpeg.Input("rtmp://192.168.0.66/live/test").
		Output("./test_videos/stream/stream_%v/stream.m3u8",
			ffmpeg.KwArgs{
			"c:a": "aac",
			"b:a": "128k",
			"ac": "2",
			"c:v": "libx264",
			"preset": "veryfast",
			"b:v":"1M",
			"f": "hls",
			"hls_time": "4",
			"hls_init_time": "12",
			"hls_list_size": "3",
			"hls_playlist_type": "event",
			"hls_segment_type": "mpegts",
			"hls_segment_filename": "./test_videos/stream/stream_%v/data%02d.ts"}).
		OverWriteOutput().Run()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

