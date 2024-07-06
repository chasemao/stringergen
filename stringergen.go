// StringerGen generates string implementations of Go structs.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/imports"
)

var (
	version = "1.0.0"
)

var (
	// source mode related
	source      = flag.String("source", "", "(source mode) Input Go source file.")
	destination = flag.String("destination", "", "(source mode) Output file; defaults to stdout, used in source mode.")

	// recusive mode related
	recursive = flag.String("recursive", "", "(recursive mode) Input directory, will handle all files recursively.")
	save      = flag.Bool("save", false, "(recursive mode) Write to file like xx_stringer.go for xx.go, used in recursive mode.")
	skipdir   = flag.String("skipdir", "", "(recursive mode) Name of directory to skip, not to generate string method within these directories; default to none.")

	// mode free flag

	exclude = flag.String("exclude", "", "Regular expression patterns for struct names to exclude from generation, separated by commas; Defaults to none.")
	method  = flag.String("method", "json", "Method for the String method generation. Supported values: json, jsoniter, fmt; Defaults to json.")

	// common flag
	debug       = flag.Bool("v", false, "Output detail information.")
	showVersion = flag.Bool("version", false, "Printf version.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	// handle common flag
	if *showVersion {
		printVersion()
		return
	}
	d.debug = *debug

	// handle mode free flagv
	excl, err := compileExcl(*exclude)
	if err != nil {
		log.Fatal("Wrong exclude regexp: ", err)
	}

	skipDirs := parseSkipDir(*skipdir)

	// handle mode
	if *source != "" {
		d.Printf(blue + "Source mode start..." + reset)
		err = genSource(*source, *destination, excl, *method)
	} else if *recursive != "" {
		d.Printf(blue + "Recursive mode start..." + reset)
		err = genRecursive(*recursive, *save, excl, *method, skipDirs)
	} else {
		usage()
		log.Fatal("You must specify source mode or recursive mode")
	}
	if err != nil {
		log.Fatalf("Generate String method failed: %v", err)
	}
}

func usage() {
	_, _ = io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `stringergen has two modes of operation: source and recursive.

Source mode generates string methods for structs from a source file.
It is enabled by using the -source flag. Other flags that
may be useful in this mode are -destination.
If destination not set, it will output to stdout.
Example:
	stringergen -source=foo.go -destination=bar.go

Recursive mode generates string methods for all structs in files in subdirectories
It is enabled by using the -recursive flag. Other flags that
may be useful in this mode are -save and -skipdir.
If save flag is set, it will output to xx_stringer.go file when handle structs in xx.go file,
otherwise it will output to stdout.
Example:
	stringergen -recursive=/path/to/directory/ -save -skipdir=/path/to/directory/skip1/,/path/to/directory/skip2/

`

func printVersion() {
	fmt.Printf("StringerGen version %s", version)
}

var d = &debuger{}

type debuger struct {
	debug bool
}

func (d *debuger) Printf(format string, args ...interface{}) {
	if d.debug {
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
}

func compileExcl(excl string) ([]*regexp.Regexp, error) {
	var exclRes []*regexp.Regexp
	if excl != "" {
		patterns := strings.Split(excl, ",")
		for _, pattern := range patterns {
			exclRe, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			}
			exclRes = append(exclRes, exclRe)
		}

		d.Printf("exclude compile success")
	}
	return exclRes, nil
}

func parseSkipDir(dir string) []string {
	return strings.Split(dir, ",")
}

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	white  = "\033[37m"
)

func genSource(source string, destination string, exclRes []*regexp.Regexp, method string) error {
	d.Printf(blue+"Handle %s start..."+reset, source)

	// if not go file, then skip
	if filepath.Ext(source) != ".go" {
		d.Printf(yellow+"NOT GO FILE: %s"+reset, source)
		return nil
	}

	// read go source file
	file, err := readGOFile(source)
	if err != nil {
		return err
	}
	d.Printf("Read Go file %s success", source)

	// parse source file, get information to generate stringerFile file
	out, err := parseFile(file, exclRes, method)
	if err != nil {
		return err
	}
	d.Printf("Parse Go file %s success get structs=%v", source, out.structNames)

	// if no struct in file, then skip
	if len(out.structNames) == 0 {
		d.Printf(yellow+"NO STRUCT IN FILE: %s"+reset, source)
		return nil
	}

	// generate stringer file
	stringerFile, err := out.gen()
	if err != nil {
		return err
	}
	d.Printf("Generate String method success for file %s", source)

	// get destination file
	dstFile, err := getDstFile(destination)
	if err != nil {
		return err
	}
	d.Printf("Going to generate destination file %s for file %s", dstFile.Name(), source)

	// write to file
	_, err = dstFile.Write(stringerFile)
	if err != nil {
		return err
	}
	d.Printf(green+"GENERATE SOURCE FILE SUCCESS: %s"+reset, source)

	return nil
}

func readGOFile(source string) (*ast.File, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, source, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed parsing source file %v: %v", source, err)
	}
	return file, nil
}

