package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// os_analyzer.go 包含操作系统指标的分析函数，检查 OS 参数、资源使用率等

// Ana_Osparameter 分析操作系统参数
func Ana_Osparameter(rule *utils.RuleInfo, osshtp *structs.OsSht, infstp *structs.InfoSht, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(infstp.Os)
	msgdata := osshtp.Osparameter.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Osparameter.Nm,
		Title:    rule.Osrule.Osparameter.Title,
		Desc:     rule.Osrule.Osparameter.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		rnm1 := regexp.MustCompile(`nproc`)
		rnm2 := regexp.MustCompile(`nofile`)
		rnm3 := regexp.MustCompile(`randomize_va_space`)
		rnm4 := regexp.MustCompile(`panic_on_oops`)
		rnm5 := regexp.MustCompile(`min_free_kbytes`)

		if strings.Contains(oS, "LINUX") {
			if rnm1.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_nproc_ne {
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s nproc 参数 %d 小于 %d，建议调整", infstp.HostName, n, rule.Osrule.Osparameter.L_nproc_ne))
					break Looop
				}
			}
			if rnm2.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_nofile_ne {
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s nofile 参数 %d 小于 %d，建议调整", infstp.HostName, n, rule.Osrule.Osparameter.L_nofile_ne))
					break Looop
				}
			}
			if rnm3.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.L_randomize_va_space {
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s randomi23_va_space 参数 %d 不等于 %d，建议调整", infstp.HostName, n, rule.Osrule.Osparameter.L_randomize_va_space))
					break Looop
				}
			}
			if rnm4.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.L_panic_on_oops {
					osshtp.Osparameter.Alarm = "G"
					entry.Minor = append(entry.Minor, fmt.Sprintf("主机 %s panic_on_oops 参数 %d 不等于 %d，建议调整", infstp.HostName, n, rule.Osrule.Osparameter.L_panic_on_oops))
					break Looop
				}
			}
			if rnm5.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.L_min_free_kbytes {
					osshtp.Osparameter.Alarm = "G"
					entry.Minor = append(entry.Minor, fmt.Sprintf("主机 %s min_free_kbytes 参数 %d 小于 %d，建议调整", infstp.HostName, n, rule.Osrule.Osparameter.L_min_free_kbytes))
					break Looop
				}
			}
		}
		if strings.Contains(oS, "SOLARIS") {
			if strings.Contains(value, "disable_ism_large_pages") {
				msg := strings.Split(value, "=")
				if !utils.Contain(msg[len(msg)-1], rule.Osrule.Osparameter.S_disable_ism_large_pages) {
					osshtp.Osparameter.Alarm = "B"
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s disable_ism_large_pages 参数 %s 不符合预期值 %v，建议调整", infstp.HostName, msg[len(msg)-1], rule.Osrule.Osparameter.S_disable_ism_large_pages))
					break Looop
				}
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ulimit 分析系统 ulimit 设置
func Ana_Ulimit(rule *utils.RuleInfo, osshtp *structs.OsSht, infstp *structs.InfoSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Ulimit.Contents
	oS := strings.ToUpper(infstp.Os)
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Ulimit.Nm,
		Title:    rule.Osrule.Ulimit.Title,
		Desc:     rule.Osrule.Ulimit.Desc,
	}
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
		if rnm1.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Open_files_ne {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s open files 限制 %d 小于 %d，建议调整", infstp.HostName, n, rule.Osrule.Ulimit.Open_files_ne))
				break Looop
			}
		}
		if rnm2.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Max_user_rocesses_ne {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("主机 %s max user processes 限制 %d 小于 %d，建议调整", infstp.HostName, n, rule.Osrule.Ulimit.Max_user_rocesses_ne))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Filesystem 分析文件系统使用率
func Ana_Filesystem(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Filesystem.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Filesystem.Nm,
		Title:    rule.Osrule.Filesystem.Title,
		Desc:     rule.Osrule.Filesystem.Desc,
	}
	// 转换阈值字符串为整数
	disk_ge1, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge1, "%"))
	disk_ge2, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge2, "%"))

	r := regexp.MustCompile(`\d+%`)
	matchs := r.FindAllString(msgdata, -1)
Looop:
	for _, p := range matchs {
		percent, _ := strconv.Atoi(strings.TrimSuffix(p, "%"))
		if percent >= disk_ge1 {
			osshtp.Filesystem.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("文件系统使用率 %d%% 超过 %d%%，建议清理或扩容", percent, disk_ge1))
			if percent >= disk_ge2 {
				osshtp.Filesystem.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("文件系统使用率 %d%% 超过 %d%%，需立即清理或扩容", percent, disk_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Inodeusage 分析索引节点使用率
func Ana_Inodeusage(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Inodeusage.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Inodeusage.Nm,
		Title:    rule.Osrule.Inodeusage.Title,
		Desc:     rule.Osrule.Inodeusage.Desc,
	}
	r := regexp.MustCompile(`\d+%`)
	matchs := r.FindAllString(msgdata, -1)
	// 转换阈值字符串为整数
	inode_ge1, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Inodeusage.Inode_ge1, "%"))
	inode_ge2, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Inodeusage.Inode_ge2, "%"))
