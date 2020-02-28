package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"

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

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Provide config for the tool")

	flag.Parse()

	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	client := resty.
		New().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Cookie", cfg.Cookies)

	log.Println("fetching all question slugs")
	_, err = getAllQuestions(client, cfg)
	if err != nil {
		log.Fatal(err)
	}

}
