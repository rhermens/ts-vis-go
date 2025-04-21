package typescript

import (
	"os"
	"path"
	"testing"
)

func TestResolvePath(t *testing.T) {
	wd, _ := os.Getwd()
	srcDir := path.Join(wd, "testdata")
	resolver := NewResolver("./testdata/src/main.ts")

	testcases := []struct {
		name string
		path string
		workingFile string
		expected string
		shouldErr bool
	} {
		{
			name: "Relative to root",
			path: "src/imported.ts",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "Relative to root with no extension",
			path: "src/imported",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "Relative non-module path",
			path: "./imported.ts",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "Relative non-module path up a directory",
			path: "../../imported.ts",
			workingFile: path.Join(srcDir, "src/mod/p/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "No extension",
			path: "./imported",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "Absolute path",
			path: path.Join(srcDir, "src/imported.ts"),
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: path.Join(srcDir, "src/imported.ts"),
			shouldErr: false,
		},
		{
			name: "Non-existent path",
			path: "./nonexistent.ts",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: "",
			shouldErr: true,
		},
		{
			name: "Module",
			path: "cache-manager",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: "",
			shouldErr: true,
		},
		{
			name: "Module in dev deps",
			path: "@swc/core",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: "",
			shouldErr: true,
		},
		{
			name: "Module sub path",
			path: "cache-manager/lib",
			workingFile: path.Join(srcDir, "src/main.ts"),
			expected: "",
			shouldErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resolved, err := resolver.ResolvePath(tc.path, tc.workingFile)
			if resolved != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, resolved)
			}

			if tc.shouldErr && err == nil {
				t.Errorf("expected error, didnt error")
			}
		})
	}
}
