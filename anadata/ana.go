package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// nproc_ne :=GetRule()
func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))

	for index, val := range strArr {
		res[index], _ = strconv.Atoi(val)
	}

	return res
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func Ana(infstp *structs.InfoSht, osshtp *structs.OsSht, dbshtp *structs.DbSht, summaryEntries *structs.SummaryEntries) {

	// log.Println("###--->Start Ana")
	rules, err := utils.GetRule()
	if err != nil {
		log.Printf("rule err: #%v", err)
	}

	//格式化展现函数(不用规则分析)
	Fmt_DbRole(infstp)
	Fmt_LogMode(infstp)
	Fmt_FlashBack(infstp)
	Fmt_DbTotalsize(infstp)
	Fmt_DbFilecount(infstp)
	Fmt_DbTblcount(infstp)

	//分析OS指标函数
	Ana_Osparameter(rules, osshtp, infstp, summaryEntries)
	Ana_Ulimit(rules, osshtp, infstp)
	Ana_Filesystem(rules, osshtp)
	Ana_Inodeusage(rules, osshtp)
	Ana_Cpustat(rules, osshtp)
	Ana_Memstat(rules, osshtp)
	Ana_Iostat(rules, osshtp)
	Ana_Thpstat(rules, osshtp)
	Ana_Numa(rules, osshtp)
	Ana_Ntp(rules, osshtp)

	//分析DB指标函数
	Ana_DbTbs(rules, dbshtp)
	Ana_DBF(rules, dbshtp)
	Ana_DBCTRF(rules, dbshtp)
	Ana_RDF(rules, dbshtp)
	Ana_RDSW(rules, dbshtp)
	Ana_RESOURCE(rules, dbshtp)
	Ana_LOADPROFILE(rules, dbshtp)
	Ana_INSTEFFICIENCY(rules, dbshtp)
	Ana_DBLSNRINFO(rules, dbshtp)
	Ana_DBTABLEPARALLEL(rules, dbshtp)
	Ana_DBINDEXPARALLEL(rules, dbshtp)
	Ana_DBINVALIDINDEX(rules, dbshtp)
	Ana_DBSEQUENCE(rules, dbshtp, infstp)
	Ana_DBFLASHRECOVERYUSEAGE(rules, dbshtp)
	Ana_DBERRLOG(rules, dbshtp)
	Ana_DBPRODUCTUSERFAILEDLOGIN(rules, dbshtp)
	Ana_DBDBAPRIV(rules, dbshtp)
	Ana_DBSYSDBA(rules, dbshtp)
	Ana_DBDGLAGCHECK(rules, dbshtp, infstp)
	Ana_DBDGERRCHECK(rules, dbshtp)
	Ana_DBRMANCHECK(rules, dbshtp)
	Ana_DBAUDITSEGMENT(rules, dbshtp)
	Ana_DBAUDITCONT(rules, dbshtp)
	Ana_DBVIRSCHECK(rules, dbshtp)
	Ana_DBSCNHEALTHCHECK(rules, dbshtp)

}
func Fmt_DbRole(infstp *structs.InfoSht) {
	msgdata := infstp.DbRole
	for index, value := range strings.Split(msgdata, "\n") { //按行拆分
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.DbRole = value
	}

}

func Fmt_LogMode(infstp *structs.InfoSht) {
	msgdata := infstp.LogMode
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.LogMode = strings.TrimSpace(value)
	}
}

func Fmt_FlashBack(infstp *structs.InfoSht) {
	msgdata := infstp.FlashBack
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.FlashBack = strings.TrimSpace(value)
	}
}

func Fmt_DbTotalsize(infstp *structs.InfoSht) {
	msgdata := infstp.DbTotalsize
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.DbTotalsize = strings.TrimSpace(value) + "  GB"
	}
}

func Fmt_DbFilecount(infstp *structs.InfoSht) {
	msgdata := infstp.DbFilecount
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.DbFilecount = strings.TrimSpace(value) + "  GB"
	}
}

func Fmt_DbTblcount(infstp *structs.InfoSht) {
	msgdata := infstp.DbTblcount
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 { //不是第三行就跳过(即只抓取第三行数据)
			continue
		}
		infstp.DbTblcount = strings.TrimSpace(value) + "  GB"
	}
}

