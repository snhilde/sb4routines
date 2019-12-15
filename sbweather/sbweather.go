package sbweather

import (
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
)

// routine is the main object for this package.
// err:    error encountered along the way, if any
// client: HTTP client to reuse for all requests out
// zip:    user-supplied zip code for where the temperature should reflect
// url:    NWS-provided URL for getting the temperature, as found during the init
// temp:   current temperature for the provided zip code
// high:   forecast high
// low:    forecast low
type routine struct {
	err    error
	client http.Client
	zip    string
	url    string
	temp   int
	high   int
	low    int
}

// Sanity-check zip code, and return new routine object.
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

// Get the current hourly temperature. Also, if first run of the session, initialize object.
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

	high, low, err := getForecast(r.client, r.url)
	if err != nil {
		r.err = err
		return
	}
	r.high = high
	r.low  = low
}

// Format and print current temperature.
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("weather: %v °F (%v/%v)", r.temp, r.high, r.low)
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
// Our value should be here: properties -> forecast.
func getURL(client http.Client, lat string, long string) (string, error) {
	type props struct {
		// Properties map[string]interface{} `json:"properties"`
		Properties struct {
			Forecast string `json: "temperature"`
		}`json:"properties"`
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

	url = p.Properties.Forecast
	if url == "" {
		return "", errors.New("Missing temperature URL")
	}

	return url, nil
}

// Get the current temperature from the NWS database.
// Our value should be here: properties -> periods -> (latest period) -> temperature.
func getTemp(client http.Client, url string) (int, error) {
	type temp struct {
		Properties struct {
			Periods []interface{} `json:"periods"`
		} `json:"properties"`
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

	// Get the list of weather readings.
	periods := t.Properties.Periods
	if len(periods) == 0 {
		return -1, errors.New("Missing hourly temperature periods")
	}

	// Use the most recent reading.
	latest := periods[0].(map[string]interface{})
	if len(latest) == 0 {
		return -1, errors.New("Missing current temperature")
	}

	// Get just the temperature reading.
	temperature := latest["temperature"].(float64)
	return int(temperature), nil
}

// Get the forecasted temperatures from the NWS database.
// Our values should be here: properties -> periods -> (chosen periods) -> temperature.
// We're going to use these rules to determine which day's forecast we want:
//   1. If it's before noon, we'll use the current day.
func getForecast(client http.Client, url string) (int, int, error) {
	var  high     float64
	var  low      float64
	var  high_s   string
	var  low_s    string

	type forecast struct {
		Properties struct {
			Periods []map[string]interface{} `json:"periods"`
		} `json:"properties"`
	}

	// Determine which day's forecast we want.
	t := time.Now()
	if (t.Hour() < 12) {
		high_s = "Today"
		low_s  = "Tonight"
	} else {
		t       = t.Add(time.Hour * 24)
		d      := t.Weekday()
		high_s  = d.String()
		low_s   = high_s + " Night"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, err
	}

	f   := forecast{}
	err  = json.Unmarshal(body, &f)
	if err != nil {
		return -1, -1, err
	}

	// Get the list of forecasts.
	periods := f.Properties.Periods
	if len(periods) == 0 {
		return -1, -1, errors.New("Missing forecast periods")
	}

	// Iterate through the list until we find the forecast for tomorrow.
	for _, f := range periods {
		fmt.Println(f)
		name := f["name"].(string)
		if name == high_s {
			// We'll get the high from here.
			high = f["temperature"].(float64)
		} else if name == low_s {
			// We'll get the low from here.
			low = f["temperature"].(float64)

			// This is all we need from the forecast, so we can exit now.
			return int(high), int(low), nil
		}
	}

	// If we're here, then we didn't find the forecast.
	return -1, -1, errors.New("Failed to determine forecast")
}
