package toxls

import (
	"autochk/structs"
	"fmt"
	"strings"

	// "hlchk/utils"
	"github.com/xuri/excelize/v2"
	// "github.com/myz-git/excelize"
	// "github.com/qax-os/excelize"
)

func Xlsx(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries, xlsnm string, colcnt int, sglf bool) {

	//open xlsx
	var newfnm string
	if sglf {
		newfnm = xlsnm + ".Done.xlsx"
	} else {
		// 保存为一个xlsx ,每个xml为一列
		newfnm = "HelthCheckReport.ALLDone.xlsx"
	}

	if colcnt == 1 { //第一次则新建文件
		NewXlsx(newfnm)
	}

	f, err := excelize.OpenFile(newfnm)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// log.Println("打开文件->", newfnm)
	//初始化工作表
	// InitSheet(f)

	///---sss---///

	PutSht_INFO(f, infstp, colcnt)

	//编辑工作表(OS)
	PutSht_OS(f, infstp, osshtp, colcnt)

	//编辑工作表(DB)
	PutSht_DB(f, infstp, dbshtp, colcnt)

	// Add Summary sheet
	PutSht_Summary(f, summaryEntries)

	//新文件名保留原来的fnm名字,后缀名为_new.xlsx
	// newcsvfnm := strings.Replace(fnm, ".xlsx", "_done.xlsx", -1)

	// //每一个xml使用保存一个xlsx
	// nowTime := time.Now().Format("20060102")
	// newcsvfnm := "_" + nowTime + ".xlsx"

	f.SaveAs(newfnm)

	///---eee---///

}

func PutSht_INFO(f *excelize.File, infstp *structs.InfoSht, colcnt int) {
	shnm := "INFO"
	// log.Printf("开始编辑工作表:%v,请等待....", shnm)

	//获取sheet有多少列
	// cols, _ := f.GetRows(shnm)
	// colnm := len(cols[0])

	// log.Printf("colcnt-------->%v,", colcnt)

	rows, _ := f.Rows(shnm) //按行的方式读取数据到二维数组中
	//[A1],[B1],[C1],[D1]....
	//[A2],[B2],[C2],[D2]....
	//....
	rowidx := 1       //行计数器
	for rows.Next() { //遍历每一行
		row, _ := rows.Columns() //将每一行赋值为一个一维数组row  :[A1],[B1],[C1],[D1]......
		// log.Printf(row[0])
		cellA := row[0] //row[0] 为每一行的第A列 , 如A1 ,A2,A3....
		//需要写入的单元格
		// v, h := excelize.CellNameToCoordinates(row[colcnt])
		aixcell, _ := excelize.CoordinatesToCellName(colcnt+1, rowidx) //定位第colcnt+1列, 第J行单元格
		// log.Println("aixcell--->", aixcell)
		switch cellA {
		case "数据库名称\nDB_UNIQUE_NAME":
			f.SetCellStr(shnm, aixcell, infstp.DbName)
		case "主机名":
			f.SetCellStr(shnm, aixcell, infstp.HostName)
		case "IP地址":
			f.SetCellStr(shnm, aixcell, infstp.Ipaddr)
		case "系统平台":
			f.SetCellStr(shnm, aixcell, infstp.Os)
		case "操作系统补丁":
			f.SetCellStr(shnm, aixcell, infstp.Relver)
		case "数据库版本":
			f.SetCellStr(shnm, aixcell, infstp.DbVer)
		case "数据库架构":
			f.SetCellStr(shnm, aixcell, infstp.DbMaa)
		case "主备库角色":
			f.SetCellStr(shnm, aixcell, infstp.DbRole)
		case "是否开启归档":
			f.SetCellStr(shnm, aixcell, infstp.LogMode)
		case "是否开启闪回":
			f.SetCellStr(shnm, aixcell, infstp.FlashBack)
		case "数据库总大小GB":
			f.SetCellStr(shnm, aixcell, infstp.DbTotalsize)
		case "数据库文件数":
			f.SetCellStr(shnm, aixcell, infstp.DbFilecount)
		case "数据库表数量":
			f.SetCellStr(shnm, aixcell, infstp.DbTblcount)
		case "数据库字符集":
			f.SetCellStr(shnm, aixcell, infstp.DbLang)
			// case "CORES\n(CPU)":
			// 	f.SetCellStr(shnm, aixcell, infstp.Cores)
		case "CPU核数":
			f.SetCellStr(shnm, aixcell, infstp.CpuCount)
		case "CPU频率":
			f.SetCellStr(shnm, aixcell, infstp.CpuMHZ)
		case "内存":
			f.SetCellStr(shnm, aixcell, infstp.MemTotal)
		case "交换分区":
			f.SetCellStr(shnm, aixcell, infstp.SwapTotal)
			// default:
			// 	f.SetCellStr(shnm, aixcell, infstp.Others)
		}
		rowidx++
	}
}

