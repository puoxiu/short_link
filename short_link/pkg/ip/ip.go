package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetLocation(ip string)  (country, region, city string, err error) {
  url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
  fmt.Println("url:", url)
	country, region, city = "unknown", "unknown", "unknown"

    resp, err := http.Get(url)
    if err != nil {
        return country, region, city, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return country, region, city, err
    }

    var result struct {
        Status      string `json:"status"`
        Country     string `json:"country"`
        RegionName  string `json:"regionName"`
        City        string `json:"city"`
    }
    err = json.Unmarshal(body, &result)
    if err != nil {
		return country, region, city, err
    }

    if result.Status != "success" {
		return country, region, city, err
    }

    return result.Country, result.RegionName, result.City, nil
}