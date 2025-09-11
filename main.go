package main

import (
	"autochk/anadata"
	"autochk/readxml"
	"autochk/structs"
	"autochk/toxls"
	"autochk/utils"
	"autochk/xmlfile"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Version   = "0.0.0-dev"
	Commit    = "unknown"
	BuildDate = ""
)

func main() {
	// 先行处理 --version / version
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("autochk version %s (commit %s, built %s)\n", Version, Commit, BuildDate)
		return
	}

	// 定义命令行参数
	custNm := flag.String("u", "", "customer name: specify customer name")

	flag.Parse()

	// 显示使用说明
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		showUsage()
		return
	}

	// 处理客户名称逻辑
	var finalCustNm string

	// 检查是否使用了-u参数（不管是否有值）
	hasUFlag := false
	for i, arg := range os.Args {
		if arg == "-u" && i < len(os.Args)-1 {
			hasUFlag = true
			break
		}
	}

	if hasUFlag {
		// 如果使用了-u参数
		if *custNm != "" {
			// 如果-u参数有值，直接使用
			finalCustNm = *custNm
			utils.LogInfof("使用命令行指定的客户名称: %s", finalCustNm)
		} else {
			// 如果-u参数没有值，提示用户输入
			finalCustNm = getCustomerName()
		}
	} else {
		// 如果没有使用-u参数，直接使用默认值
		finalCustNm = "    " // 4个空格
		utils.LogInfof("使用默认客户名称: %s", finalCustNm)
	}

	start := time.Now()

	utils.LogInfof("######---Start---######")
	utils.LogInfof("Version: %s, Commit: %s, Build: %s", Version, Commit, BuildDate)
	utils.LogInfof("客户名称: %s", finalCustNm)

	// 首先执行XML文件合并
	utils.LogInfof("开始执行XML文件合并...")
	if err := xmlfile.MergeXMLFiles(); err != nil {
		utils.LogErrorf("XML文件合并失败: %v", err)
		return
	}
	utils.LogInfof("XML文件合并完成")

	//删除*Done.xlsx文件
	ClearFile()

	files := GetXMLS("R")
	files = append(files, GetXMLS("S")...)

	//循环打开文件名为*.R.xml或*.S.xml的文件
	colcnt := 1
	for _, fnm := range files {
		utils.LogInfof("开始处理---> %s", fnm)
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
		// 传递客户名称参数，使用单文件模式
		toxls.Xlsx(&osshts, &dbsht, &instshts, summaryEntries, prex, colcnt, true, finalCustNm)
		colcnt++
	}
	elapsed := time.Since(start)
	utils.LogInfof("#####---Completed! Elapsed Time:%v---#####", elapsed)
	utils.LogInfof("Anaylze Check Data (ACD) %s", Version)
}

// 获取客户名称输入
func getCustomerName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please Input Customer Name: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 如果直接回车，默认赋值为4个空格
	if input == "" {
		return "    " // 4个空格
	}
	return input
}

// 显示使用说明
func showUsage() {
	fmt.Println("健康检查报告生成工具")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  autochk.exe                    # 不带参数执行，使用默认客户名称(4个空格)")
	fmt.Println("  autochk.exe -u 客户名称         # 指定客户名称，不提示输入")
	fmt.Println("  autochk.exe -u                  # 使用-u参数但不指定值，提示输入客户名称")
	fmt.Println("  autochk.exe --version           # 显示版本信息")
	fmt.Println("  autochk.exe -h, --help          # 显示帮助信息")
	fmt.Println("")
	fmt.Println("参数说明:")
	fmt.Println("  -u string    客户名称 (可选，不指定则提示输入)")
	fmt.Println("  -h, --help   显示帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  autochk.exe                    # 使用默认客户名称，不提示输入")
	fmt.Println("  autochk.exe -u ABC公司          # 直接使用ABC公司作为客户名称")
	fmt.Println("  autochk.exe -u                  # 使用-u参数但提示输入客户名称")
}

func GetXMLS(typ string) (xmlnms []string) {
	//遍历打开xmlfile/output_xml路径下的指定后缀的xml文件
	dirname := "xmlfile/output_xml"
	//根据传入的类型来确定按什么样的后缀遍历文件 ,如  ".DB.xml" ".OS.xml"  ".AWR.xml"
	xmltyp := "." + typ + ".xml"
	f, err := os.Open(dirname)
	if err != nil {
		utils.LogErrorf("打开目录 %s 失败: %v", dirname, err)
		return xmlnms
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		utils.LogErrorf("读取目录 %s 失败: %v", dirname, err)
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

func ClearFile() {
	//遍历打开当前路径下的文件
	dirname := "."
	xlsxTyp := ".Done.xlsx"
	htmlTyp := ".Done.html"

	f, err := os.Open(dirname)
	if err != nil {
		utils.LogErrorf("打开目录失败: %v", err)
		return
	}
	files, err := f.Readdir(-1)
	if err != nil {
		utils.LogErrorf("读取目录失败: %v", err)
		return
	}
	for _, file := range files {
		//遍历查找是否为"*.Done.xlsx"或"*.Done.html"结尾的文件,如果是则删除
		if strings.HasSuffix(file.Name(), xlsxTyp) || strings.HasSuffix(file.Name(), htmlTyp) {
			del := os.Remove(file.Name())
			if del != nil {
				utils.LogWarnf("删除文件失败: %v", del)
			}
		}
	}
}
