package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Go struct {
		Rules []struct {
			Package string   `yaml:"package"`
			Forbid  []string `yaml:"forbid"`
		} `yaml:"rules"`
	} `yaml:"go"`
	Paths struct {
		FrontendRoot string `yaml:"frontend_root"`
	} `yaml:"paths"`
}

type Violation struct {
	Package string `json:"package"`
	RelDir  string `json:"relDir"`
	Import  string `json:"import"`
	Rule    string `json:"rule"`
}

type Output struct {
	Violations []Violation `json:"violations"`
	Warnings   []string    `json:"warnings"`
}

var (
	flagConfig   string
	flagJSONOut  string
	flagGoModDir string
)

func runCmd(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("%s %v failed: %w\n%s", name, args, err, out)
	}
	return out, nil
}

func main() {
	flag.StringVar(&flagConfig, "config", ".github/bots/archcheck.config.yaml", "config path")
	flag.StringVar(&flagJSONOut, "json-out", ".github/bots/out/archcheck-go.json", "json output path")
	flag.StringVar(&flagGoModDir, "gomoddir", ".", "directory that contains go.mod (backend module root)")
	flag.Parse()

	raw, err := os.ReadFile(flagConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "archcheck-go: read config: %v\n", err)
		os.Exit(1)
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "archcheck-go: parse config: %v\n", err)
		os.Exit(1)
	}

	// go list packages in module
	outJSON, err := runCmd(flagGoModDir, "go", "list", "-json", "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "archcheck-go: go list: %v\n", err)
		os.Exit(1)
	}

	type goPkg struct {
		ImportPath string
		Dir        string
		Imports    []string
		GoFiles    []string
	}
	dec := json.NewDecoder(bytes.NewReader(outJSON))
	var pkgs []goPkg
	for {
		var p goPkg
		if err := dec.Decode(&p); err != nil {
			break
		}
		if p.ImportPath != "" {
			pkgs = append(pkgs, p)
		}
	}

	var res Output

	for _, rule := range cfg.Go.Rules {
		pkgRe, err := regexp.Compile(rule.Package)
		if err != nil {
			res.Warnings = append(res.Warnings, fmt.Sprintf("bad rule package regex: %s", rule.Package))
			continue
		}
		var forbidRes []*regexp.Regexp
		for _, s := range rule.Forbid {
			r, err := regexp.Compile(s)
			if err != nil {
				res.Warnings = append(res.Warnings, fmt.Sprintf("bad forbid regex: %s", s))
				continue
			}
			forbidRes = append(forbidRes, r)
		}

		for _, p := range pkgs {
			rel, _ := filepath.Rel(flagGoModDir, p.Dir)
			rel = filepath.ToSlash(rel)
			if !pkgRe.MatchString(rel) && !pkgRe.MatchString(p.ImportPath) {
				continue
			}
			for _, imp := range p.Imports {
				for _, fr := range forbidRes {
					if fr.MatchString(imp) {
						res.Violations = append(res.Violations, Violation{
							Package: p.ImportPath,
							RelDir:  rel,
							Import:  imp,
							Rule:    rule.Package,
						})
						break
					}
				}
			}
		}
	}

	// Quick CORS star check inside backend module
	err = filepath.WalkDir(flagGoModDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		b, _ := os.ReadFile(path)
		if bytes.Contains(b, []byte(`AllowOrigins: []string{"*"}`)) {
			rel, _ := filepath.Rel(".", path)
			res.Warnings = append(res.Warnings, fmt.Sprintf("CORS AllowOrigins \"*\" found at %s", filepath.ToSlash(rel)))
		}
		return nil
	})
	_ = err

	// write JSON
	_ = os.MkdirAll(filepath.Dir(flagJSONOut), 0o755)
	j, _ := json.MarshalIndent(res, "", "  ")
	_ = os.WriteFile(flagJSONOut, j, 0o644)

	if len(res.Violations) > 0 {
		fmt.Fprintln(os.Stderr, "❌ archcheck-go: violations detected (see JSON output)")
		os.Exit(2)
	}
	fmt.Println("✅ archcheck-go: no violations")
}
