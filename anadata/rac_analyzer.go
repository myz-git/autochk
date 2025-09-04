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

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为B级告警
		dbshtp.Ocr_info.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,OCR信息数据采集异常", dbshtp.Dbname.Contents))
	} else {
		// 有数据，进行详细分析
		lines := strings.Split(msgdata, "\n")
		onlineCount := 0
		ocrSpace := 0

		// 从第三行开始查找ONLINE状态
		for i := 2; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			// 统计ONLINE状态数量
			if strings.Contains(line, "ONLINE") {
				onlineCount++
			}

			// 查找可用空间信息
			if strings.Contains(line, "Available space (kbytes)") {
				// 提取数值
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					spaceStr := strings.TrimSpace(parts[1])
					// 尝试转换为整数
					if space, err := fmt.Sscanf(spaceStr, "%d", &ocrSpace); err == nil && space > 0 {
						// 成功解析到空间数值
					}
				}
			}
		}

		// 根据ONLINE数量判断告警级别
		switch onlineCount {
		case 3, 5:
			// 正常情况，3个或5个ONLINE
			dbshtp.Ocr_info.Alarm = ""
		case 1:
			// 1个ONLINE，存在单点故障风险
			dbshtp.Ocr_info.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,OCR存在单点故障风险,建议配置NORMAL或以上冗余模式", dbshtp.Dbname.Contents))
		case 2, 4:
			// 2个或4个ONLINE，疑似有冗余OCR盘掉线
			dbshtp.Ocr_info.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,疑似有冗余OCR盘掉线,建议尽快核查所有OCR盘是否状态正常", dbshtp.Dbname.Contents))
		default:
			// 其他情况，设置为B级告警
			dbshtp.Ocr_info.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,OCR状态异常,ONLINE数量为%d", dbshtp.Dbname.Contents, onlineCount))
		}

		// 检查OCR可用空间
		if ocrSpace > 0 && ocrSpace < 204800 {
			// 如果之前没有设置R级告警，则设置为R级
			if dbshtp.Ocr_info.Alarm != "R" {
				dbshtp.Ocr_info.Alarm = "R"
			}
			entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,OCR可用空间不足(%d KB),建议尽快扩展OCR磁盘组", dbshtp.Dbname.Contents, ocrSpace))
		}
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

	// 检查是否为空或不足3行
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		dbshtp.Asm_offset.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,ASM偏移量数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		dbshtp.Asm_offset.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,ASM偏移量数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 从第3行开始检查数据
	var hasSevereOffset, hasModerateOffset bool
	var severeDiskGroups, moderateDiskGroups []string

	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 按空格分割行数据
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		diskGroupName := fields[0] // 第1列：磁盘组名
		rangePctStr := fields[5]   // 第6列：偏移量百分比

		// 移除百分号并转换为浮点数
		rangePctStr = strings.TrimSuffix(rangePctStr, "%")
		if rangePct, err := fmt.Sscanf(rangePctStr, "%f", new(float64)); err == nil && rangePct > 0 {
			if rangePct >= 30 {
				// 偏移量>=30%，严重偏移
				hasSevereOffset = true
				severeDiskGroups = append(severeDiskGroups, diskGroupName)
			} else if rangePct >= 10 {
				// 偏移量>=10%，中等偏移
				hasModerateOffset = true
				moderateDiskGroups = append(moderateDiskGroups, diskGroupName)
			}
		}
	}

	// 根据检查结果设置告警级别
	if hasSevereOffset {
		// 存在严重偏移，设置为B级告警
		dbshtp.Asm_offset.Alarm = "B"
		for _, diskGroup := range severeDiskGroups {
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM磁盘组%s磁盘使用严重偏移,建议及时对ASM磁盘维护校准", dbshtp.Dbname.Contents, diskGroup))
		}
	} else if hasModerateOffset {
		// 存在中等偏移，设置为G级告警
		dbshtp.Asm_offset.Alarm = "G"
		for _, diskGroup := range moderateDiskGroups {
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,ASM磁盘组%s磁盘使用偏移超过10%%,建议保持关注ASM磁盘使用情况", dbshtp.Dbname.Contents, diskGroup))
		}
	} else {
		// 正常情况，不设置告警
		dbshtp.Asm_offset.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_ASM_usage 分析 ASM 使用情况
func Ana_ASM_usage(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Asm_usage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库集群",
		Nm:       rule.Dbrule.Asm_usage.Nm,
		Title:    rule.Dbrule.Asm_usage.Title,
		Desc:     rule.Dbrule.Asm_usage.Desc,
	}

	// 检查是否为空或不足3行
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为B级告警
		dbshtp.Asm_usage.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM使用情况数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		dbshtp.Asm_usage.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM使用情况数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 从第3行开始检查数据
	var hasWarning, hasCritical bool
	var warningDiskGroups, criticalDiskGroups []string

	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 按空格分割行数据
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		diskGroupName := fields[0] // 第1列：磁盘组名
		pctUsedStr := fields[5]    // 第6列：使用百分比
		status := fields[6]        // 第7列：状态

		// 检查第6列是否大于90
		if pctUsed, err := fmt.Sscanf(pctUsedStr, "%f", new(float64)); err == nil && pctUsed > 0 {
			if pctUsed > 90 {
				hasWarning = true
				warningDiskGroups = append(warningDiskGroups, diskGroupName)
			}
		}

		// 检查第7列是否包含"WARNING"
		if strings.Contains(strings.ToUpper(status), "WARNING") {
			hasWarning = true
			warningDiskGroups = append(warningDiskGroups, diskGroupName)
		}

		// 检查第7列是否包含"CRITICAL"
		if strings.Contains(strings.ToUpper(status), "CRITICAL") {
			hasCritical = true
			criticalDiskGroups = append(criticalDiskGroups, diskGroupName)
		}
	}

	// 根据检查结果设置告警级别
	if hasCritical {
		// 存在CRITICAL状态，设置为R级告警
		dbshtp.Asm_usage.Alarm = "R"
		for _, diskGroup := range criticalDiskGroups {
			entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,ASM磁盘组%s空间严重不足,已不能满足冗余需求,建议尽快扩容磁盘或数据清理", dbshtp.Dbname.Contents, diskGroup))
		}
	} else if hasWarning {
		// 存在WARNING状态或使用率超过90%，设置为B级告警
		dbshtp.Asm_usage.Alarm = "B"
		for _, diskGroup := range warningDiskGroups {
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,ASM磁盘组%s空间使用超过90%%,建议及时扩容磁盘或数据清理", dbshtp.Dbname.Contents, diskGroup))
		}
	} else {
		// 正常情况，不设置告警
		dbshtp.Asm_usage.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
