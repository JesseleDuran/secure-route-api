package crime

import (
  "log"
  "testing"
)

func TestFromGeoJSONFile(t *testing.T) {
  cc := FromGeoJSONFile("/Users/jesseleduran/Documents/secure route graph" +
    "/secure-route-api/crimes.json")

  //crimes := cc.tree.Search(4.639557220462878, -74.06363904476164, 50)
  r := cc.SetWeight(4.639557220462878, -74.06363904476164,
    4.638991123520539, -74.06295776367188)
  log.Println("aja", r)
}
