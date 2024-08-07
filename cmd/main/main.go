package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/victor247k/TerminalVideoViewer/internal/audio"
	"github.com/victor247k/TerminalVideoViewer/internal/download"
	"github.com/victor247k/TerminalVideoViewer/internal/extractvideoframes"
	"github.com/victor247k/TerminalVideoViewer/internal/render"
)

var horizontal_scale int = 4
const factor float64 = 2.5
var vertical_scale int = int(float64(horizontal_scale) * factor)
var numWorkers int = 8

func main() {
	clean()

	var link string

	fmt.Println("Paste your YouTube link:")
	fmt.Scan(&link)

	download.DownloadFromYoutubeLink(link)

	extractvideoframes.ExtractVideoFrames()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		frameCount := extractvideoframes.GetFrameCount()
		fps := download.GetFps()
		frameDuration := time.Second / time.Duration(fps)

		render.Render(frameCount, frameDuration, horizontal_scale, vertical_scale, numWorkers, fps)

	}()
	go func() {
		defer wg.Done()
		audio.PlayAudio("temp/audio.mp3")
	}()

	wg.Wait()

	clean()
}

func clean() {
	download.DeleteTempFiles()
	extractvideoframes.RemoveFrames()
	render.ClearTerminal()
}
