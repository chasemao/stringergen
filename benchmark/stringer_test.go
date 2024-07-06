package benchmark

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Default = func() *zap.SugaredLogger {
	// Configure lumberjack logger
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "output",
		MaxSize:    1, // Megabytes
		MaxBackups: 1,
	}

	// Create a zapcore.WriteSyncer using the lumberjack logger
	writeSyncer := zapcore.AddSync(lumberjackLogger)

	// Configure zap core with the lumberjack WriteSyncer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writeSyncer,
		zap.InfoLevel,
	)

	// Create a new zap logger with the core
	logger := zap.New(core)

	return logger.Sugar()
}()

type simple struct {
	F1 int
	F2 string
	F3 float64
}

func getSimple() *simple {
	return &simple{
		F1: 1221321321,
		F2: "asldhfaishdfj02318u40123",
		F3: 12312312.12321,
	}
}

func BenchmarkSimpleGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getSimple()
	}
}

//func BenchmarkSimpleSprintf(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		simpleItem := getSimple()
//		_ = fmt.Sprintf("%v", simpleItem)
//	}
//}

func BenchmarkSimpleSprintfWithPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleItem := getSimple()
		_ = fmt.Sprintf("%+v", simpleItem)
	}
}

func BenchmarkSimpleSpew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleItem := getSimple()
		_ = spew.Sprintf("%v", simpleItem)
	}
}

//func BenchmarkSimpleSpewPlus(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		simpleItem := getSimple()
//		_ = spew.Sprintf("%+v", simpleItem)
//	}
//}

func BenchmarkSimpleJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleItem := getSimple()
		v, _ := json.Marshal(simpleItem)
		_ = fmt.Sprintf("%v", string(v))
	}
}

func BenchmarkSimpleJsonIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleItem := getSimple()
		v, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(simpleItem)
		_ = string(v)
	}
}

type simpleCustomJson struct {
	F1 int
	F2 string
	F3 float64
}

func getSimpleCustomJson() *simpleCustomJson {
	return &simpleCustomJson{
		F1: 1221321321,
		F2: "asldhfaishdfj02318u40123",
		F3: 12312312.12321,
	}
}

func (s *simpleCustomJson) String() string {
	return fmt.Sprintf(`{"F1": %d, "F2": "%s", "F3": %f}`, s.F1, s.F2, s.F3)
}

func BenchmarkSimpleCustomJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleCustomJsonItem := getSimpleCustomJson()
		_ = fmt.Sprintf("%v", simpleCustomJsonItem)
	}
}

type simpleCustomJsonStringBuilder struct {
	F1 int
	F2 string
	F3 float64
}

func getSimpleCustomJsonStringBuilder() *simpleCustomJsonStringBuilder {
	return &simpleCustomJsonStringBuilder{
		F1: 1221321321,
		F2: "asldhfaishdfj02318u40123",
		F3: 12312312.12321,
	}
}

func (s *simpleCustomJsonStringBuilder) String() string {
	b := &strings.Builder{}
	b.Write([]byte(`{"F1": `))
	b.Write([]byte(strconv.Itoa(s.F1)))
	b.Write([]byte(`, "F2": "`))
	b.Write([]byte(s.F2))
	b.Write([]byte(`", "F3":`))
	b.Write([]byte(strconv.FormatFloat(s.F3, 'f', -1, 64)))
	b.Write([]byte(`}`))
	return b.String()
}

func BenchmarkSimpleCustomJsonStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleCustomJsonStringBuilderItem := getSimpleCustomJsonStringBuilder()
		_ = fmt.Sprintf("%v", simpleCustomJsonStringBuilderItem)
	}
}

type simpleStringer struct {
	F1 int
	F2 string
	F3 float64
}

func getSimpleStringer() *simpleStringer {
	return &simpleStringer{
		F1: 1221321321,
		F2: "asldhfaishdfj02318u40123",
		F3: 12312312.12321,
	}
}

func (s *simpleStringer) String() string {
	return fmt.Sprint(*s)
}

func BenchmarkSimpleStringer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simpleStringerItem := getSimpleStringer()
		_ = fmt.Sprintf("%v", simpleStringerItem)
	}
}

//func BenchmarkSimpleStringerPlus(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		simpleStringerItem := getSimpleStringer()
//		_ = fmt.Sprintf("%+v", simpleStringerItem)
//	}
//}

