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
