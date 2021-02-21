// Copyright 2021 tappythumbz development
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
)

// Log formats and logs application output
func Log(msg string, isError bool) error {
	logger.SetFormatter(&logger.JSONFormatter{})

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	standardFields := logger.Fields{
		"hostname": hostname,
		"appname":  "doggo",
	}

	if isError {
		logger.WithTime(time.Now()).WithFields(standardFields).Error(msg)
	} else {
		logger.WithTime(time.Now()).WithFields(standardFields).Info(msg)
	}

	return nil
}
