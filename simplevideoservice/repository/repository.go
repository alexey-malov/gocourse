package repository

import (
	"database/sql"
	"github.com/alexey-malov/gocourse/simplevideoservice/domain"
	log "github.com/sirupsen/logrus"
)

type videoRepository struct {
	db *sql.DB
}

type Videos interface {
	Enumerate(handler func(v domain.Video) bool) error
	Find(id string) (*domain.Video, error)
	Add(v domain.Video) error
}

func MakeVideoRepository(db *sql.DB) Videos {
	return &videoRepository{db}
}

func safeCloseRows(rr *sql.Rows) {
	if err := rr.Close(); err != nil {
		log.Errorf("Failed to close rows")
	}
}

func (r *videoRepository) Enumerate(handler func(v domain.Video) bool) error {
	rows, err := r.db.Query("SELECT video_key, title, url, thumbnail_url, duration FROM video")
	if err != nil {
		return err
	}
	defer safeCloseRows(rows)

	for rows.Next() {
		var id, title, video, screenshot string
		var duration int
		if err := rows.Scan(&id, &title, &video, &screenshot, &duration); err != nil {
			return err
		}
		if !handler(domain.MakeVideo(id, title, video, screenshot, duration)) {
			return nil
		}
	}
	return nil
}

func (r *videoRepository) Find(id string) (*domain.Video, error) {
	var title, video, screenshot string
	var duration int
	if err := r.db.QueryRow("SELECT title, url, thumbnail_url, duration FROM video WHERE video_key=?", id).
		Scan(&title, &video, &screenshot, &duration); err != nil {
		return nil, err
	}
	v := domain.MakeVideo(id, title, video, screenshot, duration)
	return &v, nil
}

func (r *videoRepository) Add(v domain.Video) error {
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
