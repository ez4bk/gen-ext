package ezgen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
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
	Desc           bool // params key sort
}

//go:embed crud.dao.tpl
var crudTemplate string

//go:embed custom.dao.tpl
var customTemplate string

func BuildParamsKey(colGo, colGoType string, unique bool) string {
	if colGoType == "string" && !unique {
		return fmt.Sprintf("%s %s // optional, likely", colGo, colGoType)
	} else {
		return fmt.Sprintf("%s %s // optional", colGo, colGoType)
	}
}

func BuildCand(colGo, columnName, colGoType string, unique bool) string {
	if colGoType == `string` && !unique {
		return fmt.Sprintf(`Scopes(ezgen.Cond(!reflect.ValueOf(params.%s).IsZero(), "%s like ?", "%%"+params.%s+"%%")).`, colGo, columnName, colGo)
	} else {
		return fmt.Sprintf(`Scopes(ezgen.Cond(!reflect.ValueOf(params.%s).IsZero(), "%s = ?", params.%s)).`, colGo, columnName, colGo)
	}
}

func BuildNullable(colGo, columnName string) string {
	return fmt.Sprintf(`Scopes(ezgen.Nullable(params.%s != nil, "%s = ?", func() any { return *params.%s })).`, colGo, columnName, colGo)
}

func Build(params *GenParams, fileName string) error {
	// 创建一个 buffer 用于存储生成的代码
	var buf bytes.Buffer

	// 解析和执行模板
	tmpl, err := template.New("dao").Parse(crudTemplate)
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

func BuildCustom(params *GenParams, fileName string) error {
	// 文件存在不生成
	if fileExists(fileName) {
		return nil
	}

	// 创建一个 buffer 用于存储生成的代码
	var buf bytes.Buffer

	// 解析和执行模板
	tmpl, err := template.New("dao-custom").Parse(customTemplate)
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

func BuildDao(modelStructNames []string, fileName string) error {
	daoNames := make([]string, 0, len(modelStructNames))
	modelNames := make([]string, 0, len(modelStructNames))
	for _, modelStructName := range modelStructNames {
		daoNames = append(daoNames, unCapitalize(modelStructName)+"Dao")
		modelNames = append(modelNames, modelStructName)
	}
	type genParams struct {
		DaoNameList   []string // dao names
		ModelNameList []string // model names
	}

	params := &genParams{
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

func Params(table, modelStructName string, columnTypes []gorm.ColumnType,
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
	}
	for _, columnType := range columnTypes {
		columnName := columnType.Name()
		colGo := SnakeToPascalCase(columnName)
		colGoType := dataMap[strings.ToLower(columnType.DatabaseTypeName())](columnType)
		unique := false

		if isPrimaryKey, ok := columnType.PrimaryKey(); ok {
			if isPrimaryKey {
				p.PrimaryField = columnName
				p.PKType = colGoType

				if comment, ok := columnType.Comment(); ok {
					if strings.Contains(comment, "asc") {
						p.Desc = false
					} else if strings.Contains(comment, "desc") {
						p.Desc = true
					}
				}

				continue
			}
		}

		if flag, ok := columnType.Unique(); ok {
			unique = flag
		}

		if columnName == "version" || columnName == "created_at" || columnName == "updated_at" || columnName == "deleted_at" {
			continue
		}
		p.ParamsKey = append(p.ParamsKey, BuildParamsKey(colGo, colGoType, unique))
		if strings.HasPrefix(colGoType, "*") {
			p.ParamsScopes = append(p.ParamsScopes, BuildNullable(colGo, columnName))
		} else {
			p.ParamsScopes = append(p.ParamsScopes, BuildCand(colGo, columnName, colGoType, unique))
		}
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
