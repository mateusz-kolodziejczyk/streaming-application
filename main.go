package main

import (
	"bufio"
	"fmt"
	ffmpeg "github.com/mateusz-kolodziejczyk/ffmpeg-go"
	"log"
	"os/exec"
)

const ServerAddress string = "127.0.0.1"
const StreamDirectory string = "stream"
const ServerDirectory string = "/usr/local/nginx/html"
func hlsStream(streamKey string, streamName string){
	err := ffmpeg.Input(fmt.Sprintf("rtmp://%s/live/%s", ServerAddress, streamKey)).
		Output(fmt.Sprintf("%s/%s/%s.m3u8", ServerDirectory, StreamDirectory, streamName),
			ffmpeg.KwArgs{
				"c:a": "aac",
				"b:a": "128k",
				"ac": "2",
				"c:v": "libx264",
				"preset": "veryfast",
				"b:v":"5M",
				"f": "hls",
				"hls_time": "5",
				"hls_wrap": "5",
				"hls_playlist_type": "event",
				"hls_segment_type": "mpegts",
				"g": "4",
				"hls_list_size": "5",
				"hls_segment_filename": fmt.Sprintf("%s/%s/%s",ServerDirectory, StreamDirectory, streamName) + "%03d.ts"}).
		OverWriteOutput().Run()
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	cmd := exec.Command("/bin/sh", "-c", "rm " + fmt.Sprintf("%s/%s/*", ServerDirectory, StreamDirectory))
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
	if err != nil{
		log.Fatal(err)
	}
	hlsStream("cool", "camera")
}

