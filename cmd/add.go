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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/RaphaelPour/blogctl/pkg/metadata"
	"github.com/fatih/color"
	shellquote "github.com/kballard/go-shellquote"
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
		slug := slug(Title)
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
			err := Open(&content)
			if err != nil {
				return fmt.Errorf("Error getting content from user: %s", err)
			}
		}

		/* Save post at the right place */
		if err := os.Mkdir(postPath, os.ModeDir|os.ModePerm); err != nil {
			rescuePost(content)
			return fmt.Errorf("Error creating post dir: %s", err)
		}

		if err := ioutil.WriteFile(GetContentFile(postPath), []byte(content), os.ModePerm); err != nil {
			rescuePost(content)
			return fmt.Errorf("Error writing post: %s", err)
		}

		/* Store metadata info */
		metadata := &metadata.Metadata{
			Title:     Title,
			CreatedAt: time.Now().Unix(),
		}
		if err := metadata.Save(postPath); err != nil {
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
	SlugRegex   = regexp.MustCompile(`[^A-Za-z0-9-]`)
)

func init() {
	postCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&Title, "title", "t", "", "Title of the post (required)")
	addCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", false, "Opens default editor with new post.")
}

func slug(s string) string {
	return SlugRegex.ReplaceAllString(strings.ReplaceAll(strings.ToLower(s), " ", "-"), "")
}

func rescuePost(content string) {
	fmt.Printf("<<< POST\n%s\n>>> POST\n", content)
}

func Open(initialValue *string) error {
	var editor string
	if val, ok := os.LookupEnv("EDITOR"); ok {
		editor = val
	}
	if val, ok := os.LookupEnv("VISUAL"); ok {
		editor = val
	}
	if editor == "" {
		return fmt.Errorf("set EDITOR or VISUAL variable")
	}

	file, err := ioutil.TempFile("", "new-post*.md")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(*initialValue); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, file.Name())

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	*initialValue = string(raw)
	return nil
}
