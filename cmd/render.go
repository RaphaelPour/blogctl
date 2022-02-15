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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/RaphaelPour/blogctl/pkg/metadata"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/feeds"
	"github.com/spf13/cobra"
)

type Post struct {
	Title            string
	Link             string
	PreviousPostLink string
	NextPostLink     string
	HomeLink         string
	Timestamp        int64
	CreatedAt        string
	Content          string
	Rendered         template.HTML
	Metadata         map[string]string
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders blog to static website",
	Long:  "Collects all posts and renders the markdown using the metadata as static website",
	RunE: func(cmd *cobra.Command, args []string) error {

		feed := &feeds.Feed{
			Title:       "evilcookie.de",
			Link:        &feeds.Link{Href: "https://evilcookie.de"},
			Description: "drunken stack developer",
			Author:      &feeds.Author{Name: "Raphael Pour"},
			Created:     time.Now(),
		}

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

		feed.Items = make([]*feeds.Item, 0)
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

			/* Overstep posts which aren't set to 'public' */
			if metadata.Status != "public" {
				continue
			}

			fmt.Printf("Rendering post #%02d: %s\n", i, dir.Name())
			slugTitle := slug(metadata.Title)

			content, err := ioutil.ReadFile(GetContentFile(postPath))
			if err != nil {
				return fmt.Errorf("Error reading post content %s: %s", postPath, err)
			}

			rendered := markdown.ToHTML(content, nil, nil)

			/* replace all IMAGE(<filename>) with valid path to filename */
			re := regexp.MustCompile(`IMAGE\(([\w\.]+)\)`)
			renderedStr := re.ReplaceAllString(string(rendered), fmt.Sprintf(`<img src="%s_$1"/>`, slugTitle))

			for _, file := range re.FindAllStringSubmatch(string(rendered), -1) {
				src := filepath.Join(BlogPath, fmt.Sprintf("%s/%s", slugTitle, file[1]))
				dst := filepath.Join(OutPath, fmt.Sprintf("%s_%s", slugTitle, file[1]))
				if err := copyFile(src, dst); err != nil {
					return fmt.Errorf("error copying '%s' to '%s': %w", src, dst, err)
				}
			}

			timestamp := time.Unix(metadata.CreatedAt, 0)
			postFileName := fmt.Sprintf(
				POST_FILE_TEMPLATE,
				slugTitle,
			)
			post := Post{
				Title:     metadata.Title,
				Link:      postFileName,
				HomeLink:  INDEX_FILE,
				Timestamp: timestamp.Unix(),
				CreatedAt: timestamp.String(),
				Content:   renderedStr,
				Rendered:  template.HTML(renderedStr),
			}
			posts = append(posts, post)
		}

		/* Sort posts */
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp > posts[j].Timestamp
		})

		/* set previous/next post of any post */
		for i := 0; i < len(posts); i++ {
			if i-1 >= 0 {
				posts[i].NextPostLink = posts[i-1].Link
			}

			if i+1 < len(posts) {
				posts[i].PreviousPostLink = posts[i+1].Link
			}
		}

		/* render all posts */
		for _, post := range posts {
			/* Render single post */
			postTemplate, err := template.New("post").Parse(POST_TEMPLATE)
			if err != nil {
				return fmt.Errorf("Error creating post file '%s': %s", post.Title, err)
			}

			postFilePath := filepath.Join(OutPath, post.Link)
			file, err := os.Create(postFilePath)
			if err != nil {
				return fmt.Errorf("Error creating post file '%s': %s", post.Title, err)
			}

			if err := postTemplate.Execute(file, post); err != nil {
				return fmt.Errorf("Error rendering post '%s': %s", post.Title, err)
			}

			if err := file.Close(); err != nil {
				return fmt.Errorf("Error closing post file '%s': %s", post.Title, err)
			}

			feed.Items = append(feed.Items, &feeds.Item{
				Title:   post.Title,
				Content: string(post.Rendered),
				Link: &feeds.Link{
					Href: fmt.Sprintf(
						"https://evilcookie.de/%s.html",
						slug(post.Title),
					),
				},
				Author:  &feeds.Author{Name: "Raphael Pour"},
				Created: time.Unix(post.Timestamp, 0),
			})
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

		if err := t.Execute(file, posts); err != nil {
			return fmt.Errorf("Error rendering posts: %s", err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("Error closing file: %s", err)
		}

		rss, err := feed.ToRss()
		if err != nil {
			return fmt.Errorf("Error generating rss feed: %w", err)
		}

		rssPath := filepath.Join(OutPath, "rss.xml")
		if err := ioutil.WriteFile(rssPath, []byte(rss), 0777); err != nil {
			return fmt.Errorf("Error writing rss.xml: %w", err)
		}

		return nil
	},
}

func copyFile(sourceFile, destinationFile string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("error open source file %s: %w", sourceFile, err)
	}
	defer src.Close()

	dst, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("error open destination file %s: %w", destinationFile, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}
	return nil
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
		<link rel="alternate" type="application/rss+xml" title="Feed" href="/rss.xml">
		<style>
			h1 { margin:0px;}
			pre { width:100%;overflow:auto}
			.date { margin-top:10px;font-size: small; color: gray; }
			.post { margin-top:10px;}
		</style>
	</head>
	<body>
	<h2>evilcookie</h2>
		Hi, I'm Raphael. I write about my software developer journey.
		You can find my stuff on <a href='https://github.com/RaphaelPour'>github</a>.
		<ul>
		{{range .}}
		<li>
		  <span class='date'>[{{.CreatedAt}}]</span>
			<a href='{{ .Link }}'>{{ .Title }}</a>
		</li>
		{{else}}<li><strong>no posts</strong></li>{{end}}
		</ul>
		<a href='rss.xml'>RSS</a>
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
			{{if .PreviousPostLink}}
			<a href='{{.PreviousPostLink}}'>&lt;</a>
			{{end}}
			<a href='{{.HomeLink}}'>up</a>
			{{if .NextPostLink}}
			<a href='{{.NextPostLink}}'>&gt;</a>
			{{end}}
			</br>

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