Looop:
	for _, p := range matchs {
		percent, _ := strconv.Atoi(strings.TrimSuffix(p, "%"))
		if percent >= inode_ge1 {
			osshtp.Inodeusage.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("索引节点使用率 %d%% 超过 %d%%，建议清理", percent, inode_ge1))
			if percent >= inode_ge2 {
				osshtp.Inodeusage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("索引节点使用率 %d%% 超过 %d%%，需立即清理", percent, inode_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Cpustat 分析 CPU 使用情况
func Ana_Cpustat(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Cpustat.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Cpustat.Nm,
		Title:    rule.Osrule.Cpustat.Title,
		Desc:     rule.Osrule.Cpustat.Desc,
	}
Looop:
	for index, value := range strings.Split(msgdata, "\n") {
		if index < 4 {
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		if rd.MatchString(value) {
			msgs := strings.Fields(value)
			data := utils.String2Int(msgs)
			switch {
			case data[6] >= rule.Osrule.Cpustat.Swap_ge2:
				osshtp.Cpustat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("内存换页 si 值 %d 超过 %d，需优化内存使用", data[6], rule.Osrule.Cpustat.Swap_ge2))
				break Looop
			case data[14] < rule.Osrule.Cpustat.Idle_le2:
				osshtp.Cpustat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("CPU 空闲率 %d%% 小于 %d%%，负载过高", data[14], rule.Osrule.Cpustat.Idle_le2))
				break Looop
			case data[6] >= rule.Osrule.Cpustat.Swap_ge1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("内存换页 si 值 %d 超过 %d，建议优化内存使用", data[6], rule.Osrule.Cpustat.Swap_ge1))
			case data[14] < rule.Osrule.Cpustat.Idle_le1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("CPU 空闲率 %d%% 小于 %d%%，建议关注负载", data[14], rule.Osrule.Cpustat.Idle_le1))
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Memstat 分析内存使用情况
func Ana_Memstat(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Memstat.Contents
	if msgdata == "" {
		log.Printf("rule.Osrule.Memstat--->[%v] No data found!!! ", "msgdata")
		return
	}
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Memstat.Nm,
		Title:    rule.Osrule.Memstat.Title,
		Desc:     rule.Osrule.Memstat.Desc,
	}
	re := regexp.MustCompile(`\d+\s+`)
	vals := re.FindAllString(msgdata, -1)
	matchs, _ := strconv.Atoi(strings.ReplaceAll(vals[5], "\n", ""))
	if matchs < rule.Osrule.Memstat.Available_le1 {
		osshtp.Memstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("可用内存 %d MB 小于 %d MB，建议增加内存或优化使用", matchs, rule.Osrule.Memstat.Available_le1))
		if matchs < rule.Osrule.Memstat.Available_le2 {
			osshtp.Memstat.Alarm = "R"
			entry.Severe = append(entry.Severe, fmt.Sprintf("可用内存 %d MB 小于 %d MB，需立即增加内存或优化使用", matchs, rule.Osrule.Memstat.Available_le2))
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Iostat 分析磁盘 IO 性能
func Ana_Iostat(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Iostat.Contents
	if msgdata == "" {
		log.Printf("rule.Osrule.Iostat--->[%v] No data found!!! ", "Iostat")
		return
	}
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Iostat.Nm,
		Title:    rule.Osrule.Iostat.Title,
		Desc:     rule.Osrule.Iostat.Desc,
	}
Looop:
	for index, row := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		re := regexp.MustCompile(`^Average.*\d+$`)
		if re.MatchString(row) {
			msgs := strings.Fields(row)
			data, err := strconv.ParseFloat(msgs[len(msgs)-1], 64)
			if err != nil {
				log.Fatal(err)
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge1 {
				osshtp.Iostat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("磁盘 %s 使用率 %.2f%% 超过 %.0f%%，建议优化 IO 负载", msgs[1], data, rule.Osrule.Iostat.Diskutil_ge1))
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge2 {
				osshtp.Iostat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("磁盘 %s 使用率 %.2f%% 超过 %f%%，需立即优化 IO 负载", msgs[1], data, rule.Osrule.Iostat.Diskutil_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Thpstat 分析透明大页使用情况
func Ana_Thpstat(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Thpstat.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Thpstat.Nm,
		Title:    rule.Osrule.Thpstat.Title,
		Desc:     rule.Osrule.Thpstat.Desc,
	}
	re := regexp.MustCompile(`\d+`)
	matchs := re.FindString(msgdata)
	data, _ := strconv.Atoi(matchs)
	if data > rule.Osrule.Thpstat.Anpages_gt {
		osshtp.Thpstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("透明大页 AnonHugePages %d kB 超过 %d kB，建议关闭透明大页", data, rule.Osrule.Thpstat.Anpages_gt))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Numa 分析 NUMA 配置
func Ana_Numa(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Numa.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Numa.Nm,
		Title:    rule.Osrule.Numa.Title,
		Desc:     rule.Osrule.Numa.Desc,
	}
Looop:
	for index, row := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		re1 := regexp.MustCompile(rule.Osrule.Numa.Flg1)
		if re1.MatchString(row) {
			osshtp.Numa.Alarm = ""
			break Looop
		}
		re2 := regexp.MustCompile(rule.Osrule.Numa.Flg2)
		if re2.MatchString(row) {
			osshtp.Numa.Alarm = ""
			break Looop
		}
		osshtp.Numa.Alarm = "B"
		entry.Moderate = append(entry.Moderate, "主机 NUMA 未关闭，建议禁用 NUMA 以提升数据库稳定性")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ntp 分析 NTP 时钟同步配置
func Ana_Ntp(rule *utils.RuleInfo, osshtp *structs.OsSht, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Ntp.Contents
	str := strings.Replace(msgdata, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	entry := structs.SummaryEntry{
		Category: "主机系统分析",
		Nm:       rule.Osrule.Ntp.Nm,
		Title:    "NTP时钟同步检查",
		Desc:     "检查主机是否配置了 NTP 时钟同步",
	}
	if str == "" {
		osshtp.Ntp.Alarm = "B"
		entry.Moderate = append(entry.Moderate, "未配置 NTP 时钟同步，建议启用 NTP 服务")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
