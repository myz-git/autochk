package xmlfile

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

// XML文件信息结构
type XMLFileInfo struct {
	FilePath    string
	FileName    string
	Date        string
	Hostname    string
	DBName      string
	InstName    string
	InstNumber  int
	IsRAC       bool
	InstPattern string // 实例名模式，如 "racdb", "oradb"
}

// RAC组结构
type RACGroup struct {
	Date       string
	DBName     string
	Files      []XMLFileInfo
	OutputName string
}

// MergeXMLFiles 合并XML文件的主函数
//
// 文件合并规则：
// 1. 单实例文件处理：
//   - XML内容不包含 <DBMAA>RAC</DBMAA> 标签，或
//   - 实例名不以数字结尾（如 oradb, mydb），或
//   - 没有相同日期、数据库名和实例名模式的其他文件
//   - 输出格式：原文件名.S.xml
//
// 2. RAC文件合并：
//   - XML内容必须包含 <DBMAA>RAC</DBMAA> 标签
//   - 实例名必须以数字结尾（如 racdb1, racdb2, racdb21, racdb22）
//   - 相同日期、数据库名和实例名模式的文件会被合并
//   - 实例名模式：racdb1, racdb2 属于 "racdb" 模式
//   - 实例名模式：racdb21, racdb22 属于 "racdb" 模式（但数字不同，不会合并）
//   - 输出格式：日期_主机名1.主机名2_数据库名.R.xml
//
// 3. 文件命名格式：yyyymmdd_hostname_dbname_instname.xml
func MergeXMLFiles() error {
	// 输入和输出目录
	inputDir := "xmlfile/input_xml"
	outputDir := "xmlfile/output_xml"

	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 获取所有XML文件
	xmlFiles, err := getXMLFiles(inputDir)
	if err != nil {
		return fmt.Errorf("获取XML文件失败: %v", err)
	}

	if len(xmlFiles) == 0 {
		log.Println("未找到XML文件")
		return nil
	}

	log.Printf("找到 %d 个XML文件", len(xmlFiles))

	// 解析文件信息
	fileInfos := parseFileInfos(xmlFiles)

	// 分组RAC文件
	racGroups := groupRACFiles(fileInfos)

	// 处理单实例文件
	processSingleInstanceFiles(fileInfos, outputDir)

	// 合并RAC文件
	for _, group := range racGroups {
		if len(group.Files) > 1 {
			log.Printf("合并RAC组: %s, 包含 %d 个文件", group.OutputName, len(group.Files))
			if err := mergeRACFiles(group, outputDir); err != nil {
				log.Printf("合并RAC组 %s 失败: %v", group.OutputName, err)
			}
		}
	}

	log.Println("XML文件处理完成")
	return nil
}

// 获取XML文件列表
func getXMLFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理文件，不处理目录
		if info.IsDir() {
			return nil
		}

		// 检查是否是XML文件
		if strings.HasSuffix(strings.ToLower(info.Name()), ".xml") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// 解析文件信息
func parseFileInfos(files []string) []XMLFileInfo {
	var fileInfos []XMLFileInfo

	for _, file := range files {
		fileName := filepath.Base(file)
		info := parseFileName(fileName)
		if info != nil {
			info.FilePath = file
			info.FileName = fileName

			// 检查XML内容是否包含RAC字样
			isRAC, err := checkRACInXML(file)
			if err != nil {
				log.Printf("检查文件 %s 的RAC信息失败: %v", fileName, err)
				isRAC = false
			}
			info.IsRAC = isRAC

			// 提取实例名模式
			info.InstPattern = extractInstPattern(info.InstName)

			fileInfos = append(fileInfos, *info)
		}
	}

	return fileInfos
}

// 解析文件名
func parseFileName(fileName string) *XMLFileInfo {
	// 匹配格式: yyyymmdd_hostname_dbname_instname.xml
	re := regexp.MustCompile(`^(\d{8})_([^_]+)_([^_]+)_([^_]+)\.xml$`)
	matches := re.FindStringSubmatch(fileName)

	if len(matches) != 5 {
		return nil
	}

	// 提取实例号（最后两位数字）
	instName := matches[4]
	instNumber := extractInstNumber(instName)

	return &XMLFileInfo{
		Date:       matches[1],
		Hostname:   matches[2],
		DBName:     matches[3],
		InstName:   instName,
		InstNumber: instNumber,
	}
}

