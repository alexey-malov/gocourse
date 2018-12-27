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
	Enumerate(handler func(v *domain.Video) bool) error
	Find(id string) (*domain.Video, error)
	EnumerateWithStatus(status domain.Status, handler func(v *domain.Video) bool) error
	SaveVideo(v *domain.Video) error
	Add(v *domain.Video) error
}

func MakeVideoRepository(db *sql.DB) Videos {
	return &videoRepository{db}
}

func safeCloseRows(rr *sql.Rows) {
	if err := rr.Close(); err != nil {
		log.Errorf("Failed to close rows")
	}
}

func (r *videoRepository) SaveVideo(v *domain.Video) error {
	_, err := r.db.Exec(`UPDATE video 
SET
    title=?,
    status=?,
    duration=?,
    url=?,
    thumbnail_url=?
    WHERE video_key=?`, v.Name(), int(v.Status()), v.Duration(), v.VideoUrl(), v.ThumbnailURL(), v.Id())
	if err != nil {
		return err
	}
	return nil
}

func (r *videoRepository) EnumerateWithStatus(status domain.Status, handler func(v *domain.Video) bool) error {
	rows, err := r.db.Query("SELECT video_key, title, url, thumbnail_url, duration, status FROM video WHERE status=?", int(status))
	if err != nil {
		return err
	}
	defer safeCloseRows(rows)

	for rows.Next() {
		var id, title, video, screenshot string
		var duration, status int
		if err := rows.Scan(&id, &title, &video, &screenshot, &duration, &status); err != nil {
			return err
		}
		if !handler(domain.MakeVideo(id, title, video, screenshot, duration, domain.Status(status))) {
			return nil
		}
	}
	return nil
}

func (r *videoRepository) Enumerate(handler func(v *domain.Video) bool) error {
	rows, err := r.db.Query("SELECT video_key, title, url, thumbnail_url, duration, status FROM video")
	if err != nil {
		return err
	}
	defer safeCloseRows(rows)

	for rows.Next() {
		var id, title, video, screenshot string
		var duration, status int
		if err := rows.Scan(&id, &title, &video, &screenshot, &duration, &status); err != nil {
			return err
		}
		if !handler(domain.MakeVideo(id, title, video, screenshot, duration, domain.Status(status))) {
			return nil
		}
	}
	return nil
}

func (r *videoRepository) Find(id string) (*domain.Video, error) {
	var title, video, screenshot string
	var duration, status int
	if err := r.db.QueryRow("SELECT title, url, thumbnail_url, duration, status FROM video WHERE video_key=?", id).
		Scan(&title, &video, &screenshot, &duration, &status); err != nil {
		return nil, err
	}
	v := domain.MakeVideo(id, title, video, screenshot, duration, domain.Status(status))
	return v, nil
}

func (r *videoRepository) Add(v *domain.Video) error {
	_, err := r.db.Exec(`INSERT INTO video
SET
    video_key = ?,
    title = ?,
    status = ?,
    duration = ?,
    url = ?,
    thumbnail_url = ?`, v.Id(), v.Name(), int(v.Status()), v.Duration(), v.VideoUrl(), v.ThumbnailURL())
	return err
}