func Ana_Osparameter(rule *utils.RuleInfo, osshtp *structs.OsSht, infstp *structs.InfoSht, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(infstp.Os)
	log.Println("rule.Osrule.Osparameter->", rule.Osrule.Osparameter)
	msgdata := osshtp.Osparameter.Contents
	entry := structs.SummaryEntry{
		Category: "OS",
		Nm:       rule.Osrule.Osparameter.Nm,
		Desc:     rule.Osrule.Osparameter.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}

		//当为LINUX系统时, 检查参项
		if strings.Contains(oS, "LINUX") {
			rd := regexp.MustCompile(`\d+$`)
			rnm1 := regexp.MustCompile(`nproc`)
			rnm2 := regexp.MustCompile(`nofile`)
			rnm3 := regexp.MustCompile(`randomize_va_space`)
			rnm4 := regexp.MustCompile(`panic_on_oops`)
			rnm5 := regexp.MustCompile(`min_free_kbytes`)

			if rnm1.MatchString(value) { // 判断是否匹配到 nproc :true or false
				matchs := rd.FindString(value) // 匹配以数字结尾的数值,  如16384
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_nproc_ne { //如果 < 16384 设定为 BLUE 重要告警级别, 且结束本次FOR循环
					log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Osparameter.L_nproc_ne)
					osshtp.Osparameter.Alarm = "B"
					// 添加到 SummaryEntries
					entry.Moderate = append(entry.Moderate, infstp.HostName)
					break Looop
				}
			}

			if rnm2.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_nofile_ne { //如果 <65536 设定为 BLUE 重要告警级别, 且结束本次FOR循环
					log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Osparameter.L_nofile_ne)
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, infstp.HostName)
					break Looop
				}
			}

			if rnm3.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.L_randomize_va_space {
					log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Osparameter.L_randomize_va_space)
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, infstp.HostName)
					break Looop
				}
			}
			if rnm4.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.L_panic_on_oops {
					log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Osparameter.L_panic_on_oops)
					osshtp.Osparameter.Alarm = "G"
					entry.Minor = append(entry.Minor, infstp.HostName)
					break Looop
				}
			}
			if rnm5.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_min_free_kbytes {
					log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Osparameter.L_min_free_kbytes)
					osshtp.Osparameter.Alarm = "G"
					entry.Minor = append(entry.Minor, infstp.HostName)
					break Looop
				}
			}
		}

		// if strings.Contains(value, "disable_ism_large_pages") && !strings.Contains(value, "0xF4") && !strings.Contains(value, "0x74") {
		// 	///log.Printf("!!Matched!! value [%v] match rule [%v]", value, rule.Osrule.Osparameter.S_disable_ism_large_pages)
		// 	osshtp.Osparameter.Alarm = "B"
		// 	break Looop
		// }

		//当为solaris系统时, 检查参项
		if strings.Contains(oS, "SOLARIS") {
			//检查disable_ism_large_pages参数设置  如:  set disable_i1sm_large_pages=0xF4
			if strings.Contains(value, "disable_ism_large_pages") {
				msg := strings.Split(value, "=") //取"="右侧数据 0xF4
				// log.Println("msg---->", msg[len(msg)-1])
				if !Contain(msg[len(msg)-1], rule.Osrule.Osparameter.S_disable_ism_large_pages) {
					//判断取到的值是否包含在 list [0xF4 0x74]中,假如没有则
					log.Printf("!!Matched!! value [%v] match rule [%v]", value, rule.Osrule.Osparameter.S_disable_ism_large_pages)
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, infstp.HostName)
					break Looop
				}
			}
		}

	}
	// log.Printf("Osparameter.Alarm->%s", osshtp.Osparameter.Alarm)
	// 如果有任何问题，添加到 SummaryEntries 中
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

func Ana_Ulimit(rule *utils.RuleInfo, osshtp *structs.OsSht, infstp *structs.InfoSht) {
	//log.Println("rule.Osrule.Ulimit->", rule.Osrule.Ulimit)
	msgdata := osshtp.Ulimit.Contents
	oS := strings.ToUpper(infstp.Os)
	// log.Println("oS--->", oS)

Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}

		rd := regexp.MustCompile(`\d+$`)
		rnm1 := regexp.MustCompile(`open files`)
		rnm2 := regexp.MustCompile(`max user processes`)
		if strings.Contains(oS, "LINUX") {
			rnm1 = regexp.MustCompile(`open files`)
			rnm2 = regexp.MustCompile(`max user processes`)
		} else if strings.Contains(oS, "SOLARIS") {
			rnm1 = regexp.MustCompile(`nofile`)
			rnm2 = regexp.MustCompile(`nproc`)
		}
		// } else if strings.Contains(oS, "AIX") {
		// 	rnm2 = regexp.MustCompile(`maxuprocs`)
		// // } else if strings.Contains(oS, "HP") {
		// 	rnm1 = regexp.MustCompile(`open files`)
		// 	rnm2 = regexp.MustCompile(`max user processes`)
		// }

		if rnm1.MatchString(value) { // 判断是否匹配到 openfiles :true or false
			matchs := rd.FindString(value) // 匹配以数字结尾的数值,  如65536
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Open_files_ne { //如果 < 65536 设定为 BLUE 重要告警级别, 且结束本次FOR循环

				///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Ulimit.Open_files_ne)
				osshtp.Ulimit.Alarm = "B"
				break Looop
			}

		}

		if rnm2.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Max_user_rocesses_ne { //如果 <16384 设定为 BLUE 重要告警级别, 且结束本次FOR循环

				///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Ulimit.Max_user_rocesses_ne)
				osshtp.Ulimit.Alarm = "B"
				break Looop
			}
		}
	}
	// log.Printf("Ulimit.Alarm->%s", osshtp.Ulimit.Alarm)
}

func Ana_Filesystem(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Filesystem->", rule.Osrule.Filesystem)
	msgdata := osshtp.Filesystem.Contents

	r := regexp.MustCompile(`\d+%`) //匹配 n%
	matchs := r.FindAllString(msgdata, -1)
	// log.Println("Filesystem->", matchs)

Looop:
	for _, p := range matchs {
		if p >= rule.Osrule.Filesystem.Disk_ge1 { //判断是否符合 重要告警级别

			///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Filesystem.Disk_ge1)
			osshtp.Filesystem.Alarm = "B"
			if p >= rule.Osrule.Filesystem.Disk_ge2 { //判断是否符合 严重告警级别

				///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Filesystem.Disk_ge2)
				osshtp.Filesystem.Alarm = "R"
				break Looop //当发现为严重告警级别 结束循环
			}
		}
	}
	// log.Printf("Filesystem.Alarm->%s", osshtp.Filesystem.Alarm)
}

