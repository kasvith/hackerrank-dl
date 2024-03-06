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
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Contest           string   `yaml:"contest"`
	Cookies           string   `yaml:"cookies"`
	Output            string   `yaml:"output"`
	OutDir            string   `yaml:"-"`
	ParallelDownloads int      `yaml:"parallelDownloads"`
	MaxWaitTime       int      `yaml:"waitTime"`
	Rate              int      `yaml:"rate"`
	SpecificQuestions []string `yaml:"specificQuestions,omitempty"`
	SpecificUsers []string `yaml:"specificUsers,omitempty"`
}


func ParseConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file, %v", err)
	}

	if cfg.MaxWaitTime == 0 {
		cfg.MaxWaitTime = 1
	}
	if cfg.ParallelDownloads == 0 {
		cfg.ParallelDownloads = 5
	}
	if cfg.Rate == 0 {
		cfg.Rate = 10
	}

	return &cfg, nil
}
