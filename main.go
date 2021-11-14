package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
)

const ServerAddress string = "127.0.0.1"
const StreamDirectory string = "stream"
const ServerDirectory string = "/usr/local/nginx/html"
var outputResolutions [][]int = [][]int{{1920, 1080}, {1280, 720}, {854, 480}, {640, 360}}

// MaxBitrate Bitrate in megabits
const MaxBitrate float64 = 1


func constructHLSArgs(hls_time string, hls_wrap string, hls_playlist_type string, hls_flags string, hls_segment_filename string, master_pl_name string, hls_segment_type string, nSplits int) string {
	s := ""
	s += fmt.Sprintf("-f hls -hls_time %s -hls_wrap %s -hls_playlist_type %s -hls_flags %s -hls_segment_type %s -master_pl_name %s %s -hls_segment_filename %s",
		hls_time, hls_wrap, hls_playlist_type, hls_flags, hls_segment_type, master_pl_name, constructValStreamMap(nSplits),  hls_segment_filename)
	s += " "
	return s
}

func constructFilterArgs(prefix string, nSplits int) string{
	s := ""
	if nSplits < 1 {
		return s
	}
	// Can only split the video for each output resolution that exists
	nSplits = validateNumberOfSplits(nSplits)
	// Add the start of the command
	s += "-filter_complex "
	s += fmt.Sprintf("\"[0:v]split=%d", nSplits)
	for i := 0; i < nSplits; i ++ {
		s += fmt.Sprintf("[%s%d]", prefix, i)
	}
	s += "; "
	for i := 0; i < nSplits; i++{
		s += fmt.Sprintf("[%s%d]scale=w=%d:h=%d[%s%dout]", prefix, i, outputResolutions[i][0], outputResolutions[i][1], prefix, i)
		if i+1 == nSplits{
			s+="\" "
			break
		}
		s += "; "
	}
	return s
}

func constructMapArgs(prefix string, nSplits int, preset string, g int) string{
	nSplits = validateNumberOfSplits(nSplits)
	s := ""
	for i := 0; i < nSplits; i++ {
		// Halve the bitrate every step
		bitrate := MaxBitrate / math.Pow(2.0, float64(i))
		bitrateString := strconv.FormatFloat(bitrate, 'f', 2, 64)
		s+=fmt.Sprintf("-map [%s%dout] -c:v:%d libx264 -b:v:%d %sM -preset %s -g %d ", prefix, i, i, i, bitrateString, preset, g)
	}
	for i := 0; i < nSplits; i++ {
		// Halve the bitrate every step
		s+=fmt.Sprintf("-map a:0 -c:a:%d aac -b:a:%d 96k -ac 2 ", i, i)
	}
	return s
}

func constructValStreamMap(nSplits int) string{
	nSplits = validateNumberOfSplits(nSplits)
	s := "-var_stream_map \""
	for i := 0; i < nSplits; i++{
		s+=fmt.Sprintf("v:%d,a:%d", i, i)
		if i+1 == nSplits{
			s +="\" "
			break
		}
		s += " "
	}
	return s
}

func validateNumberOfSplits(nSplits int) int {
	if nSplits >= len(outputResolutions) {
		return len(outputResolutions)
	}
	return nSplits
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
	nSplits := 2
	print(fmt.Sprintf("ffmpeg -i rtmp://%s/live/%s ", ServerAddress, "cool") +
		constructFilterArgs("v", nSplits) +
		constructMapArgs("v", nSplits, "ultrafast", 10) +
		constructHLSArgs("5", "5", "event",
			"independent_segments", "stream_%v\\stream%03d.ts",
			"main.m3u8", "mpegts", nSplits) +
		"stream_%v\\stream.m3u8")
}
