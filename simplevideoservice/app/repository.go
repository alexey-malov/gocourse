package app

import (
	"database/sql"
	"fmt"
	"github.com/alexey-malov/gocourse/simplevideoservice/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type VideoPersister interface {
	Videos() repository.Videos
	io.Closer
}

type repo struct {
	db     *sql.DB
	videos repository.Videos
}

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) Videos() repository.Videos {
	return r.videos
}

func MakeVideoPersister() (VideoPersister, error) {

	const dbUrlEnvVar = "SIMPLE_VIDEO_SERVICE_DB"
	dbUrl := os.Getenv(dbUrlEnvVar)
	if dbUrl == "" {
		return nil, fmt.Errorf("No %s environment variable", dbUrlEnvVar)
	}

	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		if err = db.Close(); err != nil {
			logrus.Error(err)
		}
		return nil, err
	}

	vr := repository.MakeVideoRepository(db)

	return &repo{db, vr}, nil
}
