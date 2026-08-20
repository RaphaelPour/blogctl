package site

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/RaphaelPour/blogctl/internal/common"
	"github.com/RaphaelPour/blogctl/internal/config"

	"github.com/gorilla/feeds"
)

const rssFile = "rss.xml"

// Stage is one step of the render pipeline. Stages run in order and share the
// Site, so a later stage can rely on state produced by an earlier one (e.g.
// renderPostsStage populates Published and the feed before renderIndexStage
// and generateRSSStage run). Add a new output artifact by implementing Stage
// and appending it in stages().
type Stage interface {
	Name() string
	Run(s *Site) error
}

// stages returns the ordered pipeline of build steps.
func (s *Site) stages() []Stage {
	return []Stage{
		renderPostsStage{},
		renderIndexStage{},
		generateRSSStage{},
		copyAssetsStage{},
		copyChillFilesStage{},
	}
}

// renderPostsStage writes one HTML file per post, copies referenced images and
// accumulates the published (non-static) posts into Published and the feed.
type renderPostsStage struct{}

func (renderPostsStage) Name() string { return "render posts" }

func (renderPostsStage) Run(s *Site) error {
	published := make([]Post, 0)
	for _, post := range s.Posts {
		/* copy referenced images next to the rendered post */
		for _, img := range post.images {
			if err := common.CopyFile(img.src, img.dst); err != nil {
				return fmt.Errorf("error copying '%s' to '%s': %w", img.src, img.dst, err)
			}
		}

		/* Render single post */
		tmpl := postTmpl
		if post.Metadata.Static {
			tmpl = staticTmpl
		}

		postFilePath := filepath.Join(s.Options.OutPath, post.Link)
		file, err := os.Create(postFilePath)
		if err != nil {
			return fmt.Errorf("Error creating post file '%s': %s", post.Title, err)
		}

		if err := tmpl.Execute(file, post); err != nil {
			return fmt.Errorf("Error rendering post '%s': %s", post.Title, err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("Error closing post file '%s': %s", post.Title, err)
		}

		/* skip static sites, add others to published+feed to list them on the start page */
		if post.Metadata.Static {
			continue
		}
		published = append(published, post)

		s.Feed.Items = append(s.Feed.Items, &feeds.Item{
			Title:   post.Title,
			Content: string(post.Rendered),
			Link: &feeds.Link{
				Href: fmt.Sprintf(
					"https://%s/%s.html",
					s.Config.Domain,
					common.Slug(post.Title),
				),
			},
			Author:  &feeds.Author{Name: s.Config.Author},
			Created: time.Unix(post.Timestamp, 0),
		})
	}

	s.Published = published
	return nil
}

// renderIndexStage renders the start page listing all published posts.
type renderIndexStage struct{}

func (renderIndexStage) Name() string { return "render index" }

func (renderIndexStage) Run(s *Site) error {
	sitePath := filepath.Join(s.Options.OutPath, indexFile)
	file, err := os.Create(sitePath)
	if err != nil {
		return fmt.Errorf("Error creating index file: %s", err)
	}

	indexVars := struct {
		Posts []Post
		Cfg   config.Config
	}{
		Posts: s.Published,
		Cfg:   *s.Config,
	}

	if err := indexTmpl.Execute(file, indexVars); err != nil {
		return fmt.Errorf("Error rendering posts: %s", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("Error closing file: %s", err)
	}

	return nil
}

// generateRSSStage writes the RSS feed accumulated by renderPostsStage.
type generateRSSStage struct{}

func (generateRSSStage) Name() string { return "generate rss" }

func (generateRSSStage) Run(s *Site) error {
	rss, err := s.Feed.ToRss()
	if err != nil {
		return fmt.Errorf("Error generating rss feed: %w", err)
	}

	rssPath := filepath.Join(s.Options.OutPath, rssFile)
	if err := os.WriteFile(rssPath, []byte(rss), 0777); err != nil {
		return fmt.Errorf("Error writing rss.xml: %w", err)
	}

	return nil
}

// copyAssetsStage copies the default theme's static assets (CSS, ...) from the
// embedded filesystem into the output directory, skipping templates.
type copyAssetsStage struct{}

func (copyAssetsStage) Name() string { return "copy assets" }

func (copyAssetsStage) Run(s *Site) error {
	entries, err := publicFS.ReadDir(publicDir)
	if err != nil {
		return fmt.Errorf("Error reading embedded assets: %w", err)
	}

	for _, file := range entries {
		if strings.Contains(file.Name(), "tmpl") {
			continue
		}

		inPath := path.Join(publicDir, file.Name())
		f, err := publicFS.ReadFile(inPath)
		if err != nil {
			return fmt.Errorf("Error reading %s", inPath)
		}

		outPath := filepath.Join(s.Options.OutPath, file.Name())
		if err := os.WriteFile(outPath, f, 0777); err != nil {
			return fmt.Errorf("Error writing %s: %w", file.Name(), err)
		}
	}

	return nil
}

// copyChillFilesStage copies the extra files listed in the config into the
// output directory.
type copyChillFilesStage struct{}

func (copyChillFilesStage) Name() string { return "copy chill-files" }

func (copyChillFilesStage) Run(s *Site) error {
	for _, chillFile := range s.Config.ChillFiles {
		src := chillFile
		if !filepath.IsAbs(src) {
			src = filepath.Join(s.Options.BlogPath, src)
		}
		dst := filepath.Join(s.Options.OutPath, filepath.Base(src))
		if err := common.CopyFile(src, dst); err != nil {
			return fmt.Errorf("copy chill-file %s to %s failed: %w", src, dst, err)
		}
		s.logf("copied chill-file %s to %s\n", src, dst)
	}

	return nil
}
