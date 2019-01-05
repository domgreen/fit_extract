package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/domgreen/fit_extract/pkg/fitbit"
	"github.com/spf13/pflag"
	strava "github.com/strava/go.strava"
)

func main() {
	var ft *string = pflag.String("fit_token", "", "Token for auth against Fitbit apis")
	var stravaAccessToken *string = pflag.String("strava_token", "", "Token for auth with Strava")
	pflag.Parse()

	if *ft == "" {
		os.Exit(1)
	}

	if *stravaAccessToken == "" {
		os.Exit(1)
	}

	client := strava.NewClient(*stravaAccessToken)
	service := strava.NewUploadsService(client)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	p := filepath.Join(usr.HomeDir, "/.fb_extract/afterDate")
	fmt.Println(p)
	dat, _ := ioutil.ReadFile(p)

	afterDate := "2000-01-08T00:01:01"
	if len(dat) > 0 {
		afterDate = strings.TrimSpace(string(dat))
	}

	nextUrl := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/list.json?afterDate=%s&sort=asc&offset=0&limit=20", afterDate)
	fmt.Println(nextUrl)
	for nextUrl != "" {
		req, _ := http.NewRequest("GET", nextUrl, nil)
		req.Header.Add("Authorization", "Bearer "+*ft)
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		var al = new(fitbit.ActivitiesList)
		json.Unmarshal(body, &al)
		defer resp.Body.Close()

		if len(al.Activities) == 0 {
			break
		}

		for _, v := range al.Activities {
			if v.ActivityTypeID != 3000 {
				continue
			}

			req, _ = http.NewRequest("GET", v.HeartRateLink, nil)
			req.Header.Add("Authorization", "Bearer "+*ft)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}
			body, _ = ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var a = new(fitbit.ActivityFullHeartRate)
			json.Unmarshal(body, &a)
			uploadToStrava(service, genGPX(a), a.ActivitiesHeart[0].DateTime)
			afterDate = a.ActivitiesHeart[0].DateTime + "T" + a.ActivitiesHeartIntraday.Dataset[0].Time
			fmt.Printf("uploaded activity from %s\n", afterDate)
		}

		if al.Pagination.Next == "" {
			t, err := time.Parse("2006-01-02T03:04:05", afterDate)
			if err != nil {
			}
			afterDate = t.Add(time.Second).Format("2006-01-02T03:04:05")
			fmt.Printf("Saving %s", afterDate)
			err = ioutil.WriteFile(p, []byte(afterDate), 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
		nextUrl = al.Pagination.Next
		fmt.Println(nextUrl)
		time.Sleep(time.Millisecond * 10)
	}
}

func uploadToStrava(service *strava.UploadsService, data string, date string) {
	_, err := service.Create(strava.FileDataTypes.GPX, "upload.gpx", strings.NewReader(data)).
		ActivityType(strava.ActivityTypes.Crossfit).
		Name(fmt.Sprintf("Crossfit | %s", date)).
		Do()

	if e, ok := err.(strava.Error); ok && strings.Contains(e.Message, "duplicate of activity") {
		fmt.Println(e.Message)
		return
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genGPX(a *fitbit.ActivityFullHeartRate) string {
	format := "2006-01-02T15:04:05Z"
	t, err := time.Parse("2006-01-02 03:04:05", a.ActivitiesHeart[0].DateTime+" "+a.ActivitiesHeartIntraday.Dataset[0].Time)
	if err != nil {
		fmt.Println(err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`
	<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="dom_green" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
 <metadata>
  <time>%s</time>
 </metadata>
 <trk>
  <name>%s</name>
  <trkseg>`, t.Format(format), "Crossfit | "+t.Format("2006-01-02")))
	for _, v := range a.ActivitiesHeartIntraday.Dataset {
		t, err := time.Parse("2006-01-02 03:04:05", a.ActivitiesHeart[0].DateTime+" "+v.Time)
		if err != nil {
			break
		}
		buffer.WriteString(fmt.Sprintf(`
   <trkpt>
    <ele>10</ele>
    <time>%s</time>
    <extensions>
     <gpxtpx:TrackPointExtension>
      <gpxtpx:hr>%s</gpxtpx:hr>
     </gpxtpx:TrackPointExtension>
    </extensions>
   </trkpt>`, t.Format(format), strconv.Itoa(v.Value)))
	}

	buffer.WriteString(`
  </trkseg>
 </trk>
</gpx>
`)
	return buffer.String()
}
