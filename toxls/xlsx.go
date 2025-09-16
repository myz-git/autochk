package toxls

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time" // 添加time包导入

	"github.com/xuri/excelize/v2"
)

// 定义公共的单元格样式
func getCellStyles(f *excelize.File) (styleB, styleR, styleG int) {
	styleB, _ = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#00B0F0"}, Pattern: 1},
		Font: &excelize.Font{
			Family: "Cascadia Code Light",
			Size:   8,
		},
		Alignment: &excelize.Alignment{
			Horizontal:      "left",
			Vertical:        "center",
			WrapText:        true,
			ShrinkToFit:     true,
			JustifyLastLine: true,
		},
	})
	styleR, _ = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#EBA8A5"}, Pattern: 1},
		Font: &excelize.Font{
			Family: "Cascadia Code Light",
			Size:   8,
		},
		Alignment: &excelize.Alignment{
			Horizontal:      "left",
			Vertical:        "center",
			WrapText:        true,
			ShrinkToFit:     true,
			JustifyLastLine: true,
		},
	})
	styleG, _ = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#B3DEE3"}, Pattern: 1},
		Font: &excelize.Font{
			Family: "Cascadia Code Light",
			Size:   8,
		},
		Alignment: &excelize.Alignment{
			Horizontal:      "left",
			Vertical:        "center",
			WrapText:        true,
			ShrinkToFit:     true,
			JustifyLastLine: true,
		},
	})
	return
}

func Xlsx(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool, custNm string, reportType string) {
	// 确保report目录存在
	if err := os.MkdirAll("report", 0755); err != nil {
		utils.LogErrorf("创建report目录失败: %v", err)
		return
	}

	// 确定输出文件名
	var newfnm string
	if sglf {
		if reportType == "deep" {
			newfnm = "report/" + xlsnm + "_Deep.xlsx"
		} else {
			newfnm = "report/" + xlsnm + ".xlsx"
		}
	} else {
		if reportType == "deep" {
			newfnm = "report/HealthCheckReport.ALL_Deep.xlsx"
		} else {
			newfnm = "report/HealthCheckReport.ALL.xlsx"
		}
	}

	// 加载模板文件
	tpl := "local/HealthReport.xlsx"
	if reportType == "deep" {
		tpl = "local/HealthReportDeep.xlsx"
	}
	f, err := excelize.OpenFile(tpl)
	if err != nil {
		utils.LogErrorf("打开模板文件失败: %v", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			utils.LogWarnf("关闭文件失败: %v", err)
		}
	}()

	// 填充 HealthReport Sheet 的 Server Info 和 DataBase Info
	PutSht_INFO(f, osshts, dbshtp, colcnt, custNm)

	// 填充 OS Sheet - 支持多节点
	PutSht_OS(f, osshts, summaryEntries, colcnt)

	// 填充 DB Sheet
	PutSht_DB(f, dbshtp, osshts, summaryEntries, colcnt)

	// 填充 Inst Sheet
	PutSht_Inst(f, instshts, summaryEntries)

	// 填充 HealthReport Sheet 的 Issue Summary
	PutSht_Summary(f, summaryEntries)

	// 填充 HealthReport Sheet 的 Issue List
	PutSht_Issuelist(f, summaryEntries)

	// 保存文件
	if err := f.SaveAs(newfnm); err != nil {
		utils.LogErrorf("保存文件失败: %v", err)
		return
	}
}

