package structs

// OsShts 用于存储TAG0中每个NODE的OS相关数据
type OsShts struct {
	NodeID           string // 节点ID，如"NODE1","NODE2","NODE3"等
	Hostname         Tpstrc
	Ipaddr           Tpstrc
	Os               Tpstrc
	Relver           Tpstrc
	Cores            Tpstrc
	Cpucount         Tpstrc
	Cpumhz           Tpstrc
	Memtotal         Tpstrc
	Swaptotal        Tpstrc
	Osparameter      Tpstrc
	Ulimit           Tpstrc
	Filesystem       Tpstrc
	Inodeusage       Tpstrc
	Cpustat          Tpstrc
	Memstat          Tpstrc
	Iostat           Tpstrc
	Thpstat          Tpstrc
	Hugepage         Tpstrc
	Numa             Tpstrc
	Ntp              Tpstrc
	Tmzone           Tpstrc
	Selinux          Tpstrc
	Firewall         Tpstrc
	Nsswitch         Tpstrc
	Lo_mtu           Tpstrc
	Machine_platform Tpstrc
	CPU_PERF_MODE    Tpstrc
	NOZEROCONF       Tpstrc
	RPM_PACKAGES     Tpstrc
	Oslog            Tpstrc
}

// DbSht 用于存储TAG1中NODE1的DATABASE相关信息
type DbSht struct {
	NodeID             string // 节点ID，只有一个"NODE1"
	Dbname             Tpstrc
	Dbmaa              Tpstrc
	Dbver              Tpstrc
	Dbstatus           Tpstrc
	Dblang             Tpstrc
	Logmode            Tpstrc
	Flashback          Tpstrc
	Dbcursize          Tpstrc
	Dbf_size           Tpstrc
	Dbf_cnt            Tpstrc
	Dbf_stat           Tpstrc
	Tmpfile_size       Tpstrc
	Dbtblcount         Tpstrc
	Dbrole             Tpstrc
	Dbtbsusage         Tpstrc
	Dbcontrolfile      Tpstrc
	User_info          Tpstrc
	User_size          Tpstrc
	Tab_info           Tpstrc
	Tab_parallel       Tpstrc
	Inx_parallel       Tpstrc
	Invalid_obj        Tpstrc
	Invalid_inx        Tpstrc
	Dbsequence         Tpstrc
	Db_seq_usage       Tpstrc
	Dboption           Tpstrc
	Dbfeatures         Tpstrc
	Db_expir_user      Tpstrc
	Db_password_verif  Tpstrc
	Dbdbapriv          Tpstrc
	Dbsysdba           Tpstrc
	Dbauditsegment     Tpstrc
	Dbauditcont        Tpstrc
	Db_Nosys_In_System Tpstrc
	Userfailedlogin    Tpstrc
	Dbvirscheck        Tpstrc
	Dbscnhealthcheck   Tpstrc
	Dbrmancheck        Tpstrc
	Crs_stat           Tpstrc
	Crs_stat2          Tpstrc
	Ocr_info           Tpstrc
	Ocr_bak_check      Tpstrc
	Asm_usage          Tpstrc
	Asm_offset         Tpstrc
}

// InstShts 用于存储TAG2中每个NODE的实例、日志、监听等相关数据
type InstShts struct {
	NodeID            string
	Instname          Tpstrc
	Loadprofile       Tpstrc
	Instefficiency    Tpstrc
	Topevent          Tpstrc
	Topsql_by_ela     Tpstrc
	Cursor_share_mem  Tpstrc
	Dbresource        Tpstrc
	Dbpsu             Tpstrc
	Dbpatch           Tpstrc
	Dblsnrinfo        Tpstrc
	Dbparameter       Tpstrc
	Db_parameter_file Tpstrc
	Dbredocheck       Tpstrc
	Dbredoswitch      Tpstrc
	Recovery_usage    Tpstrc
	Recovery_detail   Tpstrc
	Dberrlog          Tpstrc
	Dbdglagcheck      Tpstrc
	Dbdgerrcheck      Tpstrc
}

type Tpstrc struct {
	Contents string
	Alarm    string // 告警级别: R B G (Red, Blue, Green)
}

type SummaryEntry struct {
	Category string   // 检查类别
	Nm       string   // 检查项
	Title    string   // 检查项中文
	Desc     string   // 检查说明
	Severe   []string // 检查结果（严重）
	Moderate []string // 检查结果（一般）
	Minor    []string // 检查结果（轻微）
}

type SummaryEntries struct {
	Entries []SummaryEntry
}
