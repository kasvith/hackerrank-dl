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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func exec(cfg *Config) {
	client := CreateClient(cfg)

	log.Println("fetching all question meta data")
	qs, err := GetAllQuestions(client, cfg)
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
		s, err := GetAllSubmissions(client, cfg, q.Slug)
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
			ds, err := DownloadSubmission(client, cfg, &v)
			if err != nil {
				log.Fatal(err)
			}

			err = SaveDownload(cfg, q, ds)
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

	cfg, err := ParseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg.OutDir = filepath.Join(filepath.FromSlash(cfg.Output), time.Now().Format("2006-01-02-15-04"))

	exec(cfg)
}
