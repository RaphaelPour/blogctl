# blogctl
Static markdown blog backend as a binary. Generated [my blog](https://evilcookie.de).

## Getting started

- create new folder for the blog environment: `mkdir blog`
- initialize blog: `blogctl init --path ./blog`
- create a new post:
  - interactively with `$EDITOR`: `bloctl post add -i --title="My first blog post"`
  - one-shot: `bloctl post add --path blog --title="My first blog post"`
- render html: `blogctl render --path blog -f`
- find your ready-to-serve blog in `./out`

## How it works (for contributors)

A blog is just a directory of files. Each post is its own folder (the folder
name becomes the URL slug) holding the markdown body, a metadata sidecar and any
referenced images. `blogctl render` turns that source tree into a static site.

```
  SOURCE (--path blog)                 RENDER PIPELINE                 OUTPUT (-o out)
  ────────────────────                 ───────────────                 ───────────────

  blog/
  ├── blog.json ............ site config (domain, author, title, chill-files, ...)
  ├── robots.txt ........... "chill-files": copied verbatim
  └── my-first-post/ ....... one dir per post (dir name = slug)
      ├── content.md ....... markdown body, may contain IMAGE(pic.png) shortcodes
      ├── metadata.json .... title, status (draft|public), static, createdAt, ...
      └── pic.png .......... images referenced from content.md

                       │
                       ▼   site.New(opts)              [internal/site/load.go]
            ┌──────────────────────────────────────────────┐
            │  • config.Load(blog.json)                     │
            │  • discover post dirs, load metadata.json     │
            │  • skip anything not status:"public"          │
            │  • content.md → HTML  (gomarkdown + chroma)   │
            │  • resolve IMAGE() shortcodes                 │
            │  • sort newest-first, link prev/next nav      │
            └──────────────────────────────────────────────┘
                       │
                       ▼   site.Render()  — ordered stages  [internal/site/stages.go]
            ┌──────────────────────────────────────────────┐
            │  1. renderPosts ....→ out/<slug>.html (+images)│
            │  2. renderIndex ....→ out/index.html           │
            │  3. generateRSS ....→ out/rss.xml              │
            │  4. copyAssets .....→ out/*.css (embedded)     │
            │  5. copyChillFiles .→ out/<chill-files>        │
            └──────────────────────────────────────────────┘
                       │
                       ▼
                      out/   ← ready-to-serve static site
```

The render flow lives in `internal/site`; the `cmd/` package is a thin
[cobra](https://github.com/spf13/cobra) layer that parses flags and calls into it.

```
  cmd/ ................. CLI commands (add, publish, draft, list, render, ...)
  internal/site/ ....... the render pipeline: New() loads, Render() runs stages  ← add features here
  internal/config/ ..... blog.json
  internal/metadata/ ... per-post metadata.json
  internal/highlighter/  code-block syntax highlighting (chroma)
  internal/common/ ..... file helpers, slugs, $EDITOR integration
```

**Adding a new output artifact** (sitemap, tag pages, JSON feed, ...) is a
self-contained change: implement the `Stage` interface and append it to the list
in `Site.stages()` (`internal/site/stages.go`). A stage receives the loaded
`*Site` and writes into the output directory — no changes to the loading code or
the other stages required.
