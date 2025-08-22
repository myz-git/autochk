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

	// 设置模板文件路径为当前目录下的 chk198S.docx
	templatePath := "chk198S.docx"

	// 确保输出目录存在
	outputDir := "report"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755) // 创建目录
	}

	// 设置输出文件的路径和文件名
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_Done.docx", prefix))

	// 从结构体构造替换映射 - 完整映射所有字段
	replaceMap := docx.PlaceholderMap{
		// InfoSht 基本信息字段
		"DBNAME":      infstp.DbName,
		"DBVER":       infstp.DbVer,
		"DBROLE":      infstp.DbRole,
		"LOGMODE":     infstp.LogMode,
		"FLASHBACK":   infstp.FlashBack,
		"DBTOTALSIZE": infstp.DbTotalsize,
		"DBFILECOUNT": infstp.DbFilecount,
		"DBTBLCOUNT":  infstp.DbTblcount,
		"DBLANG":      infstp.DbLang,
		"DBMAA":       infstp.DbMaa,
		"HOSTNAME":    infstp.HostName,
		"IPADDR":      infstp.Ipaddr,
		"OS":          infstp.Os,
		"RELVER":      infstp.Relver,
		"CORES":       infstp.Cores,
		"CPUCOUNT":    infstp.CpuCount,
		"CPUMHZ":      infstp.CpuMHZ,
		"MEMTOTAL":    infstp.MemTotal,
		"SWAPTOTAL":   infstp.SwapTotal,
		"OTHERS":      infstp.Others,

		// OsSht OS相关字段 - 访问Contents字段
		"OSPARAMETER": osshtp.Osparameter.Contents,
		"ULIMIT":      osshtp.Ulimit.Contents,
		"FILESYSTEM":  osshtp.Filesystem.Contents,
		"INODEUSAGE":  osshtp.Inodeusage.Contents,
		"CPUSTAT":     osshtp.Cpustat.Contents,
		"MEMSTAT":     osshtp.Memstat.Contents,
		"IOSTAT":      osshtp.Iostat.Contents,
		"THPSTAT":     osshtp.Thpstat.Contents,
		"HUGPAGE":     osshtp.Hugpage.Contents,
		"NUMA":        osshtp.Numa.Contents,
		"NTP":         osshtp.Ntp.Contents,

		// DbSht 数据库相关字段 - 访问Contents字段
		"DBTBSUSAGE":               dbshtp.DbTbsusage.Contents,
		"DBDATAFILE":               dbshtp.Dbdatafile.Contents,
		"DBCONTROLFILE":            dbshtp.Dbcontrolfile.Contents,
		"DBUSERSIZE":               dbshtp.Dbusersize.Contents,
		"DBREDOCHECK":              dbshtp.Dbredocheck.Contents,
		"DBREDOSWITCH":             dbshtp.Dbredoswitch.Contents,
		"DBRESOURCE":               dbshtp.Dbresource.Contents,
		"LOADPROFILE":              dbshtp.Loadprofile.Contents,
		"INSTEFFICIENCY":           dbshtp.Instefficiency.Contents,
		"DBTOPEVENT":               dbshtp.Dbtopevent.Contents,
		"DBTOPSQL":                 dbshtp.DbtopSQL.Contents,
		"DBLSNRINFO":               dbshtp.Dblsnrinfo.Contents,
		"DBTABLEPARALLEL":          dbshtp.Dbtableparallel.Contents,
		"DBINDEXPARALLEL":          dbshtp.Dbindexparallel.Contents,
		"DBINVALIDINDEX":           dbshtp.Dbinvalidindex.Contents,
		"DBSEQUENCE":               dbshtp.Dbsequence.Contents,
		"DB_SEQ_USAGE":             dbshtp.Db_seq_usage.Contents,
		"DBRECOVERYDEST":           dbshtp.Dbrecoverydest.Contents,
		"DBFLASHRECOVERYUSEAGE":    dbshtp.Dbflashrecoveryuseage.Contents,
		"DBERRLOG":                 dbshtp.Dberrlog.Contents,
		"DBPRODUCTUSERFAILEDLOGIN": dbshtp.Dbproductuserfailedlogin.Contents,
		"DBDGLAGCHECK":             dbshtp.Dbdglagcheck.Contents,
		"DBDGERRCHECK":             dbshtp.Dbdgerrcheck.Contents,
		"DBRMANCHECK":              dbshtp.Dbrmancheck.Contents,
		"DBDBAPRIV":                dbshtp.Dbdbapriv.Contents,
		"DBSYSDBA":                 dbshtp.Dbsysdba.Contents,
		"DBAUDITSEGMENT":           dbshtp.Dbauditsegment.Contents,
		"DBAUDITCONT":              dbshtp.Dbauditcont.Contents,
		"DB_NOSYS_IN_SYSTEM":       dbshtp.Db_Nosys_In_System.Contents,
		"DBVIRSCHECK":              dbshtp.Dbvirscheck.Contents,
		"DBSCNHEALTHCHECK":         dbshtp.Dbscnhealthcheck.Contents,
		"DBCRSCHECK":               dbshtp.Dbcrscheck.Contents,
		"DBASMUSAGE":               dbshtp.Dbasmusage.Contents,
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
