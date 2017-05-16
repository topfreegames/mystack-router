// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

type Nginx struct{}

//Reload reloads nginx
func (*Nginx) Reload(logger logrus.FieldLogger) error {
	logger.Info("reloading nginx")

	cmd := exec.Command("nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logger.Fatal("Nginx exec error:", err)
		return err
	}

	logger.Info("nginx successfully started")
	return nil
}

//AssertConfig tests config file for correct syntax
func (*Nginx) AssertConfig(filePath string, logger logrus.FieldLogger) error {
	cmd := exec.Command("nginx", "-t", "-c", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logger.Fatal("Nginx exec error:", err)
		return err
	}

	logger.Info("correct nginx config file")
	return nil
}
