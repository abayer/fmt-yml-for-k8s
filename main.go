/*
 * The MIT License
 *
 * Copyright (c) 2019, CloudBees, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

func main() {
	inputFile := flag.String("file", "", "File to read in and format")
	outputDir := flag.String("output-dir", "", "Directory to write the formatted output to")

	flag.Parse()

	if inputFile == nil {
		fmt.Println("No --file specified")
		os.Exit(1)
	}
	if outputDir == nil {
		fmt.Println("No --output-dir specified")
		os.Exit(1)
	}

	fileOk, err := fileExists(*inputFile)
	if err != nil {
		fmt.Printf("Error checking if input file %s exists:\n%s", *inputFile, err)
		os.Exit(1)
	}
	if !fileOk {
		fmt.Printf("Input file %s does not exist\n", *inputFile)
		os.Exit(1)
	}
	dirOk, err := dirExists(*outputDir)
	if err != nil {
		fmt.Printf("Error checking if output directory %s exists:\n%s", *outputDir, err)
		os.Exit(1)
	}
	if !dirOk {
		err = os.MkdirAll(*outputDir, 0760)
		if err != nil {
			fmt.Printf("Error creating output directory %s:\n%s", *outputDir, err)
			os.Exit(1)
		}
	}

	outputFileName := filepath.Base(*inputFile)

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file %s:\n%s", *inputFile, err)
		os.Exit(1)
	}

	var yamlAsStruct interface{}

	err = yaml.Unmarshal(bytes, &yamlAsStruct)
	if err != nil {
		fmt.Printf("Could not unmarshal contents of %s as YAML:\n%s", *inputFile, err)
		os.Exit(1)
	}

	outputData, err := yaml.Marshal(&yamlAsStruct)
	if err != nil {
		fmt.Printf("Could not marshal contents of %s as formatted YAML:\n%s", *inputFile, err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(filepath.Join(*outputDir, outputFileName), outputData, 0760)
	if err != nil {
		fmt.Printf("Couldn't write formatted YAML to %s:\n%s", filepath.Join(*outputDir, outputFileName), err)
		os.Exit(1)
	}
	os.Exit(0)
}

// fileExists checks if path exists and is a file
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, errors.Wrapf(err, "failed to check if file exists %s", path)
}

// dirExists checks if path exists and is a directory
func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
