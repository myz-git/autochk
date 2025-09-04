package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// db_analyzer.go 包含与数据库状态相关的指标告警解析函数
// Ana_DB_status 分析数据库状态
func Ana_DB_status(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbstatus.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbstatus.Nm,
		Title:    rule.Dbrule.Dbstatus.Title,
		Desc:     rule.Dbrule.Dbstatus.Desc,
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		dbshtp.Dbstatus.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,数据库状态检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 检查第三行是否包含期望的状态
	dataLine := strings.TrimSpace(lines[2])
	if dataLine == "" {
		// 第三行为空，数据采集异常
		dbshtp.Dbstatus.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,数据库状态检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	if !strings.Contains(dataLine, rule.Dbrule.Dbstatus.Status) {
		dbshtp.Dbstatus.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,当前状态异常，建议尽快确认数据库是否正常", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DB_logmode 分析是否数据库开启归档
func Ana_DB_logmode(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Logmode.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Logmode.Nm,
		Title:    rule.Dbrule.Logmode.Title,
		Desc:     rule.Dbrule.Logmode.Desc,
	}

	// 检查数据库是否开启归档
	if strings.TrimSpace(msgdata) != rule.Dbrule.Logmode.Status {
		dbshtp.Logmode.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,当前日志模式为%s，建议开启归档模式(ARCHIVELOG)以确保数据安全", dbshtp.Dbname.Contents, strings.TrimSpace(msgdata)))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBF_CNT 检查数据文件数量
func Ana_DBF_CNT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbf_cnt.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbf_cnt.Nm,
		Title:    rule.Dbrule.Dbf_cnt.Title,
		Desc:     rule.Dbrule.Dbf_cnt.Desc,
	}

	// 按行分割数据
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，数据采集异常
		dbshtp.Dbf_cnt.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,数据文件数量检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 检查第三行数据
	dataLine := strings.TrimSpace(lines[2])
	if dataLine == "" {
		// 第三行为空，数据采集异常
		dbshtp.Dbf_cnt.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,数据文件数量检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 使用正则表达式提取数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(dataLine, -1)
	if len(matches) < 2 {
		return
	}

	// 解析 FILE_CNT 和 DB_FILES_LIMIT
	fileCnt, err1 := strconv.Atoi(matches[0])
	dbFilesLimit, err2 := strconv.Atoi(matches[1])
	if err1 != nil || err2 != nil {
		return
	}

	// 计算使用率
	usagePercent := float64(fileCnt) / float64(dbFilesLimit) * 100

	// 当文件数使用率超过90%时设置为严重告警
	if usagePercent >= 90 {
		dbshtp.Dbf_cnt.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库数据文件数量%d接近限制%d，使用率%.1f%%，建议及时调整增加db_files参数值", dbshtp.Dbname.Contents, fileCnt, dbFilesLimit, usagePercent))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBF_STAT 检查数据文件状态
func Ana_DBF_STAT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbf_stat.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbf_stat.Nm,
		Title:    rule.Dbrule.Dbf_stat.Title,
		Desc:     rule.Dbrule.Dbf_stat.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 1 {
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 5 {
				continue
			}
			dbfstatus := msgs[1]
			if dbfstatus != rule.Dbrule.Dbf_stat.Status {
				dbshtp.Dbf_stat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,数据文件%s当前状态为%s非AVAILABLE,需立即检查并修复", dbshtp.Dbname.Contents, msgs[0], dbfstatus))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBTbs 分析表空间使用情况
func Ana_DBTbs(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbtbsusage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbtbsusage.Nm,
		Title:    rule.Dbrule.Dbtbsusage.Title,
		Desc:     rule.Dbrule.Dbtbsusage.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 1 {
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 5 {
				continue
			}
			maxsize, _ := strconv.ParseFloat(msgs[1], 64)
			usedsize, _ := strconv.ParseFloat(msgs[3], 64)
			percent, _ := strconv.ParseFloat(msgs[5], 64)
			switch {
			case percent >= rule.Dbrule.Dbtbsusage.Tbsutil_ge2 && (maxsize-usedsize) < rule.Dbrule.Dbtbsusage.Freesize_le2:
				dbshtp.Dbtbsusage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,表空间%s使用率%.2f%%超过%.0f%%且剩余空间%.2fGB小于%.0fGB,需要及时扩容或清理数据", dbshtp.Dbname.Contents, msgs[0], percent, rule.Dbrule.Dbtbsusage.Tbsutil_ge2, maxsize-usedsize, rule.Dbrule.Dbtbsusage.Freesize_le2))
				break Looop
			case percent >= rule.Dbrule.Dbtbsusage.Tbsutil_ge1 && (maxsize-usedsize) < rule.Dbrule.Dbtbsusage.Freesize_le1:
				dbshtp.Dbtbsusage.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,表空间%s使用率%.2f%%超过%.0f%%且剩余空间%.2fGB小于%.0fGB,建议持续关注并准备扩容", dbshtp.Dbname.Contents, msgs[0], percent, rule.Dbrule.Dbtbsusage.Tbsutil_ge1, maxsize-usedsize, rule.Dbrule.Dbtbsusage.Freesize_le1))
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBCTRF 分析控制文件数量
func Ana_DBCTRF(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbcontrolfile.Contents
	index := len(strings.Split(msgdata, "\n"))
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbcontrolfile.Nm,
		Title:    rule.Dbrule.Dbcontrolfile.Title,
		Desc:     rule.Dbrule.Dbcontrolfile.Desc,
	}
	if index < rule.Dbrule.Dbcontrolfile.Cnt_le1 {
		dbshtp.Dbcontrolfile.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,控制文件数量当前%d个小于阈值%d个,建议配置至少2路冗余控制文件", dbshtp.Dbname.Contents, index, rule.Dbrule.Dbcontrolfile.Cnt_le1))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Tab_parallel 检查是否存在并行度大于0的表
func Ana_Tab_parallel(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Tab_parallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Tab_parallel.Nm,
		Title:    rule.Dbrule.Tab_parallel.Title,
		Desc:     rule.Dbrule.Tab_parallel.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Tab_parallel.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Tab_parallel.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,表%s并行度当前值%s大于0,建议设置并行度为0", dbshtp.Dbname.Contents, msgs[1], msgs[len(msgs)-1]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Inx_parallel 检查是否存在并行度大于0的索引
func Ana_Inx_parallel(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Inx_parallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Inx_parallel.Nm,
		Title:    rule.Dbrule.Inx_parallel.Title,
		Desc:     rule.Dbrule.Inx_parallel.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Inx_parallel.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Inx_parallel.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,索引%s并行度当前值%s大于0,建议设置并行度为0", dbshtp.Dbname.Contents, msgs[2], msgs[len(msgs)-1]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Invalid_obj 检查是否存在大量无效对象
func Ana_Invalid_obj(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Invalid_obj.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Invalid_obj.Nm,
		Title:    rule.Dbrule.Invalid_obj.Title,
		Desc:     rule.Dbrule.Invalid_obj.Desc,
	}

	// 从第3行开始检查第3列是否有大于10的值
	lines := strings.Split(msgdata, "\n")
	if len(lines) >= 3 {
		for i := 2; i < len(lines); i++ { // 从第3行开始（索引2）
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 3 {
				// 检查第3列（OBJ_COUNT）是否大于规则中设定的阈值
				if objCount, err := strconv.Atoi(fields[2]); err == nil {
					if threshold, err2 := strconv.Atoi(rule.Dbrule.Invalid_obj.Result); err2 == nil && objCount > threshold {
						dbshtp.Invalid_obj.Alarm = "B"
						entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,%s用户存在%d个无效对象,建议重新编译或及时清理",
							dbshtp.Dbname.Contents, fields[0], objCount, fields[1]))
					}
				}
			}
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Invalid_inx 分析无效索引
func Ana_Invalid_inx(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Invalid_inx.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Invalid_inx.Nm,
		Title:    rule.Dbrule.Invalid_inx.Title,
		Desc:     rule.Dbrule.Invalid_inx.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Invalid_inx.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Invalid_inx.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,索引%s当前状态为无效,建议重建或删除该索引", dbshtp.Dbname.Contents, msgs[2]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBSEQUENCE 分析序列配置
func Ana_DBSEQUENCE(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbsequence.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Dbsequence.Nm,
		Title:    rule.Dbrule.Dbsequence.Title,
		Desc:     rule.Dbrule.Dbsequence.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Dbsequence.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbsequence.Alarm = "G"
			msgs := strings.Fields(value)
			cache, _ := strconv.Atoi(msgs[1])
			if cache < 400 {
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,序列 %s cache %d 小于 400，建议设置为 400 或以上并使用 NOORDER", dbshtp.Dbname.Contents, msgs[0], cache))
			}
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DB_SEQ_USAGE 分析序列最大值使用情况
func Ana_DB_SEQ_USAGE(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbsequence.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Db_seq_usage.Nm,
		Title:    rule.Dbrule.Db_seq_usage.Title,
		Desc:     rule.Dbrule.Db_seq_usage.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Db_seq_usage.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbsequence.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,SEQUENCE %s 当前值达到最大值限制的80%%,建议尽快修改", dbshtp.Dbname.Contents, msgs[0]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
