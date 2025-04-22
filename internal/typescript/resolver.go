package typescript

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type tsConfig struct {
	CompilerOptions map[string]string `json:"compilerOptions"`
}

type packageFile struct {
	Dependencies map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Resolver struct {
	cfg tsConfig
	cwd string
	modules []string
}

func (r *Resolver) ResolvePath(p string, origin string) (string, error) {
	if !path.IsAbs(origin) {
		return "", errors.New("origin is not absolute")
	}

	if r.isModule(p) {
		return "", errors.New(fmt.Sprintf("%s is a module", p))
	}

	if strings.HasPrefix(p, ".") {
		return r.resolveRelative(p, origin) 
	}

	if path.IsAbs(p) {
		return r.resolveRelative(p, "/") 
	}

	return r.resolveRelative(p, filepath.Join(r.cwd, "tsconfig.json"))
}

func (r *Resolver) resolveRelative(p string, origin string) (string, error) {
	if !strings.HasSuffix(p, ".ts") {
		p = p + ".ts"
	}

	target := path.Join(path.Dir(origin), p)

	if _, err := os.Stat(target); err != nil {
		return "", err
	}

	return target, nil
}

func (r *Resolver) isModule(p string) bool {
	if slices.Contains(r.modules, p) {
		slog.Debug("is module", "path", p)
		return true
	}

	for _, module := range r.modules {
		if strings.HasPrefix(p, module) {
			slog.Debug("is module prefix", "path", p, "module", module)
			return true
		}
	}

	return false
}

func NewResolver(entrypoint string) *Resolver {
	cwd, err := resolveCwd(entrypoint)
	if err != nil {
		panic(err)
	}

	deps, _ := readDeps(cwd)
	tsconfig, _ := readTsConfig(cwd)

	return &Resolver{
		cwd: cwd,
		modules: deps,
		cfg: tsconfig,
	}
}

func readDeps(cwd string) ([]string, error) {
	file, err := os.ReadFile(path.Join(cwd, "package.json"));
	if err != nil {
		return nil, err
	}

	var pkg packageFile
	json.Unmarshal(file, &pkg)

	deps := make([]string, 0, len(pkg.Dependencies))
	for k := range pkg.Dependencies {
		deps = append(deps, k)
	}
	for k := range pkg.DevDependencies {
		deps = append(deps, k)
	}

	return deps, nil
}

func readTsConfig(cwd string) (tsConfig, error) {
	file, err := os.ReadFile(path.Join(cwd, "tsconfig.json"));
	if err != nil {
		return tsConfig{}, err
	}

	var tsconfig tsConfig
	json.Unmarshal(file, &tsconfig)

	return tsconfig, nil
}

func resolveCwd(entrypoint string) (string, error) {
	abs, _ := filepath.Abs(entrypoint)
	if _, ok := os.Stat(abs); ok != nil {
		return "", errors.New("could not find entrypoint")
	}

	cwd, _ := path.Split(path.Clean(abs))

	if _, err := os.Stat(path.Join(cwd, "tsconfig.json")); err == nil {
		return cwd, nil
	}

	if cwd == "/" {
		return "", errors.New("could not find tsconfig.json")
	}

	return resolveCwd(cwd)
}