func getDstFile(destination string) (*os.File, error) {
	if destination == "" {
		return os.Stdout, nil
	}
	return os.Create(destination)
}

func parseFile(file *ast.File, exclRes []*regexp.Regexp, method string) (*output, error) {
	out := &output{
		pkg:    file.Name.Name,
		method: method,
	}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok = ts.Type.(*ast.StructType); !ok {
				continue
			}
			name := ts.Name.Name
			if !matchExcl(name, exclRes) {
				out.structNames = append(out.structNames, name)
			} else {
				d.Printf("EXCLUDE STRUCT: %s in file %s", name, *source)
			}
		}
	}
	return out, nil
}

func matchExcl(name string, exclRes []*regexp.Regexp) bool {
	for _, exclRe := range exclRes {
		if exclRe.Match([]byte(name)) {
			return true
		}
	}
	return false
}

func genRecursive(root string, save bool, exclRes []*regexp.Regexp, method string, skipDirs []string) error {
	return filepath.WalkDir(root, func(path string, de fs.DirEntry, err error) error {
		if isInSkipDirs(de, skipDirs) {
			d.Printf(yellow+"SKIP DIR: %s"+reset, de.Name())
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if de.IsDir() {
			return nil
		}
		if !save {
			return genSource(path, "", exclRes, method)
		}
		fileName := filepath.Base(path)
		ext := filepath.Ext(fileName)
		dst := path[:len(path)-len(ext)] + "_stringer" + ext
		return genSource(path, dst, exclRes, method)
	})
}

func isInSkipDirs(d fs.DirEntry, skipDirs []string) bool {
	if !d.IsDir() {
		return false
	}
	for _, skipDir := range skipDirs {
		if d.Name() == skipDir {
			return true
		}
	}
	return false
}

type output struct {
	buf         strings.Builder
	pkg         string
	structNames []string
	method      string
}

func (o *output) gen() ([]byte, error) {
	o.buf = strings.Builder{}
	switch o.method {
	case "json":
		o.genJSON()
	case "jsoniter":
		o.genJSONIter()
	case "fmt":
		o.genFmt()
	default:
		return nil, fmt.Errorf("unknown method: %s", o.method)
	}
	res := o.buf.String()
	return imports.Process("", []byte(res), nil)
}

func (o *output) addln(s string) {
	o.buf.WriteString(s)
	o.buf.WriteByte('\n')
}

func (o *output) genJSON() {
	o.addln("package " + o.pkg)
	o.addln("")
	o.addln("import (")
	o.addln(`"encoding/json"`)
	o.addln(")")
	o.addln("")
	for i, name := range o.structNames {
		if name == "" {
			continue
		}
		if i != 0 {
			o.addln("")
		}
		n := strings.ToLower(name[0:1])
		o.addln("// String Used in fmt to generate string")
		o.addln(fmt.Sprintf("func (%s *%s) String() string {", n, name))
		o.addln(fmt.Sprintf("v, _ := json.Marshal(%s)", n))
		o.addln("return string(v)")
		o.addln("}")
	}
}

func (o *output) genJSONIter() {
	o.addln("package " + o.pkg)
	o.addln("")
	o.addln("import (")
	o.addln(`jsoniter "github.com/json-iterator/go"`)
	o.addln(")")
	o.addln("")
	for i, name := range o.structNames {
		if name == "" {
			continue
		}
		if i != 0 {
			o.addln("")
		}
		n := strings.ToLower(name[0:1])
		o.addln("// String Used in fmt to generate string")
		o.addln(fmt.Sprintf("func (%s *%s) String() string {", n, name))
		o.addln(fmt.Sprintf("v, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(%s)", n))
		o.addln("return string(v)")
		o.addln("}")
	}
}

func (o *output) genFmt() {
	o.addln("package " + o.pkg)
	o.addln("")
	o.addln("import (")
	o.addln(`"fmt"`)
	o.addln(")")
	o.addln("")
	for i, name := range o.structNames {
		if name == "" {
			continue
		}
		if i != 0 {
			o.addln("")
		}
		n := strings.ToLower(name[0:1])
		o.addln("// String Used in fmt to generate string")
		o.addln(fmt.Sprintf("func (%s *%s) String() string {", n, name))
		o.addln(fmt.Sprintf(`return fmt.Sprintf("%%+v",*%s)`, n))
		o.addln("}")
	}
}
