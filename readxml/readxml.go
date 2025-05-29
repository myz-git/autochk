package readxml

import (
	"autochk/structs"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func ReadXml(path string, infoshtp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
	}

	root := doc.SelectElement("EACHK")

	// 处理 TAG0（主机相关信息）
	for _, tag0 := range root.SelectElements("TAG0") {
		for _, tag := range tag0.ChildElements() {
			switch tag.Tag {
			// Server Info
			case "HOSTNAME":
				infoshtp.HostName = strings.TrimSpace(tag.Text())
			case "IPADDR":
				infoshtp.Ipaddr = strings.TrimSpace(tag.Text())
			case "OS":
				infoshtp.Os = strings.TrimSpace(tag.Text())
			case "RELVER":
				infoshtp.Relver = strings.TrimSpace(tag.Text())
			case "CPUCOUNT":
				infoshtp.CpuCount = strings.TrimSpace(tag.Text())
			case "CPUMHZ":
				infoshtp.CpuMHZ = strings.TrimSpace(tag.Text())
			case "MEMTOTAL":
				memKB, _ := strconv.Atoi(strings.TrimSpace(tag.Text()))
				infoshtp.MemTotal = strconv.Itoa(memKB/1024) + " Megabytes" // 转换为 MB
			// OS Sheet
			case "OSPARAMETER":
				osshtp.Osparameter.Contents = strings.TrimSpace(tag.Text())
			case "ULIMIT":
				osshtp.Ulimit.Contents = strings.TrimSpace(tag.Text())
			case "FILESYSTEM":
				osshtp.Filesystem.Contents = strings.TrimSpace(tag.Text())
			case "INODEUSAGE":
				osshtp.Inodeusage.Contents = strings.TrimSpace(tag.Text())
			case "CPUSTAT":
				osshtp.Cpustat.Contents = strings.TrimSpace(tag.Text())
			case "MEMSTAT":
				osshtp.Memstat.Contents = strings.TrimSpace(tag.Text())
			case "IOSTAT":
				osshtp.Iostat.Contents = strings.TrimSpace(tag.Text())
			case "THPSTAT":
				osshtp.Thpstat.Contents = strings.TrimSpace(tag.Text())
			case "HUGPAGE":
				osshtp.Hugpage.Contents = strings.TrimSpace(tag.Text())
			case "NUMA":
				osshtp.Numa.Contents = strings.TrimSpace(tag.Text())
			case "NTP":
				osshtp.Ntp.Contents = strings.TrimSpace(tag.Text())
			case "DBMAA":
				infoshtp.DbMaa = strings.TrimSpace(tag.Text())
			case "DBCRSCHECK":
				dbshtp.Dbcrscheck.Contents = strings.TrimSpace(tag.Text())
			case "DBASMUSAGE":
				dbshtp.Dbasmusage.Contents = strings.TrimSpace(tag.Text())
			}
		}
	}

	// 处理 TAG1（数据库相关信息）
	for _, tag1 := range root.SelectElements("TAG1") {
		for _, tag11 := range tag1.ChildElements() {
			for _, tag := range tag11.ChildElements() {
				switch tag.Tag {
				case "DBNAME":
					infoshtp.DbName = strings.TrimSpace(tag.Text())
				case "DBVER":
					infoshtp.DbVer = strings.TrimSpace(tag.Text())
				case "DBROLE":
					infoshtp.DbRole = strings.TrimSpace(tag.Text())
				case "LOGMODE":
					infoshtp.LogMode = strings.TrimSpace(tag.Text())
				case "FLASHBACK":
					infoshtp.FlashBack = strings.TrimSpace(tag.Text())
				case "DBTOTALSIZE":
					infoshtp.DbTotalsize = strings.TrimSpace(tag.Text()) + " GB"
				case "DBFILECOUNT":
					infoshtp.DbFilecount = strings.TrimSpace(tag.Text())
				case "DBTBLCOUNT":
					infoshtp.DbTblcount = strings.TrimSpace(tag.Text()) // 去除 GB 单位
				case "DBLANG":
					infoshtp.DbLang = strings.TrimSpace(tag.Text())
				case "DBTBSUSAGE":
					dbshtp.DbTbsusage.Contents = strings.TrimSpace(tag.Text())
				case "DBDATAFILE":
					dbshtp.Dbdatafile.Contents = strings.TrimSpace(tag.Text())
				case "DBCONTROLFILE":
					dbshtp.Dbcontrolfile.Contents = strings.TrimSpace(tag.Text())
				case "DBUSERSIZE":
					dbshtp.Dbusersize.Contents = strings.TrimSpace(tag.Text())
				case "DBREDOCHECK":
					dbshtp.Dbredocheck.Contents = strings.TrimSpace(tag.Text())
				case "DBREDOSWITCH":
					dbshtp.Dbredoswitch.Contents = strings.TrimSpace(tag.Text())
				case "DBRESOURCE":
					dbshtp.Dbresource.Contents = strings.TrimSpace(tag.Text())
				case "LOADPROFILE":
					dbshtp.Loadprofile.Contents = strings.TrimSpace(tag.Text())
				case "INSTEFFICIENCY":
					dbshtp.Instefficiency.Contents = strings.TrimSpace(tag.Text())
				case "TOPEVENT":
					dbshtp.Dbtopevent.Contents = strings.TrimSpace(tag.Text())
				case "TOPSQLBYELAPSTIME":
					dbshtp.DbtopSQL.Contents = strings.TrimSpace(tag.Text())
				case "DBLSNRINFO":
					dbshtp.Dblsnrinfo.Contents = strings.TrimSpace(tag.Text())
				case "DBTABLEPARALLEL":
					dbshtp.Dbtableparallel.Contents = strings.TrimSpace(tag.Text())
				case "DBINDEXPARALLEL":
					dbshtp.Dbindexparallel.Contents = strings.TrimSpace(tag.Text())
				case "DBINVALIDINDEX":
					dbshtp.Dbinvalidindex.Contents = strings.TrimSpace(tag.Text())
				case "DBSEQUENCE":
					dbshtp.Dbsequence.Contents = strings.TrimSpace(tag.Text())
				case "DBRECOVERYDEST":
					dbshtp.Dbrecoverydest.Contents = strings.TrimSpace(tag.Text())
				case "DBFLASHRECOVERYUSEAGE":
					dbshtp.Dbflashrecoveryuseage.Contents = strings.TrimSpace(tag.Text())
				case "DBERRLOG":
					dbshtp.Dberrlog.Contents = strings.TrimSpace(tag.Text())
				case "DBPRODUCTUSERFAILEDLOGIN":
					dbshtp.Dbproductuserfailedlogin.Contents = strings.TrimSpace(tag.Text())
				case "DBDGLAGCHECK":
					dbshtp.Dbdglagcheck.Contents = strings.TrimSpace(tag.Text())
				case "DBDGERRCHECK":
					dbshtp.Dbdgerrcheck.Contents = strings.TrimSpace(tag.Text())
				case "DBRMANCHECK":
					dbshtp.Dbrmancheck.Contents = strings.TrimSpace(tag.Text())
				case "DBDBAPRIV":
					dbshtp.Dbdbapriv.Contents = strings.TrimSpace(tag.Text())
				case "DBSYSDBA":
					dbshtp.Dbsysdba.Contents = strings.TrimSpace(tag.Text())
				case "DBAUDITSEGMENT":
					dbshtp.Dbauditsegment.Contents = strings.TrimSpace(tag.Text())
				case "DBAUDITCONT":
					dbshtp.Dbauditcont.Contents = strings.TrimSpace(tag.Text())
				case "DB_NOSYS_IN_SYSTEM":
					dbshtp.Db_Nosys_In_System.Contents = strings.TrimSpace(tag.Text())
				case "DBVIRSCHECK":
					dbshtp.Dbvirscheck.Contents = strings.TrimSpace(tag.Text())
				case "DBSCNHEALTHCHECK":
					dbshtp.Dbscnhealthcheck.Contents = strings.TrimSpace(tag.Text())
				}
			}
		}
	}
}
