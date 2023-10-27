/*
Copyright © 2023 Raphael Pour <info@raphaelpour.de>

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

	"github.com/RaphaelPour/blogctl/internal/metadata"
	"github.com/spf13/cobra"
)

var rewriteMetadataCmd = &cobra.Command{
	Use:   "rewrite-metadata",
	Short: "Rewrite metadata files e.g. for pretty printing legacy ones",
	RunE: func(cmd *cobra.Command, args []string) error {

		if Slug != "" {
			return rewriteMetadata(filepath.Join(BlogPath, Slug))
		}

		postDirs, err := os.ReadDir(BlogPath)
		if err != nil {
			return fmt.Errorf("Error reading blog path: %s", err)
		}

		for _, dir := range postDirs {

			if !dir.IsDir() {
				continue
			}

			postPath := filepath.Join(BlogPath, dir.Name())
			if err := rewriteMetadata(postPath); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	adminCmd.AddCommand(rewriteMetadataCmd)
	rewriteMetadataCmd.Flags().StringVarP(&Slug, "slug", "s", "", "Slug of the post. Use all if empty.")
}

func rewriteMetadata(postPath string) error {
	meta, err := metadata.Load(postPath)
	if err != nil {
		return fmt.Errorf("%s: %w", postPath, err)
	}
	if err := meta.Save(postPath); err != nil {
		return fmt.Errorf("%s: %w", postPath, err)
	}
	fmt.Printf("processed %s\n", postPath)
	return nil
}
