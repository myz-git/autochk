package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"strings"
)

// Ana_RAC_status 分析集群状态
func Ana_RAC_status(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbcrscheck.Contents
	entry := structs.SummaryEntry{
		Category: "集群检查",
		Nm:       rule.Dbrule.Dbcrscheck.Nm,
		Title:    rule.Dbrule.Dbcrscheck.Title,
		Desc:     rule.Dbrule.Dbcrscheck.Desc,
	}
	if strings.Contains(msgdata, "UNKNOWN") {
		dbshtp.Dbcrscheck.Alarm = "R"
		entry.Severe = append(entry.Severe, "集群存在 UNKNOWN 状态服务，需检查集群服务状态异常原因")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBASMUSAGE 分析 ASM 使用情况
func Ana_ASM_usage(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
}
