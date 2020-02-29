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

func GetAllQuestions(client *resty.Client, config *Config) (*Questions, error) {
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
