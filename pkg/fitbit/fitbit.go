package fitbit

import "time"

type ActivitiesList struct {
	Activities []struct {
		ActiveDuration int `json:"activeDuration"`
		ActivityLevel  []struct {
			Minutes int    `json:"minutes"`
			Name    string `json:"name"`
		} `json:"activityLevel"`
		ActivityName     string `json:"activityName"`
		ActivityTypeID   int    `json:"activityTypeId"`
		AverageHeartRate int    `json:"averageHeartRate"`
		Calories         int    `json:"calories"`
		CaloriesLink     string `json:"caloriesLink"`
		Duration         int    `json:"duration"`
		ElevationGain    int    `json:"elevationGain"`
		HeartRateLink    string `json:"heartRateLink"`
		HeartRateZones   []struct {
			Max     int    `json:"max"`
			Min     int    `json:"min"`
			Minutes int    `json:"minutes"`
			Name    string `json:"name"`
		} `json:"heartRateZones"`
		LastModified          time.Time `json:"lastModified"`
		LogID                 int64     `json:"logId"`
		LogType               string    `json:"logType"`
		ManualValuesSpecified struct {
			Calories bool `json:"calories"`
			Distance bool `json:"distance"`
			Steps    bool `json:"steps"`
		} `json:"manualValuesSpecified"`
		OriginalDuration  int       `json:"originalDuration"`
		OriginalStartTime time.Time `json:"originalStartTime"`
		StartTime         time.Time `json:"startTime"`
		Steps             int       `json:"steps,omitempty"`
		TcxLink           string    `json:"tcxLink"`
		Distance          float64   `json:"distance,omitempty"`
		DistanceUnit      string    `json:"distanceUnit,omitempty"`
		Pace              float64   `json:"pace,omitempty"`
		Source            struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"source,omitempty"`
		Speed float64 `json:"speed,omitempty"`
	} `json:"activities"`
	Pagination struct {
		AfterDate string `json:"afterDate"`
		Limit     int    `json:"limit"`
		Next      string `json:"next"`
		Offset    int    `json:"offset"`
		Previous  string `json:"previous"`
		Sort      string `json:"sort"`
	} `json:"pagination"`
}

type ActivityFullHeartRate struct {
	ActivitiesHeart []struct {
		CustomHeartRateZones []interface{} `json:"customHeartRateZones"`
		DateTime             string        `json:"dateTime"`
		HeartRateZones       []struct {
			CaloriesOut int    `json:"caloriesOut"`
			Max         int    `json:"max"`
			Min         int    `json:"min"`
			Minutes     int    `json:"minutes"`
			Name        string `json:"name"`
		} `json:"heartRateZones"`
		Value string `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday struct {
		Dataset []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-heart-intraday"`
}