func Ana_Inodeusage(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Inodeusage->", rule.Osrule.Filesystem) ///inodeusage和filesystem使用相同rule
	msgdata := osshtp.Inodeusage.Contents

	r := regexp.MustCompile(`\d+%`) //匹配 n%
	matchs := r.FindAllString(msgdata, -1)
	// log.Println("Filesystem->", matchs)

Looop:
	for _, p := range matchs {
		if p >= rule.Osrule.Filesystem.Disk_ge1 { //inodeusage和filesystem使用相同rule
			///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Filesystem.Disk_ge1)
			osshtp.Inodeusage.Alarm = "B"
			if p >= rule.Osrule.Filesystem.Disk_ge2 { //inodeusage和filesystem使用相同rule
				///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Filesystem.Disk_ge2)
				osshtp.Inodeusage.Alarm = "R"
				break Looop //当发现为严重告警级别 结束循环
			}
		}
	}
	// log.Printf("Filesystem.Alarm->%s", osshtp.Inodeusage.Alarm)
}

func Ana_Cpustat(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Cpustat->", rule.Osrule.Cpustat)
	msgdata := osshtp.Cpustat.Contents

Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 4 { //跳过vmstat的前面3行
			continue
		}

		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) { // 判断是否以数值结尾,
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			data := String2Int(msgs)      //将mst字符串数组转为数值数组

			// data[6] 对应si列 内存换页  ,data[14] 对应IDLEg列 CPU空闲使用
			switch {
			case data[6] >= rule.Osrule.Cpustat.Swap_ge2:
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data[6], rule.Osrule.Cpustat.Swap_ge2)
				osshtp.Cpustat.Alarm = "R"
				break Looop
			case data[14] < rule.Osrule.Cpustat.Idle_le2:
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data[14], rule.Osrule.Cpustat.Idle_le2)
				osshtp.Cpustat.Alarm = "R"
				break Looop
			case data[6] >= rule.Osrule.Cpustat.Swap_ge1:
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data[6], rule.Osrule.Cpustat.Swap_ge1)
				osshtp.Cpustat.Alarm = "B"
			case data[14] < rule.Osrule.Cpustat.Idle_le1:
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data[14], rule.Osrule.Cpustat.Idle_le1)
				osshtp.Cpustat.Alarm = "B"
			}
		}
	}
	// log.Printf("Cpustat.Alarm->%s", osshtp.Cpustat.Alarm)
}

func Ana_Memstat(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Memstat->", rule.Osrule.Memstat)
	msgdata := osshtp.Memstat.Contents

	if msgdata == "" { //如果为空则返回
		log.Printf("rule.Osrule.Memstat--->[%v] No data found!!! ", "msgdata")
		return
	}

	re := regexp.MustCompile(`\d+\s+`)
	vals := re.FindAllString(msgdata, -1) //获得由数字组成的字符串数组 (含了换行符)

	matchs, _ := strconv.Atoi(strings.ReplaceAll(vals[5], "\n", "")) //去除换行符,字符串转数字
	if matchs < rule.Osrule.Memstat.Available_le1 {
		///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Memstat.Available_le1)
		osshtp.Memstat.Alarm = "B"
		if matchs < rule.Osrule.Memstat.Available_le2 {
			///log.Printf("!!Matched!! value [%v] match rule [%v]", matchs, rule.Osrule.Memstat.Available_le2)
			osshtp.Memstat.Alarm = "R"
		}
	}
	// log.Printf("Memstat.Alarm->%s", osshtp.Memstat.Alarm)
}

func Ana_Iostat(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Iostat->", rule.Osrule.Iostat)
	msgdata := osshtp.Iostat.Contents
	if msgdata == "" { //如果为空则返回
		log.Printf("rule.Osrule.Iostat--->[%v] No data found!!! ", "Iostat")
		return
	}

Looop:
	for index, row := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		re := regexp.MustCompile(`^Average.*\d+$`) //匹配以Averageg开始,以数字结尾的行记录

		if re.MatchString(row) { // 判断是否匹配到
			// fmt.Println("match->", index, row)
			msgs := strings.Fields(row)                            //拆分为按空格的字符串数组
			data, err := strconv.ParseFloat(msgs[len(msgs)-1], 64) //取最右边一位数组元素 ,然后转化为float64
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println("match->", index, data)
			if data >= rule.Osrule.Iostat.Diskutil_ge1 {
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data, rule.Osrule.Iostat.Diskutil_ge1)
				osshtp.Iostat.Alarm = "B"
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge2 {
				///log.Printf("!!Matched!! value [%v] match rule [%v]", data, rule.Osrule.Iostat.Diskutil_ge2)
				osshtp.Iostat.Alarm = "R"
				break Looop
			}
		}

	}
	// log.Printf("Iostat.Alarm->%s", osshtp.Iostat.Alarm)
}

func Ana_Thpstat(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Thpstat->", rule.Osrule.Thpstat)
	msgdata := osshtp.Thpstat.Contents
	re := regexp.MustCompile(`\d+`)
	matchs := re.FindString(msgdata)
	data, _ := strconv.Atoi(matchs)
	if data > rule.Osrule.Thpstat.Anpages_gt {
		///log.Printf("!!Matched!! value [%v] match rule [%v]", data, rule.Osrule.Thpstat.Anpages_gt)
		osshtp.Thpstat.Alarm = "B"
	}
	// log.Printf("Thpstat.Alarm->%s", osshtp.Thpstat.Alarm)
}

