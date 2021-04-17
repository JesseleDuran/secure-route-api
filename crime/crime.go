package crime

import (
  "encoding/csv"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "secure-route-api/s3"
  "strconv"
  "strings"
  "time"

  "github.com/JesseleDuran/osm-graph/coordinates"
  gj "github.com/paulmach/go.geojson"
)

type Crime struct {
  Date      time.Time
  ID        int
  Lat, Lng  float64
  Type      string
  Transport string //should be type Transport
  Weapon    string
  Victim    Victim
}

type Victim struct {
  Sex string
  Age int
}

type Crimes struct {
  values []Crime
  tree   Tree
}

func FromS3(client s3.Client, b s3.Bucket) (Crimes, error) {
  crimes := make([]Crime, 0)
  files := client.GetAllObjectKeys(b.Name)

  for _, file := range files {
    log.Println("Downloading file:", file)
    if filepath.Ext(file) == ".csv" && file != "output1.csv" {
      //err := client.Get(b.Name, file, "downloads/"+file)
      //if err != nil {
      //  log.Printf("couldnt download file %s, err:", file)
      //  continue
      //}
      crimes = append(crimes, FromCSVFile("downloads/"+file)...)
    }
  }
  log.Printf("got %d crimes", len(crimes))
  return Crimes{
    values: crimes,
    tree:   MakeTree(crimes),
  }, nil
}

func FromCSVFile(path string) []Crime {
  crimes := make([]Crime, 0)
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

func fromCsvValues(record []string) (Crime, error) {
  values := strings.Split(record[0], ";")
  if len(values) >= 20 {
    c, err := coordinates.FromStrings(values[2], values[3])
    if err != nil {
      return Crime{}, fmt.Errorf("invalid lng")
    }
    t, _ := time.Parse("2006-01-02 15:04:05", values[0])
    age, _ := strconv.Atoi(values[5])
    return Crime{
      Date:      t,
      ID:        0,
      Lat:       c.Lat,
      Lng:       c.Lng,
      Type:      values[16],
      Transport: values[13],
      Weapon:    values[20],
      Victim: Victim{
        Sex: values[4],
        Age: age,
      },
    }, nil
  }
  return Crime{}, fmt.Errorf("not enough values")
}

func FromGeoJSONFile(path string) Crimes {
  f, _ := os.Open(path)
  bytes, _ := ioutil.ReadAll(f)
  return FromGeoJSON(bytes)
}

func FromGeoJSON(value []byte) Crimes {
  fc, _ := gj.UnmarshalFeatureCollection(value)
  crimes := make([]Crime, 0)
  for i, f := range fc.Features {
    if f.Geometry.IsPoint() {
      crime := fromFeature(f)
      crime.ID = i
      crimes = append(crimes, crime)
    }
  }

  return Crimes{
    values: crimes,
    tree:   MakeTree(crimes),
  }
}

func fromFeature(f *gj.Feature) Crime {
  point := f.Geometry.Point
  return Crime{
    Lat: point[1],
    Lng: point[0],
  }
}

// SetWeight sets the Weight between two nodes of the map graph.
// This function is going to be passed to the graph for using it in the
// construction.
func (cc Crimes) SetWeight(a, b coordinates.Coordinates) float64 {
  radius := float64(100)
  crimes := cc.tree.Search(a.Lat, a.Lng, radius)
  crimes1 := cc.tree.Search(b.Lat, b.Lng, radius)
  return float64(len(crimes) + len(crimes1))
}