func PutSht_INFO(f *excelize.File, osshts *[]structs.OsShts, dbshtp *structs.DbSht, colcnt int, custNm string) {
	shnm := "HealthReport"

	// 分开填写：客户名称填写到C2，报告日期填写到K2和H80
	currentDate := time.Now().Format("2006-01-02")

	// 客户名称填写到C2，格式："客户: custNm"
	f.SetCellStr(shnm, "C2", fmt.Sprintf("客户: %s", custNm))

	// 报告日期填写到K2和H80，格式："报告日期: currentDate"
	f.SetCellStr(shnm, "K2", fmt.Sprintf("报告日期: %s", currentDate))
	f.SetCellStr(shnm, "H80", fmt.Sprintf("报告日期: %s", currentDate))

	// 填充 Server Info - 支持多节点信息
	if len(*osshts) > 0 {
		// 使用第一个节点作为主要信息
		firstOs := (*osshts)[0]
		f.SetCellStr(shnm, "F4", firstOs.Hostname.Contents)
		f.SetCellStr(shnm, "F5", firstOs.Ipaddr.Contents)
		f.SetCellStr(shnm, "F6", firstOs.Os.Contents)
		f.SetCellStr(shnm, "F7", firstOs.Relver.Contents)
		f.SetCellStr(shnm, "F8", firstOs.Cpu_model.Contents)
		f.SetCellStr(shnm, "F9", firstOs.Memtotal.Contents)
		f.SetCellStr(shnm, "F10", firstOs.Machine_platform.Contents)
	}

	// 填充 DataBase Info
	f.SetCellStr(shnm, "K4", dbshtp.Dbname.Contents)
	f.SetCellStr(shnm, "K5", dbshtp.Dbver.Contents)
	f.SetCellStr(shnm, "K6", dbshtp.Dbmaa.Contents)
	f.SetCellStr(shnm, "K7", dbshtp.Dbrole.Contents)
	f.SetCellStr(shnm, "K8", dbshtp.Logmode.Contents)
	f.SetCellStr(shnm, "K9", dbshtp.Dblang.Contents)
	f.SetCellStr(shnm, "K10", dbshtp.Dbcursize.Contents)
}

// getCommentForField 根据字段名获取对应的注释内容
func getCommentForField(fieldName string, summaryEntries *structs.SummaryEntries) string {
	utils.LogDebugf("查找字段 %s 的注释，SummaryEntries 总数: %d", fieldName, len(summaryEntries.Entries))

	// 动态生成检查项名（将字段名转换为大写）
	checkName := strings.ToUpper(fieldName)
	utils.LogDebugf("字段 %s 映射到检查项: %s", fieldName, checkName)

	// 在SummaryEntries中查找对应的检查项
	for i, entry := range summaryEntries.Entries {
		utils.LogDebugf("检查项 %d: Title=%s, Nm=%s, Category=%s", i+1, entry.Title, entry.Nm, entry.Category)
		// 直接匹配Nm字段
		if entry.Nm == checkName {
			utils.LogDebugf("找到匹配的检查项: %s", checkName)
			var comments []string
			if len(entry.Severe) > 0 {
				comments = append(comments, "严重问题:\n"+strings.Join(entry.Severe, "\n"))
				utils.LogDebugf("严重问题数量: %d", len(entry.Severe))
			}
			if len(entry.Moderate) > 0 {
				comments = append(comments, "普通问题:\n"+strings.Join(entry.Moderate, "\n"))
				utils.LogDebugf("一般问题数量: %d", len(entry.Moderate))
			}
			if len(entry.Minor) > 0 {
				comments = append(comments, "轻微问题:\n"+strings.Join(entry.Minor, "\n"))
				utils.LogDebugf("轻微问题数量: %d", len(entry.Minor))
			}
			if len(comments) > 0 {
				result := strings.Join(comments, "\n\n")
				utils.LogDebugf("生成的注释内容: %s", result)
				return result
			} else {
				utils.LogDebugf("检查项 %s 没有问题描述", checkName)
			}
		}
	}
	utils.LogDebugf("未找到字段 %s 对应的检查项", fieldName)
	return ""
}

// anchorRow finds the sheet and row index for a named range (anchor at A column).
// name is the defined name (e.g., HOSTNAME). If not found, ok=false.
func anchorRow(f *excelize.File, name string) (string, int, bool) {
	for _, dn := range f.GetDefinedName() {
		if dn.Name == name {
			// dn.RefersTo example: =OS!$A$12 or 'OS'!$A$12
			refers := dn.RefersTo
			if len(refers) == 0 {
				return "", 0, false
			}
			// remove leading '='
			if refers[0] == '=' {
				refers = refers[1:]
			}
			// split sheet and cell by '!'
			parts := strings.Split(refers, "!")
			if len(parts) != 2 {
				return "", 0, false
			}
			sheet := strings.Trim(parts[0], "'")
			cell := parts[1]
			// get row from A1 notation
			_, row, err := excelize.CellNameToCoordinates(cell)
			if err != nil {
				return "", 0, false
			}
			return sheet, row, true
		}
	}
	return "", 0, false
}

