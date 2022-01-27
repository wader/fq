package ksexpr_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/wader/fq/format/kaitai/ksexpr"
	"gopkg.in/yaml.v3"
)

// TODO: const_string2 escape test

type kst struct {
	ID      string `yaml:"id"`
	Data    string `yaml:"data"`
	Asserts []struct {
		Actual   string `yaml:"actual"` // actual seems to be instance name?
		Expected string `yaml:"expected"`
	} `yaml:"asserts"`
}

// TODO: minimal just for tests, move?
type ksy struct {
	Meta struct {
		ID string `yaml:"id"`
	} `yaml:"meta"`
	Instances map[string]struct {
		Value string `yaml:"value"`
	} `yaml:"instances"`
}

type testEnv struct{}

func lookup(input any, ns string, name string) (any, error) {
	switch ns {
	case "":
		switch i := input.(type) {
		case map[string]any:
			v := i[name]
			return v, nil
		}
	default:
		return "with_namespace", nil
	}

	return nil, fmt.Errorf("failed to lookup ident %s", name)
}

func decodeYAML(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(v); err != nil {
		return err
	}

	return nil
}

func testKSY(t *testing.T, path string) {
	var ksy ksy
	var kst kst

	if err := decodeYAML(path+".ksy", &ksy); err != nil {
		t.Fatal(err)
	}
	if err := decodeYAML(path+".kst", &kst); err != nil {
		t.Fatal(err)
	}

	expectedMap := map[string]string{}
	for _, ka := range kst.Asserts {
		expectedMap[ka.Actual] = ka.Expected
	}

	// sort to keep stable
	var instanceNames []string
	for n := range ksy.Instances {
		instanceNames = append(instanceNames, n)
	}
	sort.Strings(instanceNames)

	for _, n := range instanceNames {
		ki := ksy.Instances[n]
		expectedStr := expectedMap[n]

		t.Run(n, func(t *testing.T) {
			expectedExpr, err := ksexpr.Parse(expectedStr)
			if err != nil {
				t.Fatalf("Expected parse error: %s: %q", err, expectedStr)
			}
			expected, err := expectedExpr.Eval(0)
			if err != nil {
				t.Fatalf("Expected eval error: %s", err)
			}

			valueExpr, err := ksexpr.Parse(ki.Value)
			if err != nil {
				t.Fatalf("Value parse error: %s: %q", err, ki.Value)
			}

			actual, err := valueExpr.Eval(
				// the env ksdump will have
				ksexpr.ToValue(map[string]any{
					"zeros": "\x00\x00\x00",
					"_root": map[string]any{
						"zeros": "\x00\x00\x00",
					},
				}))
			if err != nil {
				t.Fatalf("Value eval error %s", err)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected %v, got %s", expected, actual)
			}
		})

	}
}

func TestKSY(t *testing.T) {
	_ = filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".ksy" {
			return nil
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			testKSY(t, path[0:len(path)-4])
		})

		return nil
	})
}
