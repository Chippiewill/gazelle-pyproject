package gazelle_pyproject

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	pyLibraryEntrypointFilename = "__init__.py"
	pyBinaryEntrypointFilename  = "__main__.py"
	pyTestEntrypointFilename    = "__test__.py"
	pyTestEntrypointTargetname  = "__test__"
	conftestFilename            = "conftest.py"
	conftestTargetname          = "conftest"
)

var (
	buildFilenames = []string{"BUILD", "BUILD.bazel"}
)

type PyProjectToml struct {
	Project struct {
		Name string `toml:"name"`
	} `toml:"project"`

	Tool struct {
		Uv struct {
			Workspace struct {
				Members []string `toml:"members"`
			}
		} `toml:"uv"`
	} `toml:"tool"`
}

func (py *PyProject) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	// cfgs := args.Config.Exts[languageName].(pythonconfig.Configs)
	// cfg := cfgs[args.Rel]

	var result language.GenerateResult
	result.Gen = make([]*rule.Rule, 0)

	for _, f := range args.RegularFiles {
		if f == "pyproject.toml" {
			log.Printf("pyproject.toml found: %s, %s", f, args.Dir)

			var pyproject PyProjectToml
			_, err := toml.DecodeFile(args.Dir+"/"+f, &pyproject)
			if err != nil {
				log.Printf("error decoding pyproject.toml: %s", err)
			}

			if pyproject.Tool.Uv.Workspace.Members != nil {
				return result
			}

			r := rule.NewRule("py_library", pyproject.Project.Name)

			result.Gen = append(result.Gen, r)
			result.Imports = append(result.Imports, 42)
		}
	}

	return result
}
