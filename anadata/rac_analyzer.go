package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"strings"
)

// Ana_Crs_stat 分析集群状态
func Ana_Crs_stat(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Crs_stat.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Crs_stat.Nm,
		Title:    rule.Dbrule.Crs_stat.Title,
		Desc:     rule.Dbrule.Crs_stat.Desc,
	}

	// 按空行分割，得到每个服务组
	groups := strings.Split(strings.TrimSpace(msgdata), "\n\n")
	var abnormalServices []string

	for _, group := range groups {
		if strings.TrimSpace(group) == "" {
			continue
		}

		lines := strings.Split(group, "\n")
		var targetCount, stateCount int
		var serviceName string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "NAME=") {
				serviceName = strings.TrimPrefix(line, "NAME=")
			} else if strings.HasPrefix(line, "TARGET=") {
				// 统计 TARGET 行中的 ONLINE 数量
				targetLine := strings.TrimPrefix(line, "TARGET=")
				targetCount = strings.Count(targetLine, "ONLINE")
			} else if strings.HasPrefix(line, "STATE=") {
				// 统计 STATE 行中的 ONLINE 数量
				stateLine := strings.TrimPrefix(line, "STATE=")
				stateCount = strings.Count(stateLine, "ONLINE")
			}
		}

		// 如果 TARGET 和 STATE 中的 ONLINE 数量不匹配，则标记为异常
		if targetCount != stateCount {
			abnormalServices = append(abnormalServices, serviceName)
		}
	}

	// 如果存在异常服务，设置告警级别为 R
	if len(abnormalServices) > 0 {
		dbshtp.Crs_stat.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("集群服务状态异常，异常服务：%s", strings.Join(abnormalServices, ", ")))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ocr_info 分析OCR信息
func Ana_Ocr_info(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Ocr_info.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Ocr_info.Nm,
		Title:    rule.Dbrule.Ocr_info.Title,
		Desc:     rule.Dbrule.Ocr_info.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Ocr_info.Alarm = ""
	} else {
		// 有记录说明存在OCR问题，设置为B级告警
		dbshtp.Ocr_info.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,OCR信息异常,建议检查Oracle集群注册表状态", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ocr_bak_check 分析OCR备份检查
func Ana_Ocr_bak_check(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Ocr_bak_check.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Ocr_bak_check.Nm,
		Title:    rule.Dbrule.Ocr_bak_check.Title,
		Desc:     rule.Dbrule.Ocr_bak_check.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Ocr_bak_check.Alarm = ""
	} else {
		// 有记录说明存在OCR备份问题，设置为B级告警
		dbshtp.Ocr_bak_check.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,OCR备份检查异常,建议检查OCR备份配置", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Asm_offset 分析ASM偏移量
func Ana_Asm_offset(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Asm_offset.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Asm_offset.Nm,
		Title:    rule.Dbrule.Asm_offset.Title,
		Desc:     rule.Dbrule.Asm_offset.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Asm_offset.Alarm = ""
	} else {
		// 有记录说明存在ASM偏移量问题，设置为B级告警
		dbshtp.Asm_offset.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM偏移量异常,建议检查ASM磁盘对齐", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBASMUSAGE 分析 ASM 使用情况
func Ana_ASM_usage(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Asm_usage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Asm_usage.Nm,
		Title:    rule.Dbrule.Asm_usage.Title,
		Desc:     rule.Dbrule.Asm_usage.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Asm_usage.Alarm = ""
	} else {
		// 有记录说明存在ASM使用问题，设置为B级告警
		dbshtp.Asm_usage.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM使用情况异常,建议检查ASM磁盘组状态", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
