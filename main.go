package main

import (
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
)

const ServerAddress string = "192.168.0.66"
const StreamDirectory string = "./test_video/stream"
func hlsStream(streamKey string, streamName string){
	err := ffmpeg.Input(fmt.Sprintf("rtmp://%s/live/%s", ServerAddress, streamKey)).
		Output(fmt.Sprintf("%s/%s.m3u8", StreamDirectory, streamName),
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
				"hls_segment_filename": fmt.Sprintf("%s/%s", StreamDirectory, streamName) + "%02d.ts"}).
		OverWriteOutput().Run()
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	hlsStream("MyStream", "firstStream")
}