// 提取实例号
func extractInstNumber(instName string) int {
	// 提取最后两位数字
	re := regexp.MustCompile(`(\d{1,2})$`)
	matches := re.FindStringSubmatch(instName)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 999 // 默认值，表示无法解析
}

// 检查XML文件中是否包含RAC字样
// 通过正则表达式查找 <DBMAA>RAC</DBMAA> 标签
// 支持标签内可能有空白字符的情况
func checkRACInXML(filePath string) (bool, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// 查找 <DBMAA> 标签并检查是否包含 RAC 字样
	// 匹配格式：<DBMAA>RAC</DBMAA> 或 <DBMAA> RAC </DBMAA>
	re := regexp.MustCompile(`<DBMAA>\s*RAC\s*</DBMAA>`)
	return re.Match(content), nil
}

// 提取实例名模式
// 从实例名中提取数字前的部分作为模式
// 例如：racdb1 -> "racdb", racdb21 -> "racdb", oradb -> "oradb"
func extractInstPattern(instName string) string {
	// 提取数字前的部分作为模式
	re := regexp.MustCompile(`^(.+?)\d+$`)
	matches := re.FindStringSubmatch(instName)
	if len(matches) > 1 {
		return matches[1]
	}
	// 如果没有数字结尾，返回原实例名
	return instName
}

// 分组RAC文件
// 将符合RAC条件的文件按日期、数据库名和实例名模式进行分组
// RAC条件：XML内容包含RAC字样 且 实例名以数字结尾
func groupRACFiles(fileInfos []XMLFileInfo) []RACGroup {
	groups := make(map[string]RACGroup)

	for _, info := range fileInfos {
		// 只有XML内容包含RAC字样且实例名以数字结尾的文件才参与RAC分组
		if !info.IsRAC || !isNumericEnding(info.InstName) {
			continue
		}

		// 使用日期、数据库名和实例名模式作为分组键
		// 例如：20250907_racdb_racdb 表示2025年9月7日racdb数据库的racdb模式实例
		key := fmt.Sprintf("%s_%s_%s", info.Date, info.DBName, info.InstPattern)

		group, exists := groups[key]
		if !exists {
			group = RACGroup{
				Date:   info.Date,
				DBName: info.DBName,
				Files:  []XMLFileInfo{},
			}
		}

		group.Files = append(group.Files, info)
		groups[key] = group
	}

	// 转换为切片并生成输出文件名
	var result []RACGroup
	for _, group := range groups {
		if len(group.Files) > 1 {
			// 按实例号排序，确保合并时顺序正确
			sort.Slice(group.Files, func(i, j int) bool {
				return group.Files[i].InstNumber < group.Files[j].InstNumber
			})

			// 生成输出文件名：日期_主机名1.主机名2_数据库名.R.xml
			// 例如：20250907_rac19c1.rac19c2_racdb.R.xml
			var hostnames []string
			for _, file := range group.Files {
				hostnames = append(hostnames, file.Hostname)
			}
			hostnameStr := strings.Join(hostnames, ".")
			group.OutputName = fmt.Sprintf("%s_%s_%s.R.xml", group.Date, hostnameStr, group.DBName)
		}
		result = append(result, group)
	}

	return result
}

// 检查实例名是否以数字结尾
// 用于判断实例是否为RAC环境（如 racdb1, racdb2, racdb21）
func isNumericEnding(instName string) bool {
	re := regexp.MustCompile(`\d+$`)
	return re.MatchString(instName)
}

