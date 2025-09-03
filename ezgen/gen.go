package ezgen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
	"gorm.io/gorm"
)

type GenParams struct {
	ModelPackage   string   // model package
	DaoName        string   // dao name
	ModelName      string   // model name
	S              string   // the first letter(lower case) of simple Name (receiver)
	PKType         string   // primary key type
	ParamsKey      []string // params key
	ParamsScopes   []string // params scopes
	ImportPkgPaths []string
	PrimaryField   string
	Desc           bool   // params key sort
	SortField      string // sort field, default is primary key
}

//go:embed crud.dao.tpl
var crudTemplate string

//go:embed interface.dao.tpl
var interfaceTemplate string

//go:embed dao.tpl
var daoTemplate string

func BuildParamsKey(colGo, colGoType string, unique bool) string {
	if colGoType == "string" && !unique {
		return fmt.Sprintf("%s %s // optional, likely", colGo, colGoType)
	} else if strings.Contains(colGoType, "time.Time") {
		return fmt.Sprintf("%sRange ezgen.TimeRange // optional", colGo)
	} else {
		return fmt.Sprintf("%s %s // optional", colGo, colGoType)
	}
}

func BuildScope(colGo, columnName, colGoType string, unique bool) string {
	if colGoType == `string` && !unique {
		return fmt.Sprintf(`Scopes(ezgen.Cond(!reflect.ValueOf(params.%s).IsZero(), "%s like ?", "%%"+params.%s+"%%")).`, colGo, columnName, colGo)
	} else if strings.Contains(colGoType, "time.Time") {
		return fmt.Sprintf(`Scopes(ezgen.Cond(!reflect.ValueOf(params.%sRange).IsZero(), "? <= %s and %s <= ?",params.%sRange.Start, params.%sRange.End)).`, colGo, columnName, columnName, colGo, colGo)
	} else {
		return fmt.Sprintf(`Scopes(ezgen.Cond(!reflect.ValueOf(params.%s).IsZero(), "%s = ?", params.%s)).`, colGo, columnName, colGo)
	}
}

func BuildNullable(colGo, columnName string) string {
	return fmt.Sprintf(`Scopes(ezgen.Nullable(params.%s != nil, "%s = ?", func() any { return *params.%s })).`, colGo, columnName, colGo)
}

func Generate(params *GenParams, targetDir, entityName string) (err error) {
	crudFileName := filepath.Join(targetDir, entityName+".crud.go")
	interfaceFileName := filepath.Join(targetDir, entityName+".go")
	err = generateCrud(params, crudFileName)
	if err != nil {
		return err
	}
	err = generateInterface(params, interfaceFileName)
	if err != nil {
		return err
	}
	return
}

