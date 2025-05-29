package toxls

import (
	"autochk/structs"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Xlsx(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool) {
	// 确定输出文件名
	var newfnm string
	if sglf {
		newfnm = xlsnm + ".Done.xlsx"
	} else {
		newfnm = "HealthCheckReport.ALLDone.xlsx"
	}

	// 加载模板文件
	f, err := excelize.OpenFile("HealthReport.xlsx")
	if err != nil {
		fmt.Println("打开模板文件失败:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 填充 HealthReport Sheet 的 Server Info 和 DataBase Info
	PutSht_INFO(f, infstp, colcnt)

	// 填充 OS Sheet
	PutSht_OS(f, infstp, osshtp, colcnt)

	// 填充 DB Sheet
	PutSht_DB(f, infstp, dbshtp, colcnt)

	// 填充 HealthReport Sheet 的 Issue List
	PutSht_Summary(f, summaryEntries)

	// 保存文件
	if err := f.SaveAs(newfnm); err != nil {
		fmt.Println("保存文件失败:", err)
		return
	}
}

func PutSht_INFO(f *excelize.File, infstp *structs.InfoSht, colcnt int) {
	shnm := "HealthReport"

	// 填充 Server Info
	f.SetCellStr(shnm, "F4", infstp.HostName)
	f.SetCellStr(shnm, "F5", infstp.Ipaddr)
	f.SetCellStr(shnm, "F6", infstp.Os)
	f.SetCellStr(shnm, "F7", infstp.Relver)
	f.SetCellStr(shnm, "F8", infstp.CpuCount)
	f.SetCellStr(shnm, "F9", infstp.CpuMHZ)
	f.SetCellStr(shnm, "F10", infstp.MemTotal)

	// 填充 DataBase Info
	f.SetCellStr(shnm, "K4", infstp.DbName)
	f.SetCellStr(shnm, "K5", infstp.DbVer)
	f.SetCellStr(shnm, "K6", infstp.DbMaa)
	f.SetCellStr(shnm, "K7", infstp.DbRole)
	f.SetCellStr(shnm, "K8", infstp.LogMode)
	f.SetCellStr(shnm, "K9", infstp.DbLang)
	f.SetCellStr(shnm, "K10", infstp.DbTotalsize)
	f.SetCellStr(shnm, "K11", infstp.DbTblcount)
}

func PutSht_OS(f *excelize.File, infstp *structs.InfoSht, ossht *structs.OsSht, colcnt int) {
	shnm := "OS"

	// 定义单元格样式
	styleB, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4876FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})
	styleR, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFAEB9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	// 填充 OS Sheet
	f.SetCellStr(shnm, "B1", infstp.HostName)
	f.SetCellStr(shnm, "B2", infstp.Ipaddr)
	f.SetCellStr(shnm, "B3", ossht.Osparameter.Contents)
	if ossht.Osparameter.Alarm == "R" {
		f.SetCellStyle(shnm, "B3", "B3", styleR)
	} else if ossht.Osparameter.Alarm == "B" {
		f.SetCellStyle(shnm, "B3", "B3", styleB)
	}
	f.SetCellStr(shnm, "B4", ossht.Ulimit.Contents)
	if ossht.Ulimit.Alarm == "R" {
		f.SetCellStyle(shnm, "B4", "B4", styleR)
	} else if ossht.Ulimit.Alarm == "B" {
		f.SetCellStyle(shnm, "B4", "B4", styleB)
	}
	f.SetCellStr(shnm, "B5", ossht.Filesystem.Contents)
	if ossht.Filesystem.Alarm == "R" {
		f.SetCellStyle(shnm, "B5", "B5", styleR)
	} else if ossht.Filesystem.Alarm == "B" {
		f.SetCellStyle(shnm, "B5", "B5", styleB)
	}
	f.SetCellStr(shnm, "B6", ossht.Inodeusage.Contents)
	if ossht.Inodeusage.Alarm == "R" {
		f.SetCellStyle(shnm, "B6", "B6", styleR)
	} else if ossht.Inodeusage.Alarm == "B" {
		f.SetCellStyle(shnm, "B6", "B6", styleB)
	}
	f.SetCellStr(shnm, "B7", ossht.Cpustat.Contents)
	if ossht.Cpustat.Alarm == "R" {
		f.SetCellStyle(shnm, "B7", "B7", styleR)
	} else if ossht.Cpustat.Alarm == "B" {
		f.SetCellStyle(shnm, "B7", "B7", styleB)
	}
	f.SetCellStr(shnm, "B8", ossht.Memstat.Contents)
	if ossht.Memstat.Alarm == "R" {
		f.SetCellStyle(shnm, "B8", "B8", styleR)
	} else if ossht.Memstat.Alarm == "B" {
		f.SetCellStyle(shnm, "B8", "B8", styleB)
	}
	f.SetCellStr(shnm, "B9", ossht.Iostat.Contents)
	if ossht.Iostat.Alarm == "R" {
		f.SetCellStyle(shnm, "B9", "B9", styleR)
	} else if ossht.Iostat.Alarm == "B" {
		f.SetCellStyle(shnm, "B9", "B9", styleB)
	}
	f.SetCellStr(shnm, "B10", ossht.Thpstat.Contents)
	if ossht.Thpstat.Alarm == "R" {
		f.SetCellStyle(shnm, "B10", "B10", styleR)
	} else if ossht.Thpstat.Alarm == "B" {
		f.SetCellStyle(shnm, "B10", "B10", styleB)
	}
	f.SetCellStr(shnm, "B11", ossht.Hugpage.Contents)
	if ossht.Hugpage.Alarm == "R" {
		f.SetCellStyle(shnm, "B11", "B11", styleR)
	} else if ossht.Hugpage.Alarm == "B" {
		f.SetCellStyle(shnm, "B11", "B11", styleB)
	}
	f.SetCellStr(shnm, "B12", ossht.Numa.Contents)
	if ossht.Numa.Alarm == "R" {
		f.SetCellStyle(shnm, "B12", "B12", styleR)
	} else if ossht.Numa.Alarm == "B" {
		f.SetCellStyle(shnm, "B12", "B12", styleB)
	}
	f.SetCellStr(shnm, "B13", ossht.Ntp.Contents)
	if ossht.Ntp.Alarm == "R" {
		f.SetCellStyle(shnm, "B13", "B13", styleR)
	} else if ossht.Ntp.Alarm == "B" {
		f.SetCellStyle(shnm, "B13", "B13", styleB)
	}
}

func PutSht_DB(f *excelize.File, infstp *structs.InfoSht, dbsht *structs.DbSht, colcnt int) {
	shnm := "DB"

	// 定义单元格样式
	styleB, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4876FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})
	styleR, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFAEB9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})
	styleG, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#00bf5f"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	// 填充 DB Sheet
	f.SetCellStr(shnm, "B1", infstp.DbName)
	f.SetCellStr(shnm, "B2", infstp.HostName)
	f.SetCellStr(shnm, "B3", dbsht.DbTbsusage.Contents)
	if dbsht.DbTbsusage.Alarm == "R" {
		f.SetCellStyle(shnm, "B3", "B3", styleR)
	} else if dbsht.DbTbsusage.Alarm == "B" {
		f.SetCellStyle(shnm, "B3", "B3", styleB)
	}
	f.SetCellStr(shnm, "B4", dbsht.Dbdatafile.Contents)
	if dbsht.Dbdatafile.Alarm == "R" {
		f.SetCellStyle(shnm, "B4", "B4", styleR)
	} else if dbsht.Dbdatafile.Alarm == "B" {
		f.SetCellStyle(shnm, "B4", "B4", styleB)
	}
	f.SetCellStr(shnm, "B5", dbsht.Dbcontrolfile.Contents)
	if dbsht.Dbcontrolfile.Alarm == "R" {
		f.SetCellStyle(shnm, "B5", "B5", styleR)
	} else if dbsht.Dbcontrolfile.Alarm == "B" {
		f.SetCellStyle(shnm, "B5", "B5", styleB)
	}
	f.SetCellStr(shnm, "B6", dbsht.Dbusersize.Contents)
	f.SetCellStr(shnm, "B7", dbsht.Dbredocheck.Contents)
	if dbsht.Dbredocheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B7", "B7", styleR)
	} else if dbsht.Dbredocheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B7", "B7", styleB)
	}
	f.SetCellStr(shnm, "B8", dbsht.Dbredoswitch.Contents)
	if dbsht.Dbredoswitch.Alarm == "R" {
		f.SetCellStyle(shnm, "B8", "B8", styleR)
	} else if dbsht.Dbredoswitch.Alarm == "B" {
		f.SetCellStyle(shnm, "B8", "B8", styleB)
	}
	f.SetCellStr(shnm, "B9", dbsht.Dbresource.Contents)
	if dbsht.Dbresource.Alarm == "R" {
		f.SetCellStyle(shnm, "B9", "B9", styleR)
	} else if dbsht.Dbresource.Alarm == "B" {
		f.SetCellStyle(shnm, "B9", "B9", styleB)
	}
	f.SetCellStr(shnm, "B10", dbsht.Loadprofile.Contents)
	if dbsht.Loadprofile.Alarm == "R" {
		f.SetCellStyle(shnm, "B10", "B10", styleR)
	} else if dbsht.Loadprofile.Alarm == "B" {
		f.SetCellStyle(shnm, "B10", "B10", styleB)
	} else if dbsht.Loadprofile.Alarm == "G" {
		f.SetCellStyle(shnm, "B10", "B10", styleG)
	}
	f.SetCellStr(shnm, "B11", dbsht.Instefficiency.Contents)
	if dbsht.Instefficiency.Alarm == "R" {
		f.SetCellStyle(shnm, "B11", "B11", styleR)
	} else if dbsht.Instefficiency.Alarm == "B" {
		f.SetCellStyle(shnm, "B11", "B11", styleB)
	} else if dbsht.Instefficiency.Alarm == "G" {
		f.SetCellStyle(shnm, "B11", "B11", styleG)
	}
	f.SetCellStr(shnm, "B12", dbsht.Dbtopevent.Contents)
	if dbsht.Dbtopevent.Alarm == "R" {
		f.SetCellStyle(shnm, "B12", "B12", styleR)
	} else if dbsht.Dbtopevent.Alarm == "B" {
		f.SetCellStyle(shnm, "B12", "B12", styleB)
	} else if dbsht.Dbtopevent.Alarm == "G" {
		f.SetCellStyle(shnm, "B12", "B12", styleG)
	}
	f.SetCellStr(shnm, "B13", dbsht.DbtopSQL.Contents)
	if dbsht.DbtopSQL.Alarm == "R" {
		f.SetCellStyle(shnm, "B13", "B13", styleR)
	} else if dbsht.DbtopSQL.Alarm == "B" {
		f.SetCellStyle(shnm, "B13", "B13", styleB)
	} else if dbsht.DbtopSQL.Alarm == "G" {
		f.SetCellStyle(shnm, "B13", "B13", styleG)
	}
	f.SetCellStr(shnm, "B14", dbsht.Dblsnrinfo.Contents)
	if dbsht.Dblsnrinfo.Alarm == "R" {
		f.SetCellStyle(shnm, "B14", "B14", styleR)
	} else if dbsht.Dblsnrinfo.Alarm == "B" {
		f.SetCellStyle(shnm, "B14", "B14", styleB)
	} else if dbsht.Dblsnrinfo.Alarm == "G" {
		f.SetCellStyle(shnm, "B14", "B14", styleG)
	}
	f.SetCellStr(shnm, "B15", dbsht.Dbtableparallel.Contents)
	if dbsht.Dbtableparallel.Alarm == "R" {
		f.SetCellStyle(shnm, "B15", "B15", styleR)
	} else if dbsht.Dbtableparallel.Alarm == "B" {
		f.SetCellStyle(shnm, "B15", "B15", styleB)
	} else if dbsht.Dbtableparallel.Alarm == "G" {
		f.SetCellStyle(shnm, "B15", "B15", styleG)
	}
	f.SetCellStr(shnm, "B16", dbsht.Dbindexparallel.Contents)
	if dbsht.Dbindexparallel.Alarm == "R" {
		f.SetCellStyle(shnm, "B16", "B16", styleR)
	} else if dbsht.Dbindexparallel.Alarm == "B" {
		f.SetCellStyle(shnm, "B16", "B16", styleB)
	} else if dbsht.Dbindexparallel.Alarm == "G" {
		f.SetCellStyle(shnm, "B16", "B16", styleG)
	}
	f.SetCellStr(shnm, "B17", dbsht.Dbinvalidindex.Contents)
	if dbsht.Dbinvalidindex.Alarm == "R" {
		f.SetCellStyle(shnm, "B17", "B17", styleR)
	} else if dbsht.Dbinvalidindex.Alarm == "B" {
		f.SetCellStyle(shnm, "B17", "B17", styleB)
	} else if dbsht.Dbinvalidindex.Alarm == "G" {
		f.SetCellStyle(shnm, "B17", "B17", styleG)
	}
	f.SetCellStr(shnm, "B18", dbsht.Dbsequence.Contents)
	if dbsht.Dbsequence.Alarm == "R" {
		f.SetCellStyle(shnm, "B18", "B18", styleR)
	} else if dbsht.Dbsequence.Alarm == "B" {
		f.SetCellStyle(shnm, "B18", "B18", styleB)
	} else if dbsht.Dbsequence.Alarm == "G" {
		f.SetCellStyle(shnm, "B18", "B18", styleG)
	}
	f.SetCellStr(shnm, "B19", dbsht.Dbrecoverydest.Contents)
	if dbsht.Dbrecoverydest.Alarm == "R" {
		f.SetCellStyle(shnm, "B19", "B19", styleR)
	} else if dbsht.Dbrecoverydest.Alarm == "B" {
		f.SetCellStyle(shnm, "B19", "B19", styleB)
	} else if dbsht.Dbrecoverydest.Alarm == "G" {
		f.SetCellStyle(shnm, "B19", "B19", styleG)
	}
	f.SetCellStr(shnm, "B20", dbsht.Dbflashrecoveryuseage.Contents)
	if dbsht.Dbflashrecoveryuseage.Alarm == "R" {
		f.SetCellStyle(shnm, "B20", "B20", styleR)
	} else if dbsht.Dbflashrecoveryuseage.Alarm == "B" {
		f.SetCellStyle(shnm, "B20", "B20", styleB)
	} else if dbsht.Dbflashrecoveryuseage.Alarm == "G" {
		f.SetCellStyle(shnm, "B20", "B20", styleG)
	}
	f.SetCellStr(shnm, "B21", dbsht.Dberrlog.Contents)
	if dbsht.Dberrlog.Alarm == "R" {
		f.SetCellStyle(shnm, "B21", "B21", styleR)
	} else if dbsht.Dberrlog.Alarm == "B" {
		f.SetCellStyle(shnm, "B21", "B21", styleB)
	} else if dbsht.Dberrlog.Alarm == "G" {
		f.SetCellStyle(shnm, "B21", "B21", styleG)
	}
	f.SetCellStr(shnm, "B22", dbsht.Dbrmancheck.Contents)
	if dbsht.Dbrmancheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B22", "B22", styleR)
	} else if dbsht.Dbrmancheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B22", "B22", styleB)
	} else if dbsht.Dbrmancheck.Alarm == "G" {
		f.SetCellStyle(shnm, "B22", "B22", styleG)
	}
	f.SetCellStr(shnm, "B23", dbsht.Dbdbapriv.Contents)
	if dbsht.Dbdbapriv.Alarm == "R" {
		f.SetCellStyle(shnm, "B23", "B23", styleR)
	} else if dbsht.Dbdbapriv.Alarm == "B" {
		f.SetCellStyle(shnm, "B23", "B23", styleB)
	} else if dbsht.Dbdbapriv.Alarm == "G" {
		f.SetCellStyle(shnm, "B23", "B23", styleG)
	}
	f.SetCellStr(shnm, "B24", dbsht.Dbsysdba.Contents)
	if dbsht.Dbsysdba.Alarm == "R" {
		f.SetCellStyle(shnm, "B24", "B24", styleR)
	} else if dbsht.Dbsysdba.Alarm == "B" {
		f.SetCellStyle(shnm, "B24", "B24", styleB)
	} else if dbsht.Dbsysdba.Alarm == "G" {
		f.SetCellStyle(shnm, "B24", "B24", styleG)
	}
	f.SetCellStr(shnm, "B25", dbsht.Dbauditsegment.Contents)
	if dbsht.Dbauditsegment.Alarm == "R" {
		f.SetCellStyle(shnm, "B25", "B25", styleR)
	} else if dbsht.Dbauditsegment.Alarm == "B" {
		f.SetCellStyle(shnm, "B25", "B25", styleB)
	} else if dbsht.Dbauditsegment.Alarm == "G" {
		f.SetCellStyle(shnm, "B25", "B25", styleG)
	}
	f.SetCellStr(shnm, "B26", dbsht.Dbauditcont.Contents)
	if dbsht.Dbauditcont.Alarm == "R" {
		f.SetCellStyle(shnm, "B26", "B26", styleR)
	} else if dbsht.Dbauditcont.Alarm == "B" {
		f.SetCellStyle(shnm, "B26", "B26", styleB)
	} else if dbsht.Dbauditcont.Alarm == "G" {
		f.SetCellStyle(shnm, "B26", "B26", styleG)
	}
	f.SetCellStr(shnm, "B27", dbsht.Db_Nosys_In_System.Contents)
	if dbsht.Db_Nosys_In_System.Alarm == "R" {
		f.SetCellStyle(shnm, "B27", "B27", styleR)
	} else if dbsht.Db_Nosys_In_System.Alarm == "B" {
		f.SetCellStyle(shnm, "B27", "B27", styleB)
	} else if dbsht.Db_Nosys_In_System.Alarm == "G" {
		f.SetCellStyle(shnm, "B27", "B27", styleG)
	}
	f.SetCellStr(shnm, "B28", dbsht.Dbproductuserfailedlogin.Contents)
	if dbsht.Dbproductuserfailedlogin.Alarm == "R" {
		f.SetCellStyle(shnm, "B28", "B28", styleR)
	} else if dbsht.Dbproductuserfailedlogin.Alarm == "B" {
		f.SetCellStyle(shnm, "B28", "B28", styleB)
	} else if dbsht.Dbproductuserfailedlogin.Alarm == "G" {
		f.SetCellStyle(shnm, "B28", "B28", styleG)
	}
	f.SetCellStr(shnm, "B29", dbsht.Dbvirscheck.Contents)
	if dbsht.Dbvirscheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B29", "B29", styleR)
	} else if dbsht.Dbvirscheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B29", "B29", styleB)
	} else if dbsht.Dbvirscheck.Alarm == "G" {
		f.SetCellStyle(shnm, "B29", "B29", styleG)
	}
	f.SetCellStr(shnm, "B30", dbsht.Dbscnhealthcheck.Contents)
	if dbsht.Dbscnhealthcheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B30", "B30", styleR)
	} else if dbsht.Dbscnhealthcheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B30", "B30", styleB)
	} else if dbsht.Dbscnhealthcheck.Alarm == "G" {
		f.SetCellStyle(shnm, "B30", "B30", styleG)
	}
	f.SetCellStr(shnm, "B31", dbsht.Dbdglagcheck.Contents)
	if dbsht.Dbdglagcheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B31", "B31", styleR)
	} else if dbsht.Dbdglagcheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B31", "B31", styleB)
	} else if dbsht.Dbdglagcheck.Alarm == "G" {
		f.SetCellStyle(shnm, "B31", "B31", styleG)
	}
	f.SetCellStr(shnm, "B32", dbsht.Dbdgerrcheck.Contents)
	if dbsht.Dbdgerrcheck.Alarm == "R" {
		f.SetCellStyle(shnm, "B32", "B32", styleR)
	} else if dbsht.Dbdgerrcheck.Alarm == "B" {
		f.SetCellStyle(shnm, "B32", "B32", styleB)
	} else if dbsht.Dbdgerrcheck.Alarm == "G" {
		f.SetCellStyle(shnm, "B32", "B32", styleG)
	}
}

