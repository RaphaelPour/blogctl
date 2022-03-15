package highlighter

import (
	"bytes"
	"io"

	htmlFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// skip if the current node is not a code-block
	if _, ok := node.(*ast.CodeBlock); !ok {
		return ast.GoToNext, false
	}

	// get the language of the code from the top of the code-block
	codeLanguage := string(node.(*ast.CodeBlock).Info)
	lexer := lexers.Get(codeLanguage)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	formatter := htmlFormatter.New(htmlFormatter.WithClasses(true))

	// get the raw code from the node
	rawCode := string(node.AsLeaf().Literal)
	iterator, err := lexer.Tokenise(nil, rawCode)
	if err != nil {
		return ast.Terminate, false
	}

	err = formatter.Format(w, styles.AlgolNu, iterator)
	if err != nil {
		return ast.Terminate, false
	}

	return ast.GoToNext, true
}

// get a html-renderer that adds css classes to different parts
// of a code-block
func GetRenderer() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderHook,
	}

	return html.NewRenderer(opts)
}

func GetStyle() (string, error) {
	formatter := htmlFormatter.New(htmlFormatter.WithClasses(true))
	style := styles.Get("algol_nu")
	if style == nil {
		style = styles.Fallback
	}
	var buf bytes.Buffer
	err := formatter.WriteCSS(&buf, style)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
