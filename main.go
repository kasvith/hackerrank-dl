/*
 * MIT License
 *
 * Copyright (c) 2020 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"

type SubmissionMap map[string]Submission

type Config struct {
	Contest string `yaml:"contest"`
	Cookies string `yaml:"cookies"`
	Output  string `yaml:"output"`
	OutDir  string `yaml:"-"`
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
			} else if s.Score == val.Score && s.CreatedAt < val.CreatedAt {
				// if we have two submissions with same score, pick the oldest one
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

type SubmissionData struct {
	ID             int    `json:"id"`
	Language       string `json:"language"`
	Code           string `json:"code"`
	HackerID       int    `json:"hacker_id"`
	HackerUsername string `json:"hacker_username"`
}

type SubmissionDataResp struct {
	Model SubmissionData `json:"model"`
}

func buildSubmissionDownloadUrl(contest string, id int) string {
	return fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s/submissions/%d", contest, id)
}

func downloadSubmission(client *resty.Client, config *Config, submission *Submission) (*SubmissionData, error) {
	log.Printf("downloading submission %s:%d", submission.HackerUsername, submission.ID)
	resp, err := client.R().Get(buildSubmissionDownloadUrl(config.Contest, submission.ID))
	if err != nil {
		return nil, fmt.Errorf("error getting submission data for %d, %v", submission.ID, err)
	}
	var submissionData SubmissionDataResp
	err = json.Unmarshal(resp.Body(), &submissionData)
	if err != nil {
		return nil, fmt.Errorf("error parsing submission data for %d, %v", submission.ID, err)
	}
	// hr does not set username for same person so we set it explicitly
	submissionData.Model.HackerUsername = submission.HackerUsername
	return &submissionData.Model, nil
}

func createDownloadDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveDownload(config *Config, question string, data *SubmissionData) error {
	// create dir
	dirPath := filepath.Join(config.OutDir, question, data.Language)
	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", data.HackerUsername, GetExtension(data.Language)))
	err := createDownloadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error creating download directory, %v", err)
	}
	log.Printf("saving file: %s", filePath)
	err = ioutil.WriteFile(filePath, []byte(data.Code), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error saving file, %v", err)
	}
	return nil
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

	log.Printf("creating output directory: %v", cfg.OutDir)
	err = os.MkdirAll(cfg.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating output directory, %v", err)
	}

	var subs = make(map[string]SubmissionMap)
	for _, q := range qs.Models {
		s, err := getAllSubmissions(client, cfg, q.Slug)
		if err != nil {
			log.Fatalf("error fetching submissions, %v", err)
		}
		subs[q.Slug] = s
		fmt.Println("waiting 1s...")
		time.Sleep(1 * time.Second)
	}

	for q, s := range subs {
		log.Printf("downloading submissions for %s", q)
		for _, v := range s {
			ds, err := downloadSubmission(client, cfg, &v)
			if err != nil {
				log.Fatal(err)
			}

			err = saveDownload(cfg, q, ds)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("sleeping 200ms...")
			time.Sleep(200 * time.Millisecond)
		}
	}

	log.Println("finished executing")
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Provide config for the tool")
	flag.Parse()

	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg.OutDir = filepath.Join(filepath.FromSlash(cfg.Output), time.Now().Format("2006-01-02-15-04"))

	exec(cfg)
}