func Ana_Numa(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Numa->", rule.Osrule.Numa)
	msgdata := osshtp.Numa.Contents
Looop:
	for index, row := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}

		re1 := regexp.MustCompile(rule.Osrule.Numa.Flg1) //逐行查找是否有"No NUMA"关键词
		if re1.MatchString(row) {                        //当找到 "No NUMA" 表明未使用NUMA, 判断结束
			osshtp.Numa.Alarm = ""
			break Looop
		}
		re2 := regexp.MustCompile(rule.Osrule.Numa.Flg2) //逐行查找是否有"No NUMA"关键词
		if re2.MatchString(row) {                        //当找到 "No NUMA" 表明未使用NUMA, 判断结束
			osshtp.Numa.Alarm = ""
			break Looop
		}
		osshtp.Numa.Alarm = "B" //当未找到"No NUMA" ,满足ALarm条件
		// ///log.Printf("!!Matched!! value [%v] match rule [%v]", row, rule.Osrule.Numa.Flg)
	}
	// log.Printf("Numa.Alarm->%s", osshtp.Numa.Alarm)
}

func Ana_Ntp(rule *utils.RuleInfo, osshtp *structs.OsSht) {
	//log.Println("rule.Osrule.Ntp->is Not Null")
	msgdata := osshtp.Ntp.Contents
	// 去除空格
	str := strings.Replace(msgdata, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	if str == "" { //如果str为空则未配置NTP
		///log.Printf("!!Matched!! value [%v] match rule [%v]", str, "is Not Null")
		osshtp.Ntp.Alarm = "B"
	}
	// log.Printf("Ntp.Alarm->%s", osshtp.Ntp.Alarm)
}

func Ana_DbTbs(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.DbTbsusage->", rule.Dbrule.DbTbsusage)
	msgdata := dbshtp.DbTbsusage.Contents

Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 1 { //跳过前面1行
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) { // 判断是否以数值结尾,
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			if len(msgs) < 5 {            //不足5列 跳过当前行
				continue
			}
			// log.Println("msgs->", msgs)
			// msgs[4] 对应USED_RATE列 ,msgs[3]对应 USED(GB), msgs[1]对应 MAXEXTEND(GB)
			maxsize, _ := strconv.ParseFloat(msgs[1], 64)  //取左边第二个数组元素 ,然后转化为float64
			usedsize, _ := strconv.ParseFloat(msgs[3], 64) //取左边第4个数组元素 ,然后转化为float64
			// percent, _ := strconv.ParseFloat(msgs[len(msgs)-1], 64) //取最右边第一位数组元素 ,然后转化为float64
			percent, _ := strconv.ParseFloat(msgs[4], 64) //取左边第5个数组元素 ,然后转化为float64

			switch {
			//组合条件:使用率>=90%,且可用空间<=4G
			case percent >= rule.Dbrule.DbTbsusage.Tbsutil_ge2 && (maxsize-usedsize) < rule.Dbrule.DbTbsusage.Freesize_le2:
				///log.Printf("!!Matched!! value [%v],[%v] match rule [%v]", percent, (maxsize - usedsize), rule.Dbrule.DbTbsusage.Tbsutil_ge2)
				dbshtp.DbTbsusage.Alarm = "R"
				///log.Printf("!!Matched!! value [%v]", percent)
				break Looop
			//组合条件 :使用率>=80%,且可用空间<8G
			case percent >= rule.Dbrule.DbTbsusage.Tbsutil_ge1 && (maxsize-usedsize) < rule.Dbrule.DbTbsusage.Freesize_le1:
				///log.Printf("!!Matched!! value [%v],[%v] match rule [%v]", percent, (maxsize - usedsize), rule.Dbrule.DbTbsusage.Tbsutil_ge1)
				dbshtp.DbTbsusage.Alarm = "B"
				///log.Printf("!!Matched!! value [%v]", percent)

			}
		}
	}
	// log.Printf("DbTbsusage.Alarm->%s", dbshtp.DbTbsusage.Alarm)
}

func Ana_DBF(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbdatafile->", rule.Dbrule.Dbdatafile.Status)
	msgdata := dbshtp.Dbdatafile.Contents

Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 1 { //跳过前面1行
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) { // 判断是否以数值结尾,
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			if len(msgs) < 5 {            //不足5列 跳过当前行
				continue
			}

			dbfstatus := msgs[1] //取左边第1个数组元素
			if dbfstatus != rule.Dbrule.Dbdatafile.Status {
				dbshtp.Dbdatafile.Alarm = "R"
				///log.Printf("!!Matched!! value [%v]", dbfstatus)
				break Looop
			}

		}
	}
	// log.Printf("Dbdatafile.Alarm->%s", dbshtp.Dbdatafile.Alarm)
}

func Ana_DBCTRF(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbcontrolfile->", rule.Dbrule.Dbcontrolfile.Cnt_le1)
	msgdata := dbshtp.Dbcontrolfile.Contents
	index := len(strings.Split(msgdata, "\n"))
	if index < rule.Dbrule.Dbcontrolfile.Cnt_le1 { //只有一行
		dbshtp.Dbcontrolfile.Alarm = "B"
		///log.Printf("!!Matched!! value [%v]", index)
	}

	// log.Printf("Dbcontrolfile.Alarm->%s", dbshtp.Dbcontrolfile.Alarm)
}

