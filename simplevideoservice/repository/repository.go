package repository

import (
	"database/sql"
	"github.com/alexey-malov/gocourse/simplevideoservice/model"
	log "github.com/sirupsen/logrus"
)

type videoRepository struct {
	db *sql.DB
}

type VideoRepository interface {
	EnumVideos(handler func(v model.VideoItem) bool) error
	FindVideo(id string) (*model.VideoItem, error)
	AddVideo(v model.VideoItem) error
}

func MakeVideoRepository(db *sql.DB) VideoRepository {
	return &videoRepository{db}
}

func safeCloseRows(rr *sql.Rows) {
	if err := rr.Close(); err != nil {
		log.Errorf("Failed to close rows")
	}
}

func (r *videoRepository) EnumVideos(handler func(v model.VideoItem) bool) error {
	rows, err := r.db.Query("SELECT video_key, title, duration FROM video")
	if err != nil {
		return err
	}
	defer safeCloseRows(rows)

	for rows.Next() {
		var id, title string
		var duration int
		if err := rows.Scan(&id, &title, &duration); err != nil {
			return err
		}
		if !handler(model.MakeVideoItem(id, title, duration)) {
			return nil
		}
	}
	return nil
}

func (r *videoRepository) FindVideo(id string) (*model.VideoItem, error) {
	var title string
	var duration int
	if err := r.db.QueryRow("SELECT title, duration FROM video WHERE video_key=?", id).Scan(&title, &duration); err != nil {
		return nil, err
	}
	v := model.MakeVideoItem(id, title, duration)
	return &v, nil
}

func (r *videoRepository) AddVideo(v model.VideoItem) error {
	_, err := r.db.Exec(`INSERT INTO
    video
SET
    video_key = ?,
    title = ?,
    status = 3,
    duration = ?,
    url = ?,
    thumbnail_url = ?`, v.Id(), v.Name(), v.Duration(), v.VideoUrl(), v.ScreenShotUrl())
	return err
}
