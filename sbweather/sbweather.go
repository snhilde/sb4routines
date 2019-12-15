package sbweather

import (
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type routine struct {
	err    error
	client http.Client
	zip    string
	lat    string
	long   string
}

type coords struct {
	Status int                 `json:"status"`
	Output []map[string]string `json:"output"`
}

func New(zip string) *routine {
	var r routine

	if len(zip) != 5 {
		r.err = errors.New("Invalid Zip Code length")
		return &r
	}

	_, err := strconv.Atoi(zip)
	if err != nil {
		r.err = err
		return &r
	}
	r.zip = zip

	return &r
}

func (r *routine) Update() {
	if r.lat == "" || r.long == "" {
		lat, long, err := getCoords(r.client, r.zip)
		if err != nil {
			r.err = err
			return
		}
		fmt.Println(lat, long)
	}
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "weather"
}

// We should receive a response in this format:
// {"status":1,"output":[{"zip":"90210","latitude":"34.103131","longitude":"-118.416253"}]}
func getCoords(client http.Client, zip string) (string, string, error) {
	url      := "https://api.promaptools.com/service/us/zip-lat-lng/get/?zip=" + zip + "&key=17o8dysaCDrgv1c"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	c   := coords{}
	err  = json.Unmarshal(body, &c)
	if err != nil {
		return "", "", err
	}

	// Make sure we got back just one dictionary.
	if len(c.Output) != 1 {
		return "", "", errors.New("Received invalid coordinates array")
	}

	// Make sure the status is good.
	if c.Status != 1 {
		return "", "", errors.New("Coordinates request failed")
	}

	lat  := c.Output[0]["latitude"]
	long := c.Output[0]["longitude"]
	if lat == "" || long == "" {
		return "", "", errors.New("Missing coordinates in response")
	}

	return lat, long, nil
}
