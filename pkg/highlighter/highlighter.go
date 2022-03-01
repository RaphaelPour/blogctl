package highlighter

import (
	"io"

	"github.com/d4l3k/go-highlight"
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

	// get the raw code from the node
	rawCode := string(node.AsLeaf().Literal)
	highlightedCode, _ := highlight.HTML(codeLanguage, []byte(rawCode))

	// write the highlighted code
	_, err := w.Write(highlightedCode)
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