func PutSht_Summary(f *excelize.File, summaryEntries *structs.SummaryEntries) {
	shnm := "HealthReport"

	// 填充 Issue List（从 B25 开始）
	rowIndex := 25
	for i, entry := range summaryEntries.Entries {
		// 计算告警级别和描述
		var alarm, impact, desc string
		if len(entry.Severe) > 0 {
			alarm = "R"
			impact = "重要"
			desc = strings.Join(entry.Severe, "\n")
		} else if len(entry.Moderate) > 0 {
			alarm = "B"
			impact = "普通"
			desc = strings.Join(entry.Moderate, "\n")
		} else if len(entry.Minor) > 0 {
			alarm = "G"
			impact = "轻微"
			desc = strings.Join(entry.Minor, "\n")
		} else {
			alarm = ""
			impact = "PASS"
			desc = ""
		}

		var result int
		switch alarm {
		case "R":
			result = 0
		case "B":
			result = 5
		case "G":
			result = 8
		default:
			result = 10
		}

		// 填充一行
		f.SetCellInt(shnm, fmt.Sprintf("B%d", rowIndex), i+1)
		f.SetCellStr(shnm, fmt.Sprintf("C%d", rowIndex), entry.Category)
		f.SetCellStr(shnm, fmt.Sprintf("D%d", rowIndex), entry.Title)
		f.SetCellInt(shnm, fmt.Sprintf("E%d", rowIndex), result)
		f.SetCellStr(shnm, fmt.Sprintf("F%d", rowIndex), impact)
		f.SetCellStr(shnm, fmt.Sprintf("G%d", rowIndex), desc)

		rowIndex++
	}
}

