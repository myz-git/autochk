package structs

type InfoSht struct {
	DbName      string
	DbVer       string
	DbRole      string
	LogMode     string
	FlashBack   string
	DbTotalsize string
	DbFilecount string
	DbTblcount  string
	DbLang      string
	DbMaa       string
	HostName    string
	Ipaddr      string
	Os          string
	Relver      string
	Cores       string
	CpuCount    string
	CpuMHZ      string
	MemTotal    string
	SwapTotal   string

	Others string
}

type OsSht struct {
	Osparameter Tpstrc
	Ulimit      Tpstrc
	Filesystem  Tpstrc
	Inodeusage  Tpstrc
	Cpustat     Tpstrc
	Memstat     Tpstrc
	Iostat      Tpstrc
	Thpstat     Tpstrc
	Hugpage     Tpstrc
	Numa        Tpstrc
	Ntp         Tpstrc
}

type DbSht struct {
	DbTbsusage               Tpstrc
	Dbdatafile               Tpstrc
	Dbcontrolfile            Tpstrc
	Dbusersize               Tpstrc
	Dbredocheck              Tpstrc
	Dbredoswitch             Tpstrc
	Dbresource               Tpstrc
	Loadprofile              Tpstrc
	Instefficiency           Tpstrc
	Dbtopevent               Tpstrc
	DbtopSQL                 Tpstrc
	Dblsnrinfo               Tpstrc
	Dbtableparallel          Tpstrc
	Dbindexparallel          Tpstrc
	Dbinvalidindex           Tpstrc
	Dbsequence               Tpstrc
	Dbrecoverydest           Tpstrc
	Dbflashrecoveryuseage    Tpstrc
	Dberrlog                 Tpstrc
	Dbproductuserfailedlogin Tpstrc
	Dbdglagcheck             Tpstrc
	Dbdgerrcheck             Tpstrc
	Dbrmancheck              Tpstrc
	Dbdbapriv                Tpstrc
	Dbsysdba                 Tpstrc
	Dbauditsegment           Tpstrc
	Dbauditcont              Tpstrc
	Db_Nosys_In_System       Tpstrc
	Dbvirscheck              Tpstrc
	Dbscnhealthcheck         Tpstrc
	Dbcrscheck               Tpstrc
	Dbasmusage               Tpstrc
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