func PutSht_OS(f *excelize.File, infstp *structs.InfoSht, ossht *structs.OsSht, colcnt int) {
	//编辑工作表(OS)
	shnm := "OS"
	// log.Printf("开始编辑工作表:%v,请等待....", shnm)

	//设置excel styleB  单元格标记为蓝色
	styleB, _ := f.NewStyle(&excelize.Style{
		// Fill:      excelize.Fill{Type: "pattern", Color: []string{"#A020F0"}, Pattern: 1},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4876FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	//设置excel styleR  单元格标记为红色
	styleR, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFAEB9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	rows, _ := f.Rows(shnm) //按行的方式读取数据到二维数组中
	//[A1],[B1],[C1],[D1]....
	//[A2],[B2],[C2],[D2]....
	//....
	rowidx := 1       //行计数器
	for rows.Next() { //遍历每一行
		row, _ := rows.Columns() //将每一行赋值为一个一维数组row  :[A1],[B1],[C1],[D1]......
		// log.Printf(row[0])
		cellA := row[0] //row[0] 为每一行的第A列 , 如A1 ,A2,A3....
		//需要写入的单元格
		// v, h := excelize.CellNameToCoordinates(row[colcnt])
		aixcell, _ := excelize.CoordinatesToCellName(colcnt+1, rowidx) //定位第colcnt+1列, 第J行单元格
		// log.Println("aixcell--->", aixcell)
		switch cellA {
		case "主机名":
			f.SetCellStr(shnm, aixcell, infstp.HostName)
		case "IP地址":
			f.SetCellStr(shnm, aixcell, infstp.Ipaddr)
		case "主机内核参数":
			f.SetCellStr(shnm, aixcell, ossht.Osparameter.Contents)
			switch ossht.Osparameter.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
				// log.Println("主机内核参数aixcell----------------------->", aixcell)
			}
		case "主机LIMIT资源限制":
			f.SetCellStr(shnm, aixcell, ossht.Ulimit.Contents)
			switch ossht.Ulimit.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "文件系统使用率":
			f.SetCellStr(shnm, aixcell, ossht.Filesystem.Contents)
			switch ossht.Filesystem.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
			// log.Printf("Filesystem->\n%s", ossht.Filesystem.Contents)
		case "索引资源节点使用率":
			f.SetCellStr(shnm, aixcell, ossht.Inodeusage.Contents)
			switch ossht.Inodeusage.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "CPU负载":
			f.SetCellStr(shnm, aixcell, ossht.Cpustat.Contents)
			switch ossht.Cpustat.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "内存使用":
			f.SetCellStr(shnm, aixcell, ossht.Memstat.Contents)
			switch ossht.Memstat.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "磁盘IO负载检查":
			f.SetCellStr(shnm, aixcell, ossht.Iostat.Contents)
			switch ossht.Iostat.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "透明大页开启检查":
			f.SetCellStr(shnm, aixcell, ossht.Thpstat.Contents)
			switch ossht.Thpstat.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "主机大页使用检查":
			f.SetCellStr(shnm, aixcell, ossht.Hugpage.Contents)
			switch ossht.Hugpage.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "NUMA使用检查":
			f.SetCellStr(shnm, aixcell, ossht.Numa.Contents)
			switch ossht.Numa.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "NTP时钟同步检查":
			f.SetCellStr(shnm, aixcell, ossht.Ntp.Contents)
			switch ossht.Ntp.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
			// default:
			// 	f.SetCellStr(shnm, aixcell, infstp.Others)

		}
		rowidx++
	}

}

func PutSht_DB(f *excelize.File, infstp *structs.InfoSht, dbsht *structs.DbSht, colcnt int) {
	//编辑工作表(OS)
	shnm := "DB"
	// log.Printf("开始编辑工作表:%v,请等待....", shnm)

	//设置excel styleB  单元格标记为蓝色
	styleB, _ := f.NewStyle(&excelize.Style{
		// Fill:      excelize.Fill{Type: "pattern", Color: []string{"#A020F0"}, Pattern: 1},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4876FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	//设置excel styleR  单元格标记为红色
	styleR, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFAEB9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	//设置excel styleG  单元格标记为绿色
	styleG, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#00bf5f"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "top", WrapText: true},
	})

	rows, _ := f.Rows(shnm) //按行的方式读取数据到二维数组中
	//[A1],[B1],[C1],[D1]....
	//[A2],[B2],[C2],[D2]....
	//....
	rowidx := 1       //行计数器
	for rows.Next() { //遍历每一行
		row, _ := rows.Columns() //将每一行赋值为一个一维数组row  :[A1],[B1],[C1],[D1]......
		// log.Printf(row[0])
		cellA := row[0] //row[0] 为每一行的第A列 , 如A1 ,A2,A3....
		//需要写入的单元格
		// v, h := excelize.CellNameToCoordinates(row[colcnt])
		aixcell, _ := excelize.CoordinatesToCellName(colcnt+1, rowidx) //定位第colcnt+1列, 第J行单元格
		// log.Println("aixcell--->", aixcell)
		switch cellA {
		case "数据库名称\nDB_UNIQUE_NAME":
			f.SetCellStr(shnm, aixcell, infstp.DbName)
		case "主机名":
			f.SetCellStr(shnm, aixcell, infstp.HostName)
		case "表空间使用率":
			f.SetCellStr(shnm, aixcell, dbsht.DbTbsusage.Contents)
			switch dbsht.DbTbsusage.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "数据文件大小检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbdatafile.Contents)
			switch dbsht.Dbdatafile.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "控制文件检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbcontrolfile.Contents)
			switch dbsht.Dbcontrolfile.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "数据库用户大小":
			f.SetCellStr(shnm, aixcell, dbsht.Dbusersize.Contents)
		case "REDO文件性能检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbredocheck.Contents)
			switch dbsht.Dbredocheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}
		case "归档切换检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbredoswitch.Contents)
			switch dbsht.Dbredoswitch.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}

		case "数据库资源使用限制检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbresource.Contents)
			switch dbsht.Dbresource.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			}

		case "数据库性能负载分析":
			f.SetCellStr(shnm, aixcell, dbsht.Loadprofile.Contents)
			switch dbsht.Loadprofile.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "数据库性能运行效率":
			f.SetCellStr(shnm, aixcell, dbsht.Instefficiency.Contents)
			switch dbsht.Instefficiency.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "数据库Top等待":
			f.SetCellStr(shnm, aixcell, dbsht.Dbtopevent.Contents)
			switch dbsht.Dbtopevent.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "数据库Top SQL(耗时)":
			f.SetCellStr(shnm, aixcell, dbsht.DbtopSQL.Contents)
			switch dbsht.DbtopSQL.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "监听状态及日志检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dblsnrinfo.Contents)
			switch dbsht.Dblsnrinfo.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "并行度>1的表":
			f.SetCellStr(shnm, aixcell, dbsht.Dbtableparallel.Contents)
			switch dbsht.Dbtableparallel.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}

		case "并行度>1的索引":
			f.SetCellStr(shnm, aixcell, dbsht.Dbindexparallel.Contents)
			switch dbsht.Dbindexparallel.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "无效索引检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbinvalidindex.Contents)
			switch dbsht.Dbinvalidindex.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "Oracle序列检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbsequence.Contents)
			switch dbsht.Dbsequence.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "闪回区配置":
			f.SetCellStr(shnm, aixcell, dbsht.Dbrecoverydest.Contents)
			switch dbsht.Dbrecoverydest.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "FlashRecovery区使用情况":
			f.SetCellStr(shnm, aixcell, dbsht.Dbflashrecoveryuseage.Contents)
			switch dbsht.Dbflashrecoveryuseage.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "数据库日志检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dberrlog.Contents)
			switch dbsht.Dberrlog.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "数据库RMAN备份":
			f.SetCellStr(shnm, aixcell, dbsht.Dbrmancheck.Contents)
			switch dbsht.Dbrmancheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "DBA权限用户检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbdbapriv.Contents)
			switch dbsht.Dbdbapriv.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "SYSDBA权限用户检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbsysdba.Contents)
			switch dbsht.Dbsysdba.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "数据库审计空间检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbauditsegment.Contents)
			switch dbsht.Dbauditsegment.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "数据库审计对象检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbauditcont.Contents)
			switch dbsht.Dbauditcont.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "业务对象存放系统表空间":
			f.SetCellStr(shnm, aixcell, dbsht.Db_Nosys_In_System.Contents)
			switch dbsht.Db_Nosys_In_System.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "错误口令登录锁定帐户PROFILE检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbproductuserfailedlogin.Contents)
			switch dbsht.Dbproductuserfailedlogin.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "病毒勒索攻击检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbvirscheck.Contents)
			switch dbsht.Dbvirscheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "SCNHealthCheck检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbscnhealthcheck.Contents)
			switch dbsht.Dbscnhealthcheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "DataGuard同步延迟检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbdglagcheck.Contents)
			switch dbsht.Dbdglagcheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "DataGuard同步报错检查":
			f.SetCellStr(shnm, aixcell, dbsht.Dbdgerrcheck.Contents)
			switch dbsht.Dbdgerrcheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "RAC资源状态":
			f.SetCellStr(shnm, aixcell, dbsht.Dbcrscheck.Contents)
			switch dbsht.Dbcrscheck.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
		case "ASM磁盘使用":
			f.SetCellStr(shnm, aixcell, dbsht.Dbasmusage.Contents)
			switch dbsht.Dbasmusage.Alarm {
			case "R":
				f.SetCellStyle(shnm, aixcell, aixcell, styleR)
			case "B":
				f.SetCellStyle(shnm, aixcell, aixcell, styleB)
			case "G":
				f.SetCellStyle(shnm, aixcell, aixcell, styleG)
			}
			// case end
		}
		rowidx++
	}

}

