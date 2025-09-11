package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// db_performance.go 包含与数据库性能和效率相关的分析函数
// Ana_DB4031check 分析 ORA-4031 错误
func Ana_DB4031check(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dberrlog.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Db_4031check.Nm,
		Title:    rule.Dbrule.Db_4031check.Title,
		Desc:     rule.Dbrule.Db_4031check.Desc,
	}
	if strings.Contains(msgdata, "ORA-4031") {
		instshtp.Dberrlog.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s实例,检测到ORA-4031共享池内存不足错误,\n建议: 调整共享池大小或优化相关参数配置", instshtp.Instname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RESOURCE 分析数据库资源使用情况
func Ana_RESOURCE(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbresource.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Dbresource.Nm,
		Title:    rule.Dbrule.Dbresource.Title,
		Desc:     rule.Dbrule.Dbresource.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Dbresource.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,资源使用情况数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	var hasWarning bool

	// 从第3行开始检查数据（跳过标题行）
	for i := 2; i < len(lines); i++ {
		currentLine := strings.TrimSpace(lines[i])
		if currentLine == "" {
			continue
		}

		// 按空格分割行数据
		fields := strings.Fields(currentLine)
		if len(fields) >= 4 {
			resourceName := fields[0]  // 第1列：资源名称
			maxUtilStr := fields[2]    // 第3列：最大使用量
			limitValueStr := fields[3] // 第4列：限制值

			// 移除逗号并转换为数值
			maxUtilStr = strings.Replace(maxUtilStr, ",", "", -1)
			limitValueStr = strings.Replace(limitValueStr, ",", "", -1)

			if maxUtil, err := strconv.Atoi(maxUtilStr); err == nil {
				if limitValue, err := strconv.Atoi(limitValueStr); err == nil {
					// 检查使用率是否超过90%
					if limitValue > 0 {
						// 使用最大使用量计算使用率
						utilizationRate := float64(maxUtil) / float64(limitValue) * 100
						if utilizationRate > 90 {
							hasWarning = true
							instshtp.Dbresource.Alarm = "B"
							entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,%s资源近期使用%d接近最大限制%d,\n建议: 根据需要调整限制及核查资源使用增长是否预期规划", instshtp.Instname.Contents, resourceName, maxUtil, limitValue))
						}
					}
				}
			}
		}
	}

	// 如果没有告警，清空告警级别
	if !hasWarning {
		instshtp.Dbresource.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_LOADPROFILE 分析数据库负载性能
func Ana_LOADPROFILE(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Loadprofile.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Loadprofile.Nm,
		Title:    rule.Dbrule.Loadprofile.Title,
		Desc:     rule.Dbrule.Loadprofile.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Loadprofile.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,负载性能数据采集异常", instshtp.Instname.Contents))
	} else {
		// 有数据，进行详细分析
		lines := strings.Split(msgdata, "\n")
		var hasRedoWarning, hasLogonWarning bool

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 匹配 Redo size (bytes) 行
			if strings.Contains(strings.ToLower(line), "redo") && strings.Contains(strings.ToLower(line), "size") && strings.Contains(strings.ToLower(line), "bytes") {
				// 先用冒号分割，取得数值部分
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					// 取得冒号后面的部分：     57,671,680.0         228,185.1
					valuePart := strings.TrimSpace(parts[1])

					// 再以空格分割，取第一个数字
					valueFields := strings.Fields(valuePart)
					if len(valueFields) >= 1 {
						// 移除逗号并转换为浮点数
						redoSizeStr := strings.Replace(valueFields[0], ",", "", -1)

						if redoSize, err := strconv.ParseFloat(redoSizeStr, 64); err == nil {
							// 检查是否超过阈值（假设Redosize单位是MB，需要转换为字节进行比较）
							if redoSize >= rule.Dbrule.Loadprofile.Redosize*1024*1024 {
								hasRedoWarning = true
								instshtp.Loadprofile.Alarm = "B"
								entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,数据库负载较大(每秒产生redo数据量%.2f bytes),\n建议: 根据业务运行情况降低REDO压力,提升数据库性能", instshtp.Instname.Contents, redoSize))
							}
						}
					}
				}
			}

			// 匹配 Logons 行
			if strings.Contains(strings.ToLower(line), "logons:") {
				// 提取第二列数字（Per Second）
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					// 移除逗号并转换为浮点数
					logonStr := strings.Replace(fields[1], ",", "", -1)
					if logon, err := strconv.ParseFloat(logonStr, 64); err == nil {
						if logon >= rule.Dbrule.Loadprofile.Logon {
							hasLogonWarning = true
							instshtp.Loadprofile.Alarm = "B"
							entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,数据库连接压力过大(监听,每秒%.2f次连接),\n建议: 尽量避免短连接减少并发连接或启用多个监听分担压力", instshtp.Instname.Contents, logon))
						}
					}
				}
			}
		}

		// 如果没有告警，清空告警级别
		if !hasRedoWarning && !hasLogonWarning {
			instshtp.Loadprofile.Alarm = ""
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_INSTEFFICIENCY 分析数据库实例效率
func Ana_INSTEFFICIENCY(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Instefficiency.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Instefficiency.Nm,
		Title:    rule.Dbrule.Instefficiency.Title,
		Desc:     rule.Dbrule.Instefficiency.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Instefficiency.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,实例效率数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	var hasWarning bool

	// 从第3行开始检查数据（跳过标题行）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查 Buffer Hit 命中率
		if strings.Contains(line, "Buffer") && strings.Contains(line, "Hit") {
			// 提取百分比数值
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				valueFields := strings.Fields(valuePart)
				if len(valueFields) >= 1 {
					// 移除百分号并转换为浮点数
					bufferHitStr := strings.TrimSuffix(valueFields[0], "%")
					if bufferHit, err := strconv.ParseFloat(bufferHitStr, 64); err == nil {
						if bufferHit < rule.Dbrule.Instefficiency.Buffer_hit {
							hasWarning = true
							instshtp.Instefficiency.Alarm = "G"
							entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,Buffer Hit命中率%.2f%%小于%.2f%%,\n建议: 优化适当增加BufferCache或SGA,提升缓存命中率", instshtp.Instname.Contents, bufferHit, rule.Dbrule.Instefficiency.Buffer_hit))
						}
					}
				}
			}
		}

		// 检查 Library Hit 命中率
		if strings.Contains(line, "Library") && strings.Contains(line, "Hit") {
			// 提取百分比数值
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				valueFields := strings.Fields(valuePart)
				if len(valueFields) >= 1 {
					// 移除百分号并转换为浮点数
					libraryHitStr := strings.TrimSuffix(valueFields[0], "%")
					if libraryHit, err := strconv.ParseFloat(libraryHitStr, 64); err == nil {
						if libraryHit < rule.Dbrule.Instefficiency.Library_hit {
							hasWarning = true
							instshtp.Instefficiency.Alarm = "G"
							entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,Library Hit命中率%.2f%%小于%.2f%%,\n建议: 适当增加SHAREDPOOL或SGA,优化库缓存命中率", instshtp.Instname.Contents, libraryHit, rule.Dbrule.Instefficiency.Library_hit))
						}
					}
				}
			}
		}

		// 检查 Soft Parse 命中率
		if strings.Contains(line, "Soft") && strings.Contains(line, "Parse") {
			// 提取百分比数值
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				valuePart := strings.TrimSpace(parts[1])
				valueFields := strings.Fields(valuePart)
				if len(valueFields) >= 1 {
					// 移除百分号并转换为浮点数
					softParseStr := strings.TrimSuffix(valueFields[0], "%")
					if softParse, err := strconv.ParseFloat(softParseStr, 64); err == nil {
						if softParse < rule.Dbrule.Instefficiency.Soft_parse {
							hasWarning = true
							instshtp.Instefficiency.Alarm = "B"
							entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,Soft Parse命中率%.2f%%小于%.2f%%,\n建议: 改造SQL常量为绑定变量减少硬解析", instshtp.Instname.Contents, softParse, rule.Dbrule.Instefficiency.Soft_parse))
						}
					}
				}
			}
		}
	}

	// 如果没有告警，清空告警级别
	if !hasWarning {
		instshtp.Instefficiency.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBtopevent 分析顶部等待事件

// Ana_DBtopSQL 分析顶部 SQL 性能
func Ana_DBtopSQL(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Topsql_by_ela.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Topsql_by_ela.Nm,
		Title:    rule.Dbrule.Topsql_by_ela.Title,
		Desc:     rule.Dbrule.Topsql_by_ela.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Topsql_by_ela.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,Top SQL数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	var hasWarning bool

	// 从第3行开始检查数据（跳过标题行）
	for i := 2; i < len(lines); i++ {
		currentLine := strings.TrimSpace(lines[i])
		if currentLine == "" {
			continue
		}

		// 使用正则表达式匹配独立的数字（前后有空格或行首行尾）
		re := regexp.MustCompile(`(?:^|\s)(\d+(?:\.\d+)?)(?:\s|$)`)
		matches := re.FindAllStringSubmatch(currentLine, -1)

		if len(matches) >= 4 {
			// 获取SQL_ID（第一个非数字字段）
			fields := strings.Fields(currentLine)
			sqlID := fields[0]

			// 从匹配的数字中取第2个（执行次数）和第3个（平均耗时）
			executionsStr := matches[1][1] // 第2个数字：10004
			avgTimeStr := matches[2][1]    // 第3个数字：82.64

			if executions, err := strconv.Atoi(executionsStr); err == nil {
				if avgTime, err := strconv.ParseFloat(avgTimeStr, 64); err == nil {
					// 检查条件：执行次数 > 10000 且 平均耗时 > 10
					if executions > 10000 && avgTime > 10 {
						hasWarning = true
						instshtp.Topsql_by_ela.Alarm = "G"
						entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,SQL_ID:%s,执行次数%d次,平均耗时%.2f秒,\n建议: 对高频执行的SQL语句进行优化提升单次执行效率", instshtp.Instname.Contents, sqlID, executions, avgTime))
						break // 找到第一个符合条件的SQL就停止
					}
				}
			}
		}
	}

	// 如果没有告警，清空告警级别
	if !hasWarning {
		instshtp.Topsql_by_ela.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_CursorShareMem 分析游标共享内存使用情况
func Ana_CursorShareMem(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Cursor_share_mem.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Cursor_share_mem.Nm,
		Title:    rule.Dbrule.Cursor_share_mem.Title,
		Desc:     rule.Dbrule.Cursor_share_mem.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		instshtp.Cursor_share_mem.Alarm = ""
	} else {
		// 有记录说明存在游标共享内存使用过大的情况，设置为G级告警
		instshtp.Cursor_share_mem.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,游标共享内存使用过大,\n建议: 对游标共享内存使用>300M的SQL语句优化", instshtp.Instname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DB_Shp_pct 分析共享池使用率百分比
func Ana_DB_Shp_pct(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Db_shp_pct.Contents
	entry := structs.SummaryEntry{
		Category: "数据库性能",
		Nm:       rule.Dbrule.Db_shp_pct.Nm,
		Title:    rule.Dbrule.Db_shp_pct.Title,
		Desc:     rule.Dbrule.Db_shp_pct.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		instshtp.Db_shp_pct.Alarm = ""
		return
	}

	// 按行分割数据，取第3行（索引为2）的数字
	lines := strings.Split(msgdata, "\n")
	if len(lines) >= 3 {
		line3 := strings.TrimSpace(lines[2])
		if value, err := strconv.ParseFloat(line3, 64); err == nil {
			// 检查是否超过阈值
			if value >= rule.Dbrule.Db_shp_pct.Result {
				instshtp.Db_shp_pct.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,sharedpool占SGA比率%.2f%%较高,\n建议: 评估SGA中SHAREDPOOL与BUFFERCACHE分配是否合理,可考虑增加SGA", instshtp.Instname.Contents, value))
			} else {
				// 未超过阈值，正常
				instshtp.Db_shp_pct.Alarm = ""
			}
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Db_shp_size 分析共享池大小
func Ana_Db_shp_size(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	// entry := structs.SummaryEntry{
	// 	Category: "实例分析",
	// 	Nm:       rule.Dbrule.Db_shp_size.Nm,
	// 	Title:    rule.Dbrule.Db_shp_size.Title,
	// 	Desc:     rule.Dbrule.Db_shp_size.Desc,
	// }
	// if strings.Contains(msgdata, "shared_pool_size") {
	// 	re := regexp.MustCompile(`shared_pool_size\s*=\s*(\d+)`)
	// 	if matches := re.FindStringSubmatch(msgdata); len(matches) > 1 {
	// 		size, _ := strconv.Atoi(matches[1])
	// 		if size < 1000000000 {
	// 			dbshtp.Dbrecoverydest.Alarm = "B"
	// 			entry.Moderate = append(entry.Moderate, fmt.Sprintf("共享池大小 %d 字节过小,\n建议: 调整", size))
	// 		}
	// 	}
	// }
	// if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
	// 	summaryEntries.Entries = append(summaryEntries.Entries, entry)
	// }
}
