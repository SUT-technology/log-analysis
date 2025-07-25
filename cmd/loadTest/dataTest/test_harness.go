package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type EventSummary struct {
	EventName  string    `json:"event_name"`
	LastOccur  time.Time `json:"last_occur"`
	TotalCount int       `json:"total_count"`
}

type ListEventsResponse struct {
	ProjectID string         `json:"project_id"`
	Data      []EventSummary `json:"data"`
}

var (
	apiURL     = flag.String("url", "http://localhost:8000/api/logs", "Base URL of ListEvents API")
	projects   = flag.String("projects", "055337e7-37c1-4900-8048-b66ca2d6a061,2ff9f458-4921-4c12-802b-041a4a31fdcd,636abdc0-1492-4859-9c26-3771d47fbb78", "Comma-separated list of project IDs to test")
	outputFile = flag.String("output", "metrics.csv", "CSV output file")
)

func main() {
	flag.Parse()
	projectIDs := strings.Split(*projects, ",")

	// آماده‌سازی CSV
	file, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// هدر CSV
	writer.Write([]string{
		"project_id", "event_name", "searchable_keys", "page",
		"status_code", "duration_ms", "total_count", "returned_items",
	})

	scenarios := []struct {
		name string
		keys map[string]string
	}{
		{"no_filter", nil},
		{"by_event", map[string]string{"": ""}}, // فقط نام رو میاریم با event_name=login
		{"by_key", map[string]string{"k1": "v1_0"}},
		{"combo", map[string]string{"k2": "v2_0"}},
	}

	for _, pid := range projectIDs {
		for _, sc := range scenarios {
			for page := 1; page <= 5; page++ {
				q := url.Values{}
				q.Set("project_id", pid)
				q.Set("page", strconv.Itoa(page))

				if sc.name == "by_event" {
					q.Set("event_name", "login")
				}
				for k, v := range sc.keys {
					if k != "" {
						q.Set(fmt.Sprintf("searchable_keys[%s]", k), v)
					}
				}

				fullURL := *apiURL + "?" + q.Encode()
				start := time.Now()
				resp, err := http.Get(fullURL)
				duration := time.Since(start)

				status := 0
				var body ListEventsResponse
				if err == nil {
					defer resp.Body.Close()
					status = resp.StatusCode
					json.NewDecoder(resp.Body).Decode(&body)
				}

				writer.Write([]string{
					pid,
					sc.name,
					fmt.Sprintf("%v", sc.keys),
					strconv.Itoa(page),
					strconv.Itoa(status),
					strconv.FormatInt(duration.Milliseconds(), 10),
					strconv.Itoa(len(body.Data)), // در این API total_count برمی‌گردد
					strconv.Itoa(len(body.Data)),
				})
			}
		}
	}

	fmt.Println("Performance data written to", *outputFile)
}
