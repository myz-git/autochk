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
	FilePath   string
	FileName   string
	Date       string
	Hostname   string
	DBName     string
	InstName   string
	InstNumber int
}

// RAC组结构
type RACGroup struct {
	Date       string
	DBName     string
	Files      []XMLFileInfo
	OutputName string
}

// MergeXMLFiles 合并XML文件的主函数
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

// 分组RAC文件
func groupRACFiles(fileInfos []XMLFileInfo) []RACGroup {
	groups := make(map[string]RACGroup)

	for _, info := range fileInfos {
		key := fmt.Sprintf("%s_%s", info.Date, info.DBName)

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
			// 按实例号排序
			sort.Slice(group.Files, func(i, j int) bool {
				return group.Files[i].InstNumber < group.Files[j].InstNumber
			})

			// 生成输出文件名
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

// 处理单实例文件
func processSingleInstanceFiles(fileInfos []XMLFileInfo, outputDir string) {
	for _, info := range fileInfos {
		// 检查是否是单实例（没有相同日期和数据库名的其他文件）
		isSingle := true
		for _, other := range fileInfos {
			if other.Date == info.Date && other.DBName == info.DBName && other.FilePath != info.FilePath {
				isSingle = false
				break
			}
		}

		if isSingle {
			// 生成单实例输出文件名（添加_S.xml后缀）
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
func mergeRACFiles(group RACGroup, outputDir string) error {
	if len(group.Files) < 2 {
		return fmt.Errorf("RAC组文件数量不足")
	}

	// 以第一个文件为主文件
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

	// 处理其他文件
	for i, file := range group.Files[1:] {
		nodeNum := i + 2 // NODE2, NODE3, ...

		// 读取文件
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(file.FilePath); err != nil {
			log.Printf("读取文件 %s 失败: %v", file.FileName, err)
			continue
		}

		// 处理TAG0
		sourceTag0 := doc.FindElement("./EACHK/TAG0/NODE1")
		if sourceTag0 != nil {
			// 创建新的NODE元素
			newNode := tag0.CreateElement(fmt.Sprintf("NODE%d", nodeNum))
			// 复制所有子元素
			for _, child := range sourceTag0.ChildElements() {
				newNode.AddChild(child.Copy())
			}
		}

		// 处理TAG2
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
