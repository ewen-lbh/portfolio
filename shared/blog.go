package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	ortfodb "github.com/ortfo/db"
)

type BlogEntry struct {
	Date              time.Time `yaml:"date"`
	Title             string    `yaml:"title"`
	Slug              string
	Content           string
	OtherMetadata     map[string]any `yaml:",inline"`
	Pageviews         int
	RelatedWorksSlugs []string `yaml:"works"`
}

func (e *BlogEntry) RelatedWorks(db ortfodb.Database) []ortfodb.AnalyzedWork {
	var works []ortfodb.AnalyzedWork
	for _, slug := range e.RelatedWorksSlugs {
		for id, work := range db.Works() {
			if id == slug {
				works = append(works, work)
				continue
			}
		}
	}

	return works
}

func (e *BlogEntry) GetPageviews(blogRoot string) error {
	views, err := getPageviewsFor(blogRoot + "/" + e.Slug)
	if err != nil {
		return err
	}

	e.Pageviews = views
	return nil
}

func getPageviewsFor(path string) (int, error) {
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://stats.ewen.works/api/v1/stats/aggregate?site_id=ewen.works&period=12mo&metrics=visitors,pageviews&filters=event:page==%s", url.QueryEscape(path)), bytes.NewBufferString(""))
	request.Header.Set("Authorization", "Bearer "+os.Getenv("PLAUSIBLE_TOKEN"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, err
	}
	if response.StatusCode != 200 {
		return 0, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	raw, _ := io.ReadAll(response.Body)

	var data plausibleAggregate
	json.Unmarshal(raw, &data)

	return data.Results.Pageviews.Value, nil
}

type plausibleAggregate struct {
	Results struct {
		Pageviews struct {
			Value int `json:"value"`
		} `json:"pageviews"`
		Visitors struct {
			Value int `json:"value"`
		} `json:"visitors"`
	} `json:"results"`
}
