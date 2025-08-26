package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// db_monitoring.go 包含数据库错误监控、DataGuard、备份及杂项分析函数

// 错误监控

// Ana_DBERRLOG 分析数据库错误日志
func Ana_DBERRLOG(rule *utils.RuleInfo, instshtp *structs.InstShts, summaryEntries *structs.SummaryEntries) {
	msgdata := instshtp.Dberrlog.Contents
	value := strings.TrimSpace(msgdata)
	entry := structs.SummaryEntry{
		Category: "实例分析",
		Nm:       rule.Dbrule.Dberrlog.Nm,
		Title:    rule.Dbrule.Dberrlog.Title,
		Desc:     rule.Dbrule.Dberrlog.Desc,
	}
	instshtp.Dberrlog.Alarm = "B"
	if value == "" || strings.Contains(value, rule.Dbrule.Dberrlog.ResultB) {
		instshtp.Dberrlog.Alarm = ""
	} else {
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,数据库日志存在重要报错信息,建议检查并处理相关错误", instshtp.Instname.Contents))
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
	rd := regexp.MustCompile(`(?i)\.log$`)
	rd2 := regexp.MustCompile(`^Jan$|^Feb$|^Mar$|^Apr$|^May$|^Jun$|^Jul$|^Aug$|^Sept$|^Oct$|^Nov$|^Dec$`)
	rd3 := regexp.MustCompile(`^\d+$`)
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		value = strings.TrimSpace(value)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 8 {
				continue
			}
			for k := 3; k < len(msgs); k++ {
				if rd2.MatchString(msgs[k]) && rd3.MatchString(msgs[k-1]) {
					data, _ := strconv.Atoi(msgs[k-1])
					if data >= rule.Dbrule.Dblsnrinfo.Log_size {
						instshtp.Dblsnrinfo.Alarm = "G"
						entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,监听日志文件%s当前大小%d bytes超过阈值%d bytes,建议定期清理日志文件", instshtp.Instname.Contents, msgs[0], data, rule.Dbrule.Dblsnrinfo.Log_size))
						break Looop
					}
				}
			}
		}
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
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s实例,DataGuard同步延迟当前%d天超过阈值%d天,建议检查网络连接和主库性能", instshtp.Instname.Contents, vDay, rule.Dbrule.Dbdglagcheck.ResultB))
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
	if value != "" {
		instshtp.Dbdgerrcheck.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s实例,DataGuard日志存在同步错误信息,建议检查并处理相关错误", instshtp.Instname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBRECOVERYDEST 分析数据库恢复目录
func Ana_DBRECOVERYDEST(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
}

// Ana_DBFLASHRECOVERYUSEAGE 分析闪回区使用率
func Ana_DBFLASHRECOVERYUSEAGE(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbflashrecoveryuseage.Contents
	entry := structs.SummaryEntry{
		Category: "数据库备份",
		Nm:       rule.Dbrule.Dbflashrecoveryuseage.Nm,
		Title:    rule.Dbrule.Dbflashrecoveryuseage.Title,
		Desc:     rule.Dbrule.Dbflashrecoveryuseage.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 3 {
				continue
			}
			data, err := strconv.ParseFloat(msgs[len(msgs)-1], 64)
			if err != nil {
				log.Fatal(err)
			}
			if data >= rule.Dbrule.Dbflashrecoveryuseage.Useage1 {
				dbshtp.Dbflashrecoveryuseage.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,闪回区使用率当前%.2f%%超过阈值%.0f%%,建议及时清理或扩容闪回区", dbshtp.Dbname.Contents, data, rule.Dbrule.Dbflashrecoveryuseage.Useage1))
			}
			if data >= rule.Dbrule.Dbflashrecoveryuseage.Useage2 {
				dbshtp.Dbflashrecoveryuseage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s数据库,闪回区使用率当前%.2f%%超过严重阈值%.0f%%,需立即清理或扩容闪回区", dbshtp.Dbname.Contents, data, rule.Dbrule.Dbflashrecoveryuseage.Useage2))
				break Looop
			}
		}
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
	// 实现PSU分析逻辑
	if strings.TrimSpace(msgdata) != "" {
		instshtp.Dbpsu.Alarm = "G"
		entry.Minor = append(entry.Minor, "PSU补丁信息检查")
	}
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
	if strings.TrimSpace(msgdata) != "" {
		instshtp.Dbpatch.Alarm = "G"
		entry.Minor = append(entry.Minor, "数据库补丁信息检查")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
