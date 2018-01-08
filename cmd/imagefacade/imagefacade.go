/*
Copyright (C) 2016 Black Duck Software, Inc.
http://www.blackducksoftware.com/

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

// Executable that caches images in a directory as tarballs.
package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	pdocker "bitbucket.org/bdsengineering/perceptor/pkg/docker"
)

type input struct {
	fromImage string
	tag       string
	digest    string // need this so we know what to look up inthe api.
}

var in input

func init() {
	// fromImage=busyBox or fromImage=85fioh87h998hojR8h98hf.
	flag.StringVar(&in.fromImage, "fromImage", "busybox123", "imageDigest or name Will have .tar at the end.")
	flag.StringVar(&in.tag, "tag", "", "tag, empty is ok.")
}

func main() {

	flag.Parse() // Scan the arguments list

	if in.fromImage == "" {
		panic("Need -image=...")
	}

	image := pdocker.NewImage(in.fromImage, in.tag)
	err := pdocker.PullImage(*image)

	if err != nil {
		log.Errorf("Error while making tar file: %s", err)
	} else {
		log.Infof("Ready to scan !!!!! %s %s", in.fromImage, in.tag)
	}
}
