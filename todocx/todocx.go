package todocx

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"strings"

	docx "github.com/lukasjarosch/go-docx"
)

// addTp writes both Contents and Alarm placeholders
func addTp(m docx.PlaceholderMap, key string, v structs.Tpstrc) {
	m[key] = v.Contents
	m[key+"_ALARM"] = v.Alarm
}

// buildPlaceholderMap builds placeholder map for multi-node templates (NODE1_*, NODE2_*, etc.)
func buildPlaceholderMap(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts) docx.PlaceholderMap {
	m := docx.PlaceholderMap{}

	// 节点数量
	m["NODE_COUNT"] = fmt.Sprintf("%d", len(*osshts))

	// 为每个节点生成占位符 (最多4个节点)
	for i := 1; i <= 4; i++ {
		nodePrefix := fmt.Sprintf("NODE%d_", i)
		instPrefix := fmt.Sprintf("INST%d_", i)

		if i <= len(*osshts) {
			// 节点存在，填充真实数据
			osn := (*osshts)[i-1]
			m[nodePrefix+"NODEID"] = osn.NodeID
			m[nodePrefix+"HOSTNAME"] = osn.Hostname.Contents
			m[nodePrefix+"IPADDR"] = osn.Ipaddr.Contents
			m[nodePrefix+"OS"] = osn.Os.Contents
			m[nodePrefix+"RELVER"] = osn.Relver.Contents
			m[nodePrefix+"CORES"] = osn.Cores.Contents
			m[nodePrefix+"CPUCOUNT"] = osn.Cpucount.Contents
			m[nodePrefix+"CPUMHZ"] = osn.Cpumhz.Contents
			m[nodePrefix+"MEMTOTAL"] = osn.Memtotal.Contents
			m[nodePrefix+"SWAPTOTAL"] = osn.Swaptotal.Contents
			m[nodePrefix+"OSPARAMETER"] = osn.Osparameter.Contents
			m[nodePrefix+"ULIMIT"] = osn.Ulimit.Contents
			m[nodePrefix+"OSLOG"] = osn.Oslog.Contents
			m[nodePrefix+"FILESYSTEM"] = osn.Filesystem.Contents
			m[nodePrefix+"INODEUSAGE"] = osn.Inodeusage.Contents
			m[nodePrefix+"CPUSTAT"] = osn.Cpustat.Contents
			m[nodePrefix+"MEMSTAT"] = osn.Memstat.Contents
			m[nodePrefix+"IOSTAT"] = osn.Iostat.Contents
			m[nodePrefix+"THPSTAT"] = osn.Thpstat.Contents
			m[nodePrefix+"HUGEPAGE"] = osn.Hugepage.Contents
			m[nodePrefix+"NUMA"] = osn.Numa.Contents
			m[nodePrefix+"NTP"] = osn.Ntp.Contents
			m[nodePrefix+"TMZONE"] = osn.Tmzone.Contents
			m[nodePrefix+"SELINUX"] = osn.Selinux.Contents
			m[nodePrefix+"FIREWALL"] = osn.Firewall.Contents
			m[nodePrefix+"NSSWITCH"] = osn.Nsswitch.Contents
			m[nodePrefix+"LO_MTU"] = osn.Lo_mtu.Contents
			m[nodePrefix+"MACHINE_PLATFORM"] = osn.Machine_platform.Contents
			m[nodePrefix+"CPU_PERF_MODE"] = osn.CPU_PERF_MODE.Contents
			m[nodePrefix+"NOZEROCONF"] = osn.NOZEROCONF.Contents
			m[nodePrefix+"RPM_PACKAGES"] = osn.RPM_PACKAGES.Contents
		} else {
			// 节点不存在，填充空字符串
			m[nodePrefix+"NODEID"] = ""
			m[nodePrefix+"HOSTNAME"] = ""
			m[nodePrefix+"IPADDR"] = ""
			m[nodePrefix+"OS"] = ""
			m[nodePrefix+"RELVER"] = ""
			m[nodePrefix+"CORES"] = ""
			m[nodePrefix+"CPUCOUNT"] = ""
			m[nodePrefix+"CPUMHZ"] = ""
			m[nodePrefix+"MEMTOTAL"] = ""
			m[nodePrefix+"SWAPTOTAL"] = ""
			m[nodePrefix+"OSPARAMETER"] = ""
			m[nodePrefix+"ULIMIT"] = ""
			m[nodePrefix+"OSLOG"] = ""
			m[nodePrefix+"FILESYSTEM"] = ""
			m[nodePrefix+"INODEUSAGE"] = ""
			m[nodePrefix+"CPUSTAT"] = ""
			m[nodePrefix+"MEMSTAT"] = ""
			m[nodePrefix+"IOSTAT"] = ""
			m[nodePrefix+"THPSTAT"] = ""
			m[nodePrefix+"HUGPAGE"] = ""
			m[nodePrefix+"NUMA"] = ""
			m[nodePrefix+"NTP"] = ""
			m[nodePrefix+"TMZONE"] = ""
			m[nodePrefix+"SELINUX"] = ""
			m[nodePrefix+"FIREWALL"] = ""
			m[nodePrefix+"NSSWITCH"] = ""
			m[nodePrefix+"LO_MTU"] = ""
			m[nodePrefix+"MACHINE_PLATFORM"] = ""
			m[nodePrefix+"CPU_PERF_MODE"] = ""
			m[nodePrefix+"NOZEROCONF"] = ""
			m[nodePrefix+"RPM_PACKAGES"] = ""
		}

		if i <= len(*instshts) {
			// 实例存在，填充真实数据
			ins := (*instshts)[i-1]
			m[instPrefix+"INSTNAME"] = ins.Instname.Contents
			m[instPrefix+"LOADPROFILE"] = ins.Loadprofile.Contents
			m[instPrefix+"INSTEFFICIENCY"] = ins.Instefficiency.Contents
			m[instPrefix+"TOPEVENT"] = ins.Topevent.Contents
			m[instPrefix+"TOPSQL_BY_ELA"] = ins.Topsql_by_ela.Contents
			m[instPrefix+"CURSOR_SHARE_MEM"] = ins.Cursor_share_mem.Contents
			m[instPrefix+"DBRESOURCE"] = ins.Dbresource.Contents
			m[instPrefix+"DBPSU"] = ins.Dbpsu.Contents
			m[instPrefix+"DBPATCH"] = ins.Dbpatch.Contents
			m[instPrefix+"DBLSNRINFO"] = ins.Dblsnrinfo.Contents
			m[instPrefix+"DBPARAMETER"] = ins.Dbparameter.Contents
			m[instPrefix+"DB_PARAMETER_FILE"] = ins.Db_parameter_file.Contents
			m[instPrefix+"DBREDOCHECK"] = ins.Dbredocheck.Contents
			m[instPrefix+"DBREDOSWITCH"] = ins.Dbredoswitch.Contents
			m[instPrefix+"RECOVERY_USAGE"] = ins.Recovery_usage.Contents
			m[instPrefix+"RECOVERY_DETAIL"] = ins.Recovery_detail.Contents
			m[instPrefix+"DBERRLOG"] = ins.Dberrlog.Contents
			m[instPrefix+"DBDGLAGCHECK"] = ins.Dbdglagcheck.Contents
			m[instPrefix+"DBDGERRCHECK"] = ins.Dbdgerrcheck.Contents
		} else {
			// 实例不存在，填充空字符串
			m[instPrefix+"INSTNAME"] = ""
			m[instPrefix+"LOADPROFILE"] = ""
			m[instPrefix+"INSTEFFICIENCY"] = ""
			m[instPrefix+"TOPEVENT"] = ""
			m[instPrefix+"TOPSQL_BY_ELA"] = ""
			m[instPrefix+"CURSOR_SHARE_MEM"] = ""
			m[instPrefix+"DBRESOURCE"] = ""
			m[instPrefix+"DBPSU"] = ""
			m[instPrefix+"DBPATCH"] = ""
			m[instPrefix+"DBLSNRINFO"] = ""
			m[instPrefix+"DBPARAMETER"] = ""
			m[instPrefix+"DB_PARAMETER_FILE"] = ""
			m[instPrefix+"DBREDOCHECK"] = ""
			m[instPrefix+"DBREDOSWITCH"] = ""
			m[instPrefix+"RECOVERY_USAGE"] = ""
			m[instPrefix+"RECOVERY_DETAIL"] = ""
			m[instPrefix+"DBERRLOG"] = ""
			m[instPrefix+"DBDGLAGCHECK"] = ""
			m[instPrefix+"DBDGERRCHECK"] = ""
		}
	}

	// 保留原有的单节点占位符作为兼容（使用第一个节点数据）
	if len(*osshts) > 0 {
		osn := (*osshts)[0]
		m["HOSTNAME"] = osn.Hostname.Contents
		m["NODEID"] = osn.NodeID
		m["IPADDR"] = osn.Ipaddr.Contents
		m["OS"] = osn.Os.Contents
		m["RELVER"] = osn.Relver.Contents
		m["CORES"] = osn.Cores.Contents
		m["CPUCOUNT"] = osn.Cpucount.Contents
		m["CPUMHZ"] = osn.Cpumhz.Contents
		m["MEMTOTAL"] = osn.Memtotal.Contents
		m["SWAPTOTAL"] = osn.Swaptotal.Contents
		m["OSPARAMETER"] = osn.Osparameter.Contents
		m["ULIMIT"] = osn.Ulimit.Contents
		m["OSLOG"] = osn.Oslog.Contents
		m["FILESYSTEM"] = osn.Filesystem.Contents
		m["INODEUSAGE"] = osn.Inodeusage.Contents
		m["CPUSTAT"] = osn.Cpustat.Contents
		m["MEMSTAT"] = osn.Memstat.Contents
		m["IOSTAT"] = osn.Iostat.Contents
		m["THPSTAT"] = osn.Thpstat.Contents
		m["HUGEPAGE"] = osn.Hugepage.Contents
		m["NUMA"] = osn.Numa.Contents
		m["NTP"] = osn.Ntp.Contents
		m["TMZONE"] = osn.Tmzone.Contents
		m["SELINUX"] = osn.Selinux.Contents
		m["FIREWALL"] = osn.Firewall.Contents
		m["NSSWITCH"] = osn.Nsswitch.Contents
		m["LO_MTU"] = osn.Lo_mtu.Contents
		m["MACHINE_PLATFORM"] = osn.Machine_platform.Contents
		m["CPU_PERF_MODE"] = osn.CPU_PERF_MODE.Contents
		m["NOZEROCONF"] = osn.NOZEROCONF.Contents
		m["RPM_PACKAGES"] = osn.RPM_PACKAGES.Contents
	}

	// DB相关字段
	m["DBNAME"] = dbshtp.Dbname.Contents
	m["DBMAA"] = dbshtp.Dbmaa.Contents
	m["DBVER"] = dbshtp.Dbver.Contents
	m["DBSTATUS"] = dbshtp.Dbstatus.Contents
	m["DBLANG"] = dbshtp.Dblang.Contents
	m["LOGMODE"] = dbshtp.Logmode.Contents
	m["FLASHBACK"] = dbshtp.Flashback.Contents
	m["DBCURSIZE"] = dbshtp.Dbcursize.Contents
	m["DBF_SIZE"] = dbshtp.Dbf_size.Contents
	m["DBF_CNT"] = dbshtp.Dbf_cnt.Contents
	m["DBF_STAT"] = dbshtp.Dbf_stat.Contents
	m["TMPFILE_SIZE"] = dbshtp.Tmpfile_size.Contents
	m["DBTBLCOUNT"] = dbshtp.Dbtblcount.Contents
	m["DBROLE"] = dbshtp.Dbrole.Contents
	m["DBTBSUSAGE"] = dbshtp.Dbtbsusage.Contents
	m["DBCONTROLFILE"] = dbshtp.Dbcontrolfile.Contents
	m["USER_INFO"] = dbshtp.User_info.Contents
	m["USER_SIZE"] = dbshtp.User_size.Contents
	m["TAB_INFO"] = dbshtp.Tab_info.Contents
	m["TAB_PARALLEL"] = dbshtp.Tab_parallel.Contents
	m["INX_PARALLEL"] = dbshtp.Inx_parallel.Contents
	m["INVALID_OBJ"] = dbshtp.Invalid_obj.Contents
	m["INVALID_INX"] = dbshtp.Invalid_inx.Contents
	m["DBSEQUENCE"] = dbshtp.Dbsequence.Contents
	m["DB_SEQ_USAGE"] = dbshtp.Db_seq_usage.Contents
	m["DBOPTION"] = dbshtp.Dboption.Contents
	m["DBFEATURES"] = dbshtp.Dbfeatures.Contents
	m["DB_EXPIR_USER"] = dbshtp.Db_expir_user.Contents
	m["DB_PASSWORD_VERIF"] = dbshtp.Db_password_verif.Contents
	m["DBDBAPRIV"] = dbshtp.Dbdbapriv.Contents
	m["DBSYSDBA"] = dbshtp.Dbsysdba.Contents
	m["DBAUDITSEGMENT"] = dbshtp.Dbauditsegment.Contents
	m["DBAUDITCONT"] = dbshtp.Dbauditcont.Contents
	m["DB_NOSYS_IN_SYSTEM"] = dbshtp.Db_Nosys_In_System.Contents
	m["USERFAILEDLOGIN"] = dbshtp.Userfailedlogin.Contents
	m["DBVIRSCHECK"] = dbshtp.Dbvirscheck.Contents
	m["DBSCNHEALTHCHECK"] = dbshtp.Dbscnhealthcheck.Contents
	m["DBRMANCHECK"] = dbshtp.Dbrmancheck.Contents
	m["CRS_STAT"] = dbshtp.Crs_stat.Contents
	m["CRS_STAT2"] = dbshtp.Crs_stat2.Contents
	m["OCR_INFO"] = dbshtp.Ocr_info.Contents
	m["OCR_BAK_CHECK"] = dbshtp.Ocr_bak_check.Contents
	m["ASM_USAGE"] = dbshtp.Asm_usage.Contents
	m["ASM_OFFSET"] = dbshtp.Asm_offset.Contents

	// 添加README中提到的字段 - 使用现有字段
	m["DBTOTALSIZE"] = dbshtp.Dbcursize.Contents // 使用Dbcursize作为总大小
	m["DBFILECOUNT"] = dbshtp.Dbf_cnt.Contents   // 使用Dbf_cnt作为文件数量

	// 添加模板中缺失的占位符 - 使用默认值或现有字段
	m["DB_4031CHECK"] = "N/A"                // 4031错误检查
	m["FILESYSTEM_REPORT"] = "N/A"           // 文件系统报告
	m["TOPSQLBYELAPSTIME"] = "N/A"           // 按执行时间排序的SQL
	m["DBINDEXPARALLEL"] = "N/A"             // 数据库索引并行度
	m["DBDATAFILE"] = "N/A"                  // 数据文件信息
	m["DBINVALIDINDEX"] = "N/A"              // 无效索引
	m["DB_SHP_SIZE"] = "N/A"                 // 共享池大小
	m["DB_SHP_PCT"] = "N/A"                  // 共享池百分比
	m["DBCRSCHECK"] = "N/A"                  // CRS检查
	m["PROJECTNAME"] = "Oracle Health Check" // 项目名称
	m["GETRQ"] = "N/A"                       // GET请求
	m["DBTABLEPARALLEL"] = "N/A"             // 表并行度
	m["DBASMUSAGE"] = "N/A"                  // ASM使用情况
	m["DBFLASHRECOVERYUSEAGE"] = "N/A"       // 闪回恢复使用情况
	m["OWNER_SIZE"] = "N/A"                  // 所有者大小

	// 这些占位符在模板中格式不正确，但我们需要提供值

	// 使用第一个INST节点的数据
	if len(*instshts) > 0 {
		ins := (*instshts)[0]
		m["INSTNAME"] = ins.Instname.Contents
		m["LOADPROFILE"] = ins.Loadprofile.Contents
		m["INSTEFFICIENCY"] = ins.Instefficiency.Contents
		m["TOPEVENT"] = ins.Topevent.Contents
		m["TOPSQL_BY_ELA"] = ins.Topsql_by_ela.Contents
		m["CURSOR_SHARE_MEM"] = ins.Cursor_share_mem.Contents
		m["DBRESOURCE"] = ins.Dbresource.Contents
		m["DBPSU"] = ins.Dbpsu.Contents
		m["DBPATCH"] = ins.Dbpatch.Contents
		m["DBLSNRINFO"] = ins.Dblsnrinfo.Contents
		m["DBPARAMETER"] = ins.Dbparameter.Contents
		m["DB_PARAMETER_FILE"] = ins.Db_parameter_file.Contents
		m["DBREDOCHECK"] = ins.Dbredocheck.Contents
		m["DBREDOSWITCH"] = ins.Dbredoswitch.Contents
		m["RECOVERY_USAGE"] = ins.Recovery_usage.Contents
		m["RECOVERY_DETAIL"] = ins.Recovery_detail.Contents
		m["DBERRLOG"] = ins.Dberrlog.Contents
		m["DBDGLAGCHECK"] = ins.Dbdglagcheck.Contents
		m["DBDGERRCHECK"] = ins.Dbdgerrcheck.Contents
	}

	// 报告生成日期（yyyy-mm-dd）
	m["RPTDATE"] = time.Now().Format("2006-01-02")

	return m
}

