/*
Copyright © 2020 Raphael Pour <info@raphaelpour.de>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "blogctl",
	Short: "Static markdown blog backend",
	Long:  "Blogctl manages blog markdown-based posts database-less and generates a static website on-demand",
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			fmt.Println("BuildVersion: ", BuildVersion)
			fmt.Println("BuildDate: ", BuildDate)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const DEFAULT_BLOG_PATH = "./blog/"

var (
	BuildDate    string
	BuildVersion string
	BlogPath     string
	Version      bool
)

func init() {

	rootCmd.PersistentFlags().StringVar(
		&BlogPath,
		"path",
		DEFAULT_BLOG_PATH,
		fmt.Sprintf("Path to blog content. Default: %s", DEFAULT_BLOG_PATH),
	)

	rootCmd.Flags().BoolVar(
		&Version,
		"version",
		false,
		fmt.Sprintf("Print build information."),
	)
}
