// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-router/extensions"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts mystack watcher",
	Long:  `starts mystack watcher`,
	Run: func(cmd *cobra.Command, args []string) {
		w, err := extensions.NewWatcher(config)
		if err != nil {
			panic(err)
		}
		w.Start()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
