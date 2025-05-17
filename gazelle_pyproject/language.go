package gazelle_pyproject

import (
	"github.com/bazelbuild/bazel-gazelle/language"
)

type PyProject struct {
	Configurer
	Resolver
}

func NewLanguage() language.Language {
	return &PyProject{}
}