func Ana_RDF(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbredocheck->", rule.Dbrule.Dbredocheck)
	msgdata := dbshtp.Dbredocheck.Contents

Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value)
		rd := regexp.MustCompile(`^\d`) // 判断是否以数值开始的行,
		if rd.MatchString(value) {
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组

			if len(msgs) < 6 { //不足5列 跳过当前行
				continue
			}

			data := msgs[5]
			if !Contain(data, rule.Dbrule.Dbredocheck.Rdf_status_list) {
				dbshtp.Dbredocheck.Alarm = "R"
				///log.Printf("!!Matched!! value [%v]", data)
				break Looop
			}

			// matchs, _ := strconv.Atoi(msgs[3]) //去除换行符,字符串转数字
			matchs, _ := strconv.ParseFloat(msgs[3], 64) //去除换行符,字符串转数字
			// log.Println("msgs---------->", msgs)
			if matchs < rule.Dbrule.Dbredocheck.Rdf_size_lt1 {
				dbshtp.Dbredocheck.Alarm = "B"
				///log.Printf("!!Matched!! value [%v]", matchs)
			}

		}
	}
	// log.Printf("Dbredocheck.Alarm->%s", dbshtp.Dbredocheck.Alarm)
}

func Ana_RDSW(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbredoswitch->", rule.Dbrule.Dbredoswitch)
	msgdata := dbshtp.Dbredoswitch.Contents
	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:

	for index, msgs := range strings.Split(msgdata, "\n") {
		msgs = strings.TrimSpace(msgs) //去除头尾空格及空行
		if index < 2 {                 //跳过前面2行
			continue
		}

		if rd.MatchString(msgs) {
			msg := strings.Fields(msgs) //将以空格分隔的字符串转为 字符串数组

			for k, v := range msg {
				if k < 3 { //跳过前面3列
					continue
				}
				value, _ := strconv.Atoi(v)
				if value > rule.Dbrule.Dbredoswitch.Sw_cnt_ge1 {
					dbshtp.Dbredoswitch.Alarm = "B"
					///log.Printf("!!Matched!! value [%v]", value)
					break Looop
				}
			}
		}
	}
	// log.Printf("Dbredoswitch.Alarm->%s", dbshtp.Dbredoswitch.Alarm)
}

func Ana_RESOURCE(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbresource->", rule.Dbrule.Dbresource)
	msgdata := dbshtp.Dbresource.Contents

Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行
		rd := regexp.MustCompile(`\d$`)  // 判断是否以数值结尾的行
		if rd.MatchString(value) {
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			if len(msgs) < 5 {            //不足5列 跳过当前行
				continue
			}

			data1, err := strconv.Atoi(msgs[3]) //MAX_UTILIZATION
			if err != nil {
				panic(err)
			}
			data2, err := strconv.Atoi(msgs[4]) //LIMIT_VALUE
			if err != nil {
				panic(err)
			}
			// log.Println("data1,dat2,ge1", data1, data2, data2*rule.Dbrule.Dbresource.Res_use_ge1/100)
			if data1 >= data2*rule.Dbrule.Dbresource.Res_use_ge1/100 && data2 != 0 {
				dbshtp.Dbresource.Alarm = "R"
				///log.Printf("!!Matched!! value [%v]", data1)
				break Looop
			}
		}
	}
	// log.Printf("Dbresource.Alarm->%s", dbshtp.Dbresource.Alarm)
}

func Ana_LOADPROFILE(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Loadprofile->", rule.Dbrule.Loadprofile)
	msgdata := dbshtp.Loadprofile.Contents

Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行
		rd := regexp.MustCompile(`\d$`)  // 判断是否以数值结尾的行
		if rd.MatchString(value) {
			// msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			msgs := strings.Split(value, ":") //将每一行按":"分割成两个数组
			if len(msgs) < 2 {                //不足2列 跳过当前行
				continue
			}
			// log.Println("msgs--->", msgs[1])

			submsg := strings.Fields(msgs[1])
			str := strings.Replace(submsg[0], ",", "", -1) //去除原数据中的千分位逗号分隔符 "123,456.00"

			data, err := strconv.ParseFloat(str, 64)
			if err != nil {
				panic(err)
			}

			if msgs[0] == "Redo size (bytes)" && data >= rule.Dbrule.Loadprofile.Redosize_ge1*1024*1024 {
				dbshtp.Loadprofile.Alarm = "G"
				///log.Printf("!!Matched!! value [%v]", data)
				// break Looop
			}
			if msgs[0] == "Logons" && data >= rule.Dbrule.Loadprofile.Logons_ge1 {
				dbshtp.Loadprofile.Alarm = "B"
				///log.Printf("!!Matched!! value [%v]", data)
				break Looop
			}
		}
	}
	// log.Printf("Loadprofile.Alarm->%s", dbshtp.Loadprofile.Alarm)
}