// colByNode returns the column letter for node index: 0->C, 1->D, ...
func colByNode(nodeIndex int) string {
	return fmt.Sprintf("%c", 'C'+nodeIndex)
}

func PutSht_OS(f *excelize.File, osshts *[]structs.OsShts, summaryEntries *structs.SummaryEntries, colcnt int) {
	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 遍历所有OS节点
	for nodeIndex, ossht := range *osshts {
		col := colByNode(nodeIndex)

		// 使用反射遍历结构体字段
		v := reflect.ValueOf(ossht)
		t := reflect.TypeOf(ossht)

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name

			// 跳过NodeID字段
			if fieldName == "NodeID" {
				// 直接处理NodeID
				shnm, row, ok := anchorRow(f, "NODEID")
				if ok {
					cell := fmt.Sprintf("%s%d", col, row)
					f.SetCellStr(shnm, cell, ossht.NodeID)
				}
				continue
			}

			// 生成nm（将字段名转换为大写）
			nm := strings.ToUpper(fieldName)

			// 使用anchorRow定位
			shnm, row, ok := anchorRow(f, nm)
			if !ok {
				// 命名区域不存在：可能是 basic 模板缺少 deep 项，跳过
				continue
			}

			// 获取字段内容
			if tpstrc, ok := field.Interface().(structs.Tpstrc); ok {
				cell := fmt.Sprintf("%s%d", col, row)
				content := tpstrc.Contents
				alarm := tpstrc.Alarm

				// 写值
				f.SetCellStr(shnm, cell, content)

				// 注释
				if alarm != "" {
					comment := getCommentForField(fieldName, summaryEntries)
					if comment != "" {
						_ = f.AddComment(shnm, excelize.Comment{Cell: cell, Text: comment, Author: "健康检查系统", Width: 300, Height: 100})
					}
				}

				// 分数（B列，按锚点行）
				bCell := fmt.Sprintf("B%d", row)
				var score int
				switch alarm {
				case "R":
					score = 0
				case "B":
					score = 5
				case "G":
					score = 8
				default:
					score = 10
				}
				f.SetCellInt(shnm, bCell, int64(score))

				// 样式
				if alarm == "R" {
					f.SetCellStyle(shnm, cell, cell, styleR)
				} else if alarm == "B" {
					f.SetCellStyle(shnm, cell, cell, styleB)
				} else if alarm == "G" {
					f.SetCellStyle(shnm, cell, cell, styleG)
				}
			}
		}
	}
}

func PutSht_DB(f *excelize.File, dbshtp *structs.DbSht, osshts *[]structs.OsShts, summaryEntries *structs.SummaryEntries, colcnt int) {
	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 只在第一列（C列）填充DB信息（按设计 DB 只对应 NODE1）
	col := "C"

	// 使用反射遍历结构体字段
	v := reflect.ValueOf(dbshtp).Elem()
	t := reflect.TypeOf(dbshtp).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		// 跳过NodeID字段
		if fieldName == "NodeID" {
			// 直接处理NodeID
			shnm, row, ok := anchorRow(f, "NODEID")
			if ok {
				cell := fmt.Sprintf("%s%d", col, row)
				f.SetCellStr(shnm, cell, dbshtp.NodeID)
			}
			continue
		}

		// 生成nm（将字段名转换为大写）
		nm := strings.ToUpper(fieldName)

		// 使用anchorRow定位
		shnm, row, ok := anchorRow(f, nm)
		if !ok {
			continue
		}

		// 获取字段内容
		if tpstrc, ok := field.Interface().(structs.Tpstrc); ok {
			cell := fmt.Sprintf("%s%d", col, row)
			content := tpstrc.Contents
			alarm := tpstrc.Alarm

			// 写值
			f.SetCellStr(shnm, cell, content)

			// 注释
			if alarm != "" {
				comment := getCommentForField(fieldName, summaryEntries)
				if comment != "" {
					_ = f.AddComment(shnm, excelize.Comment{Cell: cell, Text: comment, Author: "健康检查系统", Width: 300, Height: 100})
				}
			}

			// 分数（B列）
			bCell := fmt.Sprintf("B%d", row)
			var score int
			switch alarm {
			case "R":
				score = 0
			case "B":
				score = 5
			case "G":
				score = 8
			default:
				score = 10
			}
			f.SetCellInt(shnm, bCell, int64(score))

			// 样式
			if alarm == "R" {
				f.SetCellStyle(shnm, cell, cell, styleR)
			} else if alarm == "B" {
				f.SetCellStyle(shnm, cell, cell, styleB)
			} else if alarm == "G" {
				f.SetCellStyle(shnm, cell, cell, styleG)
			}
		}
	}
}

