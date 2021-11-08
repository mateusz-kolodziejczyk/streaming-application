module github.com/mateusz-kolodziejczyk/streaming-application

go 1.17

// Using a fork of the original ffmpeg-go module as it gave me a broken import
require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mateusz-kolodziejczyk/ffmpeg-go v0.3.1
)

require (
	github.com/torresjeff/rtmp v0.0.0-20210303201626-9aba915a956e // indirect
	github.com/u2takey/go-utils v0.0.0-20210821132353-e90f7c6bacb5 // indirect
)
