// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import "github.com/Sirupsen/logrus"

type NginxInterface interface {
	Reload(logger logrus.FieldLogger) error
	AssertConfig(filePath string, logger logrus.FieldLogger) error
}