func PutSht_Summary(f *excelize.File, summaryEntries *structs.SummaryEntries) {
	shnm := "HealthReport"

	// ——关键：设置工作簿计算属性——
	auto := "auto"
	yes := true
	_ = f.SetCalcProps(&excelize.CalcPropsOptions{
		CalcMode:       &auto, // 自动计算
		FullCalcOnLoad: &yes,  // 打开时强制全量重算
		CalcOnSave:     &yes,  // （可选）保存时也计算
	})

	// 原有统计逻辑…
	categoryStats := make(map[string]struct{ severe, moderate, minor int })
	for _, entry := range summaryEntries.Entries {
		s := categoryStats[entry.Category]
		s.severe += len(entry.Severe)
		s.moderate += len(entry.Moderate)
		s.minor += len(entry.Minor)
		categoryStats[entry.Category] = s
	}

	categories := []string{"主机系统", "数据库分析", "实例分析", "数据库性能", "数据库集群", "DataGuard", "数据库备份", "数据库安全", "软件使用", "其他项检查"}
	rowIndex := 15
	for _, c := range categories {
		s := categoryStats[c]
		_ = f.SetCellStr(shnm, fmt.Sprintf("C%d", rowIndex), c)
		_ = f.SetCellInt(shnm, fmt.Sprintf("D%d", rowIndex), int64(s.severe))
		_ = f.SetCellInt(shnm, fmt.Sprintf("E%d", rowIndex), int64(s.moderate))
		_ = f.SetCellInt(shnm, fmt.Sprintf("F%d", rowIndex), int64(s.minor))
		rowIndex++
	}

	// 写公式（注意使用英文逗号 ,）
	// _ = f.SetCellFormula(shnm, "H13", "=\"Health status (\"&SUM($D$15:$F$24)&\" of \"&85&\")\"")
	_ = f.SetCellFormula(shnm, "L14", "=1-SUM($D$15:$F$24)/90")

	// ——关键：更新/清理缓存值（避免 Excel 以为不需要重算）——
	_ = f.UpdateLinkedValue()

	// （可选）把 L14 设成百分比格式
	// style, _ := f.NewStyle(&excelize.Style{NumFmt: 10}) // 10 = 0%
	// _ = f.SetCellStyle(shnm, "L14", "L14", style)
}

