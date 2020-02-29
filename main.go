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
	"flag"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type SaveFileRequest struct {
	s *SubmissionData
	q string
}

var saveChan = make(chan SaveFileRequest)

func saveFiles(wg *sync.WaitGroup, config *Config) {
	for {
		select {
		case sf := <-saveChan:
			wg.Add(1)
			if err := SaveDownload(config, sf.q, sf.s); err != nil {
				log.Errorf("error downloading %s:%s", sf.q, sf.s.HackerUsername)
			}
			wg.Done()
		}
	}
}

func exec(cfg *Config) {
	wg := sync.WaitGroup{}
	client := CreateClient(cfg)

	log.Infof("creating output directory: %v", cfg.OutDir)
	err := os.MkdirAll(cfg.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating output directory, %v", err)
	}

	log.Info("fetching all question meta data")
	qs, err := GetAllQuestions(client, cfg)
	if err != nil {
		log.Fatal(err)
	}

	var subs = make(map[string]SubmissionMap)
	for _, q := range qs.Models {
		s, err := GetAllSubmissions(client, cfg, q.Slug)
		if err != nil {
			log.Fatalf("error fetching submissions, %v", err)
		}
		subs[q.Slug] = s
		log.Info("waiting 1s...")
		time.Sleep(1 * time.Second)
	}

	go saveFiles(&wg, cfg)

	for q, submissionMap := range subs {
		log.Infof("downloading %d submissions for %s", q, len(submissionMap))
		for _, submission := range submissionMap {
			ds, err := DownloadSubmission(client, cfg, &submission)
			if err != nil {
				log.Fatal(err)
			}

			saveChan <- SaveFileRequest{
				s: ds,
				q: q,
			}

			log.Info("sleeping 200ms...")
			time.Sleep(200 * time.Millisecond)
		}
	}

	log.Info("waiting to finish all tasks")
	wg.Wait()

	log.Info("finished executing")
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Provide config for the tool")
	flag.Parse()

	cfg, err := ParseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// create output directory with current time
	cfg.OutDir = filepath.Join(filepath.FromSlash(cfg.Output), time.Now().Format("2006-01-02-15-04"))

	exec(cfg)
}