type complicated struct {
	F1 int
	F2 string
	F3 float64
	F4 *subComplicated
	F5 []*subComplicated
	F6 map[string]*subComplicated
	F7 map[string][]*subComplicated
}

type subComplicated struct {
	F1 int
	F2 string
	F3 float64
	S1 *subSubComplicated
}

type subSubComplicated struct {
	F1 int
	F2 string
	F3 float64
}

func genSubComplicated(s string) *subComplicated {
	return &subComplicated{
		F1: 14701271122,
		F2: s,
		F3: 102740127.12124211,
		S1: &subSubComplicated{
			F1: 14701271122,
			F2: s,
			F3: 102740127.12124211,
		},
	}
}

func getComplicated() *complicated {
	return &complicated{
		F1: 174012374,
		F2: "aoijhdfaoidsofahoshfoiasfdhjoaijsdifjaosjdfaijosdjf",
		F3: 12312.24124,
		F4: genSubComplicated("sdfadf"),
		F5: []*subComplicated{
			genSubComplicated("alsjfl"), genSubComplicated("asdfla"), genSubComplicated("dsl"), genSubComplicated("daf"), genSubComplicated("123"),
		},
		F6: map[string]*subComplicated{
			"sdlfjl":   genSubComplicated("123"),
			"s12dlfjl": genSubComplicated("12"),
			"sdl@1fjl": genSubComplicated("ajfsdijf"),
			"sdlf12jl": genSubComplicated("dsaljflsdj"),
			"sdlfj12l": genSubComplicated("daskjfl"),
		},
		F7: map[string][]*subComplicated{
			"dfsd2f":  {genSubComplicated("asdjfl"), genSubComplicated("asdfl"), genSubComplicated("asdlfj")},
			"df2sdf":  {genSubComplicated("asdlkf"), genSubComplicated("asdj"), genSubComplicated("asdlkf")},
			"dfs21df": {genSubComplicated("laskjdf"), genSubComplicated("asldjf"), genSubComplicated("qwejdl")},
			"dfsd12f": {genSubComplicated("aoifa"), genSubComplicated("14oij"), genSubComplicated("alsdkf")},
			"dfsd1f":  {genSubComplicated("123uje"), genSubComplicated("alksdfj"), genSubComplicated("uja8isd")},
		},
	}
}

func BenchmarkComplicatedGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getComplicated()
	}
}

//func BenchmarkComplicatedSprintf(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		complicatedItem := getComplicated()
//		_ = fmt.Sprintf("%v", complicatedItem)
//	}
//}

func BenchmarkComplicatedSprintfPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		_ = fmt.Sprintf("%+v", complicatedItem)
	}
}

func BenchmarkComplicatedSprintfPlusLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		v := fmt.Sprintf("%+v", complicatedItem)
		Default.Info(v)
	}
}

func BenchmarkComplicatedSpew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		_ = spew.Sprintf("%v", complicatedItem)
	}
}

func BenchmarkComplicatedSpewLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		v := spew.Sprintf("%v", complicatedItem)
		Default.Info(v)
	}
}

//func BenchmarkComplicatedSpewPlus(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		complicatedItem := getComplicated()
//		_ = spew.Sprintf("%+v", complicatedItem)
//	}
//}

func BenchmarkComplicatedJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		v, _ := json.Marshal(complicatedItem)
		_ = string(v)
	}
}

func BenchmarkComplicatedJsonLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		v, _ := json.Marshal(complicatedItem)
		Default.Info(string(v))
		// time.Sleep(time.Microsecond * 1)
	}
}

func BenchmarkComplicatedJsonIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedItem := getComplicated()
		v, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(complicatedItem)
		_ = string(v)
	}
}

type complicatedCustomJson struct {
	F1 int
	F2 string
	F3 float64
	F4 *subComplicatedCustomJson
	F5 []*subComplicatedCustomJson
	F6 map[string]*subComplicatedCustomJson
	F7 map[string][]*subComplicatedCustomJson
}

func (m *complicatedCustomJson) String() string {
	return fmt.Sprintf(`{"F1": %d, "F2": "%s", "F3": %f, "F4": %v, "F5": %v, "F6": %v, "F7": %v}`,
		m.F1, m.F2, m.F3, m.F4, m.F5, m.F6, m.F7)
}

type subComplicatedCustomJson struct {
	F1 int
	F2 string
	F3 float64
	S1 *subSubComplicatedCustomJson
}

