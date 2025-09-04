package main

import (
	"autochk/anadata"
	"autochk/readxml"
	"autochk/structs"
	"autochk/todocx"
	"autochk/toxls"
	"autochk/xmlfile"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// var arg string
	// if len(os.Args) == 1 {
	// 	fmt.Println("请输入参数")
	// 	// os.Exit(0)
	// 	flg = "0"  //0 默认一个文件
	// } else {
	// }

	//加入-s 的参数. 默认为 false
	singlefile := flag.Bool("s", true, "single file: true|false , default false")

	flag.Parse()
	// log.Println("use single file mode: ", *singlefile)

	// 判断是否 输入参数, 没有参数则退出
	// if len(os.Args) < 2 {
	// 	fmt.Println("expected parameter")
	// 	os.Exit(0)
	// }

	start := time.Now()

	log.Println("######---Start---######")

	// 首先执行XML文件合并
	log.Println("开始执行XML文件合并...")
	if err := xmlfile.MergeXMLFiles(); err != nil {
		log.Printf("XML文件合并失败: %v", err)
		return
	}
	log.Println("XML文件合并完成")

	//删除*Done.xlsx文件
	ClearFile(*singlefile)

	// if !*singlefile {
	// 	toxls.NewXlsx("xxx", *singlefile)
	// }

	files := GetXMLS("R")
	files = append(files, GetXMLS("S")...)
	// colInxp := &utils.ColInx

	//循环打开文件名为*.R.xml或*.S.xml的文件
	colcnt := 1
	for _, fnm := range files {
		if *singlefile {
			colcnt = 1
		}
		log.Println("开始处理--->", fnm)
		fileName := filepath.Base(fnm)
		prex := strings.Replace(fileName, ".R.xml", "", -1)
		prex = strings.Replace(prex, ".S.xml", "", -1)

		//初始化新的数据结构
		var osshts []structs.OsShts
		dbsht := structs.DbSht{}
		var instshts []structs.InstShts
		summaryEntries := &structs.SummaryEntries{}

		readxml.ReadXml(fnm, &osshts, &dbsht, &instshts)
		anadata.Ana(&osshts, &dbsht, &instshts, summaryEntries)
		toxls.Xlsx(&osshts, &dbsht, &instshts, summaryEntries, prex, colcnt, *singlefile)
		todocx.Todocxfunc Xlsx(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool) {
		colcnt++
	}
	elapsed := time.Since(start)
	log.Printf("#####---Completed! Elapsed Time:%v---#####", elapsed)
	log.Println("Anaylze Check Data (ACD) release 1.8")
	// fmt.Printf("执行完成,请按任意键退出...")

}

func GetXMLS(typ string) (xmlnms []string) {
	//遍历打开xmlfile/output_xml路径下的指定后缀的xml文件
	dirname := "xmlfile/output_xml"
	//根据传入的类型来确定按什么样的后缀遍历文件 ,如  ".DB.xml" ".OS.xml"  ".AWR.xml"
	xmltyp := "." + typ + ".xml"
	f, err := os.Open(dirname)
	if err != nil {
		log.Printf("打开目录 %s 失败: %v", dirname, err)
		return xmlnms
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Printf("读取目录 %s 失败: %v", dirname, err)
		return xmlnms
	}
	for _, file := range files {
		//但*AWR.xml及*OS.xml除外
		// if strings.HasSuffix(file.Name(), "*.xml") && file.Name() != ".AWR.xml" && file.Name() != ".OS.xml" {
		if strings.HasSuffix(file.Name(), xmltyp) {
			xmlnms = append(xmlnms, filepath.Join(dirname, file.Name()))
		}
	}
	return xmlnms
}

func ClearFile(sglf bool) {
	//遍历打开当前路径下的文件
	dirname := "."
	var xmltyp string
	if sglf {
		xmltyp = ".Done.xlsx"
	} else {
		xmltyp = ".ALLDone.xlsx"
	}

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		//遍历查找是否为"*.Done.xlsx"结尾的文件,如果是则删除 {
		if strings.HasSuffix(file.Name(), xmltyp) {
			del := os.Remove(file.Name())
			if del != nil {
				log.Println(del)
			}
		}
	}
}
