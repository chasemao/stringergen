package main

import (
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileExcl(t *testing.T) {
	tests := []struct {
		name    string
		excl    string
		want    []*regexp.Regexp
		wantErr bool
	}{
		{
			name:    "Valid single pattern",
			excl:    "foo",
			want:    []*regexp.Regexp{regexp.MustCompile("foo")},
			wantErr: false,
		},
		{
			name:    "Valid multiple patterns",
			excl:    "foo,bar",
			want:    []*regexp.Regexp{regexp.MustCompile("foo"), regexp.MustCompile("bar")},
			wantErr: false,
		},
		{
			name:    "Invalid pattern",
			excl:    "foo,(",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty input",
			excl:    "",
			want:    []*regexp.Regexp{},
			wantErr: false,
		},
		{
			name:    "Whitespace input",
			excl:    " ",
			want:    []*regexp.Regexp{regexp.MustCompile(" ")},
			wantErr: false,
		},
		{
			name:    "Special characters pattern",
			excl:    "^foo$,\\d+",
			want:    []*regexp.Regexp{regexp.MustCompile("^foo$"), regexp.MustCompile(`\d+`)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compileExcl(tt.excl)
			if (err != nil) != tt.wantErr {
				t.Errorf("compileExcl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("compileExcl() = %v, want %v", got, tt.want)
				return
			}
			for i, re := range got {
				if re.String() != tt.want[i].String() {
					t.Errorf("compileExcl() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestMatchExcl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		exclRes  []*regexp.Regexp
		expected bool
	}{
		{
			name:     "Match single pattern",
			input:    "foo",
			exclRes:  []*regexp.Regexp{regexp.MustCompile("foo")},
			expected: true,
		},
		{
			name:     "No match single pattern",
			input:    "bar",
			exclRes:  []*regexp.Regexp{regexp.MustCompile("foo")},
			expected: false,
		},
		{
			name:     "Match multiple patterns",
			input:    "bar",
			exclRes:  []*regexp.Regexp{regexp.MustCompile("foo"), regexp.MustCompile("bar")},
			expected: true,
		},
		{
			name:     "No match multiple patterns",
			input:    "baz",
			exclRes:  []*regexp.Regexp{regexp.MustCompile("foo"), regexp.MustCompile("bar")},
			expected: false,
		},
		{
			name:     "Empty exclusion list",
			input:    "foo",
			exclRes:  []*regexp.Regexp{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchExcl(tt.input, tt.exclRes)
			if result != tt.expected {
				t.Errorf("matchExcl() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Mock implementation of fs.DirEntry for testing purposes
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func TestIsInSkipDirs(t *testing.T) {
	tests := []struct {
		name     string
		dirEntry fs.DirEntry
		skipDirs []string
		expected bool
	}{
		{
			name:     "Directory in skip list",
			dirEntry: mockDirEntry{name: "foo", isDir: true},
			skipDirs: []string{"foo", "bar"},
			expected: true,
		},
		{
			name:     "Directory not in skip list",
			dirEntry: mockDirEntry{name: "baz", isDir: true},
			skipDirs: []string{"foo", "bar"},
			expected: false,
		},
		{
			name:     "Entry is not a directory",
			dirEntry: mockDirEntry{name: "foo", isDir: false},
			skipDirs: []string{"foo", "bar"},
			expected: false,
		},
		{
			name:     "Empty skip list",
			dirEntry: mockDirEntry{name: "foo", isDir: true},
			skipDirs: []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInSkipDirs(tt.dirEntry, tt.skipDirs)
			if result != tt.expected {
				t.Errorf("isInSkipDirs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseSkipDir(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want []string
	}{
		{
			name: "Single directory",
			dir:  "foo",
			want: []string{"foo"},
		},
		{
			name: "Multiple directories",
			dir:  "foo,bar",
			want: []string{"foo", "bar"},
		},
		{
			name: "Empty string",
			dir:  "",
			want: []string{""},
		},
		{
			name: "String with spaces",
			dir:  "foo, bar",
			want: []string{"foo", " bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSkipDir(tt.dir); !equal(got, tt.want) {
				t.Errorf("parseSkipDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices of strings
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParserFile(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		excl     string
		expected *output
		wantErr  bool
	}{
		{
			name: "Exclude None",
			src: `
package main

type MyStruct struct {
	Field1 string
	Field2 int
}
`,
			excl: "",
			expected: &output{
				pkg:         "main",
				structNames: []string{"MyStruct"},
			},
			wantErr: false,
		},
		{
			name: "Exclude one",
			src: `
package main

type MyStruct1 struct {
	Field1 string
}

type MyStruct2 struct {
	Field2 int
}
`,
			excl: "MyStruct2",
			expected: &output{
				pkg:         "main",
				structNames: []string{"MyStruct1"},
			},
			wantErr: false,
		},
		{
			name: "Exclude two",
			src: `
package main

type MyStruct1 struct {
	Field1 string
}

type MyStruct2 struct {
	Field2 int
}
`,
			excl: "MyStruct1,MyStruct2",
			expected: &output{
				pkg:         "main",
				structNames: []string{},
			},
			wantErr: false,
		},
		{
			name: "Exclude all",
			src: `
package main

type MyStruct struct {
	Field1 string
}
`,
			excl: ".*",
			expected: &output{
				pkg:         "main",
				structNames: []string{},
			},
			wantErr: false,
		},
		{
			name: "No structs",
			src: `
package main

var x int
`,
			excl: ".*",
			expected: &output{
				pkg:         "main",
				structNames: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the source code into an *ast.File
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.src, parser.AllErrors)
			if err != nil {
				t.Fatalf("parser.ParseFile() error: %v", err)
			}

			// Compile the regular expression
			excl, err := compileExcl(tt.excl)
			if err != nil {
				t.Fatalf("regexp.Compile() error: %v", err)
			}

			// Call the parseFile function
			got, err := parseFile(file, excl, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare the output
			if got.pkg != tt.expected.pkg {
				t.Errorf("parseFile() pkg = %v, want %v", got.pkg, tt.expected.pkg)
			}
			if len(got.structNames) != len(tt.expected.structNames) {
				t.Errorf("parseFile() structNames length = %v, want %v", len(got.structNames), len(tt.expected.structNames))
			}
			for i, name := range got.structNames {
				if name != tt.expected.structNames[i] {
					t.Errorf("parseFile() structNames[%d] = %v, want %v", i, name, tt.expected.structNames[i])
				}
			}
		})
	}
}

func TestOutputGen(t *testing.T) {
	o := &output{}
	_, err := o.gen()
	assert.Error(t, err)
}

func TestGenJSON(t *testing.T) {
	o := &output{
		pkg:         "main",
		structNames: []string{"MyStruct"},
	}

	o.genJSON()

	expected := `package main

import (
"encoding/json"
)

// String Used in fmt to generate string
func (m *MyStruct) String() string {
v, _ := json.Marshal(m)
return string(v)
}
`
	assert.Equal(t, expected, o.buf.String())
}

func TestGenJSONIter(t *testing.T) {
	o := &output{
		pkg:         "main",
		structNames: []string{"MyStruct"},
	}

	o.genJSONIter()

	expected := `package main

import (
jsoniter "github.com/json-iterator/go"
)

// String Used in fmt to generate string
func (m *MyStruct) String() string {
v, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(m)
return string(v)
}
`
	assert.Equal(t, expected, o.buf.String())
}

func TestGenFmt(t *testing.T) {
	o := &output{
		pkg:         "main",
		structNames: []string{"MyStruct"},
	}

	o.genFmt()

	expected := `package main

import (
"fmt"
)

// String Used in fmt to generate string
func (m *MyStruct) String() string {
return fmt.Sprintf("%+v",*m)
}
`
	assert.Equal(t, expected, o.buf.String())
}
