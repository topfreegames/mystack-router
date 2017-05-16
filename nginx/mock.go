// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import "github.com/Sirupsen/logrus"

type Mock struct {
	Err error
}

//Reload reloads nginx
func (m *Mock) Reload(logger logrus.FieldLogger) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}

//AssertConfig tests config file for correct syntax
func (m *Mock) AssertConfig(filePath string, logger logrus.FieldLogger) error {
	if m.Err != nil {
		return m.Err
	}
	return nil
}