// 处理单实例文件
// 将所有不符合RAC合并条件的文件作为单实例处理
func processSingleInstanceFiles(fileInfos []XMLFileInfo, outputDir string) {
	for _, info := range fileInfos {
		// 检查是否是单实例：
		// 1. XML内容不包含RAC字样，或
		// 2. 实例名不以数字结尾，或
		// 3. 没有相同日期、数据库名和实例名模式的其他文件
		isSingle := !info.IsRAC || !isNumericEnding(info.InstName)

		if !isSingle {
			// 检查是否有相同模式的其他文件
			// 如果存在相同模式的其他文件，则不是单实例
			for _, other := range fileInfos {
				if other.Date == info.Date &&
					other.DBName == info.DBName &&
					other.InstPattern == info.InstPattern &&
					other.FilePath != info.FilePath {
					isSingle = false
					break
				}
			}
		}

		if isSingle {
			// 生成单实例输出文件名（添加_S.xml后缀）
			// 例如：20250907_myzdb100_oradb_oradb.xml -> 20250907_myzdb100_oradb_oradb.S.xml
			baseName := strings.TrimSuffix(info.FileName, ".xml")
			outputFileName := baseName + ".S.xml"
			outputPath := filepath.Join(outputDir, outputFileName)
			if err := copyFile(info.FilePath, outputPath); err != nil {
				log.Printf("复制单实例文件 %s 失败: %v", outputFileName, err)
			} else {
				log.Printf("复制单实例文件: %s -> %s", info.FileName, outputFileName)
			}
		}
	}
}

// 合并RAC文件
// 将RAC组中的多个XML文件合并为一个文件
// 合并策略：以第一个文件为主文件，将其他文件的NODE1内容复制为NODE2、NODE3等
func mergeRACFiles(group RACGroup, outputDir string) error {
	if len(group.Files) < 2 {
		return fmt.Errorf("RAC组文件数量不足")
	}

	// 以第一个文件为主文件（NODE1）
	mainFile := group.Files[0]
	mainDoc := etree.NewDocument()
	if err := mainDoc.ReadFromFile(mainFile.FilePath); err != nil {
		return fmt.Errorf("读取主文件失败: %v", err)
	}

	// 获取EACHK根元素
	eachk := mainDoc.FindElement("./EACHK")
	if eachk == nil {
		return fmt.Errorf("未找到EACHK元素")
	}

	// 获取TAG0和TAG2元素
	tag0 := eachk.FindElement("./TAG0")
	tag2 := eachk.FindElement("./TAG2")
	if tag0 == nil || tag2 == nil {
		return fmt.Errorf("未找到TAG0或TAG2元素")
	}

	// 处理其他文件，将其NODE1内容复制为NODE2、NODE3等
	for i, file := range group.Files[1:] {
		nodeNum := i + 2 // NODE2, NODE3, ...

		// 读取文件
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(file.FilePath); err != nil {
			log.Printf("读取文件 %s 失败: %v", file.FileName, err)
			continue
		}

		// 处理TAG0：将源文件的NODE1复制为主文件的NODE2、NODE3等
		sourceTag0 := doc.FindElement("./EACHK/TAG0/NODE1")
		if sourceTag0 != nil {
			// 创建新的NODE元素
			newNode := tag0.CreateElement(fmt.Sprintf("NODE%d", nodeNum))
			// 复制所有子元素
			for _, child := range sourceTag0.ChildElements() {
				newNode.AddChild(child.Copy())
			}
		}

		// 处理TAG2：将源文件的NODE1复制为主文件的NODE2、NODE3等
		sourceTag2 := doc.FindElement("./EACHK/TAG2/NODE1")
		if sourceTag2 != nil {
			// 创建新的NODE元素
			newNode := tag2.CreateElement(fmt.Sprintf("NODE%d", nodeNum))
			// 复制所有子元素
			for _, child := range sourceTag2.ChildElements() {
				newNode.AddChild(child.Copy())
			}
		}
	}

	// 保存合并后的文件
	outputPath := filepath.Join(outputDir, group.OutputName)
	mainDoc.Indent(2)
	if err := mainDoc.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("保存合并文件失败: %v", err)
	}

	log.Printf("成功合并RAC文件: %s", group.OutputName)
	return nil
}

// 复制文件
// 用于将单实例文件复制到输出目录
func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