// Todocx generates docx from template with placeholder replacements
func Todocx(osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool) {
	startTime := time.Now()

	// 根据节点数量选择模板
	nodeCount := len(*osshts)
	var templatePath string
	if nodeCount <= 2 {
		templatePath = "2node.docx"
	} else {
		templatePath = "4node.docx"
	}

	// ensure output dir
	outputDir := "report"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		_ = os.MkdirAll(outputDir, 0755)
	}
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_Done.docx", xlsnm))

	// build map
	replaceMap := buildPlaceholderMap(osshts, dbshtp, instshts)

	// open template
	doc, err := docx.Open(templatePath)
	if err != nil {
		utils.LogErrorf("Failed to open docx template: %v", err)
		return
	}
	utils.LogDebugf("open took: %s", time.Since(startTime))

	// Debug: 列出模板占位符与映射键的差异，便于定位缺失键
	if utils.ShouldLog(utils.LevelDebug) {
		if list, err := doc.GetPlaceHoldersList(); err == nil {
			docKeys := map[string]struct{}{}
			for _, ph := range list {
				k := strings.TrimSpace(ph)
				k = strings.TrimPrefix(k, "{")
				k = strings.TrimSuffix(k, "}")
				// 处理格式错误的占位符（如 {HOSTNAME 缺少右大括号）
				if k != "" && !strings.Contains(k, "{") && !strings.Contains(k, "}") {
					docKeys[k] = struct{}{}
				}
			}
			missingInMap := []string{}
			for k := range docKeys {
				if _, ok := replaceMap[k]; !ok {
					missingInMap = append(missingInMap, k)
				}
			}
			unusedInDoc := []string{}
			for k := range replaceMap {
				if _, ok := docKeys[k]; !ok {
					unusedInDoc = append(unusedInDoc, k)
				}
			}
			utils.LogDebugf("Doc placeholders total=%d", len(docKeys))
			utils.LogDebugf("Map keys total=%d", len(replaceMap))
			if len(missingInMap) > 0 {
				utils.LogDebugf("Missing in map: %v", strings.Join(missingInMap, ", "))
			}
			if len(unusedInDoc) > 0 {
				utils.LogDebugf("Unused map keys: %v", strings.Join(unusedInDoc, ", "))
			}
			// 显示所有模板占位符
			allDocKeys := make([]string, 0, len(docKeys))
			for k := range docKeys {
				allDocKeys = append(allDocKeys, k)
			}
			utils.LogDebugf("All doc placeholders: %v", strings.Join(allDocKeys, ", "))
		} else {
			utils.LogDebugf("GetPlaceHoldersList error: %v", err)
		}
	}

	// replace
	if err := doc.ReplaceAll(replaceMap); err != nil {
		utils.LogErrorf("Failed to replace placeholders: %v", err)
		return
	}
	utils.LogDebugf("replace took: %s", time.Since(startTime))

	// write
	if err := doc.WriteToFile(outputFile); err != nil {
		utils.LogErrorf("Failed to write to docx file: %v", err)
		return
	}
	utils.LogDebugf("everything took: %s", time.Since(startTime))
	utils.LogInfof("Generated docx report at %s", outputFile)
}
