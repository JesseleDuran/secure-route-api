package config

import (
	"os"
)

var Config config

type config struct {
	Country        string
	S3BucketName   string
	S3DownloadPath string
}

func Initialize() {
	Config.Country = Country()
	Config.S3BucketName = S3BucketName()
	Config.S3DownloadPath = S3DownloadPath()
}

func Country() string {
	return os.Getenv("COUNTRY")
}

func S3BucketName() string {
	return os.Getenv("AWS_BUCKET_NAME")
}

func S3DownloadPath() string {
	path := os.Getenv("S3_DOWNLOAD_PATH")
	if path == "" {
		return "downloads/"
	}
	return path
}
