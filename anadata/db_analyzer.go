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
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,数据库状态检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 检查第三行是否包含期望的状态
	dataLine := strings.TrimSpace(lines[2])
	if dataLine == "" {
		// 第三行为空，数据采集异常
		dbshtp.Dbstatus.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,数据库状态检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	if !strings.Contains(dataLine, rule.Dbrule.Dbstatus.Status) {
		dbshtp.Dbstatus.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,当前状态异常,\n建议: 尽快确认数据库是否运行正常", dbshtp.Dbname.Contents))
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
		dbshtp.Logmode.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,当前未开启归档模式,\n建议: 生产环境需要开启归档模式以确保数据安全", dbshtp.Dbname.Contents))
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
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,数据文件数量检查数据采集异常", dbshtp.Dbname.Contents))
		return
	}

	// 检查第三行数据
	dataLine := strings.TrimSpace(lines[2])
	if dataLine == "" {
		// 第三行为空，数据采集异常
		dbshtp.Dbf_cnt.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,数据文件数量检查数据采集异常", dbshtp.Dbname.Contents))
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
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库数据文件数量%d接近限制%d，使用率%.1f%%,\n建议: 及时调整增加db_files参数值", dbshtp.Dbname.Contents, fileCnt, dbFilesLimit, usagePercent))
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

	// 空/无记录判断
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		dbshtp.Dbf_stat.Alarm = ""
		return
	}

	// 从第3行开始检查是否存在有效记录（非仅由空白和-组成的分隔线）
	lines := strings.Split(msgdata, "\n")
	sepRe := regexp.MustCompile(`^[\s-]+$`)
	for i, line := range lines {
		if i < 2 {
			continue
		}
		val := strings.TrimSpace(line)
		if val == "" {
			continue
		}
		if sepRe.MatchString(val) {
			continue
		}
		// 发现首条有效记录，直接判定为严重告警并返回（仅生成一条告警）
		dbshtp.Dbf_stat.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,存在数据文件状态异常记录,\n建议: 需立即检查并修复", dbshtp.Dbname.Contents))
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
		return
	}

	// 若第3行起未发现有效记录，则视为正常
	dbshtp.Dbf_stat.Alarm = ""
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
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,表空间%s使用率%.2f%%超过%.0f%%且剩余空间%.2fGB小于%.0fGB,\n建议: 需要及时扩容或清理数据", dbshtp.Dbname.Contents, msgs[0], percent, rule.Dbrule.Dbtbsusage.Tbsutil_ge2, maxsize-usedsize, rule.Dbrule.Dbtbsusage.Freesize_le2))
				break Looop
			case percent >= rule.Dbrule.Dbtbsusage.Tbsutil_ge1 && (maxsize-usedsize) < rule.Dbrule.Dbtbsusage.Freesize_le1:
				dbshtp.Dbtbsusage.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,表空间%s使用率%.2f%%超过%.0f%%且剩余空间%.2fGB小于%.0fGB,\n建议: 持续关注并计划扩容", dbshtp.Dbname.Contents, msgs[0], percent, rule.Dbrule.Dbtbsusage.Tbsutil_ge1, maxsize-usedsize, rule.Dbrule.Dbtbsusage.Freesize_le1))
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
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,控制文件数量当前%d个小于阈值%d个,\n建议: 配置至少2路冗余控制文件", dbshtp.Dbname.Contents, index, rule.Dbrule.Dbcontrolfile.Cnt_le1))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Tab_parallel 检查是否存在并行度大于1的表
func Ana_Tab_parallel(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Tab_parallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Tab_parallel.Nm,
		Title:    rule.Dbrule.Tab_parallel.Title,
		Desc:     rule.Dbrule.Tab_parallel.Desc,
	}

	// 检查是否为空或包含"无记录"、"no rows selected"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Tab_parallel.Alarm = ""
		return
	}

	// 解析数字
	if count, err := strconv.Atoi(strings.TrimSpace(msgdata)); err == nil {
		// 检查数字是否大于规则中设定的阈值
		if count > rule.Dbrule.Tab_parallel.Result {
			dbshtp.Tab_parallel.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,存在%d个并行度大于1的表,\n建议: 对于并行度大于1的表关闭并行属性", dbshtp.Dbname.Contents, count))
		} else {
			// 正常情况，无告警
			dbshtp.Tab_parallel.Alarm = ""
		}
	} else {
		// 数据格式异常
		dbshtp.Tab_parallel.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,表并行度检查数据格式异常,\n建议: 检查数据采集是否正常", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Inx_parallel 检查是否存在并行度大于1的索引
func Ana_Inx_parallel(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Inx_parallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库分析",
		Nm:       rule.Dbrule.Inx_parallel.Nm,
		Title:    rule.Dbrule.Inx_parallel.Title,
		Desc:     rule.Dbrule.Inx_parallel.Desc,
	}

	// 检查是否为空或包含"无记录"、"no rows selected"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Inx_parallel.Alarm = ""
		return
	}

	// 解析数字
	if count, err := strconv.Atoi(strings.TrimSpace(msgdata)); err == nil {
		// 检查数字是否大于规则中设定的阈值
		if count > rule.Dbrule.Inx_parallel.Result {
			dbshtp.Inx_parallel.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,存在%d个并行度大于1的索引,\n建议: 对于并行度大于1的索引关闭并行属性", dbshtp.Dbname.Contents, count))
		} else {
			// 正常情况，无告警
			dbshtp.Inx_parallel.Alarm = ""
		}
	} else {
		// 数据格式异常
		dbshtp.Inx_parallel.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,索引并行度检查数据格式异常,\n建议: 检查数据采集是否正常", dbshtp.Dbname.Contents))
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

	// 检查是否为空或包含"无记录"、"no rows selected"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Invalid_obj.Alarm = ""
		return
	}

	// 解析数字
	if count, err := strconv.Atoi(strings.TrimSpace(msgdata)); err == nil {
		// 检查数字是否大于规则中设定的阈值
		if count > rule.Dbrule.Invalid_obj.Result {
			dbshtp.Invalid_obj.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,存在%d个无效对象,\n建议: 对无效对象进行重新编译或及时清理", dbshtp.Dbname.Contents, count))
		} else {
			// 正常情况，无告警
			dbshtp.Invalid_obj.Alarm = ""
		}
	} else {
		// 数据格式异常
		dbshtp.Invalid_obj.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,无效对象检查数据格式异常,\n建议: 检查数据采集是否正常", dbshtp.Dbname.Contents))
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

	// 检查是否为空或包含"无记录"、"no rows selected"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Invalid_inx.Alarm = ""
		return
	}

	// 解析数字
	if count, err := strconv.Atoi(strings.TrimSpace(msgdata)); err == nil {
		// 检查数字是否大于规则中设定的阈值
		if count > rule.Dbrule.Invalid_inx.Result {
			dbshtp.Invalid_inx.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,存在%d个无效索引,\n建议: 重建或删除失效索引", dbshtp.Dbname.Contents, count))
		} else {
			// 正常情况，无告警
			dbshtp.Invalid_inx.Alarm = ""
		}
	} else {
		// 数据格式异常
		dbshtp.Invalid_inx.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,无效索引检查数据格式异常,\n建议: 检查数据采集是否正常", dbshtp.Dbname.Contents))
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
				entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,序列 %s cache %d 小于 400,\n建议: 设置CACHE为400或以上并使用NOORDER", dbshtp.Dbname.Contents, msgs[0], cache))
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
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,SEQUENCE %s 当前值达到最大值限制的80%%,\n建议: 尽快修改SEQUENCE最大值", dbshtp.Dbname.Contents, msgs[0]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// parseParameters 解析参数数据为键值对
func parseParameters(data string) map[string]string {
	params := make(map[string]string)
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "<") || strings.HasPrefix(line, "</") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			params[key] = value
		}
	}

	return params
}

