package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"

type SubmissionMap map[string]Submission

type Config struct {
	Contest string `yaml:"contest"`
	Cookies string `yaml:"cookies"`
	Output  string `yaml:"output"`
}

func parseConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file, %v", err)
	}

	return &cfg, nil
}

type Questions struct {
	Models []struct {
		ID   int    `json:"id"`
		Slug string `json:"slug"`
	} `json:"models"`
	Total int `json:"total"`
}

func buildGetQuestionsUrl(contestSlug string, offset int, limit int) string {
	return fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s/"+
		"challenges?offset=%d&filters=:true+page:1&limit=%d", contestSlug, offset, limit)
}

func getQuestions(client *resty.Client, contest string, offset int, limit int) (*Questions, error) {
	log.Printf("fetching %d qustions from offset %d", limit, offset)
	resp, err := client.R().Get(buildGetQuestionsUrl(contest, offset, limit))
	if err != nil {
		return nil, fmt.Errorf("error getting questions, %v", err)
	}
	var questions Questions
	err = json.Unmarshal(resp.Body(), &questions)
	if err != nil {
		return nil, fmt.Errorf("error parsing questions, %v", err)
	}
	return &questions, nil
}

func getAllQuestions(client *resty.Client, config *Config) (*Questions, error) {
	limit := 50
	qs, err := getQuestions(client, config.Contest, 0, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching question names, %v", err)
	}

	// we get total questions in two requests
	if qs.Total > limit {
		// ok we got more to download
		qs2, err := getQuestions(client, config.Contest, limit, qs.Total)
		if err != nil {
			return nil, fmt.Errorf("error fetching question names, %v", err)
		}
		qs.Models = append(qs.Models, qs2.Models...)
	}
	return qs, nil
}

type Submission struct {
	ID             int     `json:"id"`
	HackerID       int     `json:"hacker_id"`
	CreatedAt      int     `json:"created_at"`
	Kind           string  `json:"kind"`
	Language       string  `json:"language"`
	Score          float64 `json:"score"`
	HackerUsername string  `json:"hacker_username"`
}

type Submissions struct {
	Models []Submission `json:"models"`
	Total  int          `json:"total"`
}

func buildGetSubmissionsUrl(contest string, question string, offset int, limit int) string {
	return fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s"+
		"/judge_submissions/?offset=%d&"+
		"limit=%d&challenge_id=%s",
		contest, offset, limit, question)
}

func getContestSubmissionData(client *resty.Client, contest string, question string, offset int, limit int) (*Submissions, error) {
	log.Printf("fetching submissions:%s limit %d offset %d", question, limit, offset)
	resp, err := client.R().Get(buildGetSubmissionsUrl(contest, question, offset, limit))
	if err != nil {
		return nil, fmt.Errorf("error getting submissions, %v", err)
	}
	var submissions Submissions
	err = json.Unmarshal(resp.Body(), &submissions)
	if err != nil {
		return nil, fmt.Errorf("error parsing submissions, %v", err)
	}
	return &submissions, nil
}

// filter submissions and filter only top score, latest submission for a team
func filterSubmissions(submissions *Submissions) SubmissionMap {
	hm := make(SubmissionMap)

	for _, s := range submissions.Models {
		if val, ok := hm[s.HackerUsername]; ok {
			if s.Score > val.Score {
				hm[s.HackerUsername] = s
			} else if s.Score == val.Score && s.CreatedAt > val.CreatedAt {
				hm[s.HackerUsername] = s
			}
			continue
		}
		hm[s.HackerUsername] = s
	}

	return hm
}

func getAllSubmissions(client *resty.Client, config *Config, question string) (SubmissionMap, error) {
	limit := 50
	qs, err := getContestSubmissionData(client, config.Contest, question, 0, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching question names, %v", err)
	}

	// we get total questions in two requests
	if qs.Total > limit {
		// ok we got more to download
		qs2, err := getContestSubmissionData(client, config.Contest, question, limit, qs.Total)
		if err != nil {
			return nil, fmt.Errorf("error fetching question names, %v", err)
		}
		qs.Models = append(qs.Models, qs2.Models...)
	}
	return filterSubmissions(qs), nil
}

func exec(cfg *Config) {
	client := resty.
		New().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Cookie", cfg.Cookies)

	log.Println("fetching all question slugs")
	qs, err := getAllQuestions(client, cfg)
	if err != nil {
		log.Fatal(err)
	}

	outPath := filepath.Join(filepath.FromSlash(cfg.Output), time.Now().Format("2006-01-02-15-04"))
	log.Printf("creating output directory: %v", outPath)
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating output directory, %v", err)
	}

	for _, q := range qs.Models {
		_, err := getAllSubmissions(client, cfg, q.Slug)
		if err != nil {
			log.Fatalf("error fetching submissions, %v", err)
		}
		fmt.Println("waiting 1s...")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Provide config for the tool")
	flag.Parse()

	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	exec(cfg)
}
