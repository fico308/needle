package visitors

import (
	"fmt"
	"reflect"

	"github.com/pingcap/parser/ast"
	"github.com/rs/zerolog/log"
)

// baseVisitor - the base of visitors, providing error handling and logging.
type baseVisitor struct {
	name   string
	errors []error

	depth    int
	traceCtx []ast.Node
}

func newBaseVisitor(name string) *baseVisitor {
	return &baseVisitor{
		name: name,
	}
}

// Enter - Implements Visitor
func (b *baseVisitor) Enter(n ast.Node) (ast.Node, bool) {
	b.traceCtx = append(b.traceCtx, n)
	b.depth++
	return n, false
}

// Leave - Implements Visitor
func (b *baseVisitor) Leave(n ast.Node) (ast.Node, bool) {
	b.traceCtx = b.traceCtx[:len(b.traceCtx)-1]
	b.depth--
	return n, true
}

// IsEnteringRoot - enter in root node
func (b *baseVisitor) IsEnteringRoot() bool {
	return b.depth == 1
}

/// IsLeavingRoot - leave root node
func (b *baseVisitor) IsLeavingRoot() bool {
	return b.depth == 0
}

// FindInCtx - Returns the closest node that has @p t.
func (b *baseVisitor) FindInCtx(t ast.Node) (ast.Node, bool) {
	for i := len(b.traceCtx) - 1; i >= 0; i-- {
		c := b.traceCtx[i]
		if typeEqual(c, t) {
			return c, true
		}
	}
	return nil, false
}

// FindInCtxAnyOf return closest node with type is any of @p types.
func (b *baseVisitor) FindInCtxAnyOf(types ...ast.Node) (ast.Node, bool) {
	for i := len(b.traceCtx) - 1; i >= 0; i-- {
		c := b.traceCtx[i]
		for _, v := range types {
			if typeEqual(c, v) {
				return c, true
			}
		}
	}
	return nil, false
}

func (b *baseVisitor) Errors() []error {
	return b.errors
}

func (b *baseVisitor) AppendErr(err Error) {
	b.errors = append(b.errors, err)
}

// compiler error
func (b *baseVisitor) LogCE(f string, args ...interface{}) {
	log.Error().Msgf(fmt.Sprintf("[%s/CompilerErr]", b.name)+f, args...)
}

// user warning
func (b *baseVisitor) LogWarn(f string, args ...interface{}) {
	log.Warn().Msgf(fmt.Sprintf("[%s/Warn]", b.name)+f, args...)
}

// user info
func (b *baseVisitor) LogInfo(f string, args ...interface{}) {
	log.Warn().Msgf(fmt.Sprintf("[%s/Info]", b.name)+f, args...)
}

func typeEqual(a, b interface{}) bool {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	return v1.Type() == v2.Type()
}