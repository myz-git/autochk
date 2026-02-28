package readxml

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func ReadXml(path string, osshts *[]structs.OsShts, dbshtp *structs.DbSht, instshts *[]structs.InstShts) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		utils.LogErrorf("解析 XML 文件失败 (%s): %v", path, err)
		return
	}

	root := doc.SelectElement("EACHK")
	if root == nil {
		utils.LogErrorf("XML 根节点缺失或不是 <EACHK>: %s", path)
		return
	}

	// 处理 TAG0（主机相关信息）- 支持多节点
	for _, tag0 := range root.SelectElements("TAG0") {
		// 遍历所有NODE节点 (NODE1, NODE2, NODE3...)
		for _, node := range tag0.ChildElements() {
			if strings.HasPrefix(node.Tag, "NODE") {
				processTag0Node(node, osshts)
			}
		}
	}

	// 处理 TAG1（数据库相关信息）- 只处理NODE1
	for _, tag1 := range root.SelectElements("TAG1") {
		node1 := tag1.SelectElement("NODE1")
		if node1 != nil {
			processTag1Node(node1, dbshtp)
		}
	}

	// 处理 TAG2（数据库实例信息）- 支持多节点
	for _, tag2 := range root.SelectElements("TAG2") {
		// 遍历所有NODE节点 (NODE1, NODE2, NODE3...)
		for _, node := range tag2.ChildElements() {
			if strings.HasPrefix(node.Tag, "NODE") {
				processTag2Node(node, instshts)
			}
		}
	}
}

