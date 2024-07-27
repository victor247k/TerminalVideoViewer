package download

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"github.com/kkdai/youtube/v2"
)

func findFormatByItag(formats youtube.FormatList, itag int) *youtube.Format {
	for _, format := range formats {
		if format.ItagNo == itag {
			return &format
		}
	}
	return nil
}

func DownloadFromYoutubeLink(link string) {
	DeleteTempFiles()
	client := youtube.Client{}

	video, err := client.GetVideo(link)
	if err != nil {
        log.Fatalf("Error getting video: %v", err)
    }

    videoFormat := findFormatByItag(video.Formats, 136) // 136 is for 720p video
	if videoFormat == nil {
		log.Fatal("720p video format not found")
	}

	audioFormat := findFormatByItag(video.Formats, 140) // 140 is for audio
	if audioFormat == nil {
		log.Fatal("Audio format not found")
	}

	videoStream, _, err := client.GetStream(video, videoFormat)
	if err != nil {
		log.Fatalf("Error getting video stream: %v", err)
	}

	videoFile, err := os.Create("temp/video.mp4")
	if err != nil {
		log.Fatalf("Error creating video file: %v", err)
	}
	defer videoFile.Close()

	_, err = videoFile.ReadFrom(videoStream)
	if err != nil {
		log.Fatalf("Error downloading video: %v", err)
	}

	audioStream, _, err := client.GetStream(video, audioFormat)
	if err != nil {
		log.Fatalf("Error getting audio stream: %v", err)
	}

	audioFile, err := os.Create("temp/audio.m4a")
	if err != nil {
		log.Fatalf("Error creating audio file: %v", err)
	}
	defer audioFile.Close()

	_, err = audioFile.ReadFrom(audioStream)
	if err != nil {
		log.Fatalf("Error downloading audio: %v", err)
	}

	fmt.Println("Downloaded 720p video and audio streams")

	// Merge video and audio
	err = mergeVideoAndAudio("temp/video.mp4", "temp/audio.m4a", "temp/output.mp4")
	if err != nil {
		log.Fatalf("Error merging video and audio: %v", err)
	}

	fmt.Println("Merged video and audio into output.mp4")
}

func mergeVideoAndAudio(videoFilePath, audioFilePath, outputFilePath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoFilePath,
		"-i", audioFilePath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		outputFilePath,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func DeleteTempFiles() {
    cmd := exec.Command("rm", "temp/audio.m4a", "temp/output.mp4", "temp/video.mp4")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

    cmd.Run()
}