// mystack-router api
// https://github.com/topfreegames/mystack/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func getNginxBinary() (string, error) {
	out, err := exec.Command("which", "nginx").Output()
	if err != nil {
		return "", err
	}

	nginxBinary := string(out)
	if strings.Contains(nginxBinary, "not found") {
		return "", fmt.Errorf("nginx is not installed")
	}

	return nginxBinary, nil
}

//Start starts nginx
func Start(logger logrus.FieldLogger) error {
	logger.Info("starting nginx")

	nginxBinary, err := getNginxBinary()
	if err != nil {
		return err
	}

	cmd := exec.Command(nginxBinary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	logger.Info("nginx successfully started")
	return nil
}

//Reload reloads nginx
func Reload(logger logrus.FieldLogger) error {
	logger.Info("reloading nginx")

	nginxBinary, err := getNginxBinary()
	if err != nil {
		return err
	}

	cmd := exec.Command(nginxBinary, "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	logger.Info("nginx successfully started")
	return nil
}