// 处理TAG0中的NODE节点 - 主机相关信息
func processTag0Node(node *etree.Element, osshts *[]structs.OsShts) {
	// 创建新的OsShts结构体
	osSht := structs.OsShts{
		NodeID: node.Tag, // 设置节点ID为NODE1, NODE2等
	}

	// 初始化所有Tpstrc字段
	osSht.Hostname = structs.Tpstrc{}
	osSht.Ipaddr = structs.Tpstrc{}
	osSht.Os = structs.Tpstrc{}
	osSht.Relver = structs.Tpstrc{}
	osSht.Cpu_model = structs.Tpstrc{}
	osSht.Cpucount = structs.Tpstrc{}
	osSht.Cpumhz = structs.Tpstrc{}
	osSht.Memtotal = structs.Tpstrc{}
	osSht.Swaptotal = structs.Tpstrc{}
	osSht.Osparam_fs = structs.Tpstrc{}
	osSht.Osparam_ker = structs.Tpstrc{}
	osSht.Osparam_net = structs.Tpstrc{}
	osSht.Osparam_vm = structs.Tpstrc{}
	osSht.Ulimit = structs.Tpstrc{}
	osSht.Filesystem = structs.Tpstrc{}
	osSht.Inodeusage = structs.Tpstrc{}
	osSht.Cpustat = structs.Tpstrc{}
	osSht.Memstat = structs.Tpstrc{}
	osSht.Iostat = structs.Tpstrc{}
	osSht.Thpstat = structs.Tpstrc{}
	osSht.Hugepage = structs.Tpstrc{}

	osSht.Numa = structs.Tpstrc{}
	osSht.Ntp = structs.Tpstrc{}
	osSht.Tmzone = structs.Tpstrc{}
	osSht.Selinux = structs.Tpstrc{}
	osSht.Firewall = structs.Tpstrc{}
	osSht.Nsswitch = structs.Tpstrc{}
	osSht.Lo_mtu = structs.Tpstrc{}
	osSht.Machine_platform = structs.Tpstrc{}
	osSht.CPU_PERF_MODE = structs.Tpstrc{}
	osSht.NOZEROCONF = structs.Tpstrc{}
	osSht.RPM_PACKAGES = structs.Tpstrc{}
	osSht.Oslog = structs.Tpstrc{}

	for _, tag := range node.ChildElements() {
		switch tag.Tag {
		case "HOSTNAME":
			osSht.Hostname.Contents = strings.TrimSpace(tag.Text())
		case "IPADDR":
			osSht.Ipaddr.Contents = strings.TrimSpace(tag.Text())
		case "OS":
			osSht.Os.Contents = strings.TrimSpace(tag.Text())
		case "RELVER":
			osSht.Relver.Contents = strings.TrimSpace(tag.Text())
		case "CPU_MODEL":
			osSht.Cpu_model.Contents = strings.TrimSpace(tag.Text())
		case "CPUCOUNT":
			osSht.Cpucount.Contents = strings.TrimSpace(tag.Text())
		case "CPUMHZ":
			osSht.Cpumhz.Contents = strings.TrimSpace(tag.Text())
		case "MEMTOTAL":
			osSht.Memtotal.Contents = convertKBtoGB(strings.TrimSpace(tag.Text()))
		case "SWAPTOTAL":
			osSht.Swaptotal.Contents = convertKBtoGB(strings.TrimSpace(tag.Text()))
		case "OSPARAM_FS":
			osSht.Osparam_fs.Contents = strings.TrimSpace(tag.Text())
		case "OSPARAM_KER":
			osSht.Osparam_ker.Contents = strings.TrimSpace(tag.Text())
		case "OSPARAM_NET":
			osSht.Osparam_net.Contents = strings.TrimSpace(tag.Text())
		case "OSPARAM_VM":
			osSht.Osparam_vm.Contents = strings.TrimSpace(tag.Text())
		case "ULIMIT":
			osSht.Ulimit.Contents = strings.TrimSpace(tag.Text())
		case "FILESYSTEM":
			osSht.Filesystem.Contents = strings.TrimSpace(tag.Text())
		case "INODEUSAGE":
			osSht.Inodeusage.Contents = strings.TrimSpace(tag.Text())
		case "CPUSTAT":
			osSht.Cpustat.Contents = strings.TrimSpace(tag.Text())
		case "MEMSTAT":
			osSht.Memstat.Contents = strings.TrimSpace(tag.Text())
		case "IOSTAT":
			osSht.Iostat.Contents = strings.TrimSpace(tag.Text())
		case "THPSTAT":
			osSht.Thpstat.Contents = strings.TrimSpace(tag.Text())
		case "HUGEPAGE":
			osSht.Hugepage.Contents = strings.TrimSpace(tag.Text())
		case "NUMA":
			osSht.Numa.Contents = strings.TrimSpace(tag.Text())
		case "NTP":
			osSht.Ntp.Contents = strings.TrimSpace(tag.Text())
		case "TMZONE":
			osSht.Tmzone.Contents = strings.TrimSpace(tag.Text())
		case "SELINUX":
			osSht.Selinux.Contents = strings.TrimSpace(tag.Text())
		case "FIREWALL":
			osSht.Firewall.Contents = strings.TrimSpace(tag.Text())
		case "NSSWITCH":
			osSht.Nsswitch.Contents = strings.TrimSpace(tag.Text())
		case "LO_MTU":
			osSht.Lo_mtu.Contents = strings.TrimSpace(tag.Text())
		case "MACHINE_PLATFORM":
			osSht.Machine_platform.Contents = strings.TrimSpace(tag.Text())
		case "CPU_PERF_MODE":
			osSht.CPU_PERF_MODE.Contents = strings.TrimSpace(tag.Text())
		case "NOZEROCONF":
			osSht.NOZEROCONF.Contents = strings.TrimSpace(tag.Text())
		case "RPM_PACKAGES":
			osSht.RPM_PACKAGES.Contents = strings.TrimSpace(tag.Text())
		case "OSLOG":
			osSht.Oslog.Contents = strings.TrimSpace(tag.Text())
		}
	}

	// 将处理好的OsShts添加到数组中
	*osshts = append(*osshts, osSht)
}