func generateCrud(params *GenParams, fileName string) error {
	// 创建一个 buffer 用于存储生成的代码
	var buf bytes.Buffer
	// 解析和执行模板
	tmpl, err := template.New("crud").Parse(crudTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, params)
	if err != nil {
		return err
	}

	// 使用 go/format 包格式化代码
	formattedSource, err := imports.Process(fileName, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	// 生成的代码
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(formattedSource)
	if err != nil {
		return err
	}

	return nil
}

func generateInterface(params *GenParams, fileName string) error {
	// 不覆盖已存在的文件
	if fileExists(fileName) {
		return nil
	}

	// 创建一个 buffer 用于存储生成的代码
	var buf bytes.Buffer
	// 解析和执行模板
	tmpl, err := template.New("dao").Parse(interfaceTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, params)
	if err != nil {
		return err
	}

	// 使用 go/format 包格式化代码
	formattedSource, err := imports.Process(fileName, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	// 生成的代码
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(formattedSource)
	if err != nil {
		return err
	}

	return nil
}

func GenerateDao(modelStructNames []string, fileName string) error {
	daoNames := make([]string, 0, len(modelStructNames))
	modelNames := make([]string, 0, len(modelStructNames))
	for _, modelStructName := range modelStructNames {
		daoNames = append(daoNames, unCapitalize(modelStructName)+"Dao")
		modelNames = append(modelNames, modelStructName)
	}
	type genParams struct {
		ModelPackage  string
		DaoNameList   []string // dao names
		ModelNameList []string // model names
	}

	pkgName, _ := getModuleName()

	params := &genParams{
		ModelPackage:  pkgName,
		DaoNameList:   daoNames,
		ModelNameList: modelNames,
	}
	// 创建一个 buffer 用于存储生成的代码
	var buf bytes.Buffer

	// 解析和执行模板
	tmpl, err := template.New("daoInit").Parse(daoTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, params)
	if err != nil {
		return err
	}

	// 使用 go/format 包格式化代码
	formattedSource, err := imports.Process(fileName, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	// 生成的代码
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(formattedSource)
	if err != nil {
		return err
	}

	return nil
}

func BuildParams(table, modelStructName string, columnTypes []gorm.ColumnType,
	dataMap map[string]func(gorm.ColumnType) (dataType string)) (*GenParams, error) {
	goModel, err := getModuleName()
	if err != nil {
		return nil, err
	}

	p := &GenParams{
		ModelPackage:   goModel,
		DaoName:        unCapitalize(modelStructName) + "Dao",
		ModelName:      modelStructName,
		S:              "dao",
		PKType:         "",
		ParamsKey:      make([]string, 0),
		ParamsScopes:   make([]string, 0),
		ImportPkgPaths: nil,
		PrimaryField:   "id",
		Desc:           true,
		SortField:      "id",
	}
	sortField := ""
	for _, columnType := range columnTypes {
		columnName := columnType.Name()
		colGo := SnakeToPascalCase(columnName)
		colGoType := dataMap[strings.ToLower(columnType.DatabaseTypeName())](columnType)
		unique := false

		if isPrimaryKey, ok := columnType.PrimaryKey(); ok && isPrimaryKey {
			p.PrimaryField = columnName
			p.PKType = colGoType
			continue
		}

		if comment, ok := columnType.Comment(); ok {
			if strings.Contains(comment, "asc") {
				sortField = columnName
				p.Desc = false
			} else if strings.Contains(comment, "desc") {
				sortField = columnName
				p.Desc = true
			}
		}

		if flag, ok := columnType.Unique(); ok {
			unique = flag
		}

		if columnName == "version" || columnName == "created_at" || columnName == "updated_at" || columnName == "deleted_at" {
			continue
		}
		p.ParamsKey = append(p.ParamsKey, BuildParamsKey(colGo, colGoType, unique))
		if strings.Contains(colGoType, "time.Time") {
			p.ParamsScopes = append(p.ParamsScopes, BuildScope(colGo, columnName, colGoType, unique))
		} else if strings.HasPrefix(colGoType, "*") {
			p.ParamsScopes = append(p.ParamsScopes, BuildNullable(colGo, columnName))
		} else {
			p.ParamsScopes = append(p.ParamsScopes, BuildScope(colGo, columnName, colGoType, unique))
		}
	}
	if sortField == "" {
		p.SortField = p.PrimaryField
	} else {
		p.SortField = sortField
	}

	if p.PKType == "" {
		err = errors.New(fmt.Sprintf("table %s no primary key", table))
		return nil, err
	}

	return p, nil
}

func getModuleName() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("failed to read build info")
	}
	return info.Main.Path, nil
}

func unCapitalize(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToLower(s[:1]) + s[1:]
}

func SnakeToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	words := strings.Split(s, "_") // 将字符串按 "_" 分割成单词切片
	result := ""
	for _, word := range words {
		if len(word) > 0 {
			// 将每个单词的首字母转换为大写，并将剩余部分与首字母连接起来
			result += strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return result
}

func delPointerSym(name string) string {
	return strings.TrimLeft(name, "*")
}

// the first letter(lower case)of simple Name (receiver)
func getPureName(s string) string {
	return string(strings.ToLower(delPointerSym(s))[0])
}

func fileExists(filename string) bool {
	// 使用 os.Stat 获取文件信息
	_, err := os.Stat(filename)
	// 如果错误为 nil，表示文件存在
	if err == nil {
		return true
	}
	// 如果错误是 os.ErrNotExist，表示文件不存在
	if os.IsNotExist(err) {
		return false
	}
	// 其他错误处理
	return false
}
