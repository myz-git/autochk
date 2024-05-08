package readxml

import (
	"autochk/structs"
	"strings"

	"github.com/beevik/etree" // go get github.com/beevik/etree
)

func ReadXml(path string, infoshtp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht) {

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
	}

	//读取root tag
	root := doc.SelectElement("EACHK")

	//找到tag0
	for _, tag0 := range root.SelectElements("TAG0") {
		//遍历tag0中tag
		for _, tag := range tag0.ChildElements() {
			// log.Printf("<%s>", tag.Tag)
			// // log.Println( strings.TrimSpace(tag.Text()))
			// log.Printf("</%s>", tag.Tag)
			switch tag.Tag {
			// 准备INFO 工作表 数据
			case "HOSTNAME":
				infoshtp.HostName = strings.TrimSpace(tag.Text())
			case "IPADDR":
				infoshtp.Ipaddr = strings.TrimSpace(tag.Text())
			case "OS":
				infoshtp.Os = strings.TrimSpace(tag.Text())
			case "RELVER":
				infoshtp.Relver = strings.TrimSpace(tag.Text())
			case "CORES":
				infoshtp.Cores = strings.TrimSpace(tag.Text())
			case "CPUCOUNT":
				infoshtp.CpuCount = strings.TrimSpace(tag.Text())
			case "CPUMHZ":
				infoshtp.CpuMHZ = strings.TrimSpace(tag.Text())
			case "MEMTOTAL":
				infoshtp.MemTotal = strings.TrimSpace(tag.Text())
			case "SWAPTOTAL":
				infoshtp.SwapTotal = strings.TrimSpace(tag.Text())
				// default:
				// 	infoshtp.Others =  strings.TrimSpace(tag.Text())

			// 准备OS工作表数据
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
	//找到tag1
	for _, tag1 := range root.SelectElements("TAG1") {
		//遍历tag1中tag ,这里为动态的数据库名 如 <myzdb></myzdb>
		for _, tag11 := range tag1.ChildElements() {
			// log.Println(tag11.Tag)
			//遍历<数据库名> 内部tag
			for _, tag := range tag11.ChildElements() {
				// log.Printf("<%s>", tag.Tag)
				// // log.Println( strings.TrimSpace(tag.Text()))
				// log.Printf("</%s>", tag.Tag)
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
					infoshtp.DbTotalsize = strings.TrimSpace(tag.Text())
				case "DBFILECOUNT":
					infoshtp.DbFilecount = strings.TrimSpace(tag.Text())
				case "DBTBLCOUNT":
					infoshtp.DbTblcount = strings.TrimSpace(tag.Text())
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

					// default:
					// 	infoshtp.Others =  strings.TrimSpace(tag.Text())
				}
			}
		}
	}
}