func NewXlsx(xlsnm string) {
	f := excelize.NewFile()
	f.NewSheet("HealthReport")
	f.NewSheet("OS")
	f.NewSheet("DB")
	f.DeleteSheet("Sheet1")

	// 初始化 HealthReport Sheet
	shnm := "HealthReport"
	f.SetCellStr(shnm, "C1", "健康检查报告")
	f.SetCellStr(shnm, "G1", "Health Report")
	f.SetCellStr(shnm, "B13", "Issue summary")
	f.SetCellStr(shnm, "C14", "重要")
	f.SetCellStr(shnm, "D14", "普通")
	f.SetCellStr(shnm, "E14", "轻微")
	f.SetCellStr(shnm, "B15", "主机系统分析")
	f.SetCellStr(shnm, "B16", "数据库实例分析")
	f.SetCellStr(shnm, "B17", "数据库集群检查")
	f.SetCellStr(shnm, "B18", "DataGuard检查")
	f.SetCellStr(shnm, "B19", "数据库备份检查")
	f.SetCellStr(shnm, "B20", "数据库安全检查")
	f.SetCellStr(shnm, "B21", "软件使用分析")
	f.SetCellStr(shnm, "B22", "其他项检查")
	f.SetCellStr(shnm, "B23", "Issue list")
	f.SetCellStr(shnm, "B24", "No.")
	f.SetCellStr(shnm, "C24", "问题类别")
	f.SetCellStr(shnm, "D24", "检查项")
	f.SetCellStr(shnm, "E24", "结果")
	f.SetCellStr(shnm, "F24", "影响")
	f.SetCellStr(shnm, "G24", "问题描述及建议")

	// 初始化 OS Sheet
	shnm = "OS"
	f.SetCellStr(shnm, "A1", "主机名")
	f.SetCellStr(shnm, "A2", "IP地址")
	f.SetCellStr(shnm, "A3", "主机内核参数")
	f.SetCellStr(shnm, "A4", "主机资源限制")
	f.SetCellStr(shnm, "A5", "文件系统使用率")
	f.SetCellStr(shnm, "A6", "索引资源节点使用率")
	f.SetCellStr(shnm, "A7", "CPU负载")
	f.SetCellStr(shnm, "A8", "内存使用")
	f.SetCellStr(shnm, "A9", "磁盘IO负载检查")
	f.SetCellStr(shnm, "A10", "透明大页开启检查")
	f.SetCellStr(shnm, "A11", "主机大页使用检查")
	f.SetCellStr(shnm, "A12", "NUMA使用检查")
	f.SetCellStr(shnm, "A13", "NTP时钟同步检查")

	// 初始化 DB Sheet
	shnm = "DB"
	f.SetCellStr(shnm, "A1", "数据库名称\nDB_UNIQUE_NAME")
	f.SetCellStr(shnm, "A2", "主机名")
	f.SetCellStr(shnm, "A3", "表空间使用率")
	f.SetCellStr(shnm, "A4", "数据文件大小检查")
	f.SetCellStr(shnm, "A5", "控制文件检查")
	f.SetCellStr(shnm, "A6", "数据库用户大小")
	f.SetCellStr(shnm, "A7", "REDO文件性能检查")
	f.SetCellStr(shnm, "A8", "归档切换检查")
	f.SetCellStr(shnm, "A9", "数据库资源使用限制检查")
	f.SetCellStr(shnm, "A10", "数据库性能负载分析")
	f.SetCellStr(shnm, "A11", "数据库性能运行效率")
	f.SetCellStr(shnm, "A12", "数据库Top等待")
	f.SetCellStr(shnm, "A13", "数据库Top SQL(耗时)")
	f.SetCellStr(shnm, "A14", "监听状态及日志检查")
	f.SetCellStr(shnm, "A15", "并行度>1的表")
	f.SetCellStr(shnm, "A16", "并行度>1的索引")
	f.SetCellStr(shnm, "A17", "无效索引检查")
	f.SetCellStr(shnm, "A18", "Oracle序列检查")
	f.SetCellStr(shnm, "A19", "闪回区配置")
	f.SetCellStr(shnm, "A20", "FlashRecovery区使用情况")
	f.SetCellStr(shnm, "A21", "数据库日志检查")
	f.SetCellStr(shnm, "A22", "数据库RMAN备份")
	f.SetCellStr(shnm, "A23", "DBA权限用户检查")
	f.SetCellStr(shnm, "A24", "SYSDBA权限用户检查")
	f.SetCellStr(shnm, "A25", "数据库审计空间检查")
	f.SetCellStr(shnm, "A26", "数据库审计对象检查")
	f.SetCellStr(shnm, "A27", "业务对象存放系统表空间")
	f.SetCellStr(shnm, "A28", "错误口令登录锁定帐户PROFILE检查")
	f.SetCellStr(shnm, "A29", "病毒勒索攻击检查")
	f.SetCellStr(shnm, "A30", "SCNHealthCheck检查")
	f.SetCellStr(shnm, "A31", "DataGuard同步延迟检查")
	f.SetCellStr(shnm, "A32", "DataGuard同步报错检查")

	// 设置布局和样式
	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "right",
			Vertical:        "center",
			WrapText:        true,
			ShrinkToFit:     true,
			JustifyLastLine: true,
		},
	})
	styLeft, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#555555"},
		},
		Font: &excelize.Font{
			Bold:  true,
			Size:  11,
			Color: "#E6E6FA",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	shnms := f.GetSheetList()
	for _, shnm := range shnms {
		f.SetColWidth(shnm, "A", "A", 22)
		f.SetColStyle(shnm, "B:Z", wrapStyle)
		switch shnm {
		case "HealthReport":
			f.SetCellStyle(shnm, "B1", "B24", styLeft)
			f.SetColWidth(shnm, "B", "Z", 50)
		case "OS":
			f.SetCellStyle(shnm, "A1", "A13", styLeft)
			f.SetColWidth(shnm, "B", "Z", 80)
		case "DB":
			f.SetCellStyle(shnm, "A1", "A32", styLeft)
			f.SetColWidth(shnm, "B", "Z", 100)
		}
		f.SetPanes(shnm, &excelize.Panes{
			Freeze:      true,
			XSplit:      1,
			YSplit:      2,
			TopLeftCell: "B3",
			ActivePane:  "bottomLeft",
		})
	}
	f.SaveAs(xlsnm)
}