func NewXlsx(xlsnm string) {

	f := excelize.NewFile()
	f.NewSheet("Summary")
	f.NewSheet("INFO")
	f.NewSheet("OS")
	f.NewSheet("DB")
	f.DeleteSheet("sheet1")

	shnm := "INFO"
	//按行的方式写入数据到二维数组中,
	f.SetCellStr(shnm, "A1", "数据库名称\nDB_UNIQUE_NAME")
	f.SetCellStr(shnm, "A2", "主机名")
	f.SetCellStr(shnm, "A3", "IP地址")
	f.SetCellStr(shnm, "A4", "系统平台")
	f.SetCellStr(shnm, "A5", "操作系统补丁")
	f.SetCellStr(shnm, "A6", "数据库版本")
	f.SetCellStr(shnm, "A7", "数据库架构")
	f.SetCellStr(shnm, "A8", "主备库角色")
	f.SetCellStr(shnm, "A9", "是否开启归档")
	f.SetCellStr(shnm, "A10", "是否开启闪回")
	f.SetCellStr(shnm, "A11", "数据库总大小GB")
	f.SetCellStr(shnm, "A12", "数据库文件数")
	f.SetCellStr(shnm, "A13", "数据库表数量")
	f.SetCellStr(shnm, "A14", "数据库字符集")
	f.SetCellStr(shnm, "A15", "CPU核数")
	f.SetCellStr(shnm, "A16", "CPU频率")
	f.SetCellStr(shnm, "A17", "内存")
	f.SetCellStr(shnm, "A18", "交换分区")

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
	// f.SetCellStr(shnm, "A33", "RAC资源状态")
	// f.SetCellStr(shnm, "A34", "ASM磁盘使用")
	f.SetCellStr(shnm, "A33", "----")
	f.SetCellStr(shnm, "A34", "----")

	//设定总体布局
	sty_left, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "right",   // top bottom left right diagonalDown diagonalUp 中的一个
				Color: "#000000", // 十六进制的颜色编码
				Style: 1,         // 0-13 有对应的样式
			},
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
		}, Fill: excelize.Fill{
			Type:    "pattern",            // gradient 渐变色    pattern   填充图案
			Pattern: 1,                    // 填充样式  当类型是 pattern 0-18 填充图案  1 实体填充
			Color:   []string{"##555555"}, // 当Type = pattern 时，只有一个
			// Color:   []string{"#00F700", "#00F700"},
			Shading: 1, // 类型是 gradient 使用 0-5 横向(每种颜色横向分布) 纵向 对角向上 对角向下 有外向内 由内向外
		}, Font: &excelize.Font{
			Bold: true,
			// Italic: false,
			// Underline: "single",
			Size: 11,
			// Family: "宋体",
			// Strike:    true, // 删除线
			Color: "#E6E6FA",
		}, Alignment: &excelize.Alignment{
			Horizontal: "center", // 水平对齐方式 center left right fill(填充) justify(两端对齐)  centerContinuous(跨列居中) distributed(分散对齐)
			Vertical:   "center", // 垂直对齐方式 center top  justify distributed
			// Indent:     1,        // 缩进  只要有值就变成了左对齐 + 缩进
			// TextRotation: 30, // 旋转
			// RelativeIndent:  10,   // 好像没啥用
			// ReadingOrder:    0,    // 不知道怎么设置
			JustifyLastLine: true, // 两端分散对齐，只有 水平对齐 为 distributed 时 设置true 才有效
			WrapText:        true, // 自动换行
			ShrinkToFit:     true, // 缩小字体以填充单元格

		}})

	//设置excel style2  单元格自动行换

	wrap_style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "right",
			Indent:          1,
			JustifyLastLine: true,
			ReadingOrder:    0,
			RelativeIndent:  1,
			ShrinkToFit:     true,
			Vertical:        "center",
			WrapText:        true,
		},
	})
	shnms := f.GetSheetList()
	for _, shnm := range shnms {
		//设定第一列样式生效的范围
		f.SetColWidth(shnm, "A", "A", 22)
		f.SetColStyle(shnm, "B:Z", wrap_style)
		switch shnm {
		case "INFO":
			f.SetCellStyle(shnm, "A1", "A18", sty_left)
			f.SetColWidth(shnm, "B", "Z", 50)
		case "OS":
			f.SetCellStyle(shnm, "A1", "A13", sty_left)
			f.SetColWidth(shnm, "B", "Z", 80)
		case "DB":
			f.SetCellStyle(shnm, "A1", "A34", sty_left)
			f.SetColWidth(shnm, "B", "Z", 100)
		}
		//设置 首行首列冻结
		f.SetPanes(shnm, &excelize.Panes{
			Freeze:      true,
			XSplit:      1,    //冻结首列
			YSplit:      2,    //冻结首两行
			TopLeftCell: "B3", //设置活动元格 (不能在冻结范围内)
			ActivePane:  "bottomLeft",
		})
		// // 设置首行高度
		// f.SetRowHeight(shnm, 1, 30)
		// f.SetRowHeight(shnm, 2, 80)
		// for h := 3; h < 10; h++ {
		// 	f.SetRowHeight(shnm, h, 180)
		// }
		// for h := 14; h < 18; h++ {
		// 	f.SetRowHeight(shnm, h, 80)
		// }
	}
	f.SaveAs(xlsnm)

}

func PutSht_Summary(f *excelize.File, summaryEntries *structs.SummaryEntries) {
	shnm := "Summary"
	// Define the headers for the Summary sheet
	headers := []string{"检查类别", "检查项", "检查说明", "检查结果(严重)", "检查结果(一般),", "检查结果(轻微)"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellStr(shnm, cell, header)
	}

	rowIndex := 2
	for _, entry := range summaryEntries.Entries {
		f.SetCellStr(shnm, fmt.Sprintf("A%d", rowIndex), entry.Category)
		f.SetCellStr(shnm, fmt.Sprintf("B%d", rowIndex), entry.Nm)
		f.SetCellStr(shnm, fmt.Sprintf("C%d", rowIndex), entry.Desc)
		f.SetCellStr(shnm, fmt.Sprintf("D%d", rowIndex), strings.Join(entry.Severe, ", "))
		f.SetCellStr(shnm, fmt.Sprintf("E%d", rowIndex), strings.Join(entry.Moderate, ", "))
		f.SetCellStr(shnm, fmt.Sprintf("F%d", rowIndex), strings.Join(entry.Minor, ", "))
		rowIndex++
	}
}
