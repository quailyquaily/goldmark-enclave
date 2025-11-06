package enclave

import (
	"reflect"
	"testing"

	"github.com/quailyquaily/goldmark-enclave/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Test that an image transformed into an Enclave node preserves its
// original position within a paragraph (i.e., replaced in-place).
func TestTransformPreservesNodeOrderSimple(t *testing.T) {
	markdown := []byte("![👍](tg://emoji?id=123)\n`inline code`")

	md := goldmark.New(
		goldmark.WithExtensions(New(&core.Config{})),
	)

	reader := text.NewReader(markdown)
	doc := md.Parser().Parse(reader)

	found := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if para, ok := n.(*ast.Paragraph); ok {
			found = true
			// Check first child is Enclave (replaced in-place)
			first := para.FirstChild()
			if _, ok := first.(*core.Enclave); !ok {
				t.Fatalf("first child = %T, want *core.Enclave", first)
			}

			// Also ensure the CodeSpan still appears after the Enclave
			idxEnclave := -1
			idxCode := -1
			i := 0
			for child := para.FirstChild(); child != nil; child = child.NextSibling() {
				if _, ok := child.(*core.Enclave); ok && idxEnclave == -1 {
					idxEnclave = i
				}
				if _, ok := child.(*ast.CodeSpan); ok && idxCode == -1 {
					idxCode = i
				}
				i++
			}
			if idxEnclave == -1 {
				t.Fatalf("did not find *core.Enclave among children")
			}
			if idxCode == -1 {
				t.Fatalf("did not find *ast.CodeSpan among children")
			}
			if !(idxEnclave < idxCode) {
				t.Fatalf("unexpected order: enclave idx=%d, code idx=%d (want enclave before code)", idxEnclave, idxCode)
			}
		}
		return ast.WalkContinue, nil
	})

	if !found {
		t.Fatalf("no paragraph found in AST; types=%v", reflect.TypeOf(doc))
	}
}

// Image appears in the middle of a paragraph; order should be: Text, Enclave, Text, CodeSpan
func TestTransformPreservesNodeOrder_MiddleOfParagraph(t *testing.T) {
	markdown := []byte("Hello ![👍](tg://emoji?id=123) world `code`")

	md := goldmark.New(
		goldmark.WithExtensions(New(&core.Config{})),
	)

	reader := text.NewReader(markdown)
	doc := md.Parser().Parse(reader)

	found := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if para, ok := n.(*ast.Paragraph); ok {
			found = true
			// Expect: Text, Enclave, Text, CodeSpan
			expectedKinds := []reflect.Type{reflect.TypeOf(&ast.Text{}), reflect.TypeOf(&core.Enclave{}), reflect.TypeOf(&ast.Text{}), reflect.TypeOf(&ast.CodeSpan{})}
			var kinds []reflect.Type
			for child := para.FirstChild(); child != nil; child = child.NextSibling() {
				kinds = append(kinds, reflect.TypeOf(child))
			}
			if len(kinds) != len(expectedKinds) {
				t.Fatalf("unexpected children count: got %d, want %d", len(kinds), len(expectedKinds))
			}
			for i := range kinds {
				if kinds[i] != expectedKinds[i] {
					t.Fatalf("child[%d] kind=%v, want %v", i, kinds[i], expectedKinds[i])
				}
			}
		}
		return ast.WalkContinue, nil
	})

	if !found {
		t.Fatalf("no paragraph found in AST; types=%v", reflect.TypeOf(doc))
	}
}

// Multiple images in a single paragraph should each be replaced in-place
// and keep order: Text, Enclave, Text, Enclave, Text
func TestTransformPreservesNodeOrder_MultipleImages(t *testing.T) {
	markdown := []byte("A ![1](tg://emoji?id=1) B ![2](tg://emoji?id=2) C")

	md := goldmark.New(
		goldmark.WithExtensions(New(&core.Config{})),
	)

	reader := text.NewReader(markdown)
	doc := md.Parser().Parse(reader)

	found := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if para, ok := n.(*ast.Paragraph); ok {
			found = true
			expectedKinds := []reflect.Type{reflect.TypeOf(&ast.Text{}), reflect.TypeOf(&core.Enclave{}), reflect.TypeOf(&ast.Text{}), reflect.TypeOf(&core.Enclave{}), reflect.TypeOf(&ast.Text{})}
			var kinds []reflect.Type
			for child := para.FirstChild(); child != nil; child = child.NextSibling() {
				kinds = append(kinds, reflect.TypeOf(child))
			}
			if len(kinds) != len(expectedKinds) {
				t.Fatalf("unexpected children count: got %d, want %d", len(kinds), len(expectedKinds))
			}
			for i := range kinds {
				if kinds[i] != expectedKinds[i] {
					t.Fatalf("child[%d] kind=%v, want %v", i, kinds[i], expectedKinds[i])
				}
			}
		}
		return ast.WalkContinue, nil
	})

	if !found {
		t.Fatalf("no paragraph found in AST; types=%v", reflect.TypeOf(doc))
	}
}

// Image nested inside a link should be replaced within the link, preserving
// paragraph order: Text, Link(with Enclave child), Text
func TestTransformPreservesNodeOrder_ImageInsideLink(t *testing.T) {
	markdown := []byte("A [![1](tg://emoji?id=1)](https://example.com) Z")

	md := goldmark.New(
		goldmark.WithExtensions(New(&core.Config{})),
	)

	reader := text.NewReader(markdown)
	doc := md.Parser().Parse(reader)

	found := false
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if para, ok := n.(*ast.Paragraph); ok {
			found = true
			// Expect: Text, Link, Text
			expectedKinds := []reflect.Type{reflect.TypeOf(&ast.Text{}), reflect.TypeOf(&ast.Link{}), reflect.TypeOf(&ast.Text{})}
			var kinds []reflect.Type
			var linkNode *ast.Link
			for child := para.FirstChild(); child != nil; child = child.NextSibling() {
				kinds = append(kinds, reflect.TypeOf(child))
				if l, ok := child.(*ast.Link); ok {
					linkNode = l
				}
			}
			if len(kinds) != len(expectedKinds) {
				t.Fatalf("unexpected children count: got %d, want %d", len(kinds), len(expectedKinds))
			}
			for i := range kinds {
				if kinds[i] != expectedKinds[i] {
					t.Fatalf("child[%d] kind=%v, want %v", i, kinds[i], expectedKinds[i])
				}
			}
			if linkNode == nil {
				t.Fatalf("expected a link node in paragraph children")
			}
			if _, ok := linkNode.FirstChild().(*core.Enclave); !ok {
				t.Fatalf("link's first child = %T, want *core.Enclave", linkNode.FirstChild())
			}
		}
		return ast.WalkContinue, nil
	})

	if !found {
		t.Fatalf("no paragraph found in AST; types=%v", reflect.TypeOf(doc))
	}
}
