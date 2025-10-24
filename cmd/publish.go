/*
Copyright © 2020 Raphael Pour <raphael.pour@hetzner.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/RaphaelPour/blogctl/internal/metadata"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Set post's status to public",
	RunE: func(cmd *cobra.Command, args []string) error {
		if Slug == "" {
			return fmt.Errorf("Slug missing")
		}

		metadataPath := filepath.Join(BlogPath, Slug)
		meta, err := metadata.Load(metadataPath)
		if err != nil {
			return err
		}

		meta.Status = metadata.PUBLIC_STATUS

		if err := meta.Save(metadataPath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	postCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&Slug, "slug", "s", "", "Slug of the post (required)")
	_ = publishCmd.RegisterFlagCompletionFunc("slug", SlugCompletion)
}
