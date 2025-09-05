package toxls

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"strings"

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

func Xlsx(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool) {
	// 确定输出文件名
	var newfnm string
	if sglf {
		newfnm = xlsnm + ".Done.xlsx"
	} else {
		newfnm = "HealthCheckReport.ALLDone.xlsx"
	}

	// 加载模板文件
	f, err := excelize.OpenFile("HealthReport.xlsx")
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
	PutSht_INFO(f, osshts, dbshtp, colcnt)

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

func PutSht_INFO(f *excelize.File, osshts *[]structs.OsShts, dbshtp *structs.DbSht, colcnt int) {
	shnm := "HealthReport"

	// 填充 Server Info - 支持多节点信息
	if len(*osshts) > 0 {
		// 使用第一个节点作为主要信息
		firstOs := (*osshts)[0]
		f.SetCellStr(shnm, "F4", firstOs.Hostname.Contents)
		f.SetCellStr(shnm, "F5", firstOs.Ipaddr.Contents)
		f.SetCellStr(shnm, "F6", firstOs.Os.Contents)
		f.SetCellStr(shnm, "F7", firstOs.Relver.Contents)
		f.SetCellStr(shnm, "F8", firstOs.Cores.Contents)
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

	// 字段名到检查项Nm的映射（直接使用SummaryEntries中的Nm字段）
	fieldToCheckMap := map[string]string{
		// OS字段
		"Hostname":         "HOSTNAME",
		"Ipaddr":           "IPADDR",
		"Os":               "OS",
		"Relver":           "RELVER",
		"Cores":            "CORES",
		"Cpucount":         "CPUCOUNT",
		"Cpumhz":           "CPUMHZ",
		"Memtotal":         "MEMTOTAL",
		"Swaptotal":        "SWAPTOTAL",
		"Osparameter":      "OSPARAMETER",
		"Ulimit":           "ULIMIT",
		"Oslog":            "OSLOG",
		"Filesystem":       "FILESYSTEM",
		"Inodeusage":       "INODEUSAGE",
		"Cpustat":          "CPUSTAT",
		"Memstat":          "MEMSTAT",
		"Iostat":           "IOSTAT",
		"Thpstat":          "THPSTAT",
		"Hugepage":         "HUGEPAGE",
		"Numa":             "NUMA",
		"Ntp":              "NTP",
		"Tmzone":           "TMZONE",
		"Selinux":          "SELINUX",
		"Firewall":         "FIREWALL",
		"Nsswitch":         "NSSWITCH",
		"Lo_mtu":           "LO_MTU",
		"Machine_platform": "MACHINE_PLATFORM",
		"CPU_PERF_MODE":    "CPU_PERF_MODE",
		"NOZEROCONF":       "NOZEROCONF",
		"RPM_PACKAGES":     "RPM_PACKAGES",

		// DB字段
		"Dbname":             "DBNAME",
		"Dbmaa":              "DBMAA",
		"Dbver":              "DBVER",
		"Dbstatus":           "DBSTATUS",
		"Dblang":             "DBLANG",
		"Logmode":            "LOGMODE",
		"Flashback":          "FLASHBACK",
		"Dbcursize":          "DBCURSIZE",
		"Dbf_size":           "DBF_SIZE",
		"Dbf_cnt":            "DBF_CNT",
		"Dbf_stat":           "DBF_STAT",
		"Tmpfile_size":       "TMPFILE_SIZE",
		"Dbtblcount":         "DBTBLCOUNT",
		"Dbrole":             "DBROLE",
		"Dbtbsusage":         "DBTBSUSAGE",
		"Dbcontrolfile":      "DBCONTROLFILE",
		"User_info":          "USER_INFO",
		"User_size":          "USER_SIZE",
		"Tab_info":           "TAB_INFO",
		"Tab_parallel":       "TAB_PARALLEL",
		"Inx_parallel":       "INX_PARALLEL",
		"Invalid_obj":        "INVALID_OBJ",
		"Invalid_inx":        "INVALID_INX",
		"Dbsequence":         "DBSEQUENCE",
		"Db_seq_usage":       "DB_SEQ_USAGE",
		"Dboption":           "DBOPTION",
		"Dbfeatures":         "DBFEATURES",
		"Db_expir_user":      "DB_EXPIR_USER",
		"Db_password_verif":  "DB_PASSWORD_VERIF",
		"Dbdbapriv":          "DBDBAPRIV",
		"Dbsysdba":           "DBSYSDBA",
		"Dbauditsegment":     "DBAUDITSEGMENT",
		"Dbauditcont":        "DBAUDITCONT",
		"Db_Nosys_In_System": "DB_NOSYS_IN_SYSTEM",
		"Userfailedlogin":    "USERFAILEDLOGIN",
		"Dbvirscheck":        "DBVIRSCHECK",
		"Dbscnhealthcheck":   "DBSCNHEALTHCHECK",
		"Dbrmancheck":        "DBRMANCHECK",
		"Crs_stat":           "CRS_STAT",
		"Crs_stat2":          "CRS_STAT2",
		"Ocr_info":           "OCR_INFO",
		"Ocr_bak_check":      "OCR_BAK_CHECK",
		"Asm_usage":          "ASM_USAGE",
		"Asm_offset":         "ASM_OFFSET",

		// INST字段
		"Instname":          "INSTNAME",
		"Loadprofile":       "LOADPROFILE",
		"Instefficiency":    "INSTEFFICIENCY",
		"Topevent":          "TOPEVENT",
		"Topsql_by_ela":     "TOPSQL_BY_ELA",
		"Cursor_share_mem":  "CURSOR_SHARE_MEM",
		"Dbresource":        "DBRESOURCE",
		"Dbpsu":             "DBPSU",
		"Dbpatch":           "DBPATCH",
		"Dblsnrinfo":        "DBLSNRINFO",
		"Dbparameter":       "DBPARAMETER",
		"Db_parameter_file": "DB_PARAMETER_FILE",
		"Dbredocheck":       "DBREDOCHECK",
		"Dbredoswitch":      "DBREDOSWITCH",
		"Recovery_usage":    "RECOVERY_USAGE",
		"Recovery_detail":   "RECOVERY_DETAIL",
		"Dberrlog":          "DBERRLOG",
		"Dbdglagcheck":      "DBDGLAGCHECK",
		"Dbdgerrcheck":      "DBDGERRCHECK",
	}

	checkName := fieldToCheckMap[fieldName]
	utils.LogDebugf("字段 %s 映射到检查项: %s", fieldName, checkName)

	if checkName == "" {
		utils.LogDebugf("字段 %s 没有对应的检查项映射", fieldName)
		return ""
	}

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

func PutSht_OS(f *excelize.File, osshts *[]structs.OsShts, summaryEntries *structs.SummaryEntries, colcnt int) {
	shnm := "OS"

	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 按照OsShts结构体的字段顺序定义
	osFields := []struct {
		fieldName string
		row       int
	}{
		{"NodeID", 1},
		{"Hostname", 2},
		{"Ipaddr", 3},
		{"Os", 4},
		{"Relver", 5},
		{"Cores", 6},
		{"Cpumhz", 7},
		{"Memtotal", 8},
		{"Machine_platform", 9},
		{"Swaptotal", 10},
		{"Osparameter", 11},
		{"Ulimit", 12},
		{"Oslog", 13},
		{"Filesystem", 14},
		{"Inodeusage", 15},
		{"Cpustat", 16},
		{"Memstat", 17},
		{"Iostat", 18},
		{"Thpstat", 19},
		{"Hugepage", 20},
		{"Numa", 21},
		{"Ntp", 22},
		{"Tmzone", 23},
		{"Selinux", 24},
		{"Firewall", 25},
		{"Nsswitch", 26},
		{"Lo_mtu", 27},
		{"CPU_PERF_MODE", 28},
		{"NOZEROCONF", 29},
		{"RPM_PACKAGES", 30},
	}

	// 遍历所有OS节点
	for nodeIndex, ossht := range *osshts {
		// 确定列位置：C列对应NODE1，D列对应NODE2，以此类推
		col := fmt.Sprintf("%c", 'C'+nodeIndex)

		// 填充每个字段
		for _, field := range osFields {
			cell := fmt.Sprintf("%s%d", col, field.row)
			var content string
			var alarm string

			// 根据字段名获取对应的内容和告警级别
			switch field.fieldName {
			case "NodeID":
				content = ossht.NodeID
				alarm = ""
			case "Hostname":
				content = ossht.Hostname.Contents
				alarm = ossht.Hostname.Alarm
			case "Ipaddr":
				content = ossht.Ipaddr.Contents
				alarm = ossht.Ipaddr.Alarm
			case "Os":
				content = ossht.Os.Contents
				alarm = ossht.Os.Alarm
			case "Relver":
				content = ossht.Relver.Contents
				alarm = ossht.Relver.Alarm
			case "Cores":
				content = ossht.Cores.Contents
				alarm = ossht.Cores.Alarm
			case "Cpucount":
				content = ossht.Cpucount.Contents
				alarm = ossht.Cpucount.Alarm
			case "Cpumhz":
				content = ossht.Cpumhz.Contents
				alarm = ossht.Cpumhz.Alarm
			case "Memtotal":
				content = ossht.Memtotal.Contents
				alarm = ossht.Memtotal.Alarm
			case "Swaptotal":
				content = ossht.Swaptotal.Contents
				alarm = ossht.Swaptotal.Alarm
			case "Osparameter":
				content = ossht.Osparameter.Contents
				alarm = ossht.Osparameter.Alarm
			case "Ulimit":
				content = ossht.Ulimit.Contents
				alarm = ossht.Ulimit.Alarm
			case "Oslog":
				content = ossht.Oslog.Contents
				alarm = ossht.Oslog.Alarm
			case "Filesystem":
				content = ossht.Filesystem.Contents
				alarm = ossht.Filesystem.Alarm
			case "Inodeusage":
				content = ossht.Inodeusage.Contents
				alarm = ossht.Inodeusage.Alarm
			case "Cpustat":
				content = ossht.Cpustat.Contents
				alarm = ossht.Cpustat.Alarm
			case "Memstat":
				content = ossht.Memstat.Contents
				alarm = ossht.Memstat.Alarm
			case "Iostat":
				content = ossht.Iostat.Contents
				alarm = ossht.Iostat.Alarm
			case "Thpstat":
				content = ossht.Thpstat.Contents
				alarm = ossht.Thpstat.Alarm
			case "Hugepage":
				content = ossht.Hugepage.Contents
				alarm = ossht.Hugepage.Alarm
			case "Numa":
				content = ossht.Numa.Contents
				alarm = ossht.Numa.Alarm
			case "Ntp":
				content = ossht.Ntp.Contents
				alarm = ossht.Ntp.Alarm
			case "Tmzone":
				content = ossht.Tmzone.Contents
				alarm = ossht.Tmzone.Alarm
			case "Selinux":
				content = ossht.Selinux.Contents
				alarm = ossht.Selinux.Alarm
			case "Firewall":
				content = ossht.Firewall.Contents
				alarm = ossht.Firewall.Alarm
			case "Nsswitch":
				content = ossht.Nsswitch.Contents
				alarm = ossht.Nsswitch.Alarm
			case "Lo_mtu":
				content = ossht.Lo_mtu.Contents
				alarm = ossht.Lo_mtu.Alarm
			case "Machine_platform":
				content = ossht.Machine_platform.Contents
				alarm = ossht.Machine_platform.Alarm
			case "CPU_PERF_MODE":
				content = ossht.CPU_PERF_MODE.Contents
				alarm = ossht.CPU_PERF_MODE.Alarm
			case "NOZEROCONF":
				content = ossht.NOZEROCONF.Contents
				alarm = ossht.NOZEROCONF.Alarm
			case "RPM_PACKAGES":
				content = ossht.RPM_PACKAGES.Contents
				alarm = ossht.RPM_PACKAGES.Alarm
			}

			// 在样式设置前添加调试信息
			utils.LogDebugf("节点: %s, 字段: %s, 告警级别: %s", ossht.NodeID, field.fieldName, alarm)

			// 设置单元格内容
			f.SetCellStr(shnm, cell, content)

			// 添加单元格注释（如果有告警）
			if alarm != "" {
				utils.LogDebugf("尝试为字段 %s 添加注释，告警级别: %s", field.fieldName, alarm)
				comment := getCommentForField(field.fieldName, summaryEntries)
				utils.LogDebugf("获取到的注释内容: %s", comment)
				if comment != "" {
					// 创建带样式的注释
					commentObj := excelize.Comment{
						Cell:   cell,
						Text:   comment,
						Author: "健康检查系统",
						Width:  400, // 注释框宽度
						Height: 100, // 注释框高度
					}

					// 添加调试信息
					utils.LogDebugf("注释对象: Cell=%s, Width=%d, Height=%d, Text长度=%d",
						commentObj.Cell, commentObj.Width, commentObj.Height, len(commentObj.Text))

					err := f.AddComment(shnm, commentObj)
					if err != nil {
						utils.LogErrorf("添加注释失败: %v", err)
					} else {
						utils.LogDebugf("成功添加注释到单元格: %s", cell)
					}
				} else {
					utils.LogDebugf("未找到字段 %s 对应的注释内容", field.fieldName)
				}
			}

			// 设置B列检查结果分数 (从B2开始，B1是标题行，跳过NodeID)
			if field.fieldName != "NodeID" {
				bCell := fmt.Sprintf("B%d", field.row)
				var score int
				switch alarm {
				case "R":
					score = 0 // 严重影响
				case "B":
					score = 5 // 普通影响
				case "G":
					score = 8 // 轻微影响
				default:
					score = 10 // 正常
				}
				f.SetCellInt(shnm, bCell, int64(score))
			}

			// 设置单元格样式（根据告警级别）
			if alarm == "R" {
				f.SetCellStyle(shnm, cell, cell, styleR)
				utils.LogDebugf("应用红色样式到单元格: %s (节点: %s, 字段: %s)", cell, ossht.NodeID, field.fieldName)
			} else if alarm == "B" {
				f.SetCellStyle(shnm, cell, cell, styleB)
				utils.LogDebugf("应用蓝色样式到单元格: %s (节点: %s, 字段: %s)", cell, ossht.NodeID, field.fieldName)
			} else if alarm == "G" {
				f.SetCellStyle(shnm, cell, cell, styleG)
				utils.LogDebugf("应用绿色样式到单元格: %s (节点: %s, 字段: %s)", cell, ossht.NodeID, field.fieldName)
			}
		}
	}
}

func PutSht_DB(f *excelize.File, dbshtp *structs.DbSht, osshts *[]structs.OsShts, summaryEntries *structs.SummaryEntries, colcnt int) {
	shnm := "DB"

	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 按照DbSht结构体的字段顺序定义
	dbFields := []struct {
		fieldName string
		row       int
	}{
		{"NodeID", 1},
		{"Dbname", 2},
		{"Dbmaa", 3},
		{"Dbver", 4},
		{"Dbstatus", 5},
		{"Dblang", 6},
		{"Logmode", 7},
		{"Flashback", 8},
		{"Dbcursize", 9},
		{"Dbf_size", 10},
		{"Dbf_cnt", 11},
		{"Dbf_stat", 12},
		{"Tmpfile_size", 13},
		{"Dbtblcount", 14},
		{"Dbrole", 15},
		{"Dbtbsusage", 16},
		{"Dbcontrolfile", 17},
		{"User_info", 18},
		{"User_size", 19},
		{"Tab_info", 20},
		{"Tab_parallel", 21},
		{"Inx_parallel", 22},
		{"Invalid_obj", 23},
		{"Invalid_inx", 24},
		{"Dbsequence", 25},
		{"Db_seq_usage", 26},
		{"Dboption", 27},
		{"Dbfeatures", 28},
		{"Db_expir_user", 29},
		{"Db_password_verif", 30},
		{"Dbdbapriv", 31},
		{"Dbsysdba", 32},
		{"Dbauditsegment", 33},
		{"Dbauditcont", 34},
		{"Db_Nosys_In_System", 35},
		{"Userfailedlogin", 36},
		{"Dbvirscheck", 37},
		{"Dbscnhealthcheck", 38},
		{"Dbrmancheck", 39},
		{"Crs_stat", 40},
		{"Crs_stat2", 41},
		{"Ocr_info", 42},
		{"Ocr_bak_check", 43},
		{"Asm_usage", 44},
		{"Asm_offset", 45},
	}

	// 只在第一列（C列）填充DB信息，因为DB信息只有一个NODE1
	col := "C"

	// 填充每个字段
	for _, field := range dbFields {
		cell := fmt.Sprintf("%s%d", col, field.row)
		var content string
		var alarm string

		// 根据字段名获取对应的内容和告警级别
		switch field.fieldName {
		case "NodeID":
			content = "NODE1" // DB信息固定为NODE1
			alarm = ""
		case "Dbname":
			content = dbshtp.Dbname.Contents
			alarm = dbshtp.Dbname.Alarm
		case "Dbmaa":
			content = dbshtp.Dbmaa.Contents
			alarm = dbshtp.Dbmaa.Alarm
		case "Dbver":
			content = dbshtp.Dbver.Contents
			alarm = dbshtp.Dbver.Alarm
		case "Dbstatus":
			content = dbshtp.Dbstatus.Contents
			alarm = dbshtp.Dbstatus.Alarm
		case "Dblang":
			content = dbshtp.Dblang.Contents
			alarm = dbshtp.Dblang.Alarm
		case "Logmode":
			content = dbshtp.Logmode.Contents
			alarm = dbshtp.Logmode.Alarm
		case "Flashback":
			content = dbshtp.Flashback.Contents
			alarm = dbshtp.Flashback.Alarm
		case "Dbcursize":
			content = dbshtp.Dbcursize.Contents
			alarm = dbshtp.Dbcursize.Alarm
		case "Dbf_size":
			content = dbshtp.Dbf_size.Contents
			alarm = dbshtp.Dbf_size.Alarm
		case "Dbf_cnt":
			content = dbshtp.Dbf_cnt.Contents
			alarm = dbshtp.Dbf_cnt.Alarm
		case "Dbf_stat":
			content = dbshtp.Dbf_stat.Contents
			alarm = dbshtp.Dbf_stat.Alarm
		case "Tmpfile_size":
			content = dbshtp.Tmpfile_size.Contents
			alarm = dbshtp.Tmpfile_size.Alarm
		case "Dbtblcount":
			content = dbshtp.Dbtblcount.Contents
			alarm = dbshtp.Dbtblcount.Alarm
		case "Dbrole":
			content = dbshtp.Dbrole.Contents
			alarm = dbshtp.Dbrole.Alarm
		case "Dbtbsusage":
			content = dbshtp.Dbtbsusage.Contents
			alarm = dbshtp.Dbtbsusage.Alarm
		case "Dbcontrolfile":
			content = dbshtp.Dbcontrolfile.Contents
			alarm = dbshtp.Dbcontrolfile.Alarm
		case "User_info":
			content = dbshtp.User_info.Contents
			alarm = dbshtp.User_info.Alarm
		case "User_size":
			content = dbshtp.User_size.Contents
			alarm = dbshtp.User_size.Alarm
		case "Tab_info":
			content = dbshtp.Tab_info.Contents
			alarm = dbshtp.Tab_info.Alarm
		case "Tab_parallel":
			content = dbshtp.Tab_parallel.Contents
			alarm = dbshtp.Tab_parallel.Alarm
		case "Inx_parallel":
			content = dbshtp.Inx_parallel.Contents
			alarm = dbshtp.Inx_parallel.Alarm
		case "Invalid_obj":
			content = dbshtp.Invalid_obj.Contents
			alarm = dbshtp.Invalid_obj.Alarm
		case "Invalid_inx":
			content = dbshtp.Invalid_inx.Contents
			alarm = dbshtp.Invalid_inx.Alarm
		case "Dbsequence":
			content = dbshtp.Dbsequence.Contents
			alarm = dbshtp.Dbsequence.Alarm
		case "Db_seq_usage":
			content = dbshtp.Db_seq_usage.Contents
			alarm = dbshtp.Db_seq_usage.Alarm
		case "Dboption":
			content = dbshtp.Dboption.Contents
			alarm = dbshtp.Dboption.Alarm
		case "Dbfeatures":
			content = dbshtp.Dbfeatures.Contents
			alarm = dbshtp.Dbfeatures.Alarm
		case "Db_expir_user":
			content = dbshtp.Db_expir_user.Contents
			alarm = dbshtp.Db_expir_user.Alarm
		case "Db_password_verif":
			content = dbshtp.Db_password_verif.Contents
			alarm = dbshtp.Db_password_verif.Alarm
		case "Dbdbapriv":
			content = dbshtp.Dbdbapriv.Contents
			alarm = dbshtp.Dbdbapriv.Alarm
		case "Dbsysdba":
			content = dbshtp.Dbsysdba.Contents
			alarm = dbshtp.Dbsysdba.Alarm
		case "Dbauditsegment":
			content = dbshtp.Dbauditsegment.Contents
			alarm = dbshtp.Dbauditsegment.Alarm
		case "Dbauditcont":
			content = dbshtp.Dbauditcont.Contents
			alarm = dbshtp.Dbauditcont.Alarm
		case "Db_Nosys_In_System":
			content = dbshtp.Db_Nosys_In_System.Contents
			alarm = dbshtp.Db_Nosys_In_System.Alarm
		case "Userfailedlogin":
			content = dbshtp.Userfailedlogin.Contents
			alarm = dbshtp.Userfailedlogin.Alarm
		case "Dbvirscheck":
			content = dbshtp.Dbvirscheck.Contents
			alarm = dbshtp.Dbvirscheck.Alarm
		case "Dbscnhealthcheck":
			content = dbshtp.Dbscnhealthcheck.Contents
			alarm = dbshtp.Dbscnhealthcheck.Alarm
		case "Dbrmancheck":
			content = dbshtp.Dbrmancheck.Contents
			alarm = dbshtp.Dbrmancheck.Alarm
		case "Crs_stat":
			content = dbshtp.Crs_stat.Contents
			alarm = dbshtp.Crs_stat.Alarm
		case "Crs_stat2":
			content = dbshtp.Crs_stat2.Contents
			alarm = dbshtp.Crs_stat2.Alarm
		case "Ocr_info":
			content = dbshtp.Ocr_info.Contents
			alarm = dbshtp.Ocr_info.Alarm
		case "Ocr_bak_check":
			content = dbshtp.Ocr_bak_check.Contents
			alarm = dbshtp.Ocr_bak_check.Alarm
		case "Asm_usage":
			content = dbshtp.Asm_usage.Contents
			alarm = dbshtp.Asm_usage.Alarm
		case "Asm_offset":
			content = dbshtp.Asm_offset.Contents
			alarm = dbshtp.Asm_offset.Alarm
		}

		// 设置单元格内容
		f.SetCellStr(shnm, cell, content)

		// 添加单元格注释（如果有告警）
		if alarm != "" {
			utils.LogDebugf("尝试为DB字段 %s 添加注释，告警级别: %s", field.fieldName, alarm)
			comment := getCommentForField(field.fieldName, summaryEntries)
			utils.LogDebugf("获取到的DB注释内容: %s", comment)
			if comment != "" {
				// 创建带样式的注释
				commentObj := excelize.Comment{
					Cell:   cell,
					Text:   comment,
					Author: "健康检查系统",
					Width:  400, // 注释框宽度
					Height: 100, // 注释框高度
				}

				// 添加调试信息
				utils.LogDebugf("DB注释对象: Cell=%s, Width=%d, Height=%d, Text长度=%d",
					commentObj.Cell, commentObj.Width, commentObj.Height, len(commentObj.Text))

				err := f.AddComment(shnm, commentObj)
				if err != nil {
					utils.LogErrorf("添加DB注释失败: %v", err)
				} else {
					utils.LogDebugf("成功添加DB注释到单元格: %s", cell)
				}
			} else {
				utils.LogDebugf("未找到DB字段 %s 对应的注释内容", field.fieldName)
			}
		}

		// 设置B列检查结果分数 (从B2开始，B1是标题行，跳过NodeID)
		if field.fieldName != "NodeID" {
			bCell := fmt.Sprintf("B%d", field.row)
			var score int
			switch alarm {
			case "R":
				score = 0 // 严重影响
			case "B":
				score = 5 // 普通影响
			case "G":
				score = 8 // 轻微影响
			default:
				score = 10 // 正常
			}
			f.SetCellInt(shnm, bCell, int64(score))
		}

		// 设置单元格样式（根据告警级别）
		if alarm == "R" {
			f.SetCellStyle(shnm, cell, cell, styleR)
		} else if alarm == "B" {
			f.SetCellStyle(shnm, cell, cell, styleB)
		} else if alarm == "G" {
			f.SetCellStyle(shnm, cell, cell, styleG)
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
	_ = f.SetCellFormula(shnm, "H13", "=\"Health status (\"&SUM($D$15:$F$24)&\" of \"&85&\")\"")
	_ = f.SetCellFormula(shnm, "L14", "=1-SUM($D$15:$F$24)/85")

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
		hyperlink := getHyperlinkForProblem(item.Nm, item.Category, item.Description)
		if hyperlink != "" {
			f.SetCellFormula(shnm, fmt.Sprintf("L%d", rowIndex), hyperlink)
		}

		rowIndex++
		itemIndex++
	}
}

func PutSht_Inst(f *excelize.File, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries) {
	shnm := "Inst" // 使用专门的Inst sheet

	// 定义单元格样式
	styleB, styleR, styleG := getCellStyles(f)

	// 按照InstShts结构体的字段顺序定义
	instFields := []struct {
		fieldName string
		row       int
	}{
		{"NodeID", 1},
		{"Instname", 2},
		{"Loadprofile", 3},
		{"Instefficiency", 4},
		{"Topevent", 5},
		{"Topsql_by_ela", 6},
		{"Cursor_share_mem", 7},
		{"Dbresource", 8},
		{"Dbpsu", 9},
		{"Dbpatch", 10},
		{"Dblsnrinfo", 11},
		{"Dbparameter", 12},
		{"Db_parameter_file", 13},
		{"Dbredocheck", 14},
		{"Dbredoswitch", 15},
		{"Recovery_usage", 16},
		{"Recovery_detail", 17},
		{"Dberrlog", 18},
		{"Dbdglagcheck", 19},
		{"Dbdgerrcheck", 20},
	}

	// 遍历所有实例节点
	for nodeIndex, instsht := range *instshts {
		// 确定列位置：C列对应NODE1，D列对应NODE2，以此类推
		col := fmt.Sprintf("%c", 'C'+nodeIndex)

		// 填充每个字段
		for _, field := range instFields {
			cell := fmt.Sprintf("%s%d", col, field.row)
			var content string
			var alarm string

			// 根据字段名获取对应的内容和告警级别
			switch field.fieldName {
			case "NodeID":
				content = instsht.NodeID
				alarm = ""
			case "Instname":
				content = instsht.Instname.Contents
				alarm = instsht.Instname.Alarm
			case "Loadprofile":
				content = instsht.Loadprofile.Contents
				alarm = instsht.Loadprofile.Alarm
			case "Instefficiency":
				content = instsht.Instefficiency.Contents
				alarm = instsht.Instefficiency.Alarm
			case "Topevent":
				content = instsht.Topevent.Contents
				alarm = instsht.Topevent.Alarm
			case "Topsql_by_ela":
				content = instsht.Topsql_by_ela.Contents
				alarm = instsht.Topsql_by_ela.Alarm
			case "Cursor_share_mem":
				content = instsht.Cursor_share_mem.Contents
				alarm = instsht.Cursor_share_mem.Alarm
			case "Dbresource":
				content = instsht.Dbresource.Contents
				alarm = instsht.Dbresource.Alarm
			case "Dbpsu":
				content = instsht.Dbpsu.Contents
				alarm = instsht.Dbpsu.Alarm
			case "Dbpatch":
				content = instsht.Dbpatch.Contents
				alarm = instsht.Dbpatch.Alarm
			case "Dblsnrinfo":
				content = instsht.Dblsnrinfo.Contents
				alarm = instsht.Dblsnrinfo.Alarm
			case "Dbparameter":
				content = instsht.Dbparameter.Contents
				alarm = instsht.Dbparameter.Alarm
			case "Db_parameter_file":
				content = instsht.Db_parameter_file.Contents
				alarm = instsht.Db_parameter_file.Alarm
			case "Dbredocheck":
				content = instsht.Dbredocheck.Contents
				alarm = instsht.Dbredocheck.Alarm
			case "Dbredoswitch":
				content = instsht.Dbredoswitch.Contents
				alarm = instsht.Dbredoswitch.Alarm
			case "Recovery_usage":
				content = instsht.Recovery_usage.Contents
				alarm = instsht.Recovery_usage.Alarm
			case "Recovery_detail":
				content = instsht.Recovery_detail.Contents
				alarm = instsht.Recovery_detail.Alarm
			case "Dberrlog":
				content = instsht.Dberrlog.Contents
				alarm = instsht.Dberrlog.Alarm
			case "Dbdglagcheck":
				content = instsht.Dbdglagcheck.Contents
				alarm = instsht.Dbdglagcheck.Alarm
			case "Dbdgerrcheck":
				content = instsht.Dbdgerrcheck.Contents
				alarm = instsht.Dbdgerrcheck.Alarm
			}

			// 设置单元格内容
			f.SetCellStr(shnm, cell, content)

			// 添加单元格注释（如果有告警）
			if alarm != "" {
				utils.LogDebugf("尝试为INST字段 %s 添加注释，告警级别: %s", field.fieldName, alarm)
				comment := getCommentForField(field.fieldName, summaryEntries)
				utils.LogDebugf("获取到的INST注释内容: %s", comment)
				if comment != "" {
					// 创建带样式的注释
					commentObj := excelize.Comment{
						Cell:   cell,
						Text:   comment,
						Author: "健康检查系统",
						Width:  400, // 注释框宽度
						Height: 100, // 注释框高度
					}

					// 添加调试信息
					utils.LogDebugf("INST注释对象: Cell=%s, Width=%d, Height=%d, Text长度=%d",
						commentObj.Cell, commentObj.Width, commentObj.Height, len(commentObj.Text))

					err := f.AddComment(shnm, commentObj)
					if err != nil {
						utils.LogErrorf("添加INST注释失败: %v", err)
					} else {
						utils.LogDebugf("成功添加INST注释到单元格: %s", cell)
					}
				} else {
					utils.LogDebugf("未找到INST字段 %s 对应的注释内容", field.fieldName)
				}
			}

			// 设置B列检查结果分数 (从B2开始，B1是标题行，跳过NodeID)
			if field.fieldName != "NodeID" {
				bCell := fmt.Sprintf("B%d", field.row)
				var score int
				switch alarm {
				case "R":
					score = 0 // 严重影响
				case "B":
					score = 5 // 普通影响
				case "G":
					score = 8 // 轻微影响
				default:
					score = 10 // 正常
				}
				f.SetCellInt(shnm, bCell, int64(score))
			}

			// 设置单元格样式（根据告警级别）
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

func NewXlsx(xlsnm string) {
	f := excelize.NewFile()
	f.NewSheet("HealthReport")
	f.NewSheet("OS")
	f.NewSheet("DB")
	f.NewSheet("Inst")
	f.DeleteSheet("Sheet1")

	// 初始化 HealthReport Sheet
	shnm := "HealthReport"
	f.SetCellStr(shnm, "C1", "健康检查报告")
	// f.SetCellStr(shnm, "G1", "Health Report")
	// f.SetCellStr(shnm, "B13", "Issue summary")
	// f.SetCellStr(shnm, "C14", "重要")
	// f.SetCellStr(shnm, "D14", "普通")
	// f.SetCellStr(shnm, "E14", "轻微")
	// f.SetCellStr(shnm, "B15", "主机系统分析")
	// f.SetCellStr(shnm, "B16", "数据库实例分析")
	// f.SetCellStr(shnm, "B17", "数据库集群检查")
	// f.SetCellStr(shnm, "B18", "DataGuard检查")
	// f.SetCellStr(shnm, "B19", "数据库备份检查")
	// f.SetCellStr(shnm, "B20", "数据库安全检查")
	// f.SetCellStr(shnm, "B21", "软件使用分析")
	// f.SetCellStr(shnm, "B22", "其他项检查")
	// f.SetCellStr(shnm, "B23", "Issue list")
	// f.SetCellStr(shnm, "B24", "No.")
	// f.SetCellStr(shnm, "C24", "问题类别")
	// f.SetCellStr(shnm, "D24", "检查项")
	// f.SetCellStr(shnm, "E24", "结果")
	// f.SetCellStr(shnm, "F24", "影响")
	// f.SetCellStr(shnm, "G24", "问题描述及建议")

	// 初始化 OS Sheet
	shnm = "OS"
	// f.SetCellStr(shnm, "A1", "主机名")
	// f.SetCellStr(shnm, "A2", "IP地址")
	// f.SetCellStr(shnm, "A3", "主机内核参数")
	// f.SetCellStr(shnm, "A4", "主机资源限制")
	// f.SetCellStr(shnm, "A5", "文件系统使用率")
	// f.SetCellStr(shnm, "A6", "索引资源节点使用率")
	// f.SetCellStr(shnm, "A7", "CPU负载")
	// f.SetCellStr(shnm, "A8", "内存使用")
	// f.SetCellStr(shnm, "A9", "磁盘IO负载检查")
	// f.SetCellStr(shnm, "A10", "透明大页开启检查")
	// f.SetCellStr(shnm, "A11", "主机大页使用检查")
	// f.SetCellStr(shnm, "A12", "NUMA使用检查")
	// f.SetCellStr(shnm, "A13", "NTP时钟同步检查")

	// 初始化 DB Sheet
	shnm = "DB"
	// f.SetCellStr(shnm, "A1", "数据库名称\nDB_UNIQUE_NAME")
	// f.SetCellStr(shnm, "A2", "主机名")
	// f.SetCellStr(shnm, "A3", "表空间使用率")
	// f.SetCellStr(shnm, "A4", "数据文件大小检查")
	// f.SetCellStr(shnm, "A5", "控制文件检查")
	// f.SetCellStr(shnm, "A6", "数据库用户大小")
	// f.SetCellStr(shnm, "A7", "REDO文件性能检查")
	// f.SetCellStr(shnm, "A8", "归档切换检查")
	// f.SetCellStr(shnm, "A9", "数据库资源使用限制检查")
	// f.SetCellStr(shnm, "A10", "数据库性能负载分析")
	// f.SetCellStr(shnm, "A11", "数据库性能运行效率")
	// f.SetCellStr(shnm, "A12", "数据库Top等待")
	// f.SetCellStr(shnm, "A13", "数据库Top SQL(耗时)")
	// f.SetCellStr(shnm, "A14", "监听状态及日志检查")
	// f.SetCellStr(shnm, "A15", "并行度>1的表")
	// f.SetCellStr(shnm, "A16", "并行度>1的索引")
	// f.SetCellStr(shnm, "A17", "无效索引检查")
	// f.SetCellStr(shnm, "A18", "Oracle序列检查")
	// f.SetCellStr(shnm, "A19", "闪回区配置")
	// f.SetCellStr(shnm, "A20", "FlashRecovery区使用情况")
	// f.SetCellStr(shnm, "A21", "数据库日志检查")
	// f.SetCellStr(shnm, "A22", "数据库RMAN备份")
	// f.SetCellStr(shnm, "A23", "DBA权限用户检查")
	// f.SetCellStr(shnm, "A24", "SYSDBA权限用户检查")
	// f.SetCellStr(shnm, "A25", "数据库审计空间检查")
	// f.SetCellStr(shnm, "A26", "数据库审计对象检查")
	// f.SetCellStr(shnm, "A27", "业务对象存放系统表空间")
	// f.SetCellStr(shnm, "A28", "错误口令登录锁定帐户PROFILE检查")
	// f.SetCellStr(shnm, "A29", "病毒勒索攻击检查")
	// f.SetCellStr(shnm, "A30", "SCNHealthCheck检查")
	// f.SetCellStr(shnm, "A31", "DataGuard同步延迟检查")
	// f.SetCellStr(shnm, "A32", "DataGuard同步报错检查")

	// 初始化 Inst Sheet
	shnm = "Inst"
	// f.SetCellStr(shnm, "A1", "节点ID")
	// f.SetCellStr(shnm, "A2", "实例名称")
	// f.SetCellStr(shnm, "A3", "负载分析")
	// f.SetCellStr(shnm, "A4", "实例效率")
	// f.SetCellStr(shnm, "A5", "游标内存使用")
	// f.SetCellStr(shnm, "A6", "Top等待事件")
	// f.SetCellStr(shnm, "A7", "Top SQL(耗时)")
	// f.SetCellStr(shnm, "A8", "数据库资源使用")
	// f.SetCellStr(shnm, "A9", "PSU补丁信息")
	// f.SetCellStr(shnm, "A10", "数据库补丁信息")
	// f.SetCellStr(shnm, "A11", "监听器信息")
	// f.SetCellStr(shnm, "A12", "数据库参数")
	// f.SetCellStr(shnm, "A13", "参数文件")
	// f.SetCellStr(shnm, "A14", "REDO切换")
	// f.SetCellStr(shnm, "A15", "错误日志")

	// 设置布局和样式
	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "right",
			Vertical:        "center",
			WrapText:        true,
			ShrinkToFit:     true,
			JustifyLastLine: true,
		},
	})
	styLeft, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#555555"},
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  11,
			Color: "#E6E6FA",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	shnms := f.GetSheetList()
	for _, shnm := range shnms {
		f.SetColWidth(shnm, "A", "A", 22)
		f.SetColStyle(shnm, "B:Z", wrapStyle)
		switch shnm {
		case "HealthReport":
			f.SetCellStyle(shnm, "B1", "B24", styLeft)
			f.SetColWidth(shnm, "B", "Z", 50)
		case "OS":
			f.SetCellStyle(shnm, "A1", "A13", styLeft)
			f.SetColWidth(shnm, "B", "Z", 80)
		case "DB":
			f.SetCellStyle(shnm, "A1", "A32", styLeft)
			f.SetColWidth(shnm, "B", "Z", 100)
		case "Inst":
			f.SetCellStyle(shnm, "A1", "A15", styLeft)
			f.SetColWidth(shnm, "B", "Z", 80)
		}
		f.SetPanes(shnm, &excelize.Panes{
			Freeze:      true,
			XSplit:      1,
			YSplit:      2,
			TopLeftCell: "B3",
			ActivePane:  "bottomLeft",
		})
	}
	f.SaveAs(xlsnm)
}

// getHyperlinkForProblem 根据检查项Nm和Category生成超链接
func getHyperlinkForProblem(nm string, category string, problem string) string {
	// 字段名到行号的映射
	fieldToRowMap := map[string]int{
		// OS字段
		"HOSTNAME":         2,
		"IPADDR":           3,
		"OS":               4,
		"RELVER":           5,
		"CORES":            6,
		"CPUCOUNT":         6,
		"CPUMHZ":           7,
		"MEMTOTAL":         8,
		"SWAPTOTAL":        10,
		"OSPARAMETER":      11,
		"ULIMIT":           12,
		"OSLOG":            13,
		"FILESYSTEM":       14,
		"INODEUSAGE":       15,
		"CPUSTAT":          16,
		"MEMSTAT":          17,
		"IOSTAT":           18,
		"THPSTAT":          19,
		"HUGEPAGE":         20,
		"NUMA":             21,
		"NTP":              22,
		"TMZONE":           23,
		"SELINUX":          24,
		"FIREWALL":         25,
		"NSSWITCH":         26,
		"LO_MTU":           27,
		"MACHINE_PLATFORM": 9,
		"CPU_PERF_MODE":    28,
		"NOZEROCONF":       29,
		"RPM_PACKAGES":     30,

		// DB字段
		"DBNAME":             2,
		"DBMAA":              3,
		"DBVER":              4,
		"DBSTATUS":           5,
		"DBLANG":             6,
		"LOGMODE":            7,
		"FLASHBACK":          8,
		"DBCURSIZE":          9,
		"DBF_SIZE":           10,
		"DBF_CNT":            11,
		"DBF_STAT":           12,
		"TMPFILE_SIZE":       13,
		"DBTBLCOUNT":         14,
		"DBROLE":             15,
		"DBTBSUSAGE":         16,
		"DBCONTROLFILE":      17,
		"USER_INFO":          18,
		"USER_SIZE":          19,
		"TAB_INFO":           20,
		"TAB_PARALLEL":       21,
		"INX_PARALLEL":       22,
		"INVALID_OBJ":        23,
		"INVALID_INX":        24,
		"DBSEQUENCE":         25,
		"DB_SEQ_USAGE":       26,
		"DBOPTION":           27,
		"DBFEATURES":         28,
		"DB_EXPIR_USER":      29,
		"DB_PASSWORD_VERIF":  30,
		"DBDBAPRIV":          31,
		"DBSYSDBA":           32,
		"DBAUDITSEGMENT":     33,
		"DBAUDITCONT":        34,
		"DB_NOSYS_IN_SYSTEM": 35,
		"USERFAILEDLOGIN":    36,
		"DBVIRSCHECK":        37,
		"DBSCNHEALTHCHECK":   38,
		"DBRMANCHECK":        39,
		"CRS_STAT":           40,
		"CRS_STAT2":          41,
		"OCR_INFO":           42,
		"OCR_BAK_CHECK":      43,
		"ASM_USAGE":          44,
		"ASM_OFFSET":         45,

		// INST字段
		"INSTNAME":          2,
		"LOADPROFILE":       3,
		"INSTEFFICIENCY":    4,
		"TOPEVENT":          5,
		"TOPSQL_BY_ELA":     6,
		"CURSOR_SHARE_MEM":  7,
		"DBRESOURCE":        8,
		"DBPSU":             9,
		"DBPATCH":           10,
		"DBLSNRINFO":        11,
		"DBPARAMETER":       12,
		"DB_PARAMETER_FILE": 13,
		"DBREDOCHECK":       14,
		"DBREDOSWITCH":      15,
		"RECOVERY_USAGE":    16,
		"RECOVERY_DETAIL":   17,
		"DBERRLOG":          18,
		"DBDGLAGCHECK":      19,
		"DBDGERRCHECK":      20,
	}

	// 根据字段名确定Sheet名称
	var sheetName string

	// 检查字段是否属于OS sheet
	osFields := []string{
		"HOSTNAME", "IPADDR", "OS", "RELVER", "CORES", "CPUCOUNT", "CPUMHZ",
		"MEMTOTAL", "SWAPTOTAL", "OSPARAMETER", "ULIMIT", "OSLOG", "FILESYSTEM",
		"INODEUSAGE", "CPUSTAT", "MEMSTAT", "IOSTAT", "THPSTAT", "HUGEPAGE",
		"NUMA", "NTP", "TMZONE", "SELINUX", "FIREWALL", "NSSWITCH", "LO_MTU",
		"MACHINE_PLATFORM", "CPU_PERF_MODE", "NOZEROCONF", "RPM_PACKAGES",
	}

	// 检查字段是否属于DB sheet
	dbFields := []string{
		"DBNAME", "DBMAA", "DBVER", "DBSTATUS", "DBLANG", "LOGMODE", "FLASHBACK",
		"DBCURSIZE", "DBF_SIZE", "DBF_CNT", "DBF_STAT", "TMPFILE_SIZE", "DBTBLCOUNT",
		"DBROLE", "DBTBSUSAGE", "DBCONTROLFILE", "USER_INFO", "USER_SIZE", "TAB_INFO",
		"TAB_PARALLEL", "INX_PARALLEL", "INVALID_OBJ", "INVALID_INX", "DBSEQUENCE",
		"DB_SEQ_USAGE", "DBOPTION", "DBFEATURES", "DB_EXPIR_USER", "DB_PASSWORD_VERIF",
		"DBDBAPRIV", "DBSYSDBA", "DBAUDITSEGMENT", "DBAUDITCONT", "DB_NOSYS_IN_SYSTEM",
		"USERFAILEDLOGIN", "DBVIRSCHECK", "DBSCNHEALTHCHECK", "DBRMANCHECK", "CRS_STAT",
		"CRS_STAT2", "OCR_INFO", "OCR_BAK_CHECK", "ASM_USAGE", "ASM_OFFSET",
	}

	// 检查字段是否属于Inst sheet
	instFields := []string{
		"INSTNAME", "LOADPROFILE", "INSTEFFICIENCY", "TOPEVENT", "TOPSQL_BY_ELA",
		"CURSOR_SHARE_MEM", "DBRESOURCE", "DBPSU", "DBPATCH", "DBLSNRINFO",
		"DBPARAMETER", "DB_PARAMETER_FILE", "DBREDOCHECK", "DBREDOSWITCH",
		"RECOVERY_USAGE", "RECOVERY_DETAIL", "DBERRLOG", "DBDGLAGCHECK", "DBDGERRCHECK",
	}

	// 根据字段名确定Sheet
	if contains(osFields, nm) {
		sheetName = "OS"
	} else if contains(dbFields, nm) {
		sheetName = "DB"
	} else if contains(instFields, nm) {
		sheetName = "Inst"
	} else {
		// 默认根据Category确定
		switch category {
		case "主机系统":
			sheetName = "OS"
		case "数据库分析":
			sheetName = "DB"
		case "实例分析":
			sheetName = "Inst"
		default:
			sheetName = "OS"
		}
	}

	// 获取行号
	row, exists := fieldToRowMap[nm]
	if !exists {
		return ""
	}

	// 从问题描述中解析节点信息
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

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