func (m *subComplicatedCustomJson) String() string {
	return fmt.Sprintf(`{"F1": %d, "F2": "%s", "F3": %f, "S1": %v}`, m.F1, m.F2, m.F3, m.S1)
}

type subSubComplicatedCustomJson struct {
	F1 int
	F2 string
	F3 float64
}

func (m *subSubComplicatedCustomJson) String() string {
	return fmt.Sprintf(`{"F1": %d, "F2": "%s", "F3": %f}`, m.F1, m.F2, m.F3)
}

func genSubComplicatedCustomJson(s string) *subComplicatedCustomJson {
	return &subComplicatedCustomJson{
		F1: 14701271122,
		F2: s,
		F3: 102740127.12124211,
		S1: &subSubComplicatedCustomJson{
			F1: 14701271122,
			F2: s,
			F3: 102740127.12124211,
		},
	}
}

func getComplicatedCustomJsonItem() *complicatedCustomJson {
	return &complicatedCustomJson{
		F1: 174012374,
		F2: "aoijhdfaoidsofahoshfoiasfdhjoaijsdifjaosjdfaijosdjf",
		F3: 12312.24124,
		F4: genSubComplicatedCustomJson("sdfadf"),
		F5: []*subComplicatedCustomJson{
			genSubComplicatedCustomJson("alsjfl"), genSubComplicatedCustomJson("asdfla"), genSubComplicatedCustomJson("dsl"), genSubComplicatedCustomJson("daf"), genSubComplicatedCustomJson("123"),
		},
		F6: map[string]*subComplicatedCustomJson{
			"sdlfjl":   genSubComplicatedCustomJson("123"),
			"s12dlfjl": genSubComplicatedCustomJson("12"),
			"sdl@1fjl": genSubComplicatedCustomJson("ajfsdijf"),
			"sdlf12jl": genSubComplicatedCustomJson("dsaljflsdj"),
			"sdlfj12l": genSubComplicatedCustomJson("daskjfl"),
		},
		F7: map[string][]*subComplicatedCustomJson{
			"dfsd2f":  {genSubComplicatedCustomJson("asdjfl"), genSubComplicatedCustomJson("asdfl"), genSubComplicatedCustomJson("asdlfj")},
			"df2sdf":  {genSubComplicatedCustomJson("asdlkf"), genSubComplicatedCustomJson("asdj"), genSubComplicatedCustomJson("asdlkf")},
			"dfs21df": {genSubComplicatedCustomJson("laskjdf"), genSubComplicatedCustomJson("asldjf"), genSubComplicatedCustomJson("qwejdl")},
			"dfsd12f": {genSubComplicatedCustomJson("aoifa"), genSubComplicatedCustomJson("14oij"), genSubComplicatedCustomJson("alsdkf")},
			"dfsd1f":  {genSubComplicatedCustomJson("123uje"), genSubComplicatedCustomJson("alksdfj"), genSubComplicatedCustomJson("uja8isd")},
		},
	}
}

func BenchmarkComplicatedCustomJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedCustomJsonItem := getComplicatedCustomJsonItem()
		_ = fmt.Sprintf("%v", complicatedCustomJsonItem)
	}
}

type complicatedCustomJsonStringBuilder struct {
	F1 int
	F2 string
	F3 float64
	F4 *subComplicatedCustomJsonStringBuilder
	F5 []*subComplicatedCustomJsonStringBuilder
	F6 map[string]*subComplicatedCustomJsonStringBuilder
	F7 map[string][]*subComplicatedCustomJsonStringBuilder
}

func (m *complicatedCustomJsonStringBuilder) String() string {
	var sb strings.Builder

	sb.WriteString(`{"F1": `)
	sb.WriteString(strconv.Itoa(m.F1))
	sb.WriteString(`, "F2": "`)
	sb.WriteString(m.F2)
	sb.WriteString(`", "F3": `)
	sb.WriteString(strconv.FormatFloat(m.F3, 'f', -1, 64))
	sb.WriteString(`, "F4": `)
	sb.WriteString(fmt.Sprintf("%v", m.F4))
	sb.WriteString(`, "F5": `)
	sb.WriteString(fmt.Sprintf("%v", m.F5))
	sb.WriteString(`, "F6": `)
	sb.WriteString(fmt.Sprintf("%v", m.F6))
	sb.WriteString(`, "F7": `)
	sb.WriteString(fmt.Sprintf("%v", m.F7))
	sb.WriteString(`}`)

	return sb.String()
}