func Ana_INSTEFFICIENCY(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Instefficiency->", rule.Dbrule.Instefficiency)
	msgdata := dbshtp.Instefficiency.Contents

	rd := regexp.MustCompile(`\d$`)                                      // 匹配以数值结尾的行
	rd2 := regexp.MustCompile(`^[1-9]\d*\.\d+$|^0\.\d+$|^[1-9]\d*$|^0$`) //匹配浮点数及整数 123 1.23  0.123等,但不包括"-"负数
Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行

		if rd.MatchString(value) {
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			// msgs := strings.Split(value, ":") //将每一行按":"分割成两个数组
			if len(msgs) < 4 { //不足4列 跳过当前行, 不做分析
				continue
			}
			// log.Println("msgs--->", msgs)

			for k := 0; k < len(msgs)-3; k++ { //循环当前行拆分后的数组
				//len(msgs)-3 是因为每次循环会查看k+3个元素的值是什么,所以最后3个元素不用循环
				// log.Printf("msgdata---> %v %v: %v", msgs[k], msgs[k+1], msgs[k+3])

				if msgs[k] == "Buffer" && msgs[k+1] == "Hit" && rd2.MatchString(msgs[k+3]) {
					// 判断连续两个数组元素是否 Buffer和Hit,并且第四个词(K+3)是个数值(整型或浮点型)  如 Buffer  Hit   %:   79.57
					//注意k+3要符合前面的len(msgs) < 4 的限制,否则数组越界
					data, err := strconv.ParseFloat(msgs[k+3], 64) //取到 buffer hit命中率数值
					if err != nil {
						panic(err)
					}
					if data < rule.Dbrule.Instefficiency.Buffer_hit {
						dbshtp.Instefficiency.Alarm = "G" //因为都是G级别,只要匹配到就可以结束了
						///log.Printf("!!Matched!! value [%v]", data)
						break Looop
					}
				}

				if msgs[k] == "Library" && msgs[k+1] == "Hit" && rd2.MatchString(msgs[k+3]) {
					data, err := strconv.ParseFloat(msgs[k+3], 64)
					if err != nil {
						panic(err)
					}
					if data < rule.Dbrule.Instefficiency.Library_hit {
						dbshtp.Instefficiency.Alarm = "G"
						///log.Printf("!!Matched!! value [%v]", data)
						break Looop
					}
				}

				if msgs[k] == "Soft" && msgs[k+1] == "Parse" && rd2.MatchString(msgs[k+3]) {
					data, err := strconv.ParseFloat(msgs[k+3], 64)
					if err != nil {
						panic(err)
					}
					if data < rule.Dbrule.Instefficiency.Soft_parse {
						dbshtp.Instefficiency.Alarm = "G"
						///log.Printf("!!Matched!! value [%v]", data)
						break Looop
					}
				}
			}

		}
	}
	// log.Printf("Instefficiency.Alarm->%s", dbshtp.Instefficiency.Alarm)
}

func Ana_DBLSNRINFO(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dblsnrinfo->", rule.Dbrule.Dblsnrinfo)
	msgdata := dbshtp.Dblsnrinfo.Contents

	rd := regexp.MustCompile(`(?i)\.log$`)                                                                 // 匹配以.log结尾的行  ?i表示不分大小写 如 .LOG
	rd2 := regexp.MustCompile(`^Jan$|^Feb$|^Mar$|^Apr $|^May$|^Jun$|^Jul$|^Aug$|^Sept$|^Oct$|^Nov$|^Dec$`) //匹配12个月份
	rd3 := regexp.MustCompile(`^\d+$`)                                                                     //匹配数字
Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行

		if rd.MatchString(value) { //匹配以.log结尾的行
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			// msgs := strings.Split(value, ":") //将每一行按":"分割成两个数组
			if len(msgs) < 8 { //不足8列 跳过当前行, 不做分析
				continue
			}
			// log.Println("msgs--->", msgs)

			for k := 3; k < len(msgs); k++ { //循环当前行拆分后的数组
				//定位锚点是3字母的月份,一般位于第6个元素后面, 所以前面几个元素不用循环检查
				// log.Printf("msgdata---> %v %v: %v", msgs[k], msgs[k+1], msgs[k+3])

				if rd2.MatchString(msgs[k]) && rd3.MatchString(msgs[k-1]) {
					//匹配定位锚点3字母月份, 且 前一个元素是数值

					data, _ := strconv.Atoi(msgs[k-1]) //取到3字母月份的前一个元素,即log文件大小
					if data >= rule.Dbrule.Dblsnrinfo.Log_size {
						dbshtp.Dblsnrinfo.Alarm = "G"
						///log.Printf("!!Matched!! value [%v]", data)
						break Looop
					}
				}

			}

		}
	}
	// log.Printf("Dblsnrinfo.Alarm->%s", dbshtp.Dblsnrinfo.Alarm)
}

func Ana_DBTABLEPARALLEL(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbtableparallel->", rule.Dbrule.Dbtableparallel)
	msgdata := dbshtp.Dbtableparallel.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)                 //去除头尾空格及空行
		if value == rule.Dbrule.Dbtableparallel.Result { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbtableparallel.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbtableparallel.Alarm->%s", dbshtp.Dbtableparallel.Alarm)
}

func Ana_DBINDEXPARALLEL(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbindexparallel->", rule.Dbrule.Dbindexparallel)
	msgdata := dbshtp.Dbindexparallel.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)                 //去除头尾空格及空行
		if value == rule.Dbrule.Dbindexparallel.Result { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbindexparallel.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbindexparallel.Alarm->%s", dbshtp.Dbindexparallel.Alarm)
}

func Ana_DBINVALIDINDEX(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbinvalidindex->", rule.Dbrule.Dbinvalidindex)
	msgdata := dbshtp.Dbinvalidindex.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)                //去除头尾空格及空行
		if value == rule.Dbrule.Dbinvalidindex.Result { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbinvalidindex.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbinvalidindex.Alarm->%s", dbshtp.Dbinvalidindex.Alarm)
}

func Ana_DBSEQUENCE(rule *utils.RuleInfo, dbshtp *structs.DbSht, infstp *structs.InfoSht) {

	if !strings.Contains(infstp.DbMaa, "RAC") { //如果不是RAC环境,不做判断
		return
	}

	msgdata := dbshtp.Dbsequence.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)            //去除头尾空格及空行
		if value == rule.Dbrule.Dbsequence.Result { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbsequence.Alarm = "G"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbsequence.Alarm->%s", dbshtp.Dbsequence.Alarm)
}

