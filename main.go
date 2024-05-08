package main

import (
	"autochk/anadata"
	"autochk/readxml"
	"autochk/structs"
	"autochk/todocx"
	"autochk/toxls"
	"flag"
	"log"
	"os"
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

	// //初始化表最后一列号
	// colInxp := &utils.ColInx
	// *colInxp = 0

	//删除*Done.xlsx文件
	ClearFile(*singlefile)

	// if !*singlefile {
	// 	toxls.NewXlsx("xxx", *singlefile)
	// }

	files := GetXMLS("ALL")
	// colInxp := &utils.ColInx

	//循环打开文件名为*.ALL.xml的文件
	colcnt := 1
	for _, fnm := range files {
		if *singlefile {
			colcnt = 1
		}
		log.Println("开始处理--->", fnm)
		prex := strings.Replace(fnm, ".ALL.xml", "", -1)

		//初始化三个结构
		infosht := structs.InfoSht{}
		ossht := structs.OsSht{}
		dbsht := structs.DbSht{}

		readxml.ReadXml(fnm, &infosht, &ossht, &dbsht)
		anadata.Ana(&infosht, &ossht, &dbsht)
		toxls.Xlsx(&infosht, &ossht, &dbsht, prex, colcnt, *singlefile)
		todocx.Todocx(&infosht, &ossht, &dbsht, prex, colcnt, *singlefile)
		colcnt++
	}
	elapsed := time.Since(start)
	log.Printf("#####---Completed! Elapsed Time:%v---#####", elapsed)
	log.Println("Anaylze Check Data (ACD) release 1.8")
	// fmt.Printf("执行完成,请按任意键退出...")

}

func GetXMLS(typ string) (xmlnms []string) {
	//遍历打开当前路径下的指定后缀的xml文件
	dirname := "."
	//根据传入的类型来确定按什么样的后缀遍历文件 ,如  ".DB.xml" ".OS.xml"  ".AWR.xml"
	xmltyp := "." + typ + ".xml"
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		//但*AWR.xml及*OS.xml除外
		// if strings.HasSuffix(file.Name(), "*.xml") && file.Name() != ".AWR.xml" && file.Name() != ".OS.xml" {
		if strings.HasSuffix(file.Name(), xmltyp) {
			xmlnms = append(xmlnms, file.Name())
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