// Ana_DBparam_b 分析数据库参数组basic
func Ana_DBparam_b(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbparam_b.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbparam_b.Nm,
		Title:    rule.Dbrule.Dbparam_b.Title,
		Desc:     rule.Dbrule.Dbparam_b.Desc,
	}

	// 解析参数数据
	params := parseParameters(msgdata)

	// 检查 memory_max_target = 0
	if value, exists := params["memory_max_target"]; exists {
		if value > "0" {
			dbshtp.Dbparam_b.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,启用了AMM内存自动调整(memory_max_target参数设置为%s,\n建议: 关闭AMM内存自动调整(设置memory_max_target=%d)", dbshtp.Dbname.Contents, value, rule.Dbrule.Dbparam_b.Memory_max_target))
		}
	} else {
		// 参数未设置 正常情况，无告警
		dbshtp.Dbparam_b.Alarm = ""
	}

	// 检查 sga_max_size <= 8589934592 (8GB)
	if value, exists := params["sga_max_size"]; exists {
		if size, err := strconv.ParseInt(value, 10, 64); err == nil {
			if size <= 8589934592 {
				dbshtp.Dbparam_b.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,sga_max_size参数%.0f字节(%.2fGB)小于等于8GB,\n建议: sga_max_size设置大于8GB", dbshtp.Dbname.Contents, float64(size), float64(size)/1024/1024/1024))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_b.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,sga_max_size参数未设置,\n建议: sga_max_size设置大于8GB", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBparam_d 分析数据库参数组1
func Ana_DBparam_d(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbparam_d.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbparam_d.Nm,
		Title:    rule.Dbrule.Dbparam_d.Title,
		Desc:     rule.Dbrule.Dbparam_d.Desc,
	}

	// 解析参数数据
	params := parseParameters(msgdata)

	// 检查 _and_pruning_enabled != FALSE (大小写不敏感)
	if value, exists := params["_and_pruning_enabled"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_and_pruning_enabled参数设置为%s,\n建议: _and_pruning_enabled设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_and_pruning_enabled参数未设置,\n建议: _and_pruning_enabled设置为FALSE", dbshtp.Dbname.Contents))
	}

	// _ash_size < 67108864 (64MB)
	if value, exists := params["_ash_size"]; exists {
		if size, err := strconv.Atoi(value); err == nil {
			if size < 67108864 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_ash_size参数%d小于64MB,\n建议: _ash_size设置不少于64MB", dbshtp.Dbname.Contents, size))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_ash_size参数未设置,\n建议: _ash_size设置不少于64MB", dbshtp.Dbname.Contents))
	}

	// _bloom_filter_enabled != FALSE (大小写不敏感)
	if value, exists := params["_bloom_filter_enabled"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_bloom_filter_enabled参数设置为%s,\n建议: _bloom_filter_enabled设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_bloom_filter_enabled参数未设置,\n建议: _bloom_filter_enabled设置为FALSE", dbshtp.Dbname.Contents))
	}

	// _bloom_pruning_enabled != FALSE (大小写不敏感)
	if value, exists := params["_bloom_pruning_enabled"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_bloom_pruning_enabled参数设置为%s,\n建议: _bloom_pruning_enabled设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_bloom_pruning_enabled参数未设置,\n建议: _bloom_pruning_enabled设置为FALSE", dbshtp.Dbname.Contents))
	}

	// _cleanup_rollback_entries < 2000
	if value, exists := params["_cleanup_rollback_entries"]; exists {
		if size, err := strconv.Atoi(value); err == nil {
			if size < 2000 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_cleanup_rollback_entries参数%d小于2000,\n建议: _cleanup_rollback_entries设置不少于2000", dbshtp.Dbname.Contents, size))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_cleanup_rollback_entries参数未设置,\n建议: _cleanup_rollback_entries设置不少于2000", dbshtp.Dbname.Contents))
	}

	// _cursor_obsolete_threshold > 1024
	if value, exists := params["_cursor_obsolete_threshold"]; exists {
		if size, err := strconv.Atoi(value); err == nil {
			if size > 1024 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_cursor_obsolete_threshold参数%d大于1024,\n建议: _cursor_obsolete_threshold设置不超过1024", dbshtp.Dbname.Contents, size))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_cursor_obsolete_threshold参数未设置,\n建议: _cursor_obsolete_threshold设置不超过1024", dbshtp.Dbname.Contents))
	}

	// _optimizer_gather_feedback
	if value, exists := params["_optimizer_gather_feedback"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_optimizer_gather_feedback参数设置为%s,\n建议: _optimizer_gather_feedback设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_optimizer_gather_feedback参数未设置,\n建议: _optimizer_gather_feedback设置为FALSE", dbshtp.Dbname.Contents))
	}

	// _rowsets_enabled
	if value, exists := params["_rowsets_enabled"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_rowsets_enabled参数设置为%s,\n建议: _rowsets_enabled设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_rowsets_enabled参数未设置,\n建议: _rowsets_enabled设置为FALSE", dbshtp.Dbname.Contents))
	}

	// 检查 _shared_pool_reserved_pct <= 5
	if value, exists := params["_shared_pool_reserved_pct"]; exists {
		if pct, err := strconv.Atoi(value); err == nil {
			if pct <= 5 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_shared_pool_reserved_pct参数%d%%小于等于5%%,\n建议: 设置_shared_pool_reserved_pct为15%%", dbshtp.Dbname.Contents, pct))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_shared_pool_reserved_pct参数未设置,\n建议: 设置_shared_pool_reserved_pct为15%%", dbshtp.Dbname.Contents))
	}

	// 检查 _max_spacebg_slaves > 100
	if value, exists := params["_max_spacebg_slaves"]; exists {
		if slaves, err := strconv.Atoi(value); err == nil {
			if slaves > 100 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_max_spacebg_slaves参数%d大于100,\n建议: _max_spacebg_slaves设置不超过100", dbshtp.Dbname.Contents, slaves))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_max_spacebg_slaves参数未设置,\n建议: _max_spacebg_slaves设置不超过100", dbshtp.Dbname.Contents))
	}

	// 检查 _undo_autotune != False (大小写不敏感)
	if value, exists := params["_undo_autotune"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_undo_autotune参数设置为%s,\n建议: _undo_autotune设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_undo_autotune参数未设置,\n建议: _undo_autotune设置为FALSE", dbshtp.Dbname.Contents))
	}

	// 检查 _use_adaptive_log_file_sync != False (大小写不敏感)
	if value, exists := params["_use_adaptive_log_file_sync"]; exists {
		if strings.ToUpper(value) != "FALSE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_use_adaptive_log_file_sync参数设置为%s,\n建议: _use_adaptive_log_file_sync设置为FALSE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_use_adaptive_log_file_sync参数未设置,\n建议: _use_adaptive_log_file_sync设置为FALSE", dbshtp.Dbname.Contents))
	}

	// 检查 _use_single_log_writer != TRUE (大小写不敏感)
	if value, exists := params["_use_single_log_writer"]; exists {
		if strings.ToUpper(value) != "TRUE" {
			dbshtp.Dbparam_d.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_use_single_log_writer参数设置为%s,\n建议: _use_single_log_writer设置为TRUE", dbshtp.Dbname.Contents, value))
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,_use_single_log_writer参数未设置,\n建议: _use_single_log_writer设置为TRUE", dbshtp.Dbname.Contents))
	}

	// 检查parallel_max_servers > 128
	if value, exists := params["parallel_max_servers"]; exists {
		if servers, err := strconv.Atoi(value); err == nil {
			if servers > 128 {
				dbshtp.Dbparam_d.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,parallel_max_servers参数%d大于128,\n建议: parallel_max_servers设置不超过128", dbshtp.Dbname.Contents, servers))
			}
		}
	} else {
		// 参数未设置，也算G级告警
		dbshtp.Dbparam_d.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,parallel_max_servers参数未设置,\n建议: parallel_max_servers设置不超过128", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
