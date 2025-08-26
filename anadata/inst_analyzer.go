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
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,归档切换次数当前%d次超过阈值%d次,建议调整为15-20分钟一次", instshtp.Instname.Contents, value, rule.Dbrule.Dbredoswitch.Sw_cnt_ge1))
					break Looop
				}
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RDF 分析 REDO 文件状态和大小
func Ana_RDF(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dbredoswitch.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbredocheck.Nm,
		Title:    rule.Dbrule.Dbredocheck.Title,
		Desc:     rule.Dbrule.Dbredocheck.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		value = strings.TrimSpace(value)
		rd := regexp.MustCompile(`^\d`)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 6 {
				continue
			}
			data := msgs[5]
			if !utils.Contain(data, rule.Dbrule.Dbredocheck.Rdf_status_list) {
				instshtp.Dbredoswitch.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("REDO 文件 %s 状态 %s 异常，需检查", msgs[0], data))
				break Looop
			}
			matchs, _ := strconv.ParseFloat(msgs[3], 64)
			if matchs < rule.Dbrule.Dbredocheck.Rdf_size_lt1 {
				instshtp.Dbredoswitch.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("REDO 文件 %s 大小 %.2f MB 小于 %f MB，建议调整", msgs[0], matchs, rule.Dbrule.Dbredocheck.Rdf_size_lt1))
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBparameter 分析数据库初始化参数
func Ana_DBparameter(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbrecoverydest.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dbparameter.Nm,
		Title:    rule.Dbrule.Dbparameter.Title,
		Desc:     rule.Dbrule.Dbparameter.Desc,
	}
	if msgdata == "" {
		dbshtp.Dbrecoverydest.Alarm = "B"
		entry.Moderate = append(entry.Moderate, "数据库初始化参数未配置，建议检查")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBParameterFile 分析数据库参数文件
func Ana_DBParameterFile(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbrecoverydest.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Db_parameter_file.Nm,
		Title:    rule.Dbrule.Db_parameter_file.Title,
		Desc:     rule.Dbrule.Db_parameter_file.Desc,
	}
	if msgdata == "" {
		dbshtp.Dbrecoverydest.Alarm = "B"
		entry.Moderate = append(entry.Moderate, "数据库初始化参数文件未配置，建议检查")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBShpSize 分析共享池大小
func Ana_DBShpSize(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbrecoverydest.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Db_shp_size.Nm,
		Title:    rule.Dbrule.Db_shp_size.Title,
		Desc:     rule.Dbrule.Db_shp_size.Desc,
	}
	if strings.Contains(msgdata, "shared_pool_size") {
		re := regexp.MustCompile(`shared_pool_size\s*=\s*(\d+)`)
		if matches := re.FindStringSubmatch(msgdata); len(matches) > 1 {
			size, _ := strconv.Atoi(matches[1])
			if size < 1000000000 {
				dbshtp.Dbrecoverydest.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("共享池大小 %d 字节过小，建议调整", size))
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBShpPct 分析共享池使用率
func Ana_DBShpPct(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbrecoverydest.Contents
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Db_shp_pct.Nm,
		Title:    rule.Dbrule.Db_shp_pct.Title,
		Desc:     rule.Dbrule.Db_shp_pct.Desc,
	}
	if strings.Contains(msgdata, "shared_pool_size") {
		re := regexp.MustCompile(`shared_pool_size\s*=\s*(\d+)`)
		if matches := re.FindStringSubmatch(msgdata); len(matches) > 1 {
			size, _ := strconv.Atoi(matches[1])
			reSga := regexp.MustCompile(`sga_max_size\s*=\s*(\d+)`)
			if sgaMatches := reSga.FindStringSubmatch(msgdata); len(sgaMatches) > 1 {
				sgaSize, _ := strconv.Atoi(sgaMatches[1])
				if float64(size)/float64(sgaSize)*100 > 60 {
					dbshtp.Dbrecoverydest.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("共享池使用率 %.2f%% 超过 60%%，建议优化", float64(size)/float64(sgaSize)*100))
				}
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
