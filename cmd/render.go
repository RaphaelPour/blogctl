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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/RaphaelPour/blogctl/pkg/metadata"
	"github.com/gomarkdown/markdown"
	"github.com/spf13/cobra"
)

type Post struct {
	Title     string
	Link      string
	Timestamp int64
	CreatedAt string
	Content   string
	Rendered  template.HTML
	Metadata  map[string]string
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders blog to static website",
	Long:  "Collects all posts and renders the markdown using the metadata as static website",
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(OutPath); !os.IsNotExist(err) && !Force {
			return fmt.Errorf("Output folder already exists")
		}

		if err := os.MkdirAll(OutPath, os.ModePerm); err != nil {
			return fmt.Errorf("Error creating output folder: %s", err)
		}

		postDirs, err := ioutil.ReadDir(BlogPath)
		if err != nil {
			return fmt.Errorf("Error reading blog path: %s", err)
		}

		posts := make([]Post, 0)
		for i, dir := range postDirs {

			if !dir.IsDir() {
				continue
			}

			postPath := filepath.Join(BlogPath, dir.Name())
			files, err := ioutil.ReadDir(postPath)
			if err != nil {
				return fmt.Errorf("Error reading post path of %s: %s", postPath, err)
			}

			if len(files) != 2 {
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

			/* Overstep posts which aren't set to 'public' */
			if metadata.Status != "public" {
				continue
			}

			fmt.Printf("Rendering post #%02d: %s\n", i, dir.Name())

			content, err := ioutil.ReadFile(GetContentFile(postPath))
			if err != nil {
				return fmt.Errorf("Error reading post content %s: %s", postPath, err)
			}

			rendered := markdown.ToHTML(content, nil, nil)

			timestamp := time.Unix(metadata.CreatedAt, 0)
			postFileName := fmt.Sprintf(
				POST_FILE_TEMPLATE,
				slug(metadata.Title),
			)
			post := Post{
				Title:     metadata.Title,
				Link:      postFileName,
				Timestamp: timestamp.Unix(),
				CreatedAt: timestamp.String(),
				Content:   string(content),
				Rendered:  template.HTML(rendered),
			}
			posts = append(posts, post)

			/* Render single post */
			postTemplate, err := template.New("post").Parse(POST_TEMPLATE)
			if err != nil {
				return fmt.Errorf("Error creating post file '%s': %s", metadata.Title, err)
			}

			postFilePath := filepath.Join(OutPath, postFileName)
			file, err := os.Create(postFilePath)
			if err != nil {
				return fmt.Errorf("Error creating post file '%s': %s", metadata.Title, err)
			}

			if err := postTemplate.Execute(file, post); err != nil {
				return fmt.Errorf("Error rendering post '%s': %s", metadata.Title, err)
			}

			if err := file.Close(); err != nil {
				return fmt.Errorf("Error closing post file '%s': %s", metadata.Title, err)
			}
		}

		/* Put everything together */
		t, err := template.New("blog").Parse(INDEX_TEMPLATE)
		if err != nil {
			return fmt.Errorf("Error parsing the html template: %s", err)
		}

		/* Save site to out dir  */
		sitePath := filepath.Join(OutPath, INDEX_FILE)
		file, err := os.Create(sitePath)
		if err != nil {
			return fmt.Errorf("Error creating index file: %s", err)
		}

		/* Sort posts */
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp > posts[j].Timestamp
		})

		if err := t.Execute(file, posts); err != nil {
			return fmt.Errorf("Error rendering posts: %s", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("Error closing file: %s", err)
		}

		return nil
	},
}

const (
	DEFAULT_OUT_PATH   = "./out/"
	INDEX_FILE         = "index.html"
	POST_FILE_TEMPLATE = "%s.html"
	INDEX_TEMPLATE     = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Blog</title>
		<style>
			h1 { margin:0px;}
			pre { width:100%;overflow:auto}
			.date { margin-top:10px;font-size: small; color: gray; }
			.post { margin-top:10px;}
		</style>
	</head>
	<body>
		<ul>
		{{range .}}
		<li>
		  <span class='date'>[{{.CreatedAt}}]</span>
			<a href='{{ .Link }}'>{{ .Title }}</a>
		</li>
		{{else}}<li><strong>no posts</strong></li>{{end}}
		</ul>
	</body>
</html>`
	POST_TEMPLATE = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Blog</title>
		<style>
			h1 { margin:0px;}
			.date { margin-top:10px;font-size: small; color: gray; }
			.post { margin-top:10px;}
		</style>
	</head>
	<body>
		<div class='post'>
			<span class='date'>{{.CreatedAt}}</span>
			{{ .Rendered }}
		</div>
	</body>
</html>`
)

var (
	OutPath string
	Force   bool
)

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(
		&OutPath,
		"out",
		"o",
		DEFAULT_OUT_PATH,
		"Output folder.",
	)

	renderCmd.Flags().BoolVarP(
		&Force,
		"force",
		"f",
		false,
		"Overwrites an existing output folder.",
	)
}
