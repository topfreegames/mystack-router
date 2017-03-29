// mystack-router api
// https://github.com/topfreegames/mystack/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package nginx

import (
//"github.com/Masterminds/sprig"
//	"github.com/topfreegames/mystack/mystack-router/model"
//"text/template"
)

const (
	confTemplate = ``
)

//WriteConfig writes a new nginx file config
//func WriteConfig(routerConfig *model.RouterConfig) error {
//	tmpl, err := template.New("nginx").Funcs(sprig.TxtFuncMap()).Parse(confTemplate)
//	if err != nil {
//		return err
//	}
//
//	file, err := os.Create(filePath)
//	if err != nil {
//		return err
//	}
//	err = tmpl.Execute(file, routerConfig)
//	return err
//}
