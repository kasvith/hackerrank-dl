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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CreateDownloadDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveDownload(config *Config, question string, data *SubmissionData) error {
	// create dir
	dirPath := filepath.Join(config.OutDir, question, data.Language)
	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", data.HackerUsername, GetExtension(data.Language)))
	err := CreateDownloadDir(dirPath)
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
