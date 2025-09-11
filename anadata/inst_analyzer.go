package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// db_instance.go 包含与数据库实例状态和存储相关的分析函数

// Ana_RDF 分析 REDO 文件状态和大小
func Ana_RDF(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbredocheck.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbredocheck.Nm,
		Title:    rule.Dbrule.Dbredocheck.Title,
		Desc:     rule.Dbrule.Dbredocheck.Desc,
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		instshtp.Dbredocheck.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,REDO文件检查数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 从第3行开始检查数据（跳过标题行）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 检查行是否以数字开头（数据行）
		rd := regexp.MustCompile(`^\d`)
		if !rd.MatchString(line) {
			continue
		}

		// 按空格分割字段
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		// 1. 检查第6列 STATUS，是否状态不在允许范围内
		status := fields[5]
		if !utils.Contain(status, rule.Dbrule.Dbredocheck.Rdf_status_list) {
			instshtp.Dbredocheck.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("%s实例,REDO文件%s状态%s异常,\n建议: 立即核查REDO文件状态是否正常", instshtp.Instname.Contents, fields[0], status))
			break // 发现状态异常就退出，不需要继续检查
		}

		// 2. 检查第4列 MB，是否有小于阈值的情况
		if size, err := strconv.ParseFloat(fields[3], 64); err == nil {
			if size < rule.Dbrule.Dbredocheck.Rdf_size {
				instshtp.Dbredocheck.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,REDO文件%s大小%.1fMB小于阈值%.1fMB,\n建议: 生产环境每个redo文件在200M-4G之间", instshtp.Instname.Contents, fields[0], size, rule.Dbrule.Dbredocheck.Rdf_size))
				break // 发现大小异常就退出，不需要继续检查
			}
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RDSW 分析归档切换次数
func Ana_RDSW(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbredoswitch.Contents
	rd := regexp.MustCompile(` \d+$`)
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbredoswitch.Nm,
		Title:    rule.Dbrule.Dbredoswitch.Title,
		Desc:     rule.Dbrule.Dbredoswitch.Desc,
	}
Looop:
	for index, msgs := range strings.Split(msgdata, "\n") {
		msgs = strings.TrimSpace(msgs)
		if index < 2 {
			continue
		}
		if rd.MatchString(msgs) {
			msg := strings.Fields(msgs)
			for k, v := range msg {
				if k < 3 {
					continue
				}
				value, _ := strconv.Atoi(v)
				if value > rule.Dbrule.Dbredoswitch.Sw_cnt_ge1 {
					instshtp.Dbredoswitch.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,归档切换每小时%d次超过阈值%d次,\n建议: 增加redo大小及数量期望每10-20分钟切换一次", instshtp.Instname.Contents, value, rule.Dbrule.Dbredoswitch.Sw_cnt_ge1))
					break Looop
				}
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RECOVERY_USAGE 分析闪回区空间使用情况
func Ana_RECOVERY_USAGE(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Recovery_usage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Recovery_usage.Nm,
		Title:    rule.Dbrule.Recovery_usage.Title,
		Desc:     rule.Dbrule.Recovery_usage.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		instshtp.Recovery_usage.Alarm = ""
		return
	}

	// 检查数据行数是否足够
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		instshtp.Recovery_usage.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,闪回区空间使用检查数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 获取第3行（索引2）
	thirdLine := strings.TrimSpace(lines[2])
	if thirdLine == "" {
		// 第3行为空，数据采集异常
		instshtp.Recovery_usage.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,闪回区空间使用检查数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 按空格分割第3行，获取第4列（索引3）
	fields := strings.Fields(thirdLine)
	if len(fields) < 4 {
		// 字段数不足，数据采集异常
		instshtp.Recovery_usage.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,闪回区空间使用检查数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 尝试将第4列转换为数字
	if usedPercent, err := strconv.ParseFloat(fields[3], 64); err == nil {
		// 如果 >= rule.Dbrule.Recovery_usage.Result[1]，则为R级告警
		if usedPercent >= float64(rule.Dbrule.Recovery_usage.Result[1]) {
			instshtp.Recovery_usage.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s实例,闪回区使用率当前%.2f%%超过严重阈值%.0f%%,\n建议: 尽快清理或扩容闪回区", instshtp.Instname.Contents, usedPercent, rule.Dbrule.Recovery_usage.Result[1]))
		} else if usedPercent >= float64(rule.Dbrule.Recovery_usage.Result[0]) {
			// 如果 >= rule.Dbrule.Recovery_usage.Result[0]，则为B级告警
			instshtp.Recovery_usage.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,闪回区使用率当前%.2f%%超过阈值%.0f%%,\n建议: 及时清理或扩容闪回区", instshtp.Instname.Contents, usedPercent, rule.Dbrule.Recovery_usage.Result[0]))
		} else {
			// 使用率在正常范围内
			instshtp.Recovery_usage.Alarm = ""
		}
	} else {
		// 第4列不是数字，数据采集异常
		instshtp.Recovery_usage.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,闪回区空间使用检查数据采集异常", instshtp.Instname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBparameter 分析数据库初始化参数
func Ana_DBparameter(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
}
