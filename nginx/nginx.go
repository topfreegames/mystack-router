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
)

//Start starts nginx
func Start(logger logrus.FieldLogger) error {
	logger.Info("starting nginx")

	cmd := exec.Command("nginx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Nginx exec error:", err)
		return err
	}

	logger.Info("nginx successfully started")
	return nil
}

//Reload reloads nginx
func Reload(logger logrus.FieldLogger) error {
	logger.Info("reloading nginx")

	cmd := exec.Command("nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	logger.Info("nginx successfully started")
	return nil
}
