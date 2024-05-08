package todocx

import (
	"autochk/structs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/lukasjarosch/go-docx"
)

// createPlaceholderMap 通过反射自动生成替换映射
func createPlaceholderMap(data interface{}) docx.PlaceholderMap {
	result := make(docx.PlaceholderMap)
	val := reflect.ValueOf(data)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 使用字段名作为占位符键，字段值作为映射值
		key := strings.ToUpper(fieldType.Name)        // 转换为大写作为占位符
		value := fmt.Sprintf("%v", field.Interface()) // 将字段值转换为字符串

		// 将生成的键值对添加到结果映射中
		result[key] = value
	}
	return result
}

func Todocx(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, prefix string, colcnt int, sglf bool) {
	startTime := time.Now()

	// 设置模板和输出文件路径
	templatePath := "chk197S.docx"
	outputDir := "report"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_Done.docx", prefix))

	// 通过反射自动生成替换映射
	replaceMap := createPlaceholderMap(infstp)
	for k, v := range createPlaceholderMap(osshtp) {
		replaceMap[k] = v
	}
	for k, v := range createPlaceholderMap(dbshtp) {
		replaceMap[k] = v
	}

	doc, err := docx.Open(templatePath)
	if err != nil {
		log.Fatalf("Failed to open docx template: %v", err)
	}

	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		log.Fatalf("Failed to replace placeholders: %v", err)
	}

	err = doc.WriteToFile(outputFile)
	if err != nil {
		log.Fatalf("Failed to write to docx file: %v", err)
	}

	log.Printf("Document processing completed in: %s", time.Since(startTime))
	log.Printf("Generated docx report at %s", outputFile)
}
