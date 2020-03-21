// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/fanux/lvscare/care"
	"github.com/spf13/cobra"
)

// health checks
var (
	HealthPath  string
	HealthSchem string // http or https
)

// careCmd represents the care command
var careCmd = &cobra.Command{
	Use:   "care",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		care.VsAndRsCare(VirtualServer, RealServer, 5, HealthPath, HealthSchem)
	},
}

func init() {
	rootCmd.AddCommand(careCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// careCmd.PersistentFlags().String("foo", "", "A help for foo")
	careCmd.Flags().StringVar(&VirtualServer, "vs", "", "virturl server like 10.54.0.2:6443")
	careCmd.Flags().StringSliceVar(&RealServer, "rs", []string{}, "virturl server like 192.168.0.2:6443")

	careCmd.Flags().StringVar(&HealthPath, "health-path", "/healthz", "health check path")
	careCmd.Flags().StringVar(&HealthSchem, "health-schem", "https", "health check schem")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// careCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
