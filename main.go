package main

import (
  "secure-route-api/actor"
  "secure-route-api/api"
)

func main() {
  //s3Client := s3.GetClient()
  //crimes, _ := crime.FromS3(s3Client, s3.Bucket{Name: os.Getenv("AWS_BUCKET_NAME")})
  //system := actor.NewSystem(map[int]city.City{
  //  1: {
  //    ID:   1,
  //    Name: "Medellin",
  //  },
  //}, &crime.Crimes{})
  api.New(actor.System{}).Run()
}
