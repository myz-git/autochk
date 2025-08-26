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

// db_security.go 包含与数据库安全检查相关的分析函数

// Ana_DBExpirUser 分析用户密码过期情况
func Ana_DBExpirUser(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Db_expir_user.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全检查",
		Nm:       rule.Dbrule.Db_expir_user.Nm,
		Title:    rule.Dbrule.Db_expir_user.Title,
		Desc:     rule.Dbrule.Db_expir_user.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 1 {
			continue
		}
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			if len(msgs) < 3 {
				continue
			}
			days, _ := strconv.Atoi(msgs[len(msgs)-1])
			if days < 30 {
				dbshtp.Db_expir_user.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,用户%s口令将在%d天后过期,建议提前处理密码更新", dbshtp.Dbname.Contents, msgs[0], days))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DB_PASSWORD_VERIF 分析密码复杂性验证
func Ana_DB_PASSWORD_VERIF(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Db_password_verif.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全检查",
		Nm:       rule.Dbrule.Db_password_verif.Nm,
		Title:    rule.Dbrule.Db_password_verif.Title,
		Desc:     rule.Dbrule.Db_password_verif.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Db_password_verif.Alarm = ""
	} else {
		// 有记录说明存在未设置密码复杂性校验，设置为B级告警
		dbshtp.Db_password_verif.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,存在未设置密码复杂性校验的用户,建议启用密码复杂性验证以增强安全性", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBPRODUCTUSERFAILEDLOGIN 分析用户登录失败限制
func Ana_DBPRODUCTUSERFAILEDLOGIN(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbproductuserfailedlogin.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbproductuserfailedlogin.Nm,
		Title:    rule.Dbrule.Dbproductuserfailedlogin.Title,
		Desc:     rule.Dbrule.Dbproductuserfailedlogin.Desc,
	}
	rd := regexp.MustCompile(` \d+$`)
Looop:
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)
		if value == rule.Dbrule.Dbproductuserfailedlogin.Result {
			break Looop
		}
		if rd.MatchString(value) {
			dbshtp.Dbproductuserfailedlogin.Alarm = "B"
			msgs := strings.Fields(value)
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,用户%s错误登录次数限制当前为%s,建议调整为有限值以增强安全性", dbshtp.Dbname.Contents, msgs[0], msgs[len(msgs)-1]))
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBDBAPRIV 分析 DBA 权限
func Ana_DBDBAPRIV(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbdbapriv.Contents
	value := strings.TrimSpace(msgdata)
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbdbapriv.Nm,
		Title:    rule.Dbrule.Dbdbapriv.Title,
		Desc:     rule.Dbrule.Dbdbapriv.Desc,
	}
	dbshtp.Dbdbapriv.Alarm = "G"
	if value == "" || strings.Contains(value, rule.Dbrule.Dbdbapriv.ResultG) {
		dbshtp.Dbdbapriv.Alarm = ""
	} else {
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,存在具有DBA权限的业务账户,建议收回不必要的DBA权限以增强安全性", dbshtp.Dbname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBSYSDBA 分析 SYSDBA 权限用户
func Ana_DBSYSDBA(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbsysdba.Contents
	value := strings.TrimSpace(msgdata)
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbsysdba.Nm,
		Title:    rule.Dbrule.Dbsysdba.Title,
		Desc:     rule.Dbrule.Dbsysdba.Desc,
	}
	dbshtp.Dbsysdba.Alarm = "B"
	if value == "" || strings.Contains(value, rule.Dbrule.Dbsysdba.ResultB) {
		dbshtp.Dbsysdba.Alarm = ""
	} else {
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s数据库,存在非必要SYSDBA权限用户,建议检查并收回不必要的SYSDBA权限", dbshtp.Dbname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBAUDITSEGMENT 分析审计段
func Ana_DBAUDITSEGMENT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbauditsegment.Contents
	value := strings.TrimSpace(msgdata)
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbauditsegment.Nm,
		Title:    rule.Dbrule.Dbauditsegment.Title,
		Desc:     rule.Dbrule.Dbauditsegment.Desc,
	}
	if value != "" { //判断是否空, 为空正常, 非空则标记后退出
		dbshtp.Dbauditsegment.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,审计段存在异常信息,建议检查审计配置", dbshtp.Dbname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}

}

// Ana_DBAUDITCONT 分析审计内容
func Ana_DBAUDITCONT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := strings.TrimSpace(dbshtp.Dbauditcont.Contents) //去除头尾空格及空行
	rd := regexp.MustCompile(` \d+$`)                         //匹配以空格+数字结尾
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbauditcont.Nm,
		Title:    rule.Dbrule.Dbauditcont.Title,
		Desc:     rule.Dbrule.Dbauditcont.Desc,
	}
Looop:
	for index, row := range strings.Split(msgdata, "\n") { //按行分割
		if index < 2 { //跳过前面2行 (如 column head  和 ----)
			continue
		}
		if rd.MatchString(row) { //匹配以"空格+数字"结尾的行
			msgs := strings.Fields(row) // 以空格分隔的字符串转为 字符串数组
			// msgs := strings.Split(value, ":") //将每一行按":"分割成两个数组
			log.Println("msgs[0]--------->", msgs[0])
			// if len(msgs) < 8 { //不足8列 跳过当前行, 不做分析
			if len(msgs) > 1 { //这里正常只应有一列, 超过1列则数据有问题,跳过当前行,不做分析
				continue
			}
			data, _ := strconv.Atoi(msgs[0]) //定位需要匹配的列是当前行拆分后转换的字符串数组第几个元素 ,这里取第一列
			if data >= rule.Dbrule.Dbauditcont.ResultG {
				dbshtp.Dbauditcont.Alarm = "G"
				log.Printf("!!Matched!! value [%v]", data)
				entry.Minor = append(entry.Minor, fmt.Sprintf("%s数据库,审计内容数量当前%d超过阈值%d,建议检查审计配置并清理历史审计数据", dbshtp.Dbname.Contents, data, rule.Dbrule.Dbauditcont.ResultG))
				break Looop
			}

		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBNosysInSystem 分析系统用户
func Ana_DBNosysInSystem(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
}

// Ana_DBVIRSCHECK 分析病毒检查
func Ana_DBVIRSCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	//log.Println("rule.Dbrule.Dbvirscheck->", rule.Dbrule.Dbvirscheck)
	msgdata := dbshtp.Dbvirscheck.Contents
	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbvirscheck.Nm,
		Title:    rule.Dbrule.Dbvirscheck.Title,
		Desc:     rule.Dbrule.Dbvirscheck.Desc,
	}
Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)              //去除头尾空格及空行
		if value == rule.Dbrule.Dbvirscheck.ResultR { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbvirscheck.Alarm = "R"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
	// log.Printf("Dbvirscheck.Alarm->%s", dbshtp.Dbvirscheck.Alarm)
}

// Ana_DBRMANCHECK 分析 RMAN 备份配置
func Ana_DBRMANCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	// log.Println("rule.Dbrule.Dbrmancheck->", rule.Dbrule.Dbrmancheck)
	msgdata := dbshtp.Dbrmancheck.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbrmancheck.Nm,
		Title:    rule.Dbrule.Dbrmancheck.Title,
		Desc:     rule.Dbrule.Dbrmancheck.Desc,
	}
	if strings.TrimSpace(msgdata) == "" { //判断是否空, 为空则标记后退出
		dbshtp.Dbrmancheck.Alarm = "G"
		return
	}
Looop:
	//按行分割
	for _, row := range strings.Split(msgdata, "\n") {
		row = strings.TrimSpace(row) //去除头尾空格及空行

		re1 := regexp.MustCompile(rule.Dbrule.Dbrmancheck.ResultR) //逐行查找是否有"ERROR"关键词
		if re1.MatchString(strings.ToUpper(row)) {
			dbshtp.Dbrmancheck.Alarm = "R"
			break Looop
		}
		re2 := regexp.MustCompile(rule.Dbrule.Dbrmancheck.ResultB) //逐行查找是否有"WARNINGS"关键词
		if re2.MatchString(row) {
			dbshtp.Dbrmancheck.Alarm = "B"
			break Looop
		}

	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
	// log.Printf("Dbrmancheck.Alarm->%s", dbshtp.Dbrmancheck.Alarm)
}

// Ana_DBSCNHEALTHCHECK 分析 SCN 健康状态
func Ana_DBSCNHEALTHCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	//log.Println("rule.Dbrule.Dbscnhealthcheck->", rule.Dbrule.Dbscnhealthcheck)
	msgdata := dbshtp.Dbscnhealthcheck.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbscnhealthcheck.Nm,
		Title:    rule.Dbrule.Dbscnhealthcheck.Title,
		Desc:     rule.Dbrule.Dbscnhealthcheck.Desc,
	}
Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行
		// rda := regexp.MustCompile(`^Result: A`) // 判断是否检测到结果A
		rdv19 := regexp.MustCompile(`^Version:\s+19\.`)          // 判断是否检测版本
		rdv1124 := regexp.MustCompile(`^Version:\s+11\.2\.0\.4`) // 判断是否检测版本
		rdb := regexp.MustCompile(`^Result: B`)                  // 判断是否检测到结果B
		rdc := regexp.MustCompile(`^Result: C`)                  // 判断是否检测到结果C

		if rdv19.MatchString(value) || rdv1124.MatchString(value) {
			dbvstr := strings.Split(value, ":")
			log.Println("dbv->", (strings.TrimSpace(dbvstr[1])))
			break Looop
		}

		if rdc.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbscnhealthcheck.Alarm = "R"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
		if rdb.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbscnhealthcheck.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
	// log.Printf("Dbscnhealthcheck.Alarm->%s", dbshtp.Dbscnhealthcheck.Alarm)
}
