package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	ortfodb "github.com/ortfo/db"
)

const WAKATIME_API_URL = "https://wakatime.com/api/v1"

func wakatimeRequest(path string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", WAKATIME_API_URL+"/"+path, nil)
	if err != nil {
		err = fmt.Errorf("while creating request: %w", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(os.Getenv("WAKATIME_API_KEY")))))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("while sending request: %w", err)
		return
	}

	// bodyString, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Got response %s\n", bodyString)
	return
}

var timeSpentOnProjects struct {
	mu         sync.Mutex
	durations  map[string]time.Duration
	computedAt map[string]time.Time
} = struct {
	mu         sync.Mutex
	durations  map[string]time.Duration
	computedAt map[string]time.Time
}{
	durations:  make(map[string]time.Duration),
	computedAt: make(map[string]time.Time),
}

const WAKATIME_CACHE_LIFETIME = 24 * time.Hour

var timeSpentOnTechs map[string]time.Duration = make(map[string]time.Duration)

func TimeSpentOnProject(work ortfodb.AnalyzedWork) time.Duration {
	// In dev, don't calculate times, it just slows everything down
	// if IsDev() {
	// 	return 0, nil
	// }

	timeSpentOnProjects.mu.Lock()
	defer timeSpentOnProjects.mu.Unlock()

	for id, duration := range timeSpentOnProjects.durations {
		if id == work.ID && time.Since(timeSpentOnProjects.computedAt[id]) < WAKATIME_CACHE_LIFETIME {
			return duration
		}
	}

	var data struct {
		Data wakatimeProjectStats `json:"data"`
	}
	resp, err := wakatimeRequest("users/current/all_time_since_today?project=" + work.ID)
	if err != nil {
		color.Red("[!!] Could not get time spent on %s: while fetching wakatime API: %w", work.ID, err)
		return 0
	}
	if resp.StatusCode == 404 {
		color.Yellow("[!!] Could not get time spent on %s: not found on wakatime", work.ID)
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)
	duration := time.Duration(data.Data.TotalSeconds) * time.Second
	if duration.Seconds() != 0 {
		fmt.Printf("[  ] Time spent on %s is %s, via wakatime\n", work.ID, duration)
	}
	timeSpentOnProjects.durations[work.ID] = duration
	timeSpentOnProjects.computedAt[work.ID] = time.Now()
	return duration
}

type wakatimeBestDay struct {
	Date         string  `json:"date"`
	TotalSeconds float64 `json:"total_seconds"`
	Text         string  `json:"text"`
}

type wakatimeCategory struct {
	Name          string  `json:"name"`
	TotalSeconds  float64 `json:"total_seconds"`
	Percent       float64 `json:"percent"`
	Digital       string  `json:"digital"`
	Decimal       string  `json:"decimal"`
	Text          string  `json:"text"`
	Hours         int64   `json:"hours"`
	Minutes       int64   `json:"minutes"`
	MachineNameID *string `json:"machine_name_id,omitempty"`
}

type wakatimeProjectStats struct {
	TotalSeconds float64 `json:"total_seconds"`
	IsUpToDate   bool    `json:"is_up_to_date"`
	Range        struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"range"`
}

type wakatimeUserStats struct {
	ID                                              string             `json:"id"`
	UserID                                          string             `json:"user_id"`
	Range                                           string             `json:"range"`
	Start                                           string             `json:"start"`
	End                                             string             `json:"end"`
	Timeout                                         int64              `json:"timeout"`
	WritesOnly                                      bool               `json:"writes_only"`
	Timezone                                        string             `json:"timezone"`
	Holidays                                        int64              `json:"holidays"`
	Status                                          string             `json:"status"`
	CreatedAt                                       string             `json:"created_at"`
	ModifiedAt                                      string             `json:"modified_at"`
	HumanReadableDailyAverage                       string             `json:"human_readable_daily_average"`
	BestDay                                         wakatimeBestDay    `json:"best_day"`
	HumanReadableTotalIncludingOtherLanguage        string             `json:"human_readable_total_including_other_language"`
	Machines                                        []wakatimeCategory `json:"machines"`
	Projects                                        []wakatimeCategory `json:"projects"`
	OperatingSystems                                []wakatimeCategory `json:"operating_systems"`
	IsUpToDatePendingFuture                         bool               `json:"is_up_to_date_pending_future"`
	DaysIncludingHolidays                           int64              `json:"days_including_holidays"`
	TotalSecondsIncludingOtherLanguage              float64            `json:"total_seconds_including_other_language"`
	TotalSeconds                                    float64            `json:"total_seconds"`
	HumanReadableDailyAverageIncludingOtherLanguage string             `json:"human_readable_daily_average_including_other_language"`
	IsAlreadyUpdating                               bool               `json:"is_already_updating"`
	Editors                                         []wakatimeCategory `json:"editors"`
	DaysMinusHolidays                               int64              `json:"days_minus_holidays"`
	IsStuck                                         bool               `json:"is_stuck"`
	PercentCalculated                               int64              `json:"percent_calculated"`
	DailyAverageIncludingOtherLanguage              float64            `json:"daily_average_including_other_language"`
	Dependencies                                    []wakatimeCategory `json:"dependencies"`
	Categories                                      []wakatimeCategory `json:"categories"`
	IsUpToDate                                      bool               `json:"is_up_to_date"`
	DailyAverage                                    float64            `json:"daily_average"`
	Languages                                       []wakatimeCategory `json:"languages"`
	HumanReadableTotal                              string             `json:"human_readable_total"`
	IsCached                                        bool               `json:"is_cached"`
	Username                                        string             `json:"username"`
	IsIncludingToday                                bool               `json:"is_including_today"`
	HumanReadableRange                              string             `json:"human_readable_range"`
	IsCodingActivityVisible                         bool               `json:"is_coding_activity_visible"`
	IsOtherUsageVisible                             bool               `json:"is_other_usage_visible"`
}
