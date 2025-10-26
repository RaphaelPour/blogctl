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
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/RaphaelPour/blogctl/internal/metadata"
)

type RegexpSet []*regexp.Regexp

func (r RegexpSet) Match(s string) bool {
	for _, pattern := range r {
		if !pattern.MatchString(s) {
			return false
		}
	}
	return true
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available posts",
	RunE: func(cmd *cobra.Command, args []string) error {

		postDirs, err := os.ReadDir(BlogPath)
		if err != nil {
			return fmt.Errorf("Error reading blog path: %s", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Creation date", "Status", "Static", "Title")

		// Compile title regex pattern, if any
		titlePattern := make(RegexpSet, len(titleFilter))
		for i, pattern := range titleFilter {
			if titlePattern[i], err = regexp.Compile(pattern); err != nil {
				return fmt.Errorf("error parsing title regex %q: %w", pattern, err)
			}
		}

		for _, dir := range postDirs {
			if !dir.IsDir() {
				continue
			}

			postPath := filepath.Join(BlogPath, dir.Name())
			files, err := os.ReadDir(postPath)
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

			if statusFilter != nil && !slices.Contains(statusFilter, metadata.Status) {
				continue
			}

			if titlePattern != nil && !titlePattern.Match(metadata.Title) {
				continue
			}

			err = table.Append([]string{
				metadata.Date(),
				metadata.Status,
				fmt.Sprintf("%t", metadata.Static),
				metadata.Title,
			})

			if err != nil {
				return err
			}
		}

		if err := table.Render(); err != nil {

		}

		return nil
	},
}

var (
	statusFilter []string
	titleFilter  []string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVar(&statusFilter, "status", nil, "Filter for status. Comma-separated or multiple flags allowed.")
	listCmd.Flags().StringSliceVar(&titleFilter, "title", nil, "Filter for title. Comma-separated or multiple flags with regex allowed.")
}
