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
	"time"

	"github.com/spf13/cobra"

	"github.com/RaphaelPour/blogctl/internal/common"
	"github.com/RaphaelPour/blogctl/internal/metadata"
	"github.com/fatih/color"
)

const CONTENT_FILE = "content.md"

func GetContentFile(postPath string) string {
	return filepath.Join(postPath, CONTENT_FILE)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new post",
	RunE: func(cmd *cobra.Command, args []string) error {

		if Title == "" {
			return fmt.Errorf("Title missing")
		}

		/* Generate slug */
		slug := common.Slug(Title)
		postPath := filepath.Join(BlogPath, slug)

		/* Check if slug already exists */

		/* Ask user for input */
		content := fmt.Sprintf("# %s\n\n", Title)

		/* Abort if file already exist to prevent data loss when checking
		 * it afterwards.
		 */
		if _, err := os.Stat(postPath); err == nil {
			return fmt.Errorf("Error adding new post '%s': File already exists", slug)
		}

		if Interactive {
			if err := common.Editor(&content); err != nil {
				return fmt.Errorf("Error getting content from user: %s", err)
			}
		}

		/* Save post at the right place */
		if err := os.Mkdir(postPath, os.ModeDir|os.ModePerm); err != nil {
			rescuePost(content)
			return fmt.Errorf("Error creating post dir: %s", err)
		}

		if err := os.WriteFile(GetContentFile(postPath), []byte(content), os.ModePerm); err != nil {
			rescuePost(content)
			return fmt.Errorf("Error writing post: %s", err)
		}

		/* Store metadata info */
		meta := &metadata.Metadata{
			Title:     Title,
			Status:    metadata.DRAFT_STATUS,
			Static:    Static,
			CreatedAt: time.Now().Unix(),
		}
		if err := meta.Save(postPath); err != nil {
			rescuePost(content)
			return err
		}

		color.Green(GetContentFile(postPath))
		return nil
	},
}

var (
	Title       string
	Interactive bool
	Static      bool
)

func init() {
	postCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&Title, "title", "t", "", "Title of the post (required)")
	addCmd.Flags().BoolVarP(&Static, "static", "s", false, "Creates static sites that are unlisted and manually used (like Impressum)")
	addCmd.Flags().BoolVarP(&Interactive, "interactive", "i", false, "Opens default editor with new post.")
}

func rescuePost(content string) {
	fmt.Printf("<<< POST\n%s\n>>> POST\n", content)
}
