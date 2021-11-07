module github.com/mateusz-kolodziejczyk/streaming-application

go 1.17

// Using a fork of the original ffmpeg module as it had a broken import
require (
	github.com/mateusz-kolodziejczyk/ffmpeg-go v0.3.1
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

require github.com/u2takey/go-utils v0.0.0-20210821132353-e90f7c6bacb5 // indirect
