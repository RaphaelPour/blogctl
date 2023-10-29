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
	"os"
	"path/filepath"
	"regexp"
	"sort"
	plainTemplate "text/template"
	"time"

	_ "embed"

	"github.com/RaphaelPour/blogctl/internal/config"
	"github.com/RaphaelPour/blogctl/internal/highlighter"
	"github.com/RaphaelPour/blogctl/internal/metadata"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/feeds"
	"github.com/spf13/cobra"
)

var (
	//go:embed index.tmpl
	indexTemplate string

	//go:embed post.tmpl
	postTemplate string

	//go:embed static.tmpl
	staticTemplate string
)

type Post struct {
	Title            string
	Link             string
	PermaLink        string
	PreviousPostLink string
	NextPostLink     string
	HomeLink         string
	Timestamp        int64
	CreatedAt        string
	Content          string
	FeaturedImage    string
	Rendered         template.HTML
	Metadata         *metadata.Metadata
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders blog to static website",
	Long:  "Collects all posts and renders the markdown using the metadata as static website",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(BlogPath)
		if err != nil {
			return err
		}

		feed := &feeds.Feed{
			Title:       cfg.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s", cfg.Domain)},
			Description: cfg.Description,
			Author:      &feeds.Author{Name: cfg.Author},
			Created:     time.Now(),
		}

		if _, err := os.Stat(OutPath); !os.IsNotExist(err) && !Force {
			return fmt.Errorf("Output folder already exists")
		}

		if err := os.MkdirAll(OutPath, os.ModePerm); err != nil {
			return fmt.Errorf("Error creating output folder: %s", err)
		}

		postDirs, err := os.ReadDir(BlogPath)
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

			meta, err := metadata.Load(postPath)
			if err != nil {
				return err
			}

			/* Overstep posts which aren't set to 'public' */
			if meta.Status != metadata.PUBLIC_STATUS {
				continue
			}

			fmt.Printf("Rendering post #%02d: %s\n", i, dir.Name())
			slugTitle := slug(meta.Title)

			content, err := os.ReadFile(GetContentFile(postPath))
			if err != nil {
				return fmt.Errorf("Error reading post content %s: %s", postPath, err)
			}

			rendered := markdown.ToHTML(
				content, parser.NewWithExtensions(parser.CommonExtensions|parser.Footnotes),
				highlighter.GetRenderer(),
			)

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

			postFileName := fmt.Sprintf(
				POST_FILE_TEMPLATE,
				slugTitle,
			)
			var featuredImage string
			if len(meta.FeaturedImage) > 0 {
				featuredImage = fmt.Sprintf("https://%s/%s_%s", cfg.Domain, slugTitle, meta.FeaturedImage)
			}
			post := Post{
				Title:         meta.Title,
				Link:          postFileName,
				PermaLink:     fmt.Sprintf("https://%s/%s.html", cfg.Domain, slugTitle),
				HomeLink:      INDEX_FILE,
				Timestamp:     meta.CreatedAt,
				CreatedAt:     meta.Date(),
				Content:       renderedStr,
				Rendered:      template.HTML(renderedStr),
				FeaturedImage: featuredImage,
				Metadata:      meta,
			}
			posts = append(posts, post)
		}

		/* Sort posts */
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp > posts[j].Timestamp
		})

		/* set previous/next post of any post */
		nextPost := -1
		for i := 0; i < len(posts); i++ {
			if posts[i].Metadata.Static {
				continue
			}

			if nextPost >= 0 {
				posts[i].NextPostLink = posts[nextPost].Link
				posts[nextPost].PreviousPostLink = posts[i].Link

				fmt.Println(posts[i].Title, "<->", posts[nextPost].Title)
			}
			nextPost = i
		}

		/* render all posts */
		publishedPosts := make([]Post, 0)
		for _, post := range posts {
			/* Render single post */
			templateString := postTemplate
			if post.Metadata.Static {
				templateString = staticTemplate
			}
			postTemplate, err := template.New("post").Parse(templateString)
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

			/* skip static sites, add other to published+feed in order to list it at the start page */
			if post.Metadata.Static {
				continue
			}
			publishedPosts = append(publishedPosts, post)

			feed.Items = append(feed.Items, &feeds.Item{
				Title:   post.Title,
				Content: string(post.Rendered),
				Link: &feeds.Link{
					Href: fmt.Sprintf(
						"https://%s/%s.html",
						cfg.Domain,
						slug(post.Title),
					),
				},
				Author:  &feeds.Author{Name: cfg.Author},
				Created: time.Unix(post.Timestamp, 0),
			})
		}

		/* Put everything together */
		t, err := plainTemplate.New("blog").Parse(indexTemplate)
		if err != nil {
			return fmt.Errorf("Error parsing the html template: %s", err)
		}

		/* Save site to out dir  */
		sitePath := filepath.Join(OutPath, INDEX_FILE)
		file, err := os.Create(sitePath)
		if err != nil {
			return fmt.Errorf("Error creating index file: %s", err)
		}

		var indexVars = struct {
			Posts []Post
			Cfg   config.Config
		}{
			Posts: publishedPosts,
			Cfg:   *cfg,
		}

		if err := t.Execute(file, indexVars); err != nil {
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
		if err := os.WriteFile(rssPath, []byte(rss), 0777); err != nil {
			return fmt.Errorf("Error writing rss.xml: %w", err)
		}

		/* copy chill-files to output dir */
		for _, chillFile := range cfg.ChillFiles {
			src := chillFile
			if !filepath.IsAbs(src) {
				src = filepath.Join(BlogPath, src)
			}
			dst := filepath.Join(OutPath, filepath.Base(src))
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("copy chill-file %s to %s failed: %w", src, dst, err)
			}
			fmt.Printf("copied chill-file %s to %s\n", src, dst)
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
