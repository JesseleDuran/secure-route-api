package main

import (
  "os"
  "secure-route-api/actor"
  "secure-route-api/api"
  "secure-route-api/city"
  "secure-route-api/crime"
  "secure-route-api/s3"
)

func main() {
  s3Client := s3.GetClient()
  crimes, _ := crime.FromS3(s3Client, s3.Bucket{Name: os.Getenv("AWS_BUCKET_NAME")})
  system := actor.NewSystem(map[int]city.City{
    1: {
      ID:   1,
      Name: "Medellin",
    },
  }, &crimes)
  api.New(system).Run()
}
