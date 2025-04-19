package typescript

import (
	"encoding/json"
	"errors"
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
}

type Resolver struct {
	cfg tsConfig
	cwd string
	modules []string
}

func (r *Resolver) ResolvePath(p string, origin string) (string, error) {
	if slices.Contains(r.modules, p) {
		return "", errors.New("Is module")
	}

	if !strings.HasSuffix(p, ".ts") {
		p = p + ".ts"
	}

	if path.IsAbs(p) {
		return p, nil
	}

	if strings.HasPrefix(p, ".") {
		return path.Join(path.Dir(origin), p), nil
	}

	return path.Join(r.cwd, p), nil
}

func NewResolver(entrypoint string) *Resolver {
	cwd, err := resolveCwd(entrypoint)
	if err != nil {
		panic(err)
	}

	deps, err := readDeps(cwd)
	tsconfig, err := readTsConfig(cwd)

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
	cwd, _ := path.Split(path.Clean(abs))

	if _, err := os.Stat(path.Join(cwd, "tsconfig.json")); err == nil {
		return cwd, nil
	}

	if cwd == "/" {
		return "", errors.New("could not find tsconfig.json")
	}

	return resolveCwd(cwd)
}
