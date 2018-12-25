package processor

import (
	"github.com/alexey-malov/gocourse/simplevideoservice/daemon/task"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	"github.com/alexey-malov/gocourse/simplevideoservice/storage"
	"github.com/alexey-malov/gocourse/simplevideoservice/usecases"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func processVideo(videos repository.Videos, stg storage.Storage, video domain.Video) {
	stats := usecases.MakeStats(stg)

	folder, err := stg.GetFolder(video.Id())
	defer func() {
		err := videos.SaveVideo(video)
		if err != nil {
			logrus.Error(err)
		}
	}()

	reportError := func(err error) {
		video.SetStatus(domain.StatusError)
		logrus.Error(err)
	}

	if err != nil {
		reportError(err)
		return
	}

	videoPath, err := stats.VideoPath(video)
	if err != nil {
		reportError(err)
		return
	}

	duration, err := getVideoDuration(videoPath)
	if err != nil {
		reportError(err)
		return
	}

	const thumbName = "thumbnail.jpg"
	thumbnailPath := folder.GetAbsFilePath(thumbName)
	if err = createVideoThumbnail(videoPath, thumbnailPath, 0); err != nil {
		reportError(err)
		return
	}

	video.SetStatus(domain.StatusReady)
	video.SetDuration(int(duration))
	video.SetThumbnailURL(folder.GetRelFilePath(thumbName))
}

func MakeVideoProcessor(videos repository.Videos, stg storage.Storage) task.TaskGenerator {
	return func() task.Task {
		var videoToProcess *domain.Video
		if err := videos.EnumerateWithStatus(domain.StatusUploaded, func(v *domain.Video) bool {
			videoToProcess = v
			return false
		}); err != nil {
			logrus.Error(err)
			return nil
		}

		if videoToProcess == nil {
			return nil
		}

		videoToProcess.SetStatus(domain.StatusProcessing)
		if err := videos.SaveVideo(*videoToProcess); err != nil {
			logrus.Error("err")
			return nil
		}

		return func() {
			processVideo(videos, stg, *videoToProcess)
		}
	}
}

func getVideoDuration(videoPath string) (float64, error) {
	result, err := exec.Command(`ffprobe`, `-v`, `error`, `-show_entries`, `format=duration`, `-of`, `default=noprint_wrappers=1:nokey=1`, videoPath).Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func ffmpegTimeFromSeconds(seconds int64) string {
	return time.Unix(seconds, 0).UTC().Format(`15:04:05.000000`)
}

func createVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset int64) error {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-ss`, ffmpegTimeFromSeconds(thumbnailOffset), `-vframes`, `1`, thumbnailPath).Run()
}