type subComplicatedCustomJsonStringBuilder struct {
	F1 int
	F2 string
	F3 float64
	S1 *subSubComplicatedCustomJsonStringBuilder
}

func (m *subComplicatedCustomJsonStringBuilder) String() string {
	var sb strings.Builder

	sb.WriteString(`{"F1": `)
	sb.WriteString(strconv.Itoa(m.F1))
	sb.WriteString(`, "F2": "`)
	sb.WriteString(m.F2)
	sb.WriteString(`", "F3": `)
	sb.WriteString(strconv.FormatFloat(m.F3, 'f', -1, 64))
	sb.WriteString(`, "S1": `)
	sb.WriteString(fmt.Sprintf("%v", m.S1))
	sb.WriteString(`}`)

	return sb.String()
}

type subSubComplicatedCustomJsonStringBuilder struct {
	F1 int
	F2 string
	F3 float64
}

func (m *subSubComplicatedCustomJsonStringBuilder) String() string {
	var sb strings.Builder

	sb.WriteString(`{"F1": `)
	sb.WriteString(strconv.Itoa(m.F1))
	sb.WriteString(`, "F2": "`)
	sb.WriteString(m.F2)
	sb.WriteString(`", "F3": `)
	sb.WriteString(strconv.FormatFloat(m.F3, 'f', -1, 64))
	sb.WriteString(`}`)

	return sb.String()
}

func genSubComplicatedCustomJsonStringBuilder(s string) *subComplicatedCustomJsonStringBuilder {
	return &subComplicatedCustomJsonStringBuilder{
		F1: 14701271122,
		F2: s,
		F3: 102740127.12124211,
		S1: &subSubComplicatedCustomJsonStringBuilder{
			F1: 14701271122,
			F2: s,
			F3: 102740127.12124211,
		},
	}
}

func getComplicatedCustomJsonStringBuilderItem() *complicatedCustomJsonStringBuilder {
	return &complicatedCustomJsonStringBuilder{
		F1: 174012374,
		F2: "aoijhdfaoidsofahoshfoiasfdhjoaijsdifjaosjdfaijosdjf",
		F3: 12312.24124,
		F4: genSubComplicatedCustomJsonStringBuilder("sdfadf"),
		F5: []*subComplicatedCustomJsonStringBuilder{
			genSubComplicatedCustomJsonStringBuilder("alsjfl"), genSubComplicatedCustomJsonStringBuilder("asdfla"), genSubComplicatedCustomJsonStringBuilder("dsl"), genSubComplicatedCustomJsonStringBuilder("daf"), genSubComplicatedCustomJsonStringBuilder("123"),
		},
		F6: map[string]*subComplicatedCustomJsonStringBuilder{
			"sdlfjl":   genSubComplicatedCustomJsonStringBuilder("123"),
			"s12dlfjl": genSubComplicatedCustomJsonStringBuilder("12"),
			"sdl@1fjl": genSubComplicatedCustomJsonStringBuilder("ajfsdijf"),
			"sdlf12jl": genSubComplicatedCustomJsonStringBuilder("dsaljflsdj"),
			"sdlfj12l": genSubComplicatedCustomJsonStringBuilder("daskjfl"),
		},
		F7: map[string][]*subComplicatedCustomJsonStringBuilder{
			"dfsd2f":  {genSubComplicatedCustomJsonStringBuilder("asdjfl"), genSubComplicatedCustomJsonStringBuilder("asdfl"), genSubComplicatedCustomJsonStringBuilder("asdlfj")},
			"df2sdf":  {genSubComplicatedCustomJsonStringBuilder("asdlkf"), genSubComplicatedCustomJsonStringBuilder("asdj"), genSubComplicatedCustomJsonStringBuilder("asdlkf")},
			"dfs21df": {genSubComplicatedCustomJsonStringBuilder("laskjdf"), genSubComplicatedCustomJsonStringBuilder("asldjf"), genSubComplicatedCustomJsonStringBuilder("qwejdl")},
			"dfsd12f": {genSubComplicatedCustomJsonStringBuilder("aoifa"), genSubComplicatedCustomJsonStringBuilder("14oij"), genSubComplicatedCustomJsonStringBuilder("alsdkf")},
			"dfsd1f":  {genSubComplicatedCustomJsonStringBuilder("123uje"), genSubComplicatedCustomJsonStringBuilder("alksdfj"), genSubComplicatedCustomJsonStringBuilder("uja8isd")},
		},
	}
}

func BenchmarkComplicatedCustomJsonStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedCustomJsonStringBuilderItem := getComplicatedCustomJsonStringBuilderItem()
		_ = fmt.Sprintf("%v", complicatedCustomJsonStringBuilderItem)
	}
}

type complicatedStringer struct {
	F1 int
	F2 string
	F3 float64
	F4 *subComplicatedStringer
	F5 []*subComplicatedStringer
	F6 map[string]*subComplicatedStringer
	F7 map[string][]*subComplicatedStringer
}

func (c *complicatedStringer) String() string {
	return fmt.Sprintf("%+v", *c)
}

type subComplicatedStringer struct {
	F1 int
	F2 string
	F3 float64
	S1 *subSubComplicatedStringer
}

func (c *subComplicatedStringer) String() string {
	return fmt.Sprintf("%+v", *c)
}

type subSubComplicatedStringer struct {
	F1 int
	F2 string
	F3 float64
}

func (c *subSubComplicatedStringer) String() string {
	return fmt.Sprintf("%+v", *c)
}

func genSubComplicatedStringer(s string) *subComplicatedStringer {
	return &subComplicatedStringer{
		F1: 14701271122,
		F2: s,
		F3: 102740127.12124211,
		S1: &subSubComplicatedStringer{
			F1: 14701271122,
			F2: s,
			F3: 102740127.12124211,
		},
	}
}

func getComplicatedStringer() *complicatedStringer {
	return &complicatedStringer{
		F1: 174012374,
		F2: "aoijhdfaoidsofahoshfoiasfdhjoaijsdifjaosjdfaijosdjf",
		F3: 12312.24124,
		F4: genSubComplicatedStringer("sdfadf"),
		F5: []*subComplicatedStringer{
			genSubComplicatedStringer("alsjfl"), genSubComplicatedStringer("asdfla"), genSubComplicatedStringer("dsl"), genSubComplicatedStringer("daf"), genSubComplicatedStringer("123"),
		},
		F6: map[string]*subComplicatedStringer{
			"sdlfjl":   genSubComplicatedStringer("123"),
			"s12dlfjl": genSubComplicatedStringer("12"),
			"sdl@1fjl": genSubComplicatedStringer("ajfsdijf"),
			"sdlf12jl": genSubComplicatedStringer("dsaljflsdj"),
			"sdlfj12l": genSubComplicatedStringer("daskjfl"),
		},
		F7: map[string][]*subComplicatedStringer{
			"dfsd2f":  {genSubComplicatedStringer("asdjfl"), genSubComplicatedStringer("asdfl"), genSubComplicatedStringer("asdlfj")},
			"df2sdf":  {genSubComplicatedStringer("asdlkf"), genSubComplicatedStringer("asdj"), genSubComplicatedStringer("asdlkf")},
			"dfs21df": {genSubComplicatedStringer("laskjdf"), genSubComplicatedStringer("asldjf"), genSubComplicatedStringer("qwejdl")},
			"dfsd12f": {genSubComplicatedStringer("aoifa"), genSubComplicatedStringer("14oij"), genSubComplicatedStringer("alsdkf")},
			"dfsd1f":  {genSubComplicatedStringer("123uje"), genSubComplicatedStringer("alksdfj"), genSubComplicatedStringer("uja8isd")},
		},
	}
}

func BenchmarkComplicatedStringer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedStringerItem := getComplicatedStringer()
		_ = fmt.Sprintf("%v", complicatedStringerItem)
	}
}

//func BenchmarkComplicatedStringerPlus(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		complicatedStringerItem := getComplicatedStringer()
//		_ = fmt.Sprintf("%+v", complicatedStringerItem)
//	}
//}

type complicatedStringerRetype struct {
	F1 int
	F2 string
	F3 float64
	F4 *subComplicatedStringerRetype
	F5 []*subComplicatedStringerRetype
	F6 map[string]*subComplicatedStringerRetype
	F7 map[string][]*subComplicatedStringerRetype
}

type complicatedStringerRetypeTarget complicatedStringerRetype

func (c *complicatedStringerRetype) String() string {
	return fmt.Sprintf("%+v", (*complicatedStringerRetypeTarget)(c))
}

