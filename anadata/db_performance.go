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

// db_performance.go 包含与数据库性能和效率相关的分析函数
// Ana_DB4031check 分析 ORA-4031 错误
func Ana_DB4031check(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dberrlog.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Db_4031check.Nm,
		Title:    rule.Dbrule.Db_4031check.Title,
		Desc:     rule.Dbrule.Db_4031check.Desc,
	}
	if strings.Contains(msgdata, "ORA-4031") {
		dbshtp.Dberrlog.Alarm = "R"
		entry.Severe = append(entry.Severe, "检测到 ORA-4031 错误，需调整共享池或清理 LRU 列表")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RESOURCE 分析数据库资源使用情况
func Ana_RESOURCE(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbresource.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Dbresource.Nm,
		Title:    rule.Dbrule.Dbresource.Title,
		Desc:     rule.Dbrule.Dbresource.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		value = strings.TrimSpace(value)
		rd := regexp.MustCompile(`\d$`)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 5 {
				continue
			}
			data1, err := strconv.Atoi(msgs[3])
			if err != nil {
				log.Fatal(err)
			}
			data2, err := strconv.Atoi(msgs[4])
			if err != nil {
				log.Fatal(err)
			}
			if data1 >= data2*rule.Dbrule.Dbresource.Res_use_ge1/100 && data2 != 0 {
				dbshtp.Dbresource.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("资源 %s 使用率 %d 接近限制值 %d，需优化", msgs[1], data1, data2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_LOADPROFILE 分析数据库负载性能
func Ana_LOADPROFILE(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Loadprofile.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Loadprofile.Nm,
		Title:    rule.Dbrule.Loadprofile.Title,
		Desc:     rule.Dbrule.Loadprofile.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		value = strings.TrimSpace(value)
		rd := regexp.MustCompile(`\d$`)
		if rd.MatchString(value) {
			msgs := strings.Split(value, ":")
			if len(msgs) < 2 {
				continue
			}
			submsg := strings.Fields(msgs[1])
			str := strings.Replace(submsg[0], ",", "", -1)
			data, err := strconv.ParseFloat(str, 64)
			if err != nil {
				log.Fatal(err)
			}
			if msgs[0] == "Redo size (bytes)" && data >= rule.Dbrule.Loadprofile.Redosize_ge1*1024*1024 {
				dbshtp.Loadprofile.Alarm = "G"
				entry.Minor = append(entry.Minor, fmt.Sprintf("Redo size %.2f bytes 超过 %f MB，建议优化", data, rule.Dbrule.Loadprofile.Redosize_ge1))
			}
			if msgs[0] == "Logons" && data >= rule.Dbrule.Loadprofile.Logons_ge1 {
				dbshtp.Loadprofile.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("登录数 %.2f 超过 %f，建议优化连接管理", data, rule.Dbrule.Loadprofile.Logons_ge1))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_INSTEFFICIENCY 分析数据库实例效率
func Ana_INSTEFFICIENCY(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Instefficiency.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Instefficiency.Nm,
		Title:    rule.Dbrule.Instefficiency.Title,
		Desc:     rule.Dbrule.Instefficiency.Desc,
	}
	rd := regexp.MustCompile(`\d$`)
	rd2 := regexp.MustCompile(`^[1-9]\d*\.\d+$|^0\.\d+$|^[1-9]\d*$|^0$`)
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		value = strings.TrimSpace(value)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 4 {
				continue
			}
			for k := 0; k < len(msgs)-3; k++ {
				if msgs[k] == "Buffer" && msgs[k+1] == "Hit" && rd2.MatchString(msgs[k+3]) {
					data, err := strconv.ParseFloat(msgs[k+3], 64)
					if err != nil {
						log.Fatal(err)
					}
					if data < rule.Dbrule.Instefficiency.Buffer_hit {
						dbshtp.Instefficiency.Alarm = "G"
						entry.Minor = append(entry.Minor, fmt.Sprintf("Buffer Hit 命中率 %.2f%% 小于 %f%%，建议优化", data, rule.Dbrule.Instefficiency.Buffer_hit))
						break Looop
					}
				}
				if msgs[k] == "Library" && msgs[k+1] == "Hit" && rd2.MatchString(msgs[k+3]) {
					data, err := strconv.ParseFloat(msgs[k+3], 64)
					if err != nil {
						log.Fatal(err)
					}
					if data < rule.Dbrule.Instefficiency.Library_hit {
						dbshtp.Instefficiency.Alarm = "G"
						entry.Minor = append(entry.Minor, fmt.Sprintf("Library Hit 命中率 %.2f%% 小于 %f%%，建议优化", data, rule.Dbrule.Instefficiency.Library_hit))
						break Looop
					}
				}
				if msgs[k] == "Soft" && msgs[k+1] == "Parse" && rd2.MatchString(msgs[k+3]) {
					data, err := strconv.ParseFloat(msgs[k+3], 64)
					if err != nil {
						log.Fatal(err)
					}
					if data < rule.Dbrule.Instefficiency.Soft_parse {
						dbshtp.Instefficiency.Alarm = "G"
						entry.Minor = append(entry.Minor, fmt.Sprintf("Soft Parse 命中率 %.2f%% 小于 %f%%，建议优化", data, rule.Dbrule.Instefficiency.Soft_parse))
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

// Ana_DBtopevent 分析顶部等待事件
func Ana_DBtopevent(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbtopevent.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Topevent.Nm,
		Title:    rule.Dbrule.Topevent.Title,
		Desc:     rule.Dbrule.Topevent.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 3 {
				continue
			}
			executions, _ := strconv.Atoi(msgs[1])
			avgTime, _ := strconv.ParseFloat(msgs[2], 64)
			if executions > 1000 && avgTime > 2 {
				dbshtp.DbtopSQL.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("SQL 执行次数 %d 次，平均耗时 %.2f 秒，建议优化", executions, avgTime))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBtopSQL 分析顶部 SQL 性能
func Ana_DBtopSQL(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.DbtopSQL.Contents
	entry := structs.SummaryEntry{
		Category: "数据库实例分析",
		Nm:       rule.Dbrule.Topsqlbyelapstime.Nm,
		Title:    rule.Dbrule.Topsqlbyelapstime.Title,
		Desc:     rule.Dbrule.Topsqlbyelapstime.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 {
			continue
		}
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 3 {
				continue
			}
			executions, _ := strconv.Atoi(msgs[1])
			avgTime, _ := strconv.ParseFloat(msgs[2], 64)
			if executions > 1000 && avgTime > 2 {
				dbshtp.DbtopSQL.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("SQL 执行次数 %d 次，平均耗时 %.2f 秒，建议优化", executions, avgTime))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