func PutSht_Issuelist(f *excelize.File, summaryEntries *structs.SummaryEntries) {
	shnm := "HealthReport"

	// 调试输出：显示所有SummaryEntry的内容
	utils.LogDebugf("=== SummaryEntry 调试信息 ===")
	utils.LogDebugf("总共有 %d 个 SummaryEntry", len(summaryEntries.Entries))

	for i, entry := range summaryEntries.Entries {
		utils.LogDebugf("--- SummaryEntry %d ---", i+1)
		utils.LogDebugf("Category: %s", entry.Category)
		utils.LogDebugf("Nm: %s", entry.Nm)
		utils.LogDebugf("Title: %s", entry.Title)
		utils.LogDebugf("Desc: %s", entry.Desc)
		utils.LogDebugf("Severe 问题数量: %d", len(entry.Severe))
		if len(entry.Severe) > 0 {
			utils.LogDebugf("Severe 问题: %v", entry.Severe)
		}
		utils.LogDebugf("Moderate 问题数量: %d", len(entry.Moderate))
		if len(entry.Moderate) > 0 {
			utils.LogDebugf("Moderate 问题: %v", entry.Moderate)
		}
		utils.LogDebugf("Minor 问题数量: %d", len(entry.Minor))
		if len(entry.Minor) > 0 {
			utils.LogDebugf("Minor 问题: %v", entry.Minor)
		}
	}
	utils.LogDebugf("=== 调试信息结束 ===")

	// 定义问题级别优先级（数字越小优先级越高）
	severityOrder := map[string]int{
		"严重": 1,
		"普通": 2,
		"轻微": 3,
		"正常": 4,
	}

	// 定义问题类别优先级（数字越小优先级越高）
	categoryOrder := map[string]int{
		"主机系统":      1,
		"数据库分析":     2,
		"实例分析":      3,
		"数据库性能":     4,
		"数据库集群":     5,
		"DataGuard": 6,
		"数据库备份":     7,
		"数据库安全":     8,
		"软件使用":      9,
		"其他项检查":     10,
	}

	// 创建问题项结构体用于排序
	type ProblemItem struct {
		Category    string
		Title       string
		Nm          string
		Severity    string
		Score       int
		Description string
	}

	var problemItems []ProblemItem

	// 收集所有问题项
	for i := range summaryEntries.Entries {
		entry := &summaryEntries.Entries[i]

		// 收集严重问题
		for _, problem := range entry.Severe {
			problemItems = append(problemItems, ProblemItem{
				Category:    entry.Category,
				Title:       entry.Title,
				Nm:          entry.Nm,
				Severity:    "严重",
				Score:       0,
				Description: problem,
			})
		}

		// 收集普通问题
		for _, problem := range entry.Moderate {
			problemItems = append(problemItems, ProblemItem{
				Category:    entry.Category,
				Title:       entry.Title,
				Nm:          entry.Nm,
				Severity:    "普通",
				Score:       5,
				Description: problem,
			})
		}

		// 收集轻微问题
		for _, problem := range entry.Minor {
			problemItems = append(problemItems, ProblemItem{
				Category:    entry.Category,
				Title:       entry.Title,
				Nm:          entry.Nm,
				Severity:    "轻微",
				Score:       8,
				Description: problem,
			})
		}

		// 如果没有问题，添加正常项
		if len(entry.Severe) == 0 && len(entry.Moderate) == 0 && len(entry.Minor) == 0 {
			problemItems = append(problemItems, ProblemItem{
				Category:    entry.Category,
				Title:       entry.Title,
				Nm:          entry.Nm,
				Severity:    "正常",
				Score:       10,
				Description: "检查通过，无问题",
			})
		}
	}

	// 对问题项进行排序
	// 首先按问题重要程度排序，然后按问题类别排序
	for i := 0; i < len(problemItems)-1; i++ {
		for j := i + 1; j < len(problemItems); j++ {
			// 比较问题重要程度
			severityI := severityOrder[problemItems[i].Severity]
			severityJ := severityOrder[problemItems[j].Severity]

			if severityI > severityJ {
				// 交换位置
				problemItems[i], problemItems[j] = problemItems[j], problemItems[i]
			} else if severityI == severityJ {
				// 问题重要程度相同，按问题类别排序
				categoryI := categoryOrder[problemItems[i].Category]
				categoryJ := categoryOrder[problemItems[j].Category]

				if categoryI > categoryJ {
					// 交换位置
					problemItems[i], problemItems[j] = problemItems[j], problemItems[i]
				}
			}
		}
	}

	// 填充 Issue List（从 B29 开始）
	rowIndex := 29
	itemIndex := 1

	// 按照排序后的顺序填充
	for _, item := range problemItems {
		f.SetCellInt(shnm, fmt.Sprintf("B%d", rowIndex), int64(itemIndex))  // B列：序号
		f.SetCellStr(shnm, fmt.Sprintf("C%d", rowIndex), item.Category)     // C列：问题类型
		f.SetCellStr(shnm, fmt.Sprintf("D%d", rowIndex), item.Title)        // D列：检查项
		f.SetCellInt(shnm, fmt.Sprintf("F%d", rowIndex), int64(item.Score)) // F列：检查结果分数
		f.SetCellStr(shnm, fmt.Sprintf("G%d", rowIndex), item.Severity)     // G列：问题级别
		f.SetCellStr(shnm, fmt.Sprintf("H%d", rowIndex), item.Description)  // H列：检查项说明

		// L列：添加超链接到具体问题位置
		hyperlink := getHyperlinkForProblem(item.Nm, item.Category, item.Description, f)
		if hyperlink != "" {
			f.SetCellFormula(shnm, fmt.Sprintf("L%d", rowIndex), hyperlink)
		}

		rowIndex++
		itemIndex++
	}

	// 删除最后一个问题所在行到75行之间的空白行
	lastProblemRow := rowIndex - 1 // 最后一个问题所在的行号
	utils.LogDebugf("最后一个问题所在行: %d", lastProblemRow)

	if lastProblemRow < 75 {
		// 需要删除的行范围：lastProblemRow+1 到 75
		deleteStartRow := lastProblemRow + 1
		deleteEndRow := 75

		utils.LogDebugf("删除空白行范围: %d 到 %d", deleteStartRow, deleteEndRow)

		// 删除空白行
		for i := deleteEndRow; i >= deleteStartRow; i-- {
			err := f.RemoveRow(shnm, i)
			if err != nil {
				utils.LogWarnf("删除第 %d 行失败: %v", i, err)
			}
		}
	}
}

