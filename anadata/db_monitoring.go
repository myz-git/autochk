package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// db_monitoring.go 包含数据库错误监控、DataGuard、备份及杂项分析函数

// 错误监控

// Ana_DBERRLOG 分析数据库错误日志
func Ana_DBERRLOG(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dberrlog.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dberrlog.Nm,
		Title:    rule.Dbrule.Dberrlog.Title,
		Desc:     rule.Dbrule.Dberrlog.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		instshtp.Dberrlog.Alarm = ""
	} else {
		// 有内容说明存在错误日志，设置为G级告警
		instshtp.Dberrlog.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,近期数据库日志存在重要报错信息,\n建议: 根据错误信息定位问题处理相关错误", instshtp.Instname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBLSNRINFO 分析监听日志文件大小
func Ana_DBLSNRINFO(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dblsnrinfo.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dblsnrinfo.Nm,
		Title:    rule.Dbrule.Dblsnrinfo.Title,
		Desc:     rule.Dbrule.Dblsnrinfo.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Dblsnrinfo.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,监听信息数据采集异常", instshtp.Instname.Contents))
		return
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	var hasWarning bool

	// 检查第一行是否有 tnslsnr 字样
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if !strings.Contains(firstLine, "tnslsnr") {
			instshtp.Dblsnrinfo.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,未检测到监听,\n建议: 核查监听是否正常运行", instshtp.Instname.Contents))
			hasWarning = true
		}
	}

	// 查找以 size= 开头的行，解析日志文件大小
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否以 size= 开头
		if strings.HasPrefix(line, "size=") {
			// 解析 size=xxx bytes 格式
			// 例如: size=170909890 bytes, file=...
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				// 提取 size 值，去掉 "size=" 前缀
				sizeStr := strings.TrimPrefix(parts[0], "size=")

				// 转换为整数
				if size, err := strconv.Atoi(sizeStr); err == nil {
					// 检查是否超过阈值
					if size >= rule.Dbrule.Dblsnrinfo.Log_size {
						// 如果之前没有设置G级告警，则设置为G级
						if instshtp.Dblsnrinfo.Alarm != "G" {
							instshtp.Dblsnrinfo.Alarm = "G"
						}
						entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,监听日志文件较大(%d bytes),影响监听响应性能,\n建议: 定期清理或归档监听日志保持在2G以下", instshtp.Instname.Contents, size))
						hasWarning = true
					}
				}
			}
		}
	}

	// 如果没有告警，清空告警级别
	if !hasWarning {
		instshtp.Dblsnrinfo.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// DataGuard 和备份

// Ana_DBDGLAGCHECK 分析 DataGuard 同步延迟
func Ana_DBDGLAGCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	if strings.Contains(dbshtp.Dbrole.Contents, "PRIMARY") {
		return
	}
	msgdata := instshtp.Dbdglagcheck.Contents
	rdok := regexp.MustCompile(`^apply lag(.*)\+00 00:00:00$`)
	rd := regexp.MustCompile(`^apply lag(.*)\+(.*):\d+$`)
	instshtp.Dbdglagcheck.Alarm = "G"
	entry := structs.SummaryEntry{
		Category: "DataGuard",
		Nm:       rule.Dbrule.Dbdglagcheck.Nm,
		Title:    rule.Dbrule.Dbdglagcheck.Title,
		Desc:     rule.Dbrule.Dbdglagcheck.Desc,
	}
Looop:
	for _, row := range strings.Split(msgdata, "\n") {
		row = strings.TrimSpace(row)
		if rdok.MatchString(row) {
			instshtp.Dbdglagcheck.Alarm = ""
			break Looop
		}
		if rd.MatchString(row) {
			values1 := strings.Split(row, "+")
			values2 := strings.Fields(values1[1])
			vDay, _ := strconv.Atoi(values2[0])
			if vDay >= rule.Dbrule.Dbdglagcheck.ResultB {
				instshtp.Dbdglagcheck.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,DataGuard同步延迟当前%d超过阈值%d,\n建议: 检查网络连接和数据库状态", instshtp.Instname.Contents, vDay, rule.Dbrule.Dbdglagcheck.ResultB))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBDGERRCHECK 分析 DataGuard 同步错误
func Ana_DBDGERRCHECK(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbdgerrcheck.Contents
	value := strings.TrimSpace(msgdata)
	entry := structs.SummaryEntry{
		Category: "DataGuard",
		Nm:       rule.Dbrule.Dbdgerrcheck.Nm,
		Title:    rule.Dbrule.Dbdgerrcheck.Title,
		Desc:     rule.Dbrule.Dbdgerrcheck.Desc,
	}
	// if value != "" {
	if value != "" && !strings.Contains(value, "no rows selected") && !strings.Contains(value, "无记录") {
		instshtp.Dbdgerrcheck.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,DataGuard日志存在同步错误信息,\n建议: 检查并处理相关错误", instshtp.Instname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// 杂项检查

// Ana_DBPSU 分析 PSU 使用情况
func Ana_DBPSU(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbpsu.Contents
	entry := structs.SummaryEntry{
		Category: "软件使用",
		Nm:       rule.Dbrule.Dbpsu.Nm,
		Title:    rule.Dbrule.Dbpsu.Title,
		Desc:     rule.Dbrule.Dbpsu.Desc,
	}

	// 检查是否为空
	if strings.TrimSpace(msgdata) == "" {
		// 数据采集异常，设置为G级告警
		instshtp.Dbpsu.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,PSU使用情况数据采集异常", instshtp.Instname.Contents))
	} else {
		// 按行分割数据
		lines := strings.Split(msgdata, "\n")

		// 检查行数：如果只有3行，说明只有标题行+1行数据，未安装PSU或RU
		if len(lines) <= 3 {
			// 只有标题行+1行数据，说明未安装PSU或RU
			instshtp.Dbpsu.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,未安装PSU或RU,\n建议: 生产系统定期安装和更新PSU或RU", instshtp.Instname.Contents))
		} else {
			// 从第4行开始检查日期（跳过标题行）
			var maxDate time.Time
			var hasValidDate bool

			for i := 3; i < len(lines); i++ {
				currentLine := strings.TrimSpace(lines[i])
				if currentLine == "" {
					continue
				}

				// 按空格分割行数据
				fields := strings.Fields(currentLine)
				if len(fields) >= 1 {
					dateStr := fields[0]

					// 尝试解析日期格式 "2025-07-05"
					if date, err := time.Parse("2006-01-02", dateStr); err == nil {
						if !hasValidDate || date.After(maxDate) {
							maxDate = date
							hasValidDate = true
						}
					}
				}
			}

			// 如果找到了有效日期，检查是否超过2年
			if hasValidDate {
				now := time.Now()
				yearsDiff := now.Sub(maxDate).Hours() / 24 / 365.25

				if yearsDiff > 2 {
					instshtp.Dbpsu.Alarm = "G"
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s实例,PSU/RU最后更新日期为%s,超过2年未更新,\n建议: 重要生产系统最少两年更新一次PSU或RU", instshtp.Instname.Contents, maxDate.Format("2006-01-02")))
				} else {
					// 正常情况，不设置告警
					instshtp.Dbpsu.Alarm = ""
				}
			} else {
				// 没有找到有效日期，设置为B级告警
				instshtp.Dbpsu.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,无法解析PSU/RU更新日期,\n建议: 重要生产系统最少两年更新一次PSU或RU", instshtp.Instname.Contents))
			}
		}
	}

	// 如果有告警信息，添加到汇总中
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBPATCH 分析补丁使用情况
func Ana_DBPATCH(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbpatch.Contents
	entry := structs.SummaryEntry{
		Category: "软件使用",
		Nm:       rule.Dbrule.Dbpatch.Nm,
		Title:    rule.Dbrule.Dbpatch.Title,
		Desc:     rule.Dbrule.Dbpatch.Desc,
	}
	// 实现补丁分析逻辑
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		instshtp.Dbpatch.Alarm = "G"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s实例,未检测到补丁安装记录,\n建议: 生产系统定期安装和更新补丁", instshtp.Instname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
