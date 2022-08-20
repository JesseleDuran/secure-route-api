package s3

import (
	"encoding/csv"
	"fmt"
	api "github.com/JesseleDuran/secure-route-api"
	"github.com/JesseleDuran/secure-route-api/config"
	"github.com/golang/geo/s2"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CrimesInMemoryProvider struct {
	client api.S3Client
}

func NewCrimesInMemoryProvider(client api.S3Client) api.CrimeInMemoryProvider {
	return CrimesInMemoryProvider{
		client: client,
	}
}

// Fetch builds a graph from a serialized file.
func (cp CrimesInMemoryProvider) Fetch() map[uint64]api.Crime {
	files := cp.client.GetAllObjectKeys(config.Config.S3BucketName)
	result := make(map[uint64]api.Crime, 0)
	for _, file := range files {
		log.Println("Downloading file:", file)
		if filepath.Ext(file) == ".csv" && file != "output1.csv" {
			err := cp.client.Get(config.Config.S3BucketName, file, "downloads/"+file)
			if err != nil {
				log.Printf("couldnt download file %s, err:", file)
				continue
			}
			crimes := FromCSVFile("downloads/" + file)
			for _, c := range crimes {
				result[c.ID] = c
			}
		}
	}
	return result
}
func FromCSVFile(path string) []api.Crime {
	crimes := make([]api.Crime, 0)
	f, _ := os.Open(path)
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error FromCSVFile: ", err.Error(), path)
			continue
		}
		crime, err := fromCsvValues(record)
		if err == nil {
			crimes = append(crimes, crime)
		}
	}
	return crimes
}

func fromCsvValues(record []string) (api.Crime, error) {
	values := strings.Split(record[0], ";")
	if len(values) >= 20 {
		age, _ := strconv.Atoi(values[5])
		t, _ := time.Parse("2006-01-02 15:04:05", values[0])
		lat, err := strconv.ParseFloat(values[2], 32)
		if err != nil || math.IsNaN(lat) {
			return api.Crime{}, fmt.Errorf("invalid lat")
		}
		lng, err := strconv.ParseFloat(values[3], 32)
		if err != nil || math.IsNaN(lng) {
			return api.Crime{}, fmt.Errorf("invalid lng")
		}
		if lat == 0 || lng == 0 {
			return api.Crime{}, fmt.Errorf("invalid lat and lng")
		}
		id := uint64(t.Unix()) + uint64(s2.CellFromLatLng(s2.LatLngFromDegrees(lat, lng)).ID())
		return api.Crime{
			Date:      t,
			ID:        id,
			Lat:       lat,
			Lng:       lng,
			Type:      values[16],
			Transport: values[13],
			Weapon:    values[20],
			Victim: api.Victim{
				Sex: values[4],
				Age: age,
			},
		}, nil
	}
	return api.Crime{}, fmt.Errorf("not enough values")
}