func Ana_DBFLASHRECOVERYUSEAGE(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbflashrecoveryuseage->", rule.Dbrule.Dbflashrecoveryuseage)
	msgdata := dbshtp.Dbflashrecoveryuseage.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value) //去除头尾空格及空行
		// log.Println("value--------->", value)
		if rd.MatchString(value) {
			msgs := strings.Fields(value) //将以空格分隔的字符串转为 字符串数组
			// msgs := strings.Split(value, "    ") //将每一行按":"分割成两个数组
			// log.Println("msgs--------->", msgs)
			// log.Println("msgs[1]--------->", msgs[1])
			if len(msgs) < 3 { //不足4列 跳过当前行, 不做分析
				continue
			}

			data, err := strconv.ParseFloat(msgs[len(msgs)-1], 64)
			if err != nil {
				panic(err)

			}

			// log.Println("data--------->", data)
			if data >= rule.Dbrule.Dbflashrecoveryuseage.Useage1 {
				dbshtp.Dbflashrecoveryuseage.Alarm = "B"
				///log.Printf("!!Matched!! value [%v]", data)
				// break Looop
			}
			if data >= rule.Dbrule.Dbflashrecoveryuseage.Useage2 {
				dbshtp.Dbflashrecoveryuseage.Alarm = "R"
				///log.Printf("!!Matched!! value [%v]", data)
				break Looop
			}
		}

	}
	// log.Printf("Dbflashrecoveryuseage.Alarm->%s", dbshtp.Dbflashrecoveryuseage.Alarm)
}

func Ana_DBERRLOG(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := dbshtp.Dberrlog.Contents
	value := strings.TrimSpace(msgdata)
	dbshtp.Dberrlog.Alarm = "B"
	if value == "" || strings.Contains(value, rule.Dbrule.Dberrlog.ResultB) { //空 或匹配到"no rows selected" ,则去除告警标记后退出
		dbshtp.Dberrlog.Alarm = ""
		return
	}
}

func Ana_DBPRODUCTUSERFAILEDLOGIN(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbproductuserfailedlogin->", rule.Dbrule.Dbproductuserfailedlogin)
	msgdata := dbshtp.Dbproductuserfailedlogin.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)                          //去除头尾空格及空行
		if value == rule.Dbrule.Dbproductuserfailedlogin.Result { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbproductuserfailedlogin.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbproductuserfailedlogin.Alarm->%s", dbshtp.Dbproductuserfailedlogin.Alarm)
}

func Ana_DBDGLAGCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht, infstp *structs.InfoSht) {

	if strings.Contains(infstp.DbRole, "PRIMARY") { //如果是主库,不做判断
		return
	}

	msgdata := dbshtp.Dbdglagcheck.Contents
	// rd1 := regexp.MustCompile(`:\d+$`)      //匹配以:数字结尾
	rdok := regexp.MustCompile(`^apply lag(.*)\+00 00:00:00$`) //匹配以"apply lag     +00 00:00:00"
	rd := regexp.MustCompile(`^apply lag(.*)\+(.*):\d+$`)      //匹配以apply lag开头,包含+数字 且以:数字结尾
	dbshtp.Dbdglagcheck.Alarm = "G"                            //初始化为绿
Looop:
	//按行分割
	for _, row := range strings.Split(msgdata, "\n") {
		row = strings.TrimSpace(row) //去除头尾空格及空行
		if rdok.MatchString(row) {   //匹配到apply lag   且为  +00 00:00:00 ,则去掉告警, 并退出
			dbshtp.Dbdglagcheck.Alarm = ""
			break Looop
		}
		if rd.MatchString(row) { //匹配到apply lag   且不为+00 00:00:00
			values1 := strings.Split(row, "+")            //按+号进行分组   [apply lag]     +[12 00:00:00]
			values2 := strings.Fields(values1[1])         // 对values[1] 即 “12 00:00:00” 再次按空格分组 [12][00:00:00]
			vDay, _ := strconv.Atoi(values2[0])           //  12
			if vDay >= rule.Dbrule.Dbdglagcheck.ResultB { //延迟超过阀值(>=1天)
				dbshtp.Dbdglagcheck.Alarm = "B"
				// log.Printf("!!Matched!! value [%v]", vDay)
				break Looop
			}
		}
	}
	// log.Printf("Dbvirscheck.Alarm->%s", dbshtp.Dbvirscheck.Alarm)
}

func Ana_DBDGERRCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := dbshtp.Dbdgerrcheck.Contents
	value := strings.TrimSpace(msgdata)
	if value != "" { //判断是否空, 为空正常, 非空则标记后退出
		dbshtp.Dbdgerrcheck.Alarm = "G"
	}

}

func Ana_DBRMANCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	// log.Println("rule.Dbrule.Dbrmancheck->", rule.Dbrule.Dbrmancheck)
	msgdata := dbshtp.Dbrmancheck.Contents
	if strings.TrimSpace(msgdata) == "" { //判断是否空, 为空则标记后退出
		dbshtp.Dbrmancheck.Alarm = "G"
		return
	}
