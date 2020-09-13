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
	"regexp"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new post",
	RunE: func(cmd *cobra.Command, args []string) error {

		if Title == "" {
			return fmt.Errorf("Title missing")
		}
		/* Generate slug */

		/* Check if sluf already exists */

		return nil
	},
}

var (
	Title     string
	SlugRegex = regexp.MustCompile(`[^A-Za-z0-9-]`)
)

func init() {
	postCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVar(&Title, "title", "", "Title of the post (required)")
}

func slug(s string) string {
	return SlugRegex.ReplaceAllString(strings.ReplaceAll(strings.ToLower(s), " ", "-"), "")
}
