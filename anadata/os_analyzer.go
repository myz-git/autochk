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
func Ana_Osparameter(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	log.Printf("分析节点 %s 的OS参数", osshtp.NodeID)
	log.Printf("OS参数内容长度: %d", len(osshtp.Osparameter.Contents))

	oS := strings.ToUpper(osshtp.Os.Contents)
	msgdata := osshtp.Osparameter.Contents
	log.Printf("操作系统类型: %s", oS)
	log.Printf("OS参数内容: %s", msgdata)

	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Osparameter.Nm,
		Title:    rule.Osrule.Osparameter.Title,
		Desc:     rule.Osrule.Osparameter.Desc,
	}

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		log.Printf("处理第%d行: %s", index, value)
		rd := regexp.MustCompile(`\d+$`)
		rm_file_max := regexp.MustCompile(`file-max`)
		rm_aio_max_nr := regexp.MustCompile(`aio-max-nr`)
		// rm_shmmni := regexp.MustCompile(`shmmni`)
		// rm_shmmax := regexp.MustCompile(`shmmax`)
		// rm_shmall := regexp.MustCompile(`shmall`)
		rm_sem := regexp.MustCompile(`sem`)
		rm_panic_on_oops := regexp.MustCompile(`panic_on_oops`)
		rm_randomize_va_space := regexp.MustCompile(`randomize_va_space`)
		// rm_numa_balancing := regexp.MustCompile(`numa_balancing`)
		// rm_tcp_keepalive_time := regexp.MustCompile(`tcp_keepalive_time`)
		// rm_dirty_ratio := regexp.MustCompile(`dirty_ratio`)
		// rm_dirty_background_ratio := regexp.MustCompile(`dirty_background_ratio`)
		// rm_dirty_expire_centisecs := regexp.MustCompile(`dirty_expire_centisecs`)
		// rm_dirty_writeback_centisecs := regexp.MustCompile(`dirty_writeback_centisecs`)
		rm_rp_filter_all := regexp.MustCompile(`rp_filter_all`)
		rm_rp_filter_default := regexp.MustCompile(`rp_filter_default`)
		rm_ipfrag_high_thresh := regexp.MustCompile(`ipfrag_high_thresh`)
		rm_ipfrag_low_thresh := regexp.MustCompile(`ipfrag_low_thresh`)
		rm_ip_local_port_range := regexp.MustCompile(`ip_local_port_range`)
		rm_rmem_default := regexp.MustCompile(`rmem_default`)
		rm_rmem_max := regexp.MustCompile(`rmem_max`)
		rm_wmem_default := regexp.MustCompile(`wmem_default`)
		rm_wmem_max := regexp.MustCompile(`wmem_max`)
		rm_swappiness := regexp.MustCompile(`swappiness`)
		rm_min_free_kbytes := regexp.MustCompile(`min_free_kbytes`)
		// rm_disable_ism_large_pages := regexp.MustCompile(`disable_ism_large_pages`)

		if strings.Contains(oS, "LINUX") {
			log.Printf("检测到Linux系统")
			if rm_file_max.MatchString(value) {
				// log.Printf("匹配到file_max: %s", value)
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				log.Printf("file_max值: %d, 阈值: %d", n, rule.Osrule.Osparameter.File_max)
				if n < rule.Osrule.Osparameter.File_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,file_max参数当前值%d小于阈值%d,建议设置file_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.File_max, rule.Osrule.Osparameter.File_max))
				}
			}
			if rm_aio_max_nr.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Aio_max_nr {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,aio_max_nr参数当前值%d小于阈值%d,建议设置aio_max_nr=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Aio_max_nr, rule.Osrule.Osparameter.Aio_max_nr))
				}
			}
			if rm_sem.MatchString(value) {
				// 对于sem参数，需要解析四个值，这里只检查第二个值
				msgs := strings.Fields(value)
				if len(msgs) >= 4 {
					sem2, _ := strconv.Atoi(msgs[len(msgs)-3])
					if len(rule.Osrule.Osparameter.Sem) > 1 {
						expectedSem2, _ := strconv.Atoi(rule.Osrule.Osparameter.Sem[1])
						if sem2 < expectedSem2 {
							entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,sem参数总信号量当前%d小于阈值%s,建议调整sem参数", osshtp.Hostname.Contents, sem2, rule.Osrule.Osparameter.Sem))
						}
					}
				}
			}
			if rm_panic_on_oops.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.Panic_on_oops {
					entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,panic_on_oops参数当前值%d不等于期望值%d,建议设置panic_on_oops=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Panic_on_oops, rule.Osrule.Osparameter.Panic_on_oops))
				}
			}
			if rm_randomize_va_space.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.Randomize_va_space {
					entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,randomize_va_space参数当前值%d不等于期望值%d,建议设置randomize_va_space=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Randomize_va_space, rule.Osrule.Osparameter.Randomize_va_space))
				}
			}
			if rm_rp_filter_all.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.Rp_filter_all {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,rp_filter_all参数当前值%d不等于期望值%d,建议设置rp_filter_all=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Rp_filter_all, rule.Osrule.Osparameter.Rp_filter_all))
				}
			}
			if rm_rp_filter_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparameter.Rp_filter_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,rp_filter_default参数当前值%d不等于期望值%d,建议设置rp_filter_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Rp_filter_default, rule.Osrule.Osparameter.Rp_filter_default))
				}
			}
			if rm_ip_local_port_range.MatchString(value) {
				// 对于ip_local_port_range，需要解析两个值
				msgs := strings.Fields(value)
				if len(msgs) >= 2 && len(rule.Osrule.Osparameter.Ip_local_port_range) >= 2 {
					start, _ := strconv.Atoi(msgs[len(msgs)-2])
					end, _ := strconv.Atoi(msgs[len(msgs)-1])
					expectedStart, _ := strconv.Atoi(rule.Osrule.Osparameter.Ip_local_port_range[0])
					expectedEnd, _ := strconv.Atoi(rule.Osrule.Osparameter.Ip_local_port_range[1])
					if start < expectedStart || end > expectedEnd {
						entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,ip_local_port_range参数当前值%d-%d不在期望范围%d-%d内,建议调整端口范围", osshtp.Hostname.Contents, start, end, expectedStart, expectedEnd))
					}
				}
			}
			if rm_ipfrag_high_thresh.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Ipfrag_high_thresh {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,ipfrag_high_thresh参数当前值%d小于阈值%d,建议设置ipfrag_high_thresh=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Ipfrag_high_thresh, rule.Osrule.Osparameter.Ipfrag_high_thresh))
				}
			}
			if rm_ipfrag_low_thresh.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Ipfrag_low_thresh {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,ipfrag_low_thresh参数当前值%d小于阈值%d,建议设置ipfrag_low_thresh=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Ipfrag_low_thresh, rule.Osrule.Osparameter.Ipfrag_low_thresh))
				}
			}
			if rm_rmem_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Rmem_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,rmem_default参数当前值%d小于阈值%d,建议设置rmem_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Rmem_default, rule.Osrule.Osparameter.Rmem_default))
				}
			}
			if rm_rmem_max.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Rmem_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,rmem_max参数当前值%d小于阈值%d,建议设置rmem_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Rmem_max, rule.Osrule.Osparameter.Rmem_max))
				}
			}
			if rm_wmem_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Wmem_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,wmem_default参数当前值%d小于阈值%d,建议设置wmem_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Wmem_default, rule.Osrule.Osparameter.Wmem_default))
				}
			}
			if rm_wmem_max.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Wmem_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,wmem_max参数当前值%d小于阈值%d,建议设置wmem_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Wmem_max, rule.Osrule.Osparameter.Wmem_max))
				}
			}
			if rm_swappiness.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n > rule.Osrule.Osparameter.Swappiness {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,swappiness参数当前值%d大于阈值%d,建议设置swappiness=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Swappiness, rule.Osrule.Osparameter.Swappiness))
				}
			}
			if rm_min_free_kbytes.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparameter.Min_free_kbytes {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,min_free_kbytes参数当前值%d小于阈值%d,建议设置min_free_kbytes=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparameter.Min_free_kbytes, rule.Osrule.Osparameter.Min_free_kbytes))
				}
			}
		}
		if strings.Contains(oS, "SOLARIS") {
			if strings.Contains(value, "disable_ism_large_pages") {
				msg := strings.Split(value, "=")
				if !utils.Contain(msg[len(msg)-1], rule.Osrule.Osparameter.Disable_ism_large_pages) {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,disable_ism_large_pages参数当前值%s不符合期望值%v,建议调整", osshtp.Hostname.Contents, msg[len(msg)-1], rule.Osrule.Osparameter.Disable_ism_large_pages))
				}
			}
		}
	}

	// 根据问题级别设置告警级别（优先级：严重 > 一般 > 轻微）
	if len(entry.Severe) > 0 {
		osshtp.Osparameter.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Osparameter.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Osparameter.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ulimit 分析系统 ulimit 设置
func Ana_Ulimit(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Ulimit.Contents
	oS := strings.ToUpper(osshtp.Os.Contents)
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Ulimit.Nm,
		Title:    rule.Osrule.Ulimit.Title,
		Desc:     rule.Osrule.Ulimit.Desc,
	}

	// 获取内存大小用于memlock检查
	memTotalKB := 0
	if memStr := strings.TrimSpace(osshtp.Memtotal.Contents); memStr != "" {
		// 从 "X GB" 格式中提取数字
		if memGB, err := strconv.Atoi(strings.Fields(memStr)[0]); err == nil {
			memTotalKB = memGB * 1024 * 1024 // 转换为KB
		}
	}

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		rd := regexp.MustCompile(`\d+$`)
		rnm1 := regexp.MustCompile(`open files`)
		rnm2 := regexp.MustCompile(`max user processes`)
		rnm3 := regexp.MustCompile(`max locked memory`)

		if strings.Contains(oS, "LINUX") {
			rnm1 = regexp.MustCompile(`open files`)
			rnm2 = regexp.MustCompile(`max user processes`)
			rnm3 = regexp.MustCompile(`max locked memory`)
		} else if strings.Contains(oS, "SOLARIS") {
			rnm1 = regexp.MustCompile(`nofile`)
			rnm2 = regexp.MustCompile(`nproc`)
			rnm3 = regexp.MustCompile(`memlock`)
		}

		if rnm1.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Open_files {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,open files限制当前值%d小于阈值%d,建议设置open files=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Open_files, rule.Osrule.Ulimit.Open_files))
			}
		}
		if rnm2.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Max_user_rocesses {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,max user processes限制当前值%d小于阈值%d,建议设置max user processes=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Max_user_rocesses, rule.Osrule.Ulimit.Max_user_rocesses))
			}
		}
		if rnm3.MatchString(value) {
			// 检查是否包含"unlimited"
			if strings.Contains(strings.ToLower(value), "unlimited") {
				// unlimited是正常值，不需要告警
				continue
			}

			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)

			// 检查memlock是否不等于-1或小于OS内存的80%
			if n != rule.Osrule.Ulimit.Memlock {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,memlock限制当前值%d不等于期望值%d,建议设置memlock=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Memlock, rule.Osrule.Ulimit.Memlock))
			} else if memTotalKB > 0 && n < int(float64(memTotalKB)*0.8) {
				// 如果memlock小于OS内存的80%，也判定为问题
				expectedMemlock := int(float64(memTotalKB) * 0.8)
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,memlock限制当前值%d小于OS内存80%%(%d),建议设置memlock=%d", osshtp.Hostname.Contents, n, expectedMemlock, expectedMemlock))
			}
		}
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Ulimit.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Ulimit.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Ulimit.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Filesystem 分析文件系统使用率
func Ana_Filesystem(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	log.Printf("分析节点 %s 的文件系统使用率", osshtp.NodeID)
	msgdata := osshtp.Filesystem.Contents
	log.Printf("文件系统内容: %s", msgdata)
	log.Printf("文件系统内容长度: %d", len(msgdata))

	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Filesystem.Nm,
		Title:    rule.Osrule.Filesystem.Title,
		Desc:     rule.Osrule.Filesystem.Desc,
	}
	// 转换阈值字符串为整数
	disk_ge1, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge1, "%"))
	disk_ge2, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge2, "%"))
	log.Printf("文件系统阈值: disk_ge1=%d%%, disk_ge2=%d%%", disk_ge1, disk_ge2)

	r := regexp.MustCompile(`\d+%`)
	matchs := r.FindAllString(msgdata, -1)
	log.Printf("找到的百分比匹配: %v", matchs)
