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
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update post content",
	RunE: func(cmd *cobra.Command, args []string) error {

		if Slug == "" {
			return fmt.Errorf("Slug missing")
		}

		contentFile := GetContentFile(filepath.Join(BlogPath, Slug))
		if _, err := os.Stat(contentFile); os.IsNotExist(err) {
			return fmt.Errorf("Error updating post '%s': Not existing", Slug)
		}

		/* Read in post content */
		rawContent, err := os.ReadFile(contentFile)
		if err != nil {
			return fmt.Errorf("Error reading post '%s': %s", contentFile, err)
		}

		content := string(rawContent)

		/* Append stdin to content if append flag is set. This option is especially
		 * useful for two cases:
		 * 1. For testing without having a doubt about the interactive editor
		 * 2. Pipe output of a post-related command to the post directly
		 */
		if AppendStdin {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				content += scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				rescuePost(content)
				return fmt.Errorf("Error appending stdin: %s", err)
			}
		}

		/* Open editor with content if interactive flag is set*/
		if Interactive {
			if err := common.Open(&content); err != nil {
				rescuePost(content)
				return fmt.Errorf("Error getting content from user: %s", err)
			}
		}

		/* Write content back to file */
		if err := os.WriteFile(contentFile, []byte(content), os.ModePerm); err != nil {
			rescuePost(content)
			return fmt.Errorf("Error writing post: %s", err)
		}
		return nil
	},
}

var (
	Slug        string
	AppendStdin bool
)

func init() {
	postCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&Slug, "slug", "s", "", "Slug of the post (required)")
	updateCmd.Flags().BoolVarP(&AppendStdin, "append", "a", false, "Append stdin to post")
	updateCmd.Flags().BoolVarP(&Interactive, "interactive", "i", false, "Opend default editor with existing post")

	_ = updateCmd.RegisterFlagCompletionFunc("slug", SlugCompletion)
}