type subComplicatedStringerRetype struct {
	F1 int
	F2 string
	F3 float64
	S1 *subSubComplicatedStringerRetype
}

type subComplicatedStringerRetypeTarget subComplicatedStringerRetype

func (c *subComplicatedStringerRetype) String() string {
	return fmt.Sprintf("%+v", (*subComplicatedStringerRetypeTarget)(c))
}

type subSubComplicatedStringerRetype struct {
	F1 int
	F2 string
	F3 float64
}

type subSubComplicatedStringerRetypeTarget subSubComplicatedStringerRetype

func (c *subSubComplicatedStringerRetype) String() string {
	return fmt.Sprintf("%+v", (*subSubComplicatedStringerRetypeTarget)(c))
}

func genSubComplicatedStringerRetype(s string) *subComplicatedStringerRetype {
	return &subComplicatedStringerRetype{
		F1: 14701271122,
		F2: s,
		F3: 102740127.12124211,
		S1: &subSubComplicatedStringerRetype{
			F1: 14701271122,
			F2: s,
			F3: 102740127.12124211,
		},
	}
}

func getComplicatedStringerRetype() *complicatedStringerRetype {
	return &complicatedStringerRetype{
		F1: 174012374,
		F2: "aoijhdfaoidsofahoshfoiasfdhjoaijsdifjaosjdfaijosdjf",
		F3: 12312.24124,
		F4: genSubComplicatedStringerRetype("sdfadf"),
		F5: []*subComplicatedStringerRetype{
			genSubComplicatedStringerRetype("alsjfl"), genSubComplicatedStringerRetype("asdfla"), genSubComplicatedStringerRetype("dsl"), genSubComplicatedStringerRetype("daf"), genSubComplicatedStringerRetype("123"),
		},
		F6: map[string]*subComplicatedStringerRetype{
			"sdlfjl":   genSubComplicatedStringerRetype("123"),
			"s12dlfjl": genSubComplicatedStringerRetype("12"),
			"sdl@1fjl": genSubComplicatedStringerRetype("ajfsdijf"),
			"sdlf12jl": genSubComplicatedStringerRetype("dsaljflsdj"),
			"sdlfj12l": genSubComplicatedStringerRetype("daskjfl"),
		},
		F7: map[string][]*subComplicatedStringerRetype{
			"dfsd2f":  {genSubComplicatedStringerRetype("asdjfl"), genSubComplicatedStringerRetype("asdfl"), genSubComplicatedStringerRetype("asdlfj")},
			"df2sdf":  {genSubComplicatedStringerRetype("asdlkf"), genSubComplicatedStringerRetype("asdj"), genSubComplicatedStringerRetype("asdlkf")},
			"dfs21df": {genSubComplicatedStringerRetype("laskjdf"), genSubComplicatedStringerRetype("asldjf"), genSubComplicatedStringerRetype("qwejdl")},
			"dfsd12f": {genSubComplicatedStringerRetype("aoifa"), genSubComplicatedStringerRetype("14oij"), genSubComplicatedStringerRetype("alsdkf")},
			"dfsd1f":  {genSubComplicatedStringerRetype("123uje"), genSubComplicatedStringerRetype("alksdfj"), genSubComplicatedStringerRetype("uja8isd")},
		},
	}
}

func BenchmarkComplicatedStringerRetype(b *testing.B) {
	for i := 0; i < b.N; i++ {
		complicatedStringerRetypeItem := getComplicatedStringerRetype()
		_ = fmt.Sprintf("%v", complicatedStringerRetypeItem)
	}
}

func TestLen(t *testing.T) {
	complicatedItem := getComplicated()
	v1 := fmt.Sprintf("%+v", complicatedItem)
	t.Logf("SprintfPlus len=%d", len(v1))
	v2, _ := json.Marshal(complicatedItem)
	t.Logf("JSON len=%d", len(v2))
	//=== RUN   TestLen
	//sprint_test.go:391: SprintfPlus len=513
	//sprint_test.go:393: JSON len=3248
	//--- PASS: TestLen (0.00s)
	//PASS
}

type outer struct {
	Inner *inner
}

type inner struct {
	Outer *outer
}

func TestCycle(t *testing.T) {
	o := &outer{Inner: &inner{}}
	_, err := json.Marshal(o)
	assert.NoError(t, err)
	o.Inner.Outer = o
	_, err = json.Marshal(o)
	assert.Error(t, err)
}
