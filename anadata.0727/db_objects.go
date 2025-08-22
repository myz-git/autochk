package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// db_objects.go 包含与数据库对象（表、索引、序列）管理相关的分析函数

// Ana_DBTABLEPARALLEL 分析表并行度
func Ana_DBTABLEPARALLEL(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbtableparallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Dbtableparallel.Nm,
		Title:    rule.Dbrule.Dbtableparallel.Title,
		Desc:     rule.Dbrule.Dbtableparallel.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Dbtableparallel.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbtableparallel.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("表 %s 并行度 %s 大于 0，建议设置为 0", msgs[1], msgs[len(msgs)-1]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBINDEXPARALLEL 分析索引并行度
func Ana_DBINDEXPARALLEL(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbindexparallel.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Dbindexparallel.Nm,
		Title:    rule.Dbrule.Dbindexparallel.Title,
		Desc:     rule.Dbrule.Dbindexparallel.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Dbindexparallel.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbindexparallel.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("索引 %s 并行度 %s 大于 0，建议设置为 0", msgs[2], msgs[len(msgs)-1]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBINVALIDINDEX 分析无效索引
func Ana_DBINVALIDINDEX(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbinvalidindex.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Dbinvalidindex.Nm,
		Title:    rule.Dbrule.Dbinvalidindex.Title,
		Desc:     rule.Dbrule.Dbinvalidindex.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Dbinvalidindex.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbinvalidindex.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("索引 %s 无效，建议重建或删除", msgs[2]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBSEQUENCE 分析序列配置
func Ana_DBSEQUENCE(rule *utils.RuleInfo, dbshtp *structs.DbSht, infstp *structs.InfoSht, summaryEntries *structs.SummaryEntries) {
	if !strings.Contains(infstp.DbMaa, "RAC") {
		return
	}
	msgdata := dbshtp.Dbsequence.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
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
				entry.Minor = append(entry.Minor, fmt.Sprintf("序列 %s cache %d 小于 400，建议设置为 400 或以上并使用 NOORDER", msgs[0], cache))
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
	msgdata := dbshtp.Db_seq_usage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
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
			dbshtp.Db_seq_usage.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("SEQUENCE %s 当前值达到最大值限制的80%%", msgs[0]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
