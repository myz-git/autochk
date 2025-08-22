package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"log"
)

// Ana 是分析的主入口函数，协调格式化和 OS/DB 指标分析
func Ana(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {
	rules, err := utils.GetRule()
	if err != nil {
		log.Printf("rule err: #%v", err)
		return
	}

	// 格式化 InfoSht 字段
	Fmt_DbRole(infstp)
	Fmt_LogMode(infstp)
	Fmt_FlashBack(infstp)
	Fmt_DbTotalsize(infstp)
	Fmt_DbFilecount(infstp)
	Fmt_DbTblcount(infstp)

	// 分析 OS 指标
	Ana_Osparameter(rules, osshtp, infstp, summaryEntries)
	Ana_Ulimit(rules, osshtp, infstp, summaryEntries)
	Ana_Filesystem(rules, osshtp, summaryEntries)
	Ana_Inodeusage(rules, osshtp, summaryEntries)
	Ana_Cpustat(rules, osshtp, summaryEntries)
	Ana_Memstat(rules, osshtp, summaryEntries)
	Ana_Iostat(rules, osshtp, summaryEntries)
	Ana_Thpstat(rules, osshtp, summaryEntries)
	Ana_Numa(rules, osshtp, summaryEntries)
	Ana_Ntp(rules, osshtp, summaryEntries)

	// 数据库实例分析
	Ana_DBstatus(rules, dbshtp, summaryEntries)
	Ana_DBTbs(rules, dbshtp, summaryEntries)
	Ana_DBF(rules, dbshtp, summaryEntries)
	Ana_DBCTRF(rules, dbshtp, summaryEntries)
	Ana_DBusersize(rules, dbshtp, summaryEntries)
	Ana_RDSW(rules, dbshtp, summaryEntries)
	Ana_RDF(rules, dbshtp, summaryEntries)
	Ana_DBparameter(rules, dbshtp, summaryEntries)
	Ana_DBParameterFile(rules, dbshtp, summaryEntries)
	Ana_DBShpSize(rules, dbshtp, summaryEntries)
	Ana_DBShpPct(rules, dbshtp, summaryEntries)

	// 数据库对象分析
	Ana_DBTABLEPARALLEL(rules, dbshtp, summaryEntries)
	Ana_DBINDEXPARALLEL(rules, dbshtp, summaryEntries)
	Ana_DBINVALIDINDEX(rules, dbshtp, summaryEntries)
	Ana_DBSEQUENCE(rules, dbshtp, infstp, summaryEntries)
	Ana_DB_SEQ_USAGE(rules, dbshtp, summaryEntries)

	// 数据库性能分析
	Ana_DB4031check(rules, dbshtp, summaryEntries)
	Ana_RESOURCE(rules, dbshtp, summaryEntries)
	Ana_LOADPROFILE(rules, dbshtp, summaryEntries)
	Ana_INSTEFFICIENCY(rules, dbshtp, summaryEntries)
	Ana_DBtopevent(rules, dbshtp, summaryEntries)
	Ana_DBtopSQL(rules, dbshtp, summaryEntries)

	// 数据库安全检查
	Ana_DBExpirUser(rules, dbshtp, summaryEntries)
	Ana_DBPRODUCTUSERFAILEDLOGIN(rules, dbshtp, summaryEntries)
	Ana_DBDBAPRIV(rules, dbshtp, summaryEntries)
	Ana_DBSYSDBA(rules, dbshtp, summaryEntries)
	Ana_DBAUDITSEGMENT(rules, dbshtp, summaryEntries)
	Ana_DBAUDITCONT(rules, dbshtp, summaryEntries)
	Ana_DBNosysInSystem(rules, dbshtp, summaryEntries)
	Ana_DBVIRSCHECK(rules, dbshtp, summaryEntries)
	Ana_DBRMANCHECK(rules, dbshtp, summaryEntries)
	Ana_DBSCNHEALTHCHECK(rules, dbshtp, summaryEntries)

	// 数据库监控、DataGuard、备份及杂项分析
	Ana_DBERRLOG(rules, dbshtp, summaryEntries)
	Ana_DBLSNRINFO(rules, dbshtp, summaryEntries)
	Ana_DBDGLAGCHECK(rules, dbshtp, infstp, summaryEntries)
	Ana_DBDGERRCHECK(rules, dbshtp, summaryEntries)
	Ana_DBRECOVERYDEST(rules, dbshtp, summaryEntries)
	Ana_DBFLASHRECOVERYUSEAGE(rules, dbshtp, summaryEntries)
	Ana_DBCRSCHECK(rules, dbshtp, summaryEntries)
	Ana_DBASMUSAGE(rules, dbshtp, summaryEntries)
	Ana_DBPSU(rules, dbshtp, summaryEntries)
	Ana_DBPATCH(rules, dbshtp, summaryEntries)
}
