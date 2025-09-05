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

// db_security.go 包含与数据库安全相关的分析函数

// Ana_DBExpirUser 分析用户密码过期情况
func Ana_DBExpirUser(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := strings.TrimSpace(dbshtp.Db_expir_user.Contents)
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Db_expir_user.Nm,
		Title:    rule.Dbrule.Db_expir_user.Title,
		Desc:     rule.Dbrule.Db_expir_user.Desc,
	}

	// 没有记录则表示正常（无问题，不追加告警）
	if msgdata == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		dbshtp.Db_expir_user.Alarm = ""
		return
	}

	// 只要有记录，即存在口令将过期的用户，判定为普通告警(B)
	dbshtp.Db_expir_user.Alarm = "G"
	entry.Moderate = append(entry.Moderate,
		fmt.Sprintf("问题: %s数据库,存在密码即将过期的用户,\n建议: 提前更新用户密码", dbshtp.Dbname.Contents))

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DB_PASSWORD_VERIF 分析密码复杂性验证
func Ana_DB_PASSWORD_VERIF(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Db_password_verif.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Db_password_verif.Nm,
		Title:    rule.Dbrule.Db_password_verif.Title,
		Desc:     rule.Dbrule.Db_password_verif.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(msgdata, "no rows selected") {
		// 正常情况，无告警
		dbshtp.Db_password_verif.Alarm = ""
	} else {
		// 有记录说明存在未设置密码复杂性校验，设置为G级告警
		dbshtp.Db_password_verif.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,存在未设置密码复杂性校验的PROFILE,\n建议: 启用密码复杂性验证以增强安全性", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Userfailedlogin 密码错误用户锁定检查
func Ana_Userfailedlogin(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Userfailedlogin.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Userfailedlogin.Nm,
		Title:    rule.Dbrule.Userfailedlogin.Title,
		Desc:     rule.Dbrule.Userfailedlogin.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Userfailedlogin.Alarm = ""
	} else {
		// 有记录说明存在登录失败锁定限制，设置为G级告警
		dbshtp.Userfailedlogin.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,存在用户登录失败锁定限制,\n建议: 对于业务用户不做登录失败锁定限制", dbshtp.Dbname.Contents))
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

	// 检查是否为空或包含"无记录"
	if value == "" || strings.Contains(value, "无记录") || strings.Contains(strings.ToLower(value), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbdbapriv.Alarm = ""
	} else {
		// 有记录，检查记录数量
		lines := strings.Split(value, "\n")
		recordCount := 0

		// 从第3行开始统计（跳过标题行和分隔线）
		for i := 2; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line != "" {
				recordCount++
			}
		}

		// 如果存在两条（包含）以上记录则为G级告警
		if recordCount >= 2 {
			dbshtp.Dbdbapriv.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,存在%d个DBA权限账户,\n建议: 只保留一个DBA账户，收回其他DBA用户权限以增强安全性", dbshtp.Dbname.Contents, recordCount))
		} else {
			// 只有一条记录，正常情况
			dbshtp.Dbdbapriv.Alarm = ""
		}
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

	// 检查是否为空或包含"无记录"
	if value == "" || strings.Contains(value, "无记录") || strings.Contains(strings.ToLower(value), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbsysdba.Alarm = ""
	} else {
		// 有记录说明存在非必要SYSDBA权限用户，设置为B级告警
		dbshtp.Dbsysdba.Alarm = "Minor"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,存在非必要SYSDBA权限用户,\n建议: 对涉及用户收回SYSDBA权限", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBAUDITSEGMENT 分析审计段
func Ana_DBAUDITSEGMENT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbauditsegment.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbauditsegment.Nm,
		Title:    rule.Dbrule.Dbauditsegment.Title,
		Desc:     rule.Dbrule.Dbauditsegment.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbauditsegment.Alarm = ""
	} else {
		// 有记录说明审计数据占用过多存储空间，设置为G级告警
		dbshtp.Dbauditsegment.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,审计数据占用过多存储空间,\n建议: 在满足审计监管要求下定期清理或归档数据。", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBAUDITCONT 分析审计内容
func Ana_DBAUDITCONT(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := strings.TrimSpace(dbshtp.Dbauditcont.Contents)
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbauditcont.Nm,
		Title:    rule.Dbrule.Dbauditcont.Title,
		Desc:     rule.Dbrule.Dbauditcont.Desc,
	}

	// 检查是否为空或包含"无记录"
	if msgdata == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbauditcont.Alarm = ""
		return
	}

	// 检查数据行数是否足够
	lines := strings.Split(msgdata, "\n")
	if len(lines) < 3 {
		// 数据行数不足，不做告警处理
		dbshtp.Dbauditcont.Alarm = ""
		return
	}

	// 获取第3行（索引2）
	thirdLine := strings.TrimSpace(lines[2])
	if thirdLine == "" {
		// 第3行为空，不做告警处理
		dbshtp.Dbauditcont.Alarm = ""
		return
	}

	// 尝试将第3行转换为数字
	if data, err := strconv.Atoi(thirdLine); err == nil {
		// 如果该值 > rule.Dbrule.Dbauditcont.Result，则判断为G
		if data > rule.Dbrule.Dbauditcont.Result {
			dbshtp.Dbauditcont.Alarm = "G"
			entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,审计内容数量当前%d超过阈值%d,\n建议: 在满足审计监管要求下定期清理或归档数据。", dbshtp.Dbname.Contents, data, rule.Dbrule.Dbauditcont.Result))
		} else {
			// 数值在阈值范围内，正常情况
			dbshtp.Dbauditcont.Alarm = ""
		}
	} else {
		// 第3行不是数字，不做告警处理
		dbshtp.Dbauditcont.Alarm = ""
		return
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBNosysInSystem 分析系统用户
func Ana_DBNosysInSystem(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Db_Nosys_In_System.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Db_Nosys_In_System.Nm,
		Title:    rule.Dbrule.Db_Nosys_In_System.Title,
		Desc:     rule.Dbrule.Db_Nosys_In_System.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Db_Nosys_In_System.Alarm = ""
	} else {
		// 有记录说明存在业务账户对象存储在SYSTEM、SYSAUX表空间，设置为G级告警
		dbshtp.Db_Nosys_In_System.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s数据库,存在业务账户对象存储在SYSTEM、SYSAUX表空间,\n建议: 将业务对象迁移到专用表空间", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBVIRSCHECK 分析病毒检查
func Ana_DBVIRSCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbvirscheck.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbvirscheck.Nm,
		Title:    rule.Dbrule.Dbvirscheck.Title,
		Desc:     rule.Dbrule.Dbvirscheck.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbvirscheck.Alarm = ""
	} else {
		// 有记录说明存在病毒植入风险，设置为R级告警
		dbshtp.Dbvirscheck.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,存在病毒植入风险,\n建议: 立即深入安全检查核实病毒风险", dbshtp.Dbname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBRMANCHECK 分析 RMAN 备份检查
func Ana_DBRMANCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbrmancheck.Contents
	entry := structs.SummaryEntry{
		Category: "数据库备份",
		Nm:       rule.Dbrule.Dbrmancheck.Nm,
		Title:    rule.Dbrule.Dbrmancheck.Title,
		Desc:     rule.Dbrule.Dbrmancheck.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 无记录说明数据库没有RMAN备份，设置为B级告警
		dbshtp.Dbrmancheck.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,未发现RMAN备份记录,\n建议: 确认是否有其他有效备份,尽快部署RMAN备份策略", dbshtp.Dbname.Contents))
	} else {
		// 有记录，逐行检查关键词
		lines := strings.Split(msgdata, "\n")
		hasError := false
		hasWarning := false

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 检查是否包含ERROR关键词
			if strings.Contains(strings.ToUpper(line), "ERROR") {
				hasError = true
			}

			// 检查是否包含WARNINGS关键词
			if strings.Contains(strings.ToUpper(line), "WARNINGS") {
				hasWarning = true
			}
		}

		// 根据检查结果设置告警级别
		if hasError {
			dbshtp.Dbrmancheck.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,近7天RMAN备份检查发现错误,\n建议: 检查备份作业执行情况,确保备份作业正常运行", dbshtp.Dbname.Contents))
		} else if hasWarning {
			dbshtp.Dbrmancheck.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,近7天RMAN备份检查发现警告,\n建议: 检查备份作业执行情况,确保备份作业正常运行", dbshtp.Dbname.Contents))
		} else {
			// 无错误无警告，正常情况
			dbshtp.Dbrmancheck.Alarm = ""
		}
	}

	// 如果有告警信息，添加到汇总中
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_DBSCNHEALTHCHECK 分析 SCN 健康状态
func Ana_DBSCNHEALTHCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	msgdata := dbshtp.Dbscnhealthcheck.Contents
	entry := structs.SummaryEntry{
		Category: "数据库安全",
		Nm:       rule.Dbrule.Dbscnhealthcheck.Nm,
		Title:    rule.Dbrule.Dbscnhealthcheck.Title,
		Desc:     rule.Dbrule.Dbscnhealthcheck.Desc,
	}

	// 检查是否为空或包含"无记录"
	if strings.TrimSpace(msgdata) == "" || strings.Contains(msgdata, "无记录") || strings.Contains(strings.ToLower(msgdata), "no rows selected") {
		// 正常情况，无告警
		dbshtp.Dbscnhealthcheck.Alarm = ""
		return
	}

	// 按行分割进行详细检查
	lines := strings.Split(msgdata, "\n")
	scnHealthStatus := ""

	// 从第3行开始检查（跳过标题行和分隔线）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// 判断是否检测到版本信息
		rdv19 := regexp.MustCompile(`^Version:\s+19\.`)          // 判断是否检测版本19
		rdv1124 := regexp.MustCompile(`^Version:\s+11\.2\.0\.4`) // 判断是否检测版本11.2.0.4

		// 判断是否检测到结果
		rdb := regexp.MustCompile(`^Result: B`) // 判断是否检测到结果B
		rdc := regexp.MustCompile(`^Result: C`) // 判断是否检测到结果C

		if rdv19.MatchString(line) || rdv1124.MatchString(line) {
			// 检测到版本信息，记录版本
			dbvstr := strings.Split(line, ":")
			if len(dbvstr) > 1 {
				log.Printf("检测到数据库版本: %s", strings.TrimSpace(dbvstr[1]))
			}
			continue
		}

		if rdc.MatchString(line) { // 匹配到结果C
			scnHealthStatus = "C"
			dbshtp.Dbscnhealthcheck.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s数据库,SCN健康检查结果为C，SCN增长异常,\n建议: 尽快核查数据库是否SCN增长异常", dbshtp.Dbname.Contents))
			break
		}

		if rdb.MatchString(line) { // 匹配到结果B
			scnHealthStatus = "B"
			dbshtp.Dbscnhealthcheck.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s数据库,SCN健康检查结果为B，SCN增长较快,\n建议: 关注SCN增长速度是否过快", dbshtp.Dbname.Contents))
			break
		}
	}

	// 如果没有检测到B或C结果，说明SCN健康状态良好
	if scnHealthStatus == "" {
		dbshtp.Dbscnhealthcheck.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
