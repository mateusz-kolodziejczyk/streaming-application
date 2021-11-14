package main

import (
	"fmt"
	ffmpeg "github.com/mateusz-kolodziejczyk/ffmpeg-go"
	"log"
	"math"
	"os/exec"
	"strconv"
)

const ServerAddress string = "192.168.0.66"
const StreamDirectory string = "stream"
const ServerDirectory string = "test_videos"
var outputResolutions [][]int = [][]int{{1920, 1080}, {1280, 720}, {854, 480}, {640, 360}}
// Bitrate in megabits
const MaxBitrate float64 = 5

// Ffmpeg commands are modified versions of ones from https://ottverse.com/hls-packaging-using-ffmpeg-live-vod/
func hlsStream(streamKey string, streamName string) {
	err := ffmpeg.Input(fmt.Sprintf("rtmp://%s/live/%s", ServerAddress, streamKey)).
		Output(fmt.Sprintf("%s/%s/%s.m3u8", ServerDirectory, StreamDirectory, streamName),
			ffmpeg.KwArgs{
				"filter_complex":       "\"[0:v]split=3[v1][v2][v3];\\ [v1]copy[v1out];\\ [v2]scale=w=1280:h=720[v2out];\\ [v3]scale=w=640:h=360[v3out]\"",
				"map [v1out]":          "-c:v:0 libx264 -x264-params \"nal-hrd=cbr:force-cfr=1\" -b:v:0 5M",
				"maxrate:v:0":          "5M -minrate:v:0 5M -bufsize:v:0 10M -preset slow -g 48 -sc_threshold 0 -keyint_min 48",
				"map [v2out]":          "-c:v:1 libx264 -x264-params \"nal-hrd=cbr:force-cfr=1\" -b:v:1 3M",
				"maxrate:v:1":          "3M -minrate:v:1 3M -bufsize:v:1 3M -preset slow -g 48 -sc_threshold 0 -keyint_min 48",
				"map [v3out]":          "-c:v:2 libx264 -x264-params \"nal-hrd=cbr:force-cfr=1\" -b:v:2 1M",
				"maxrate:v:2":          "1M -minrate:v:2 1M -bufsize:v:2 1M -preset slow -g 48 -sc_threshold",
				"map a:0 -c:a:0":       "aac -b:a:0 96k -ac 2",
				"map a:1 -c:a:1":       "aac -b:a:1 96k -ac 2",
				"map a:2 -c:a:2":       "aac -b:a:2 96k -ac 2",
				"f":                    "hls",
				"hls_time":             "5",
				"hls_wrap":             "5",
				"hls_playlist_type":    "event",
				"hls_segment_type":     "mpegts",
				"g":                    "4",
				"hls_list_size":        "5",
				"hls_segment_filename": fmt.Sprintf("%s/%s/%s", ServerDirectory, StreamDirectory, streamName) + "%03d.ts",
				"master_pl_name":       "main.m3u8",
				"var_stream_map":       "\"v:0,a:0 v:1,a:1 v:2,a:2\" stream_%v.m3u8",
			}).
		OverWriteOutput().Run()
	if err != nil {
		log.Fatal(err)
	}
}
func split(streamKey string, streamName string){
	cmd := exec.Command("ffmpeg")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

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
	/*cmd := exec.Command("/bin/sh", "-c", "rm " + fmt.Sprintf("%s/%s/*", ServerDirectory, StreamDirectory))
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
	if err != nil{
		log.Fatal(err)
	}*/
	nSplits := 2
	print(fmt.Sprintf("ffmpeg -i rtmp://%s/live/%s ", ServerAddress, "cool"))
	print(constructFilterArgs("v", nSplits))
	print(constructMapArgs("v", nSplits, "fast", 10))
	print(constructHLSArgs("5", "5", "event", "independent_segments", "stream_%v\\stream%03d.ts", "main.m3u8", "mpegts", nSplits))
	print("stream_%v\\stream.m3u8")
}