Looop:
	for _, p := range matchs {
		percent, _ := strconv.Atoi(strings.TrimSuffix(p, "%"))
		log.Printf("处理百分比: %s -> %d%%", p, percent)
		if percent >= disk_ge1 {
			log.Printf("触发蓝色告警: %d%% >= %d%%", percent, disk_ge1)
			osshtp.Filesystem.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,文件系统使用率当前值%d%%超过阈值%d%%,建议及时清理或扩容", osshtp.Hostname.Contents, percent, disk_ge1))
			if percent >= disk_ge2 {
				log.Printf("触发红色告警: %d%% >= %d%%", percent, disk_ge2)
				osshtp.Filesystem.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,文件系统使用率当前值%d%%超过严重阈值%d%%,需立即清理或扩容", osshtp.Hostname.Contents, percent, disk_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Inodeusage 分析索引节点使用率
func Ana_Inodeusage(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Inodeusage.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
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
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,索引节点使用率当前值%d%%超过阈值%d%%,建议及时清理", osshtp.Hostname.Contents, percent, inode_ge1))
			if percent >= inode_ge2 {
				osshtp.Inodeusage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,索引节点使用率当前值%d%%超过严重阈值%d%%,需立即清理", osshtp.Hostname.Contents, percent, inode_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Cpustat 分析 CPU 使用情况
func Ana_Cpustat(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Cpustat.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
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
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,内存换页si值当前%d超过阈值%d,建议优化内存使用或增加内存", osshtp.Hostname.Contents, data[6], rule.Osrule.Cpustat.Swap_ge2))
				break Looop
			case data[14] < rule.Osrule.Cpustat.Idle_le2:
				osshtp.Cpustat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,CPU空闲率当前%d%%小于阈值%d%%,建议优化进程或增加CPU资源", osshtp.Hostname.Contents, data[14], rule.Osrule.Cpustat.Idle_le2))
				break Looop
			case data[6] >= rule.Osrule.Cpustat.Swap_ge1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,内存换页 si 值 %d 超过 %d，建议优化内存使用", osshtp.Hostname.Contents, data[6], rule.Osrule.Cpustat.Swap_ge1))
			case data[14] < rule.Osrule.Cpustat.Idle_le1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,CPU 空闲率 %d%% 小于 %d%%，建议关注负载", osshtp.Hostname.Contents, data[14], rule.Osrule.Cpustat.Idle_le1))
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Memstat 分析内存使用情况
func Ana_Memstat(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Memstat.Contents
	if msgdata == "" {
		log.Printf("rule.Osrule.Memstat--->[%v] No data found!!! ", "msgdata")
		return
	}
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Memstat.Nm,
		Title:    rule.Osrule.Memstat.Title,
		Desc:     rule.Osrule.Memstat.Desc,
	}

	// 解析内存数据
	var totalMem, usedMem, availableMem int

	// 按行解析数据
	lines := strings.Split(msgdata, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Mem:") {
			// 解析 Mem 行: "Mem: 7479 5932 260 7 1287 1171"
			fields := strings.Fields(line)
			if len(fields) >= 7 {
				totalMem, _ = strconv.Atoi(fields[1])
				usedMem, _ = strconv.Atoi(fields[2])
			}
		} else if strings.HasPrefix(line, "available=") {
			// 解析 available 行: "available=1171 MB"
			re := regexp.MustCompile(`available=(\d+)\s+MB`)
			if match := re.FindStringSubmatch(line); len(match) > 1 {
				availableMem, _ = strconv.Atoi(match[1])
			}
		}
	}

	// 计算内存使用率
	var memoryUsagePercent float64
	if totalMem > 0 {
		memoryUsagePercent = float64(usedMem) / float64(totalMem) * 100
	}

	// 检查内存使用率
	if memoryUsagePercent > 90 {
		osshtp.Memstat.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,内存使用率当前%.1f%%超过90%%,建议立即优化内存使用或增加内存", osshtp.Hostname.Contents, memoryUsagePercent))
	} else if memoryUsagePercent > 80 {
		osshtp.Memstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,内存使用率当前%.1f%%超过80%%,建议关注内存使用情况", osshtp.Hostname.Contents, memoryUsagePercent))
	}

	// 检查可用内存
	if availableMem < rule.Osrule.Memstat.Available_le2 {
		osshtp.Memstat.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,可用内存当前%d MB小于严重阈值%d MB,需立即增加内存或优化内存使用", osshtp.Hostname.Contents, availableMem, rule.Osrule.Memstat.Available_le2))
	} else if availableMem < rule.Osrule.Memstat.Available_le1 {
		osshtp.Memstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,可用内存当前%d MB小于阈值%d MB,建议增加内存或优化内存使用", osshtp.Hostname.Contents, availableMem, rule.Osrule.Memstat.Available_le1))
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Memstat.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Memstat.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Memstat.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Iostat 分析磁盘 IO 性能
func Ana_Iostat(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Iostat.Contents
	if msgdata == "" {
		log.Printf("rule.Osrule.Iostat--->[%v] No data found!!! ", "Iostat")
		return
	}
	entry := structs.SummaryEntry{
		Category: "主机系统",
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
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,磁盘%s使用率当前值%.2f%%超过阈值%.0f%%,建议优化IO负载", osshtp.Hostname.Contents, msgs[1], data, rule.Osrule.Iostat.Diskutil_ge1))
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge2 {
				osshtp.Iostat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,磁盘%s使用率当前值%.2f%%超过严重阈值%.0f%%,需立即优化IO负载", osshtp.Hostname.Contents, msgs[1], data, rule.Osrule.Iostat.Diskutil_ge2))
				break Looop
			}
		}
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Thpstat 分析透明大页使用情况
func Ana_Thpstat(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Thpstat.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Thpstat.Nm,
		Title:    rule.Osrule.Thpstat.Title,
		Desc:     rule.Osrule.Thpstat.Desc,
	}

	// 检查透明大页配置，当不是 [never] 时，则为B普通告警
	// 支持格式: "always madvise [never]" 或 "[always] madvise never"
	lines := strings.Split(msgdata, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否包含 [never]，如果包含则说明透明大页已关闭（正常配置）
		if strings.Contains(line, "[never]") {
			// 透明大页已关闭，正常配置，不需要告警
			continue
		}

		// 检查是否包含 always 或 madvise，如果包含且不是 [never]，则需要告警
		if strings.Contains(line, "always") || strings.Contains(line, "madvise") {
			osshtp.Thpstat.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,透明大页当前配置为%s,建议关闭透明大页功能以提升系统稳定性", osshtp.Hostname.Contents, strings.TrimSpace(line)))
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Hugpage 分析大页配置情况
func Ana_Hugpage(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Hugpage.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "HUGEPAGE",
		Title:    "大页配置检查",
		Desc:     "检查大页配置情况，如果配置了大页但使用率过低说明配置未生效",
	}

	var totalHugePages, freeHugePages int

	// 解析大页数据
	lines := strings.Split(msgdata, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "HugePages_Total:") {
			// 解析总大页数
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				totalHugePages, _ = strconv.Atoi(fields[1])
			}
		} else if strings.HasPrefix(line, "HugePages_Free:") {
			// 解析空闲大页数
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				freeHugePages, _ = strconv.Atoi(fields[1])
			}
		}
	}

	// 判断大页配置情况
	if totalHugePages == 0 {
		// 未启用大页，普通告警
		osshtp.Hugpage.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,未配置大页内存,建议根据数据库SGA大小配置适当的大页内存", osshtp.Hostname.Contents))
	} else if totalHugePages > 0 {
		// 已配置大页，检查使用率
		if freeHugePages > 0 {
			usagePercent := float64(freeHugePages) / float64(totalHugePages) * 100
			if usagePercent > 70 {
				// 空闲大页超过70%，严重告警
				osshtp.Hugpage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("%s主机,大页配置总数为%d,空闲数为%d,使用率仅为%.1f%%,说明大页配置未生效,建议检查数据库配置", osshtp.Hostname.Contents, totalHugePages, freeHugePages, 100-usagePercent))
			}
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Numa 分析 NUMA 配置
func Ana_Numa(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Numa.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Numa.Nm,
		Title:    rule.Osrule.Numa.Title,
		Desc:     rule.Osrule.Numa.Desc,
	}

	// 检查NUMA状态
	msgdata = strings.TrimSpace(msgdata)
	if strings.Contains(msgdata, "NUMA turned off") {
		// NUMA已关闭，正常状态
		osshtp.Numa.Alarm = ""
	} else if strings.Contains(msgdata, "NUMA turned on") {
		// NUMA已开启，普通告警
		osshtp.Numa.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NUMA功能当前已开启,建议禁用NUMA以提升数据库稳定性", osshtp.Hostname.Contents))
	} else {
		// 其他情况，普通告警
		osshtp.Numa.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NUMA功能状态未知,建议检查并禁用NUMA以提升数据库稳定性", osshtp.Hostname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Ntp 分析 NTP 时钟同步配置
func Ana_Ntp(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Ntp.Contents
	str := strings.Replace(msgdata, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Ntp.Nm,
		Title:    "时钟同步检查",
		Desc:     "检查主机NTP 或chronyd 时钟同步服务是否启用",
	}
	if strings.Contains(str, "not running") {
		osshtp.Ntp.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,未配置时钟同步服务,建议启用NTP或chronyd时钟同步服务", osshtp.Hostname.Contents))
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Selinux 分析 SELinux 状态
func Ana_Selinux(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Selinux.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "SELINUX",
		Title:    "SELinux状态检查",
		Desc:     "检查SELinux安全模块是否已禁用",
	}

	// 检查SELinux状态，如果为enabled，则为B普通告警
	if strings.Contains(msgdata, "SELinux status:") {
		if strings.Contains(msgdata, "enabled") {
			// SELinux已启用，普通告警
			osshtp.Selinux.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,SELinux安全模块当前已启用,建议禁用SELinux以提升数据库性能", osshtp.Hostname.Contents))
		} else if strings.Contains(msgdata, "disabled") {
			// SELinux已禁用，正常状态
			osshtp.Selinux.Alarm = ""
		} else {
			// 其他状态，普通告警
			osshtp.Selinux.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,SELinux安全模块状态未知,建议检查并禁用SELinux以提升数据库性能", osshtp.Hostname.Contents))
		}
	} else {
		// 未找到SELinux状态信息，普通告警
		osshtp.Selinux.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,未找到SELinux状态信息,建议检查并禁用SELinux以提升数据库性能", osshtp.Hostname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Firewall 分析防火墙状态
func Ana_Firewall(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Firewall.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "FIREWALL",
		Title:    "防火墙状态检查",
		Desc:     "检查防火墙服务是否已启用",
	}

	// 检查防火墙状态，如果为is running，则为B普通告警
	if strings.Contains(msgdata, "is running") {
		// 防火墙正在运行，普通告警
		osshtp.Firewall.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,防火墙服务当前正在运行,建议禁用防火墙服务以提升数据库性能", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "not running") {
		// 防火墙未运行，正常状态
		osshtp.Firewall.Alarm = ""
	} else {
		// 其他状态，普通告警
		osshtp.Firewall.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,防火墙服务状态未知,建议检查并禁用防火墙服务以提升数据库性能", osshtp.Hostname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Oslog 分析操作系统日志
func Ana_Oslog(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Oslog.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "OSLOG",
		Title:    "操作系统日志检查",
		Desc:     "检查操作系统日志状态",
	}
	// 实现日志分析逻辑
	if strings.TrimSpace(msgdata) != "" {
		osshtp.Oslog.Alarm = "G"
		entry.Minor = append(entry.Minor, "操作系统日志检查")
	}
	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Tmzone 分析时区设置
func Ana_Tmzone(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Tmzone.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "TMZONE",
		Title:    "时区配置检查",
		Desc:     "检查主机时区配置是否为东八区",
	}

	// 检查时区配置，如果不是东八区(+0800)，则为G轻微告警
	if strings.Contains(msgdata, "+0800") {
		// 东八区，正常状态
		osshtp.Tmzone.Alarm = ""
	} else {
		// 非东八区，轻微告警
		osshtp.Tmzone.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,当前时区为%s,建议调整为东八区(+0800)以保持时间一致性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Nsswitch 分析NSSwitch配置
func Ana_Nsswitch(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Nsswitch.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "NSSWITCH",
		Title:    "NSSwitch配置检查",
		Desc:     "检查NSSwitch配置是否包含NIS配置",
	}

	// 检查NSSwitch配置，如果发现NIS configuration，则为B普通告警
	if strings.Contains(msgdata, "NIS configuration found") {
		// 发现NIS配置，普通告警
		osshtp.Nsswitch.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NSSwitch配置中发现NIS配置,建议移除NIS配置以提升系统安全性", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "No NIS configuration found") {
		// 未发现NIS配置，正常状态
		osshtp.Nsswitch.Alarm = ""
	} else {
		// 其他状态，普通告警
		osshtp.Nsswitch.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NSSwitch配置状态未知,建议检查并移除NIS配置以提升系统安全性", osshtp.Hostname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Lo_mtu 分析LO_MTU配置
func Ana_Lo_mtu(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Lo_mtu.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "LO_MTU",
		Title:    "LO网卡MTU配置检查",
		Desc:     "检查LO网卡的MTU值是否超过16384",
	}

	// 检查LO网卡MTU值，如果大于16384，则为B普通告警
	lines := strings.Split(msgdata, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mtu ") {
			// 解析MTU值
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if mtu, err := strconv.Atoi(fields[1]); err == nil {
					if mtu > 16384 {
						// MTU值大于16384，普通告警
						osshtp.Lo_mtu.Alarm = "B"
						entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,LO网卡MTU值当前为%d,超过16384,建议调整为16384或更小值以提升网络性能", osshtp.Hostname.Contents, mtu))
					} else {
						// MTU值正常，不设置告警
						osshtp.Lo_mtu.Alarm = ""
					}
					break
				}
			}
		}
	}

	// 如果没有找到MTU配置或解析失败，设置为普通告警
	if osshtp.Lo_mtu.Alarm == "" && !strings.Contains(msgdata, "mtu ") {
		osshtp.Lo_mtu.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,未找到LO网卡MTU配置信息,建议检查并设置合适的MTU值", osshtp.Hostname.Contents))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Machine_platform 分析主机平台配置
func Ana_Machine_platform(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Machine_platform.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "MACHINE_PLATFORM",
		Title:    "主机平台类型检查",
		Desc:     "检查主机是物理机还是虚拟服务器",
	}

	// 该指标仅做记录，不做告警判断
	if strings.Contains(msgdata, "Virtual Machine") {
		// 虚拟服务器，记录信息但不告警
		osshtp.Machine_platform.Alarm = ""
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,当前运行在虚拟服务器环境中", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "Physical Server") {
		// 物理服务器，记录信息但不告警
		osshtp.Machine_platform.Alarm = ""
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,当前运行在物理服务器环境中", osshtp.Hostname.Contents))
	} else {
		// 其他情况，记录信息但不告警
		osshtp.Machine_platform.Alarm = ""
		entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,主机平台类型未知: %s", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_CPU_PERF_MODE 分析CPU性能模式配置
func Ana_CPU_PERF_MODE(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.CPU_PERF_MODE.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "CPU_PERF_MODE",
		Title:    "CPU性能模式检查",
		Desc:     "检查CPU是否在高性能模式下运行",
	}

	// 检查是否为虚拟化平台，如果是则跳过检查
	if strings.Contains(msgdata, "Virtual machine detected - CPU performance mode check skipped") {
		osshtp.CPU_PERF_MODE.Alarm = ""
		// entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,虚拟化平台检测到,CPU性能模式检查已跳过", osshtp.Hostname.Contents))
		return
	}

	// 当为物理主机时,检查CPU是否在性能模式下运行
	if strings.Contains(msgdata, "CPU is in performance mode") {
		// CPU在性能模式下运行，正常
		osshtp.CPU_PERF_MODE.Alarm = ""
		// entry.Minor = append(entry.Minor, fmt.Sprintf("%s主机,CPU当前在性能模式下运行,状态正常", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "CPU is NOT in performance mode") {
		// CPU不在性能模式下运行，普通告警
		osshtp.CPU_PERF_MODE.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,CPU未在性能模式下运行,建议检查并启用CPU性能模式以提升性能", osshtp.Hostname.Contents))
	} else {
		// 其他情况，设置为普通告警
		osshtp.CPU_PERF_MODE.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,CPU性能模式状态未知: %s,建议检查并确保CPU在性能模式下运行", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_NOZEROCONF 分析NOZEROCONF配置
func Ana_NOZEROCONF(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.NOZEROCONF.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "NOZEROCONF",
		Title:    "NOZEROCONF配置检查",
		Desc:     "检查NOZEROCONF是否已正确配置",
	}

	// 检查NOZEROCONF配置
	if strings.Contains(msgdata, "NOZEROCONF=yes is configured") {
		// NOZEROCONF已正确配置，正常状态
		osshtp.NOZEROCONF.Alarm = ""
		// 正常状态不添加到任何告警列表中
	} else if strings.Contains(msgdata, "NOZEROCONF is not") {
		// NOZEROCONF未配置或未明确配置，普通告警
		osshtp.NOZEROCONF.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NOZEROCONF未配置或未明确配置,建议设置NOZEROCONF=yes以提升网络安全性", osshtp.Hostname.Contents))
	} else {
		// 其他情况，设置为普通告警
		osshtp.NOZEROCONF.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,NOZEROCONF配置状态未知: %s,建议检查并设置NOZEROCONF=yes以提升网络安全性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_RPM_PACKAGES 分析RPM_PACKAGES配置
func Ana_RPM_PACKAGES(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.RPM_PACKAGES.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       "RPM_PACKAGES",
		Title:    "RPM包安全检查",
		Desc:     "检查是否安装了存在安全风险的RPM包",
	}

	// 检查是否有安装rpm包，如果msgdata不是"No item detected"则说明有安装存在风险漏洞的rpm包
	if strings.TrimSpace(msgdata) != "" && !strings.Contains(msgdata, "No item detected") {
		// 检测到安装了RPM包，普通告警
		osshtp.RPM_PACKAGES.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("%s主机,检测到安装了存在安全风险的RPM包: %s,建议及时更新或移除这些包以提升系统安全性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	} else {
		// 未检测到RPM包，正常状态
		osshtp.RPM_PACKAGES.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
