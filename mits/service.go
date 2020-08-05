/*
   Copyright 2020 SUSE

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package mits

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/onsi/gomega/gexec"

	"github.com/SUSE/minibroker-integration-tests/mits/config"
)

const serviceKey = "test-credentials"

// Service represents a service instance to ease its manipulation from tests.
type Service struct {
	name       string
	brokerName string
	stdout     io.Writer
	stderr     io.Writer

	guid        string
	credentials map[string]interface{}
}

// NewService instantiates a new Service.
func NewService(
	name string,
	brokerName string,
	stdout io.Writer,
	stderr io.Writer,
) *Service {
	return &Service{
		name:        name,
		brokerName:  brokerName,
		stdout:      stdout,
		stderr:      stderr,
		credentials: nil,
	}
}

// Create creates the service instance on CF.
func (service *Service) Create(testConfig config.TestConfig, params map[string]interface{}, timeout time.Duration) error {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to create service instance: %w", err)
	}
	session := cf.Cf("create-service", "-b", service.brokerName, testConfig.Class, testConfig.Plan, service.name, "-c", string(paramsBytes)).
		Wait(timeout)
	if exitCode := session.ExitCode(); exitCode != 0 {
		return fmt.Errorf("failed to create service instance: cf create-service %s exited with code %d", service.name, exitCode)
	}
	var guidBuilder strings.Builder
	cmd := exec.Command("cf", "service", "--guid", service.name)
	session, err = gexec.Start(cmd, &guidBuilder, service.stderr)
	if err != nil {
		return fmt.Errorf("failed to create service instance: %w", err)
	}
	if exitCode := session.Wait(timeout).ExitCode(); exitCode != 0 {
		return fmt.Errorf("failed to create service instance: failed to get service instance guid")
	}
	service.guid = strings.TrimSpace(guidBuilder.String())
	return nil
}

// WaitForCreate waits for the creation of the service instance.
func (service *Service) WaitForCreate(timeout time.Duration) error {
	cond := conditions{
		progress:  "create in progress",
		completed: "create succeeded",
	}
	return service.waitForCondition(cond, timeout)
}

// WaitForDelete waits for the deletion of the service instance.
func (service *Service) WaitForDelete(timeout time.Duration) error {
	cond := conditions{
		progress:  "delete in progress",
		completed: "delete succeeded",
	}
	return service.waitForCondition(cond, timeout)
}

func (service *Service) waitForCondition(cond conditions, timeout time.Duration) error {
	timeLimit := time.Now().Add(timeout)
	for {
		if time.Now().After(timeLimit) {
			return fmt.Errorf("failed to wait for service instance: timed out")
		}

		cmd := exec.Command("cf", "service", service.name)
		pipeReader, pipeWriter := io.Pipe()
		session, err := gexec.Start(cmd, pipeWriter, service.stderr)
		if err != nil {
			return fmt.Errorf("failed to wait for service instance: %w", err)
		}
		go func() {
			defer pipeWriter.Close()
			session.Wait(timeout)
		}()

		scanner := bufio.NewScanner(pipeReader)
		var status string
		const statusPrefix = "status:"
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, statusPrefix) {
				status = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, statusPrefix)))
			}
		}

		if exitCode := session.ExitCode(); exitCode != 0 {
			return fmt.Errorf("failed to wait for service instance: cf service %s exited with code %d", service.name, exitCode)
		}

		if status == cond.progress || status == "" {
			time.Sleep(time.Second)
			continue
		} else if status == cond.completed {
			return nil
		} else {
			return fmt.Errorf("failed to wait for service instance: the service status is %q", status)
		}
	}
}

// Credentials creates a service-key in order to extract credentials for the service instance.
// It's useful for calculating the values of the security-group.
func (service *Service) Credentials(timeout time.Duration) (map[string]interface{}, error) {
	if service.credentials != nil {
		return service.credentials, nil
	}

	session := cf.Cf("create-service-key", service.name, serviceKey).Wait(timeout)
	if session.ExitCode() != 0 {
		return nil, fmt.Errorf("failed to get credentials for service instance: failed to create service key")
	}

	var guidBuilder strings.Builder
	cmd := exec.Command("cf", "service-key", "--guid", service.name, serviceKey)
	session, err := gexec.Start(cmd, &guidBuilder, service.stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for service instance: %w", err)
	}
	if exitCode := session.Wait(timeout).ExitCode(); exitCode != 0 {
		return nil, fmt.Errorf("failed to get credentials for service instance: failed to get service key guid")
	}
	serviceKeyGUID := strings.TrimSpace(guidBuilder.String())

	cmd = exec.Command("cf", "curl", "--fail", "/v2/service_keys/"+serviceKeyGUID)
	pipeReader, pipeWriter := io.Pipe()
	session, err = gexec.Start(cmd, pipeWriter, service.stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for service instance: %w", err)
	}
	go func() {
		defer pipeWriter.Close()
		session.Wait(timeout)
	}()

	var body map[string]interface{}
	decoder := json.NewDecoder(pipeReader)
	if err := decoder.Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to get credentials for service instance: %w", err)
	}

	// cf curl exits with code 22 if an error is reported back.
	// On success, it is returning non-zero codes, hence the unusual check below.
	if exitCode := session.ExitCode(); exitCode == 22 {
		return nil, fmt.Errorf("failed to get credentials for service instance: cf service-key %s %s exited with code %d", service.name, serviceKey, exitCode)
	}

	service.credentials = body["entity"].(map[string]interface{})["credentials"].(map[string]interface{})
	return service.credentials, nil
}

// Bind binds the service instance to an app.
func (service *Service) Bind(appName string, timeout time.Duration) error {
	session := cf.Cf("bind-service", appName, service.name).Wait(timeout)
	if exitCode := session.ExitCode(); exitCode != 0 {
		return fmt.Errorf("failed to bind service instance: cf bind-service %s %s exited with code %d", appName, service.name, exitCode)
	}
	return nil
}

// Unbind unbinds the service instance from an app.
func (service *Service) Unbind(appName string, timeout time.Duration) error {
	session := cf.Cf("unbind-service", appName, service.name).Wait(timeout)
	if exitCode := session.ExitCode(); exitCode != 0 {
		return fmt.Errorf("failed to unbind service instance: cf unbind-service %s %s exited with code %d", appName, service.name, exitCode)
	}
	return nil
}

// Destroy destroys all the created resources linked to the service instance.
func (service *Service) Destroy(timeout time.Duration) {
	cf.Cf("delete-service-key", service.name, serviceKey, "-f").Wait(timeout)
	cf.Cf("delete-service", service.name, "-f").Wait(timeout)
	service.WaitForDelete(timeout)
}

type conditions struct {
	progress  string
	completed string
}