func PutSht_Inst(f *excelize.File, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries) {
	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 遍历所有实例节点
	for nodeIndex, instsht := range *instshts {
		col := colByNode(nodeIndex)

		// 使用反射遍历结构体字段
		v := reflect.ValueOf(instsht)
		t := reflect.TypeOf(instsht)

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name

			// 跳过NodeID字段
			if fieldName == "NodeID" {
				// 直接处理NodeID
				shnm, row, ok := anchorRow(f, "NODEID")
				if ok {
					cell := fmt.Sprintf("%s%d", col, row)
					f.SetCellStr(shnm, cell, instsht.NodeID)
				}
				continue
			}

			// 生成nm（将字段名转换为大写）
			nm := strings.ToUpper(fieldName)

			// 使用anchorRow定位
			shnm, row, ok := anchorRow(f, nm)
			if !ok {
				continue
			}

			// 获取字段内容
			if tpstrc, ok := field.Interface().(structs.Tpstrc); ok {
				cell := fmt.Sprintf("%s%d", col, row)
				content := tpstrc.Contents
				alarm := tpstrc.Alarm

				// 写入值
				f.SetCellStr(shnm, cell, content)

				// 注释
				if alarm != "" {
					comment := getCommentForField(fieldName, summaryEntries)
					if comment != "" {
						_ = f.AddComment(shnm, excelize.Comment{Cell: cell, Text: comment, Author: "健康检查系统", Width: 300, Height: 100})
					}
				}

				// 分数（B列）
				bCell := fmt.Sprintf("B%d", row)
				var score int
				switch alarm {
				case "R":
					score = 0
				case "B":
					score = 5
				case "G":
					score = 8
				default:
					score = 10
				}
				f.SetCellInt(shnm, bCell, int64(score))

				// 样式
				if alarm == "R" {
					f.SetCellStyle(shnm, cell, cell, styleR)
				} else if alarm == "B" {
					f.SetCellStyle(shnm, cell, cell, styleB)
				} else if alarm == "G" {
					f.SetCellStyle(shnm, cell, cell, styleG)
				}
			}
		}
	}
}

// getHyperlinkForProblem 根据检查项Nm生成超链接到对应的工作表位置
// 使用anchorRow函数获取准确的行号，避免硬编码
func getHyperlinkForProblem(nm string, category string, problem string, f *excelize.File) string {
	// 根据Category确定Sheet名称（简化逻辑）
	var sheetName string
	switch category {
	case "主机系统":
		sheetName = "OS"
	case "数据库分析", "数据库性能", "数据库集群", "数据库备份", "数据库安全":
		sheetName = "DB"
	case "实例分析", "DataGuard":
		sheetName = "Inst"
	default:
		sheetName = "OS" // 默认为OS
	}

	// 使用anchorRow函数获取准确的行号
	_, row, ok := anchorRow(f, nm)
	if !ok {
		// 如果找不到命名区域，不生成超链接
		return ""
	}

	// 从问题描述中解析节点信息，确定列号
	// 问题描述格式通常是："问题: NODE1主机,xxx" 或 "问题: NODE2主机,xxx"
	var col string = "C" // 默认C列（NODE1）

	// 检查问题描述中是否包含节点信息
	if strings.Contains(problem, "NODE1") {
		col = "C"
	} else if strings.Contains(problem, "NODE2") {
		col = "D"
	} else if strings.Contains(problem, "NODE3") {
		col = "E"
	} else if strings.Contains(problem, "NODE4") {
		col = "F"
	}

	// 生成超链接公式
	cellRef := fmt.Sprintf("%s%d", col, row)
	return fmt.Sprintf("=HYPERLINK(\"#%s!%s\",\"Detail\")", sheetName, cellRef)
}