Looop:
	//按行分割
	for _, row := range strings.Split(msgdata, "\n") {
		row = strings.TrimSpace(row) //去除头尾空格及空行

		re1 := regexp.MustCompile(rule.Dbrule.Dbrmancheck.ResultR) //逐行查找是否有"ERROR"关键词
		if re1.MatchString(strings.ToUpper(row)) {
			dbshtp.Dbrmancheck.Alarm = "R"
			break Looop
		}
		re2 := regexp.MustCompile(rule.Dbrule.Dbrmancheck.ResultB) //逐行查找是否有"WARNINGS"关键词
		if re2.MatchString(row) {
			dbshtp.Dbrmancheck.Alarm = "B"
			break Looop
		}

	}

	// log.Printf("Dbrmancheck.Alarm->%s", dbshtp.Dbrmancheck.Alarm)
}

func Ana_DBAUDITSEGMENT(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := dbshtp.Dbauditsegment.Contents
	value := strings.TrimSpace(msgdata)
	if value != "" { //判断是否空, 为空正常, 非空则标记后退出
		dbshtp.Dbauditsegment.Alarm = "G"
	}

}

func Ana_DBAUDITCONT(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := strings.TrimSpace(dbshtp.Dbauditcont.Contents) //去除头尾空格及空行
	rd := regexp.MustCompile(` \d+$`)                         //匹配以空格+数字结尾

Looop:
	for index, row := range strings.Split(msgdata, "\n") { //按行分割
		if index < 2 { //跳过前面2行 (如 column head  和 ----)
			continue
		}
		if rd.MatchString(row) { //匹配以"空格+数字"结尾的行
			msgs := strings.Fields(row) // 以空格分隔的字符串转为 字符串数组
			// msgs := strings.Split(value, ":") //将每一行按":"分割成两个数组
			log.Println("msgs[0]--------->", msgs[0])
			// if len(msgs) < 8 { //不足8列 跳过当前行, 不做分析
			if len(msgs) > 1 { //这里正常只应有一列, 超过1列则数据有问题,跳过当前行,不做分析
				continue
			}
			data, _ := strconv.Atoi(msgs[0]) //定位需要匹配的列是当前行拆分后转换的字符串数组第几个元素 ,这里取第一列
			if data >= rule.Dbrule.Dbauditcont.ResultG {
				dbshtp.Dbauditcont.Alarm = "G"
				log.Printf("!!Matched!! value [%v]", data)
				break Looop
			}

		}
	}
}

func Ana_DBVIRSCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbvirscheck->", rule.Dbrule.Dbvirscheck)
	msgdata := dbshtp.Dbvirscheck.Contents

	rd := regexp.MustCompile(` \d+$`) //匹配以空格+数字结尾

Looop:
	//按行分割
	for _, value := range strings.Split(msgdata, "\n") {
		value = strings.TrimSpace(value)              //去除头尾空格及空行
		if value == rule.Dbrule.Dbvirscheck.ResultR { //匹配到"no rows selected" ,或者没有记录则结束循环
			break Looop
		}
		if rd.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbvirscheck.Alarm = "R"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbvirscheck.Alarm->%s", dbshtp.Dbvirscheck.Alarm)
}

func Ana_DBSCNHEALTHCHECK(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	//log.Println("rule.Dbrule.Dbscnhealthcheck->", rule.Dbrule.Dbscnhealthcheck)
	msgdata := dbshtp.Dbscnhealthcheck.Contents

Looop:
	//按行分割
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 2 { //跳过前面2行
			continue
		}
		value = strings.TrimSpace(value) //去除头尾空格及空行
		// rda := regexp.MustCompile(`^Result: A`) // 判断是否检测到结果A
		rdv19 := regexp.MustCompile(`^Version:\s+19\.`)          // 判断是否检测版本
		rdv1124 := regexp.MustCompile(`^Version:\s+11\.2\.0\.4`) // 判断是否检测版本
		rdb := regexp.MustCompile(`^Result: B`)                  // 判断是否检测到结果B
		rdc := regexp.MustCompile(`^Result: C`)                  // 判断是否检测到结果C

		if rdv19.MatchString(value) || rdv1124.MatchString(value) {
			dbvstr := strings.Split(value, ":")
			log.Println("dbv->", (strings.TrimSpace(dbvstr[1])))
			break Looop
		}

		if rdc.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbscnhealthcheck.Alarm = "R"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
		if rdb.MatchString(value) { //匹配以"空格+数字"结尾的行
			dbshtp.Dbscnhealthcheck.Alarm = "B"
			///log.Printf("!!Matched!! value [%v]", value)
			break Looop
		}
	}
	// log.Printf("Dbscnhealthcheck.Alarm->%s", dbshtp.Dbscnhealthcheck.Alarm)
}

func Ana_DBDBAPRIV(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := dbshtp.Dbdbapriv.Contents
	value := strings.TrimSpace(msgdata)
	dbshtp.Dbdbapriv.Alarm = "G"
	if value == "" || strings.Contains(value, rule.Dbrule.Dbdbapriv.ResultG) { //空 或匹配到"no rows selected" ,则去除告警标记后退出
		dbshtp.Dbdbapriv.Alarm = ""
		return
	}
}

func Ana_DBSYSDBA(rule *utils.RuleInfo, dbshtp *structs.DbSht) {
	msgdata := dbshtp.Dbsysdba.Contents
	value := strings.TrimSpace(msgdata)
	dbshtp.Dbsysdba.Alarm = "B"
	if value == "" || strings.Contains(value, rule.Dbrule.Dbsysdba.ResultB) { //空 或匹配到"no rows selected" ,则去除告警标记后退出
		dbshtp.Dbsysdba.Alarm = ""
		return
	}
}
