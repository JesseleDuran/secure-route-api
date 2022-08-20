package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	Initialize()
	os.Exit(m.Run())
}
func Test_GetCountry(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		got := Country()
		assert.Equal(t, "test", got)
	})
}

func Test_GetS3BucketName(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		got := S3BucketName()
		assert.Equal(t, "test-bucket", got)
	})
}

func Test_GetS3DownloadPath(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		got := S3DownloadPath()
		assert.Equal(t, "./testdata/", got)
	})
}
