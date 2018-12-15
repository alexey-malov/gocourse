package handlers

import "database/sql"

type videoRepository struct {
	db *sql.DB
}

type VideoRepository interface {
	enumVideos(handler func(v videoItem) bool) error
	findVideo(id string) (*videoItem, error)
	addVideo(v videoItem) error
}

func makeVideoRepository(db *sql.DB) VideoRepository {
	return &videoRepository{db}
}

func (r *videoRepository) enumVideos(handler func(v videoItem) bool) error {
	rows, err := r.db.Query("SELECT video_key, title, duration FROM video")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, title string
		var duration int
		if err := rows.Scan(&id, &title, &duration); err != nil {
			return err
		}
		if !handler(videoItem{id, title, duration}) {
			return nil
		}
	}
	return nil
}

func (r *videoRepository) findVideo(id string) (*videoItem, error) {
	var title string
	var duration int
	if err := r.db.QueryRow("SELECT title, duration FROM video WHERE video_key=?", id).Scan(&title, &duration); err != nil {
		return nil, err
	}
	return &videoItem{id, title, duration}, nil
}

func (r *videoRepository) addVideo(v videoItem) error {
	_, err := r.db.Exec(`INSERT INTO
    video
SET
    video_key = ?,
    title = ?,
    status = 3,
    duration = ?,
    url = ?,
    thumbnail_url = ?`, v.id, v.name, v.duration, v.videoUrl(), v.screenShotUrl())
	return err
}
