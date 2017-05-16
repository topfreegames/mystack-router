// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-router/extensions"
	"github.com/topfreegames/mystack-router/models"
	"github.com/topfreegames/mystack-router/nginx"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts mystack watcher",
	Long:  `starts mystack watcher`,
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		w, err := extensions.NewWatcher(config, nil)
		if err != nil {
			panic(err)
		}
		fs := models.NewRealFS()
		err = w.CreateConfigFile(fs)
		if err != nil {
			panic(err)
		}
		ng := &nginx.Nginx{}
		err = w.Start(ng, fs)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
