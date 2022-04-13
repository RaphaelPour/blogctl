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
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/RaphaelPour/blogctl/pkg/metadata"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available posts",
	RunE: func(cmd *cobra.Command, args []string) error {

		postDirs, err := ioutil.ReadDir(BlogPath)
		if err != nil {
			return fmt.Errorf("Error reading blog path: %s", err)
		}

		fmt.Println("Creation date                  | Status  |Title")
		fmt.Println("-------------------------------+---------+---------------")

		for _, dir := range postDirs {
			if !dir.IsDir() {
				continue
			}

			postPath := filepath.Join(BlogPath, dir.Name())
			files, err := ioutil.ReadDir(postPath)
			if err != nil {
				return fmt.Errorf("Error reading post path of %s: %s", postPath, err)
			}

			if len(files) < 2 {
				return fmt.Errorf(
					"Unexpected count of files in post path %s. Found: %d",
					postPath,
					len(files),
				)
			}

			metadata, err := metadata.Load(postPath)
			if err != nil {
				return err
			}

			fmt.Printf("%s | %-6s | %s\n",
				time.Unix(metadata.CreatedAt, 0).String(),
				metadata.Status,
				metadata.Title,
			)

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