// 处理TAG1中的NODE1节点 - 数据库基本信息
func processTag1Node(node *etree.Element, dbshtp *structs.DbSht) {
	// 设置节点ID
	dbshtp.NodeID = "NODE1"

	// 初始化所有Tpstrc字段
	dbshtp.Dbname = structs.Tpstrc{}
	dbshtp.Dbmaa = structs.Tpstrc{}
	dbshtp.Dbver = structs.Tpstrc{}
	dbshtp.Dbstatus = structs.Tpstrc{}
	dbshtp.Logmode = structs.Tpstrc{}
	dbshtp.Dbrole = structs.Tpstrc{}
	dbshtp.Flashback = structs.Tpstrc{}
	dbshtp.Dbcursize = structs.Tpstrc{}
	dbshtp.Dbf_size = structs.Tpstrc{}
	dbshtp.Dbf_cnt = structs.Tpstrc{}
	dbshtp.Dbf_stat = structs.Tpstrc{}
	dbshtp.Tmpfile_size = structs.Tpstrc{}
	// dbshtp.Dbtblcount = structs.Tpstrc{}
	dbshtp.Dblang = structs.Tpstrc{}
	dbshtp.Dbtbsusage = structs.Tpstrc{}
	dbshtp.Dbcontrolfile = structs.Tpstrc{}
	dbshtp.User_info = structs.Tpstrc{}
	dbshtp.User_size = structs.Tpstrc{}
	dbshtp.Tab_info = structs.Tpstrc{}
	dbshtp.Tab_parallel = structs.Tpstrc{}
	dbshtp.Inx_parallel = structs.Tpstrc{}
	dbshtp.Invalid_obj = structs.Tpstrc{}
	dbshtp.Invalid_inx = structs.Tpstrc{}
	dbshtp.Dbsequence = structs.Tpstrc{}
	dbshtp.Db_seq_usage = structs.Tpstrc{}
	// dbshtp.Dboption = structs.Tpstrc{}
	// dbshtp.Dbfeatures = structs.Tpstrc{}
	dbshtp.Db_expir_user = structs.Tpstrc{}
	dbshtp.Db_password_verif = structs.Tpstrc{}
	dbshtp.Userfailedlogin = structs.Tpstrc{}
	dbshtp.Dbrmancheck = structs.Tpstrc{}
	dbshtp.Dbdbapriv = structs.Tpstrc{}
	dbshtp.Dbsysdba = structs.Tpstrc{}
	dbshtp.Dbauditsegment = structs.Tpstrc{}
	dbshtp.Dbauditcont = structs.Tpstrc{}
	dbshtp.Db_Nosys_In_System = structs.Tpstrc{}
	dbshtp.Dbvirscheck = structs.Tpstrc{}
	dbshtp.Dbscnhealthcheck = structs.Tpstrc{}
	dbshtp.Dbparam_b = structs.Tpstrc{}
	dbshtp.Dbparam_d = structs.Tpstrc{}
	dbshtp.Crs_stat = structs.Tpstrc{}
	dbshtp.Crs_stat2 = structs.Tpstrc{}
	dbshtp.Ocr_info = structs.Tpstrc{}
	dbshtp.Ocr_bak_check = structs.Tpstrc{}
	dbshtp.Asm_usage = structs.Tpstrc{}
	dbshtp.Asm_offset = structs.Tpstrc{}

	for _, tag := range node.ChildElements() {
		switch tag.Tag {
		case "DBNAME":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbname.Contents = "无记录"
			} else {
				dbshtp.Dbname.Contents = content
			}
		case "DBMAA":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbmaa.Contents = "无记录"
			} else {
				dbshtp.Dbmaa.Contents = content
			}
		case "DBVER":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbver.Contents = "无记录"
			} else {
				dbshtp.Dbver.Contents = content
			}
		case "DBSTATUS":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbstatus.Contents = "无记录"
			} else {
				dbshtp.Dbstatus.Contents = content
			}
		case "DBROLE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbrole.Contents = "无记录"
			} else {
				dbshtp.Dbrole.Contents = content
			}
		case "LOGMODE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Logmode.Contents = "无记录"
			} else {
				dbshtp.Logmode.Contents = content
			}
		case "FLASHBACK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Flashback.Contents = "无记录"
			} else {
				dbshtp.Flashback.Contents = content
			}
		case "DBCURSIZE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbcursize.Contents = "无记录"
			} else {
				dbshtp.Dbcursize.Contents = content
			}
		case "DBF_SIZE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbf_size.Contents = "无记录"
			} else {
				dbshtp.Dbf_size.Contents = content
			}
		case "DBF_CNT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbf_cnt.Contents = "无记录"
			} else {
				dbshtp.Dbf_cnt.Contents = content
			}
		case "DBF_STAT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbf_stat.Contents = "无记录"
			} else {
				dbshtp.Dbf_stat.Contents = content
			}
		case "TMPFILE_SIZE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Tmpfile_size.Contents = "无记录"
			} else {
				dbshtp.Tmpfile_size.Contents = content
			}
		case "DBLANG":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dblang.Contents = "无记录"
			} else {
				dbshtp.Dblang.Contents = content
			}
		case "DBTBSUSAGE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbtbsusage.Contents = "无记录"
			} else {
				dbshtp.Dbtbsusage.Contents = content
			}
		case "DBCONTROLFILE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbcontrolfile.Contents = "无记录"
			} else {
				dbshtp.Dbcontrolfile.Contents = content
			}
		case "USER_INFO":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.User_info.Contents = "无记录"
			} else {
				dbshtp.User_info.Contents = content
			}
		case "USER_SIZE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.User_size.Contents = "无记录"
			} else {
				dbshtp.User_size.Contents = content
			}
		case "TAB_INFO":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Tab_info.Contents = "无记录"
			} else {
				dbshtp.Tab_info.Contents = content
			}
		case "TAB_PARALLEL":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Tab_parallel.Contents = "无记录"
			} else {
				dbshtp.Tab_parallel.Contents = content
			}
		case "INX_PARALLEL":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Inx_parallel.Contents = "无记录"
			} else {
				dbshtp.Inx_parallel.Contents = content
			}
		case "INVALID_OBJ":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Invalid_obj.Contents = "无记录"
			} else {
				dbshtp.Invalid_obj.Contents = content
			}
		case "INVALID_INX":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Invalid_inx.Contents = "无记录"
			} else {
				dbshtp.Invalid_inx.Contents = content
			}
		case "DBSEQUENCE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbsequence.Contents = "无记录"
			} else {
				dbshtp.Dbsequence.Contents = content
			}
		case "DB_SEQ_USAGE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Db_seq_usage.Contents = "无记录"
			} else {
				dbshtp.Db_seq_usage.Contents = content
			}
		case "DB_EXPIR_USER":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Db_expir_user.Contents = "无记录"
			} else {
				dbshtp.Db_expir_user.Contents = content
			}
		case "DB_PASSWORD_VERIF":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Db_password_verif.Contents = "无记录"
			} else {
				dbshtp.Db_password_verif.Contents = content
			}
		case "Userfailedlogin":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Userfailedlogin.Contents = "无记录"
			} else {
				dbshtp.Userfailedlogin.Contents = content
			}
		case "DBRMANCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbrmancheck.Contents = "无记录"
			} else {
				dbshtp.Dbrmancheck.Contents = content
			}
		case "DBDBAPRIV":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbdbapriv.Contents = "无记录"
			} else {
				dbshtp.Dbdbapriv.Contents = content
			}
		case "DBSYSDBA":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbsysdba.Contents = "无记录"
			} else {
				dbshtp.Dbsysdba.Contents = content
			}
		case "DBAUDITSEGMENT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbauditsegment.Contents = "无记录"
			} else {
				dbshtp.Dbauditsegment.Contents = content
			}
		case "DBAUDITCONT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbauditcont.Contents = "无记录"
			} else {
				dbshtp.Dbauditcont.Contents = content
			}
		case "DB_NOSYS_IN_SYSTEM":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Db_Nosys_In_System.Contents = "无记录"
			} else {
				dbshtp.Db_Nosys_In_System.Contents = content
			}
		case "DBVIRSCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbvirscheck.Contents = "无记录"
			} else {
				dbshtp.Dbvirscheck.Contents = content
			}
		case "DBSCNHEALTHCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbscnhealthcheck.Contents = "无记录"
			} else {
				dbshtp.Dbscnhealthcheck.Contents = content
			}
		case "DBPARAM_B":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbparam_b.Contents = "无记录"
			} else {
				dbshtp.Dbparam_b.Contents = content
			}
		case "DBPARAM_D":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				dbshtp.Dbparam_d.Contents = "无记录"
			} else {
				dbshtp.Dbparam_d.Contents = content
			}
		case "CRS_STAT":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Crs_stat.Contents = content
		case "CRS_STAT2":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Crs_stat2.Contents = content
		case "OCR_INFO":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Ocr_info.Contents = content
		case "OCR_BAK_CHECK":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Ocr_bak_check.Contents = content
		case "ASM_USAGE":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Asm_usage.Contents = content
		case "ASM_OFFSET":
			content := strings.TrimSpace(tag.Text())
			dbshtp.Asm_offset.Contents = content

		}
	}
}

