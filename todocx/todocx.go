package todocx

import (
	"autochk/structs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lukasjarosch/go-docx"
)

func Todocx(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, prefix string, colcnt int, sglf bool) {
	startTime := time.Now()

	// 设置模板文件路径为当前目录下的 chk197S.docx
	templatePath := "chk197S.docx"

	// 确保输出目录存在
	outputDir := "report"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755) // 创建目录
	}

	// 设置输出文件的路径和文件名
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_Done.docx", prefix))

	// 从结构体构造替换映射
	replaceMap := docx.PlaceholderMap{
		"DBNAME":        infstp.DbName,
		"DBVER":         infstp.DbVer,
		"DBROLE":        infstp.DbRole,
		"LOGMODE":       infstp.LogMode,
		"FLASHBACK":     infstp.FlashBack,
		"DBTOTALSIZE":   infstp.DbTotalsize,
		"DBFILECOUNT":   infstp.DbFilecount,
		"DBTBLCOUNT":    infstp.DbTblcount,
		"DBLANG":        infstp.DbLang,
		"HOSTNAME":      infstp.HostName,
		"IPADDR":        infstp.Ipaddr,
		"OS":            infstp.Os,
		"RELVER":        infstp.Relver,
		"CORES":         infstp.Cores,
		"CPUCOUNT":      infstp.CpuCount,
		"CPUMHZ":        infstp.CpuMHZ,
		"MEMTOTAL":      infstp.MemTotal,
		"SWAPTOTAL":     infstp.SwapTotal,
		"OSPARAMETER":   osshtp.Osparameter,
		"DBTBSUSAGE":    dbshtp.DbTbsusage,
		"DBDATAFILE":    dbshtp.Dbdatafile,
		"DBCONTROLFILE": dbshtp.Dbcontrolfile,
	}

	// 打开模板文件
	doc, err := docx.Open(templatePath)
	if err != nil {
		log.Fatalf("Failed to open docx template: %v", err)
	}

	log.Printf("open took: %s", time.Since(startTime))

	// 替换模板中的占位符
	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		log.Fatalf("Failed to replace placeholders: %v", err)
	}

	log.Printf("replace took: %s", time.Since(startTime))

	// 保存填充后的文档
	err = doc.WriteToFile(outputFile)
	if err != nil {
		log.Fatalf("Failed to write to docx file: %v", err)
	}

	log.Printf("everything took: %s", time.Since(startTime))
	log.Printf("Generated docx report at %s", outputFile)
}
