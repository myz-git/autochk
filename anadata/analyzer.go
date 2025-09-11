package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"strings"
)

// Ana 是分析的主入口函数，协调格式化和 OS/DB 指标分析
func Ana(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries) {
	rules, err := utils.GetRule()
	if err != nil {
		utils.LogErrorf("rule err: #%v", err)
		return
	}

	// 调试信息（默认不输出）
	utils.LogDebugf("开始分析 - OS节点数: %d, 实例数: %d", len(*osshts), len(*instshts))
	if len(*osshts) > 0 {
		utils.LogDebugf("第一个OS节点: %s, 主机名: %s", (*osshts)[0].NodeID, (*osshts)[0].Hostname.Contents)
	}

	// 分析 OS 指标 - 遍历所有节点
	for i := range *osshts {
		utils.LogDebugf("分析节点 %s (索引: %d)", (*osshts)[i].NodeID, i)

		// 为每个节点创建独立的分析上下文
		Ana_Osparameter(rules, &(*osshts)[i], summaryEntries)
		Ana_Ulimit(rules, &(*osshts)[i], summaryEntries)
		Ana_Filesystem(rules, &(*osshts)[i], summaryEntries)
		Ana_Inodeusage(rules, &(*osshts)[i], summaryEntries)
		Ana_Cpustat(rules, &(*osshts)[i], summaryEntries)
		Ana_Memstat(rules, &(*osshts)[i], summaryEntries)
		Ana_Iostat(rules, &(*osshts)[i], summaryEntries)
		Ana_Thpstat(rules, &(*osshts)[i], summaryEntries)
		Ana_Hugepage(rules, &(*osshts)[i], summaryEntries)
		Ana_Numa(rules, &(*osshts)[i], summaryEntries)
		Ana_Ntp(rules, &(*osshts)[i], summaryEntries)
		Ana_Selinux(rules, &(*osshts)[i], summaryEntries)
		Ana_Firewall(rules, &(*osshts)[i], summaryEntries)
		Ana_Nsswitch(rules, &(*osshts)[i], summaryEntries)
		Ana_Lo_mtu(rules, &(*osshts)[i], summaryEntries)
		Ana_CPU_PERF_MODE(rules, &(*osshts)[i], summaryEntries)
		Ana_NOZEROCONF(rules, &(*osshts)[i], summaryEntries)
		Ana_RPM_PACKAGES(rules, &(*osshts)[i], summaryEntries)
	}

	// 数据库基本信息分析
	Ana_DB_status(rules, dbshtp, summaryEntries)
	Ana_DB_logmode(rules, dbshtp, summaryEntries)
	Ana_DBF_CNT(rules, dbshtp, summaryEntries)
	Ana_DBF_STAT(rules, dbshtp, summaryEntries)
	Ana_DBTbs(rules, dbshtp, summaryEntries)
	Ana_DBCTRF(rules, dbshtp, summaryEntries)
	Ana_Tab_parallel(rules, dbshtp, summaryEntries)
	Ana_Inx_parallel(rules, dbshtp, summaryEntries)

	// 数据库对象分析
	Ana_Invalid_inx(rules, dbshtp, summaryEntries)
	Ana_Invalid_obj(rules, dbshtp, summaryEntries)
	Ana_DBSEQUENCE(rules, dbshtp, summaryEntries)
	Ana_DB_SEQ_USAGE(rules, dbshtp, summaryEntries)

	// 数据库安全检查
	Ana_DBExpirUser(rules, dbshtp, summaryEntries)
	Ana_DB_PASSWORD_VERIF(rules, dbshtp, summaryEntries)
	Ana_Userfailedlogin(rules, dbshtp, summaryEntries)
	Ana_DBDBAPRIV(rules, dbshtp, summaryEntries)
	Ana_DBSYSDBA(rules, dbshtp, summaryEntries)
	Ana_DBAUDITSEGMENT(rules, dbshtp, summaryEntries)
	Ana_DBAUDITCONT(rules, dbshtp, summaryEntries)
	Ana_DBNosysInSystem(rules, dbshtp, summaryEntries)
	Ana_DBVIRSCHECK(rules, dbshtp, summaryEntries)
	Ana_DBRMANCHECK(rules, dbshtp, summaryEntries)
	Ana_DBSCNHEALTHCHECK(rules, dbshtp, summaryEntries)

	// 分析实例相关指标 - 遍历所有实例
	for i := range *instshts {
		utils.LogDebugf("分析实例 %s (索引: %d)", (*instshts)[i].NodeID, i)

		// 实例相关的分析
		Ana_RDF(rules, &(*instshts)[i], summaryEntries)
		Ana_RDSW(rules, &(*instshts)[i], summaryEntries)
		Ana_RECOVERY_USAGE(rules, &(*instshts)[i], summaryEntries)
		Ana_LOADPROFILE(rules, &(*instshts)[i], summaryEntries)
		Ana_INSTEFFICIENCY(rules, &(*instshts)[i], summaryEntries)
		// Ana_DBtopevent(rules, &(*instshts)[i], summaryEntries)
		Ana_DBtopSQL(rules, &(*instshts)[i], summaryEntries)
		Ana_CursorShareMem(rules, &(*instshts)[i], summaryEntries)
		Ana_RESOURCE(rules, &(*instshts)[i], summaryEntries)
		Ana_Db_shp_size(rules, &(*instshts)[i], summaryEntries)
		Ana_DB_Shp_pct(rules, &(*instshts)[i], summaryEntries)
		Ana_DB4031check(rules, &(*instshts)[i], summaryEntries)
		Ana_DBPSU(rules, &(*instshts)[i], summaryEntries)
		Ana_DBPATCH(rules, &(*instshts)[i], summaryEntries)
		Ana_DBLSNRINFO(rules, &(*instshts)[i], summaryEntries)
		Ana_DBparameter(rules, &(*instshts)[i], summaryEntries)
		Ana_DBERRLOG(rules, &(*instshts)[i], summaryEntries)

		// DataGuard相关检查 - 只有当数据库角色为STANDBY且是第一个Node时才执行
		if strings.Contains(dbshtp.Dbrole.Contents, "STANDBY") && i == 0 {
			Ana_DBDGLAGCHECK(rules, dbshtp, &(*instshts)[i], summaryEntries)
			Ana_DBDGERRCHECK(rules, &(*instshts)[i], summaryEntries)
		}
	}

	//集群相关检查
	if dbshtp.Dbmaa.Contents == "RAC" {
		Ana_Crs_stat(rules, dbshtp, summaryEntries)
		Ana_Ocr_info(rules, dbshtp, summaryEntries)
		Ana_ASM_usage(rules, dbshtp, summaryEntries)
		Ana_Asm_offset(rules, dbshtp, summaryEntries)
	}
}
