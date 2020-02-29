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
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

type SubmissionMap map[string]Submission

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

// filter submissions and filter only top score, oldest submission for a team
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

func GetAllSubmissions(client *resty.Client, config *Config, question string) (SubmissionMap, error) {
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

func DownloadSubmission(client *resty.Client, config *Config, submission *Submission) (*SubmissionData, error) {
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
