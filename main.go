package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"os/exec"
	"runtime"
	"syscall"

	//"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)
const StreamServerDirectory string = "C:\\Users\\MK\\GolandProjects\\streamingApplication\\build"
const ServerAddress string = "127.0.0.1"
const WinServerAddress string = "192.168.0.66"
const StreamDirectory string = "stream"
const ServerDirectory string = "/usr/local/nginx/html"
const TimeoutMicroSeconds int = 5000000
var outputResolutions [][]int = [][]int{{1920, 1080}, {1280, 720}, {854, 480}, {640, 360}}

// MaxBitrate Bitrate in megabits
const MaxBitrate float64 = 5


func constructHLSArgs(hls_time string, hls_wrap string, hls_playlist_type string, hls_flags string, hls_segment_filename string, master_pl_name string, hls_segment_type string, nSplits int) string {
	s := ""
	s += fmt.Sprintf("-f hls -hls_time %s -hls_wrap %s -hls_playlist_type %s -hls_flags %s -hls_segment_type %s -master_pl_name %s -hls_segment_filename %s",
		hls_time, hls_wrap, hls_playlist_type, hls_flags, hls_segment_type, master_pl_name, hls_segment_filename)
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
	s += fmt.Sprintf("[0:v]split=%d", nSplits)
	for i := 0; i < nSplits; i ++ {
		s += fmt.Sprintf("[%s%d]", prefix, i)
	}
	s += ";"
	for i := 0; i < nSplits; i++{
		s += fmt.Sprintf("[%s%d]scale=w=%d:h=%d[%s%dout]", prefix, i, outputResolutions[i][0], outputResolutions[i][1], prefix, i)
		if i+1 == nSplits{
			s+=" "
			break
		}
		s += ";"
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

func constructVarStreamMap(nSplits int) string{
	nSplits = validateNumberOfSplits(nSplits)
	s := ""
	for i := 0; i < nSplits; i++{
		s+=fmt.Sprintf("v:%d,a:%d", i, i)
		if i+1 == nSplits{
			s+= ""
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

// probeRTMPStream Checks if an rtmp stream exists at the given rtmp server address and streamkey. Returns false if it does not exist and true if it does.
func probeRTMPStream(streamKey string, address string) bool{
	cmd := exec.Command("ffprobe", "-v", "error", "-rw_timeout", strconv.Itoa(TimeoutMicroSeconds), fmt.Sprintf("rtmp://%s/live/%s", address, streamKey))


	// Code modified from https://stackoverflow.com/questions/10385551/get-exit-code-go
	// This code checks the exit code. Ffprobe will return 1 if no stream was found, and 0 if one was.
	if err := cmd.Start(); err != nil {
		log.Printf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Printf("cmd.Wait: %v", err)
		}
		return false
	}

	return true
}

func transcodeToHLS(nSplits int, streamPath string){
	command := fmt.Sprintf("ffmpeg -rw_timeout %d -i rtmp://%s/live/%s ", TimeoutMicroSeconds / 10, WinServerAddress, "cool") +
		constructFilterArgs("v", nSplits) +
		constructMapArgs("v", nSplits, "ultrafast", 10) +
		constructHLSArgs("5", "5", "event",
			"independent_segments", fmt.Sprintf("%s/stream%%v_%%03d.ts", streamPath),
			"main.m3u8", "mpegts", nSplits)
	var cmd *exec.Cmd
	// Run either using bash when using linux, or cmd  when using windows
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(`ffmpeg`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.CmdLine = fmt.Sprintf(`%s -var_stream_map "%s" `, command, constructVarStreamMap(nSplits)) + fmt.Sprintf("%s/stream%%v.m3u8", streamPath)
		//cmd = exec.Command("cmd", "/C", command, "-var_stream_map", fmt.Sprintf(`"%s"`, constructVarStreamMap(nSplits)), fmt.Sprintf("%s/stream%%v.m3u8", streamPath))
	case "linux":
		cmd = exec.Command("/usr/bin/bash", "-c", command, "-val_stream_map", fmt.Sprintf(`"%s"`, constructVarStreamMap(nSplits)), fmt.Sprintf("%s/stream%%v.m3u8", streamPath))
	}
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan(){
		fmt.Println(scanner.Text())
	}
	if err != nil {
		return
	}
}

func startHLSStream(nSplits int, streamPath string, streamKey string){
	// Clear the directory containing the stream
	streamPath = streamPath + "/" + streamKey
	clearDirectory(streamPath)
	// Check if rtmp stream exists
	// If it does, start the hls conversion process
	if probeRTMPStream(streamKey, WinServerAddress){
		go transcodeToHLS(nSplits, streamPath)
	}
}

// clearDirectory Deletes all files within a directory. It does this by removing the directory and recreating it afterwards.
func clearDirectory(directory string){
	err := os.RemoveAll(directory)
	if err != nil {
		println(err)
	}
	err = os.Mkdir(directory, os.ModePerm)
	if err != nil {
		println(err)
	}
}

func main() {
	// Get the path to server
	serverPath := "build"
	localpath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	streamPath := fmt.Sprintf("%s\\%s\\stream", localpath, serverPath)

	go startHLSStream(3, streamPath, "cool")


	r := mux.NewRouter()
	r.PathPrefix("/stream/").Handler(http.StripPrefix("/stream/", http.FileServer(http.Dir(filepath.Join(streamPath)))))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET"}, // Allowing only get, just an example
	})
	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}


