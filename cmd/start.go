// mystack-router api
// https://github.com/topfreegames/mystack/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack/mystack-router/extensions"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts mystack watcher",
	Long:  `starts mystack watcher`,
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		w, err := extensions.NewWatcher(config)
		if err != nil {
			panic(err)
		}
		err = w.Start()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