// 处理TAG2中的NODE节点 - 数据库实例信息
func processTag2Node(node *etree.Element, instshts *[]structs.InstShts) {
	// 创建新的InstShts结构体
	instSht := structs.InstShts{
		NodeID: node.Tag, // 设置节点ID为NODE1, NODE2等
	}

	// 初始化所有Tpstrc字段
	instSht.Instname = structs.Tpstrc{}
	instSht.Loadprofile = structs.Tpstrc{}
	instSht.Instefficiency = structs.Tpstrc{}
	instSht.Topevent = structs.Tpstrc{}
	instSht.Topsql_by_ela = structs.Tpstrc{}
	instSht.Cursor_share_mem = structs.Tpstrc{}
	instSht.Db_shp_size = structs.Tpstrc{}
	instSht.Db_shp_pct = structs.Tpstrc{}
	instSht.Dbresource = structs.Tpstrc{}
	instSht.Dbpsu = structs.Tpstrc{}
	instSht.Dbpatch = structs.Tpstrc{}
	instSht.Dblsnrinfo = structs.Tpstrc{}
	instSht.Dbredocheck = structs.Tpstrc{}
	instSht.Dbredoswitch = structs.Tpstrc{}
	instSht.Recovery_usage = structs.Tpstrc{}
	instSht.Recovery_detail = structs.Tpstrc{}
	instSht.Dbparameter = structs.Tpstrc{}
	instSht.Db_parameter_file = structs.Tpstrc{}
	instSht.Dberrlog = structs.Tpstrc{}
	instSht.Dbdglagcheck = structs.Tpstrc{}
	instSht.Dbdgerrcheck = structs.Tpstrc{}

	for _, tag := range node.ChildElements() {
		switch tag.Tag {
		case "INSTNAME":
			instSht.Instname.Contents = strings.TrimSpace(tag.Text())
		case "LOADPROFILE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Loadprofile.Contents = "无记录"
			} else {
				instSht.Loadprofile.Contents = content
			}
		case "INSTEFFICIENCY":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Instefficiency.Contents = "无记录"
			} else {
				instSht.Instefficiency.Contents = content
			}
		case "TOPEVENT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Topevent.Contents = "无记录"
			} else {
				instSht.Topevent.Contents = content
			}
		case "TOPSQL_BY_ELA":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Topsql_by_ela.Contents = "无记录"
			} else {
				instSht.Topsql_by_ela.Contents = content
			}
		case "CURSOR_SHARE_MEM":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Cursor_share_mem.Contents = "无记录"
			} else {
				instSht.Cursor_share_mem.Contents = content
			}
		case "DB_SHP_SIZE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Db_shp_size.Contents = "无记录"
			} else {
				instSht.Db_shp_size.Contents = content
			}
		case "DB_SHP_PCT":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Db_shp_pct.Contents = "无记录"
			} else {
				instSht.Db_shp_pct.Contents = content
			}
		case "DBRESOURCE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbresource.Contents = "无记录"
			} else {
				instSht.Dbresource.Contents = content
			}
		case "DBPSU":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbpsu.Contents = "无记录"
			} else {
				instSht.Dbpsu.Contents = content
			}
		case "DBPATCH":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbpatch.Contents = "无记录"
			} else {
				instSht.Dbpatch.Contents = content
			}
		case "DBPARAMETER":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbparameter.Contents = "无记录"
			} else {
				instSht.Dbparameter.Contents = content
			}
		case "DB_PARAMETER_FILE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Db_parameter_file.Contents = "无记录"
			} else {
				instSht.Db_parameter_file.Contents = content
			}
		case "DBLSNRINFO":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dblsnrinfo.Contents = "无记录"
			} else {
				instSht.Dblsnrinfo.Contents = content
			}
		case "DBREDOCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbredocheck.Contents = "无记录"
			} else {
				instSht.Dbredocheck.Contents = content
			}
		case "DBREDOSWITCH":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbredoswitch.Contents = "无记录"
			} else {
				instSht.Dbredoswitch.Contents = content
			}
		case "RECOVERY_USAGE":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Recovery_usage.Contents = "无记录"
			} else {
				instSht.Recovery_usage.Contents = content
			}
		case "RECOVERY_DETAIL":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Recovery_detail.Contents = "无记录"
			} else {
				instSht.Recovery_detail.Contents = content
			}
		case "DBERRLOG":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dberrlog.Contents = "无记录"
			} else {
				instSht.Dberrlog.Contents = content
			}
		case "DBDGLAGCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbdglagcheck.Contents = "无记录"
			} else {
				instSht.Dbdglagcheck.Contents = content
			}
		case "DBDGERRCHECK":
			content := strings.TrimSpace(tag.Text())
			if content == "" {
				instSht.Dbdgerrcheck.Contents = "无记录"
			} else {
				instSht.Dbdgerrcheck.Contents = content
			}
		}
	}

	// 将处理好的InstShts添加到数组中
	*instshts = append(*instshts, instSht)
}

// convertKBtoGB 将KB单位转换为GB单位并取整
func convertKBtoGB(kbStr string) string {
	// 移除所有空格和单位标识
	cleanStr := strings.TrimSpace(kbStr)
	cleanStr = strings.ReplaceAll(cleanStr, " ", "")
	cleanStr = strings.ReplaceAll(cleanStr, "kB", "")
	cleanStr = strings.ReplaceAll(cleanStr, "KB", "")
	cleanStr = strings.ReplaceAll(cleanStr, "kb", "")

	// 尝试解析数字
	if kb, err := strconv.ParseFloat(cleanStr, 64); err == nil {
		// 转换为GB (1 GB = 1024 * 1024 KB)
		gb := kb / (1024 * 1024)
		// 取整
		gbInt := int(gb)
		return fmt.Sprintf("%d GB", gbInt)
	}

	// 如果解析失败，返回原字符串
	return kbStr
}
