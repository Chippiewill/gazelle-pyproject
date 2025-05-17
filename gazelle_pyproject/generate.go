package gazelle_pyproject

import (
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

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

	pyProjectFilename = "pyproject.toml"
)

type PyProjectToml struct {
	Project struct {
		Name         string   `toml:"name"`
		Dependencies []string `toml:"dependencies"`
	} `toml:"project"`

	Tool struct {
		Uv struct {
			Workspace struct {
				Members []string `toml:"members"`
			}
		} `toml:"uv"`
	} `toml:"tool"`
}

func parsePyProjectToml(path string) (*PyProjectToml, error) {
	var pyproject PyProjectToml
	_, err := toml.DecodeFile(path, &pyproject)
	if err != nil {
		return nil, err
	}

	return &pyproject, nil
}

func (py *PyProject) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	// cfgs := args.Config.Exts[languageName].(pythonconfig.Configs)
	// cfg := cfgs[args.Rel]

	slog.SetLogLoggerLevel(slog.LevelDebug)

	var result language.GenerateResult
	result.Gen = make([]*rule.Rule, 0)

	//slog.Debug("Generating rules for", "dir", args.Dir)

	// Find pyproject.toml in the directory tree

	if !slices.Contains(args.RegularFiles, pyProjectFilename) {
		//slog.Debug("pyproject.toml not found", "dir", args.Dir)
		return result
	}

	slog.Debug("pyproject.toml found", "dir", args.Dir)
	var pyproject, err = parsePyProjectToml(args.Dir + "/" + pyProjectFilename)
	if err != nil {
		slog.Error("error parsing pyproject.toml", "error", err)
		return result
	}

	if pyproject.Tool.Uv.Workspace.Members != nil {
		slog.Debug("uv.workspace.members found in pyproject.toml, not a leaf package")
		return result
	}

	var requirements = ExtractPackageNamesFromRequirements(pyproject.Project.Dependencies)

	var hasLibrary bool = slices.Contains(args.Subdirs, "src")
	var hasMain bool = slices.Contains(args.RegularFiles, "main.py")

	// If there's a src directory create a py_library
	if hasLibrary {
		slog.Debug("src directory found", "dir", args.Dir)
		r := rule.NewRule("py_library", pyproject.Project.Name)
		r.SetAttr("visibility", []string{"//visibility:public"})
		// list all python files in the src directory as relative paths
		srcs, _ := filepath.Glob(args.Dir + "/src/**/*.py")
		for i, src := range srcs {
			// remove the current directory from the path
			srcs[i] = strings.TrimPrefix(src, args.Dir+"/")
		}
		r.SetAttr("srcs", srcs)
		r.SetAttr("imports", []string{"src"})
		result.Gen = append(result.Gen, r)
		result.Imports = append(result.Imports, requirements)
	}

	// todo: this should cover all possible entrypoints really, not just main.py
	if hasMain {
		slog.Debug("main.py found", "dir", args.Dir)

		r := rule.NewRule("py_binary", pyproject.Project.Name)
		r.SetAttr("visibility", []string{"//visibility:public"})
		r.SetAttr("srcs", []string{"main.py"})
		r.SetAttr("main", "main.py")
		if hasLibrary {
			r.SetPrivateAttr(resolvedDepsKey, []string{":" + pyproject.Project.Name})
		}
		result.Gen = append(result.Gen, r)
		result.Imports = append(result.Imports, requirements)
	}

	return result
}
