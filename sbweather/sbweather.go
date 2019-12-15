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
	url    string
	temp   int
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
	if r.url == "" {
		// Get coordinates.
		lat, long, err := getCoords(r.client, r.zip)
		if err != nil {
			r.err = err
			return
		}

		// Get forecast URL.
		url, err := getURL(r.client, lat, long)
		if err != nil {
			r.err = err
			return
		}
		r.url = url
	}

	// Get hourly temperature.
	temp, err := getTemp(r.client, r.url + "/hourly")
	if err != nil {
		r.err = err
		return
	}
	r.temp = temp
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("weather: %v °F", r.temp)
}

// Get the geographic coordinates for the provided zip code.
// We should receive a response in this format:
// {"status":1,"output":[{"zip":"90210","latitude":"34.103131","longitude":"-118.416253"}]}
func getCoords(client http.Client, zip string) (string, string, error) {
	type coords struct {
		Status int                 `json:"status"`
		Output []map[string]string `json:"output"`
	}

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

// Query the NWS to determine which URL we should be using for getting the weather forecast.
// Our value should be within the "properties" field, under key "forecast".
func getURL(client http.Client, lat string, long string) (string, error) {
	type props struct {
		Properties map[string]interface{} `json:"properties"`
	}

	url      := "https://api.weather.gov/points/" + lat + "," + long
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	p   := props{}
	err  = json.Unmarshal(body, &p)
	if err != nil {
		return "", err
	}

	if len(p.Properties) == 0 {
		return "", errors.New("Received invalid properties map")
	}

	url = p.Properties["forecast"].(string)
	if url == "" {
		return "", errors.New("Missing forecast URL in response")
	}

	return url, nil
}

func getTemp(client http.Client, url string) (int, error) {
	type temp struct {
		Properties struct {
			Periods []interface{} `json:"periods"`
		} `json:"properties"`

		// Properties map[string]interface{} `json:"properties"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	t   := temp{}
	err  = json.Unmarshal(body, &t)
	if err != nil {
		return -1, err
	}

	periods := t.Properties.Periods
	if len(periods) == 0 {
		return -1, errors.New("Missing hourly temperature periods")
	}

	latest := periods[0].(map[string]interface{})
	if len(latest) == 0 {
		return -1, errors.New("Missing current temperature")
	}

	temperature := latest["temperature"].(float64)
	return int(temperature), nil
}
