package anadata

import (
	"autochk/structs"
	"autochk/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// os_analyzer.go 包含操作系统指标的分析函数，检查 OS 参数、资源使用率等

// Ana_Osparam_fs 分析文件系统类OS参数
func Ana_Osparam_fs(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(osshtp.Os.Contents)
	msgdata := osshtp.Osparam_fs.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Osparam_fs.Nm,
		Title:    rule.Osrule.Osparam_fs.Title,
		Desc:     rule.Osrule.Osparam_fs.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
	rm_file_max := regexp.MustCompile(`file-max`)
	rm_aio_max_nr := regexp.MustCompile(`aio-max-nr`)

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		if strings.Contains(oS, "LINUX") {
			if rm_file_max.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_fs.File_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,file_max参数当前值%d小于阈值%d,\n建议: 设置file_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_fs.File_max, rule.Osrule.Osparam_fs.File_max))
				}
			}
			if rm_aio_max_nr.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_fs.Aio_max_nr {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,aio_max_nr参数当前值%d小于阈值%d,\n建议: 设置aio_max_nr=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_fs.Aio_max_nr, rule.Osrule.Osparam_fs.Aio_max_nr))
				}
			}
		}
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Osparam_fs.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Osparam_fs.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Osparam_fs.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Osparam_ker 分析内核类OS参数
func Ana_Osparam_ker(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(osshtp.Os.Contents)
	msgdata := osshtp.Osparam_ker.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Osparam_ker.Nm,
		Title:    rule.Osrule.Osparam_ker.Title,
		Desc:     rule.Osrule.Osparam_ker.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
	rm_sem := regexp.MustCompile(`sem`)
	rm_panic_on_oops := regexp.MustCompile(`panic_on_oops`)
	// rm_randomize_va_space := regexp.MustCompile(`randomize_va_space`)

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		if strings.Contains(oS, "LINUX") {
			if rm_sem.MatchString(value) {
				msgs := strings.Fields(value)
				if len(msgs) >= 4 {
					sem2, _ := strconv.Atoi(msgs[len(msgs)-3])
					if len(rule.Osrule.Osparam_ker.Sem) > 1 {
						expectedSem2, _ := strconv.Atoi(rule.Osrule.Osparam_ker.Sem[1])
						if sem2 < expectedSem2 {
							entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,sem参数总信号量当前%d小于阈值%s,\n建议: 调整sem参数", osshtp.Hostname.Contents, sem2, rule.Osrule.Osparam_ker.Sem))
						}
					}
				}
			}
			if rm_panic_on_oops.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparam_ker.Panic_on_oops {
					entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,panic_on_oops参数当前值%d不等于期望值%d,\n建议: 设置panic_on_oops=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_ker.Panic_on_oops, rule.Osrule.Osparam_ker.Panic_on_oops))
				}
			}
			// if rm_randomize_va_space.MatchString(value) {
			// 	matchs := rd.FindString(value)
			// 	n, _ := strconv.Atoi(matchs)
			// 	if n != rule.Osrule.Osparam_ker.Randomize_va_space {
			// 		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,randomize_va_space参数当前值%d不等于期望值%d,\n建议: 设置randomize_va_space=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_ker.Randomize_va_space, rule.Osrule.Osparam_ker.Randomize_va_space))
			// 	}
			// }
		}
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Osparam_ker.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Osparam_ker.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Osparam_ker.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Osparam_net 分析网络类OS参数
func Ana_Osparam_net(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(osshtp.Os.Contents)
	msgdata := osshtp.Osparam_net.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Osparam_net.Nm,
		Title:    rule.Osrule.Osparam_net.Title,
		Desc:     rule.Osrule.Osparam_net.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
	rm_rp_filter_all := regexp.MustCompile(`rp_filter_all`)
	rm_rp_filter_default := regexp.MustCompile(`rp_filter_default`)
	rm_ip_local_port_range := regexp.MustCompile(`ip_local_port_range`)
	rm_ipfrag_high_thresh := regexp.MustCompile(`ipfrag_high_thresh`)
	rm_ipfrag_low_thresh := regexp.MustCompile(`ipfrag_low_thresh`)
	rm_rmem_default := regexp.MustCompile(`rmem_default`)
	rm_rmem_max := regexp.MustCompile(`rmem_max`)
	rm_wmem_default := regexp.MustCompile(`wmem_default`)
	rm_wmem_max := regexp.MustCompile(`wmem_max`)

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		if strings.Contains(oS, "LINUX") {
			if rm_rp_filter_all.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparam_net.Rp_filter_all {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,rp_filter_all参数当前值%d不等于期望值%d,\n建议: 设置rp_filter_all=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Rp_filter_all, rule.Osrule.Osparam_net.Rp_filter_all))
				}
			}
			if rm_rp_filter_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n != rule.Osrule.Osparam_net.Rp_filter_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,rp_filter_default参数当前值%d不等于期望值%d,\n建议: 设置rp_filter_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Rp_filter_default, rule.Osrule.Osparam_net.Rp_filter_default))
				}
			}
			if rm_ip_local_port_range.MatchString(value) {
				msgs := strings.Fields(value)
				if len(msgs) >= 2 && len(rule.Osrule.Osparam_net.Ip_local_port_range) >= 2 {
					start, _ := strconv.Atoi(msgs[len(msgs)-2])
					end, _ := strconv.Atoi(msgs[len(msgs)-1])
					expectedStart, _ := strconv.Atoi(rule.Osrule.Osparam_net.Ip_local_port_range[0])
					expectedEnd, _ := strconv.Atoi(rule.Osrule.Osparam_net.Ip_local_port_range[1])
					if start < expectedStart || end > expectedEnd {
						entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,ip_local_port_range参数当前值%d-%d不在期望范围%d-%d内,\n建议: 调整端口范围", osshtp.Hostname.Contents, start, end, expectedStart, expectedEnd))
					}
				}
			}
			if rm_ipfrag_high_thresh.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Ipfrag_high_thresh {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,ipfrag_high_thresh参数当前值%d小于阈值%d,\n建议: 设置ipfrag_high_thresh=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Ipfrag_high_thresh, rule.Osrule.Osparam_net.Ipfrag_high_thresh))
				}
			}
			if rm_ipfrag_low_thresh.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Ipfrag_low_thresh {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,ipfrag_low_thresh参数当前值%d小于阈值%d,\n建议: 设置ipfrag_low_thresh=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Ipfrag_low_thresh, rule.Osrule.Osparam_net.Ipfrag_low_thresh))
				}
			}
			if rm_rmem_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Rmem_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,rmem_default参数当前值%d小于阈值%d,\n建议: 设置rmem_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Rmem_default, rule.Osrule.Osparam_net.Rmem_default))
				}
			}
			if rm_rmem_max.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Rmem_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,rmem_max参数当前值%d小于阈值%d,\n建议: 设置rmem_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Rmem_max, rule.Osrule.Osparam_net.Rmem_max))
				}
			}
			if rm_wmem_default.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Wmem_default {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,wmem_default参数当前值%d小于阈值%d,\n建议: 设置wmem_default=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Wmem_default, rule.Osrule.Osparam_net.Wmem_default))
				}
			}
			if rm_wmem_max.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_net.Wmem_max {
					entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,wmem_max参数当前值%d小于阈值%d,\n建议: 设置wmem_max=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_net.Wmem_max, rule.Osrule.Osparam_net.Wmem_max))
				}
			}
		}
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Osparam_net.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Osparam_net.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Osparam_net.Alarm = "G"
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Osparam_vm 分析虚拟内存类OS参数
func Ana_Osparam_vm(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	oS := strings.ToUpper(osshtp.Os.Contents)
	msgdata := osshtp.Osparam_vm.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Osparam_vm.Nm,
		Title:    rule.Osrule.Osparam_vm.Title,
		Desc:     rule.Osrule.Osparam_vm.Desc,
	}
	rd := regexp.MustCompile(`\d+$`)
	rm_swappiness := regexp.MustCompile(`swappiness`)
	rm_min_free_kbytes := regexp.MustCompile(`min_free_kbytes`)

	for index, value := range strings.Split(msgdata, "\n") {
		if index == 0 {
			continue
		}
		if strings.Contains(oS, "LINUX") {
			if rm_swappiness.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n > rule.Osrule.Osparam_vm.Swappiness {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,swappiness参数当前值%d大于阈值%d,\n建议: 设置swappiness=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_vm.Swappiness, rule.Osrule.Osparam_vm.Swappiness))
				}
			}
			if rm_min_free_kbytes.MatchString(value) {
				matchs := rd.FindString(value)
				n, _ := strconv.Atoi(matchs)
				if n < rule.Osrule.Osparam_vm.Min_free_kbytes {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,min_free_kbytes参数当前值%d小于阈值%d,\n建议: 设置min_free_kbytes=%d", osshtp.Hostname.Contents, n, rule.Osrule.Osparam_vm.Min_free_kbytes, rule.Osrule.Osparam_vm.Min_free_kbytes))
				}
			}
		}
		if strings.Contains(oS, "SOLARIS") {
			if strings.Contains(value, "disable_ism_large_pages") {
				msg := strings.Split(value, "=")
				if !utils.Contain(msg[len(msg)-1], rule.Osrule.Osparam_vm.Disable_ism_large_pages) {
					entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,disable_ism_large_pages参数当前值%s不符合期望值%v,\n建议: 调整", osshtp.Hostname.Contents, msg[len(msg)-1], rule.Osrule.Osparam_vm.Disable_ism_large_pages))
				}
			}
		}
	}

	// 根据找到的问题设置最高级别的告警
	if len(entry.Severe) > 0 {
		osshtp.Osparam_vm.Alarm = "R"
	} else if len(entry.Moderate) > 0 {
		osshtp.Osparam_vm.Alarm = "B"
	} else if len(entry.Minor) > 0 {
		osshtp.Osparam_vm.Alarm = "G"
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
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,open files限制当前值%d小于阈值%d,\n建议: 设置open files=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Open_files, rule.Osrule.Ulimit.Open_files))
			}
		}
		if rnm2.MatchString(value) {
			matchs := rd.FindString(value)
			n, _ := strconv.Atoi(matchs)
			if n < rule.Osrule.Ulimit.Max_user_rocesses {
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,max user processes限制当前值%d小于阈值%d,\n建议: 设置max user processes=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Max_user_rocesses, rule.Osrule.Ulimit.Max_user_rocesses))
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
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,memlock限制当前值%d不等于期望值%d,\n建议: 设置memlock=%d", osshtp.Hostname.Contents, n, rule.Osrule.Ulimit.Memlock, rule.Osrule.Ulimit.Memlock))
			} else if memTotalKB > 0 && n < int(float64(memTotalKB)*0.8) {
				// 如果memlock小于OS内存的80%，也判定为问题
				expectedMemlock := int(float64(memTotalKB) * 0.8)
				osshtp.Ulimit.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,memlock限制当前值%d小于OS内存80%%(%d),\n建议: 设置memlock=%d", osshtp.Hostname.Contents, n, expectedMemlock, expectedMemlock))
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

// isExcludedMountPoint 检查挂载点是否应该被排除
func isExcludedMountPoint(mountPoint string) bool {
	excludedPaths := []string{
		"/media", "/mnt", "/run", "/tmp", "/dev", "/proc", "/sys",
		"/boot/efi", "/sys/fs/cgroup", "/dev/shm", "/run/user",
	}

	for _, excluded := range excludedPaths {
		if strings.HasPrefix(mountPoint, excluded) {
			return true
		}
	}
	return false
}

// Ana_Filesystem 分析文件系统使用率
func Ana_Filesystem(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	utils.LogDebugf("分析节点 %s 的文件系统使用率", osshtp.NodeID)
	msgdata := osshtp.Filesystem.Contents
	utils.LogDebugf("文件系统内容: %s", msgdata)
	utils.LogDebugf("文件系统内容长度: %d", len(msgdata))

	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Filesystem.Nm,
		Title:    rule.Osrule.Filesystem.Title,
		Desc:     rule.Osrule.Filesystem.Desc,
	}
	// 转换阈值字符串为整数
	disk_ge1, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge1, "%"))
	disk_ge2, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Filesystem.Disk_ge2, "%"))
	utils.LogDebugf("文件系统阈值: disk_ge1=%d%%, disk_ge2=%d%%", disk_ge1, disk_ge2)

	// 解析文件系统数据，过滤掉不需要检查的挂载点
	lines := strings.Split(msgdata, "\n")
	var filteredData strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "文件系统") || strings.Contains(line, "Filesystem") {
			// 保留标题行
			filteredData.WriteString(line + "\n")
			continue
		}

		// 提取挂载点（最后一列）
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			mountPoint := fields[len(fields)-1]
			if !isExcludedMountPoint(mountPoint) {
				filteredData.WriteString(line + "\n")
				utils.LogDebugf("保留文件系统行: %s (挂载点: %s)", line, mountPoint)
			} else {
				utils.LogDebugf("排除文件系统行: %s (挂载点: %s)", line, mountPoint)
			}
		} else {
			// 如果格式异常，保留该行
			filteredData.WriteString(line + "\n")
		}
	}

	filteredMsgData := filteredData.String()
	utils.LogDebugf("过滤后的文件系统内容: %s", filteredMsgData)

	r := regexp.MustCompile(`\d+%`)
	matchs := r.FindAllString(filteredMsgData, -1)
	utils.LogDebugf("找到的百分比匹配: %v", matchs)
Looop:
	for _, p := range matchs {
		percent, _ := strconv.Atoi(strings.TrimSuffix(p, "%"))
		utils.LogDebugf("处理百分比: %s -> %d%%", p, percent)
		if percent >= disk_ge1 {
			utils.LogDebugf("触发蓝色告警: %d%% >= %d%%", percent, disk_ge1)
			osshtp.Filesystem.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,文件系统使用率当前值%d%%超过阈值%d%%,\n建议: 及时清理或扩容", osshtp.Hostname.Contents, percent, disk_ge1))
			if percent >= disk_ge2 {
				utils.LogDebugf("触发红色告警: %d%% >= %d%%", percent, disk_ge2)
				osshtp.Filesystem.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,文件系统使用率当前值%d%%超过严重阈值%d%%,\n建议: 需尽快清理或扩容", osshtp.Hostname.Contents, percent, disk_ge2))
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

	// 解析索引节点数据，过滤掉不需要检查的挂载点
	lines := strings.Split(msgdata, "\n")
	var filteredData strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "文件系统") || strings.Contains(line, "Filesystem") || strings.Contains(line, "Inode") {
			// 保留标题行
			filteredData.WriteString(line + "\n")
			continue
		}

		// 提取挂载点（最后一列）
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			mountPoint := fields[len(fields)-1]
			if !isExcludedMountPoint(mountPoint) {
				filteredData.WriteString(line + "\n")
				utils.LogDebugf("保留索引节点行: %s (挂载点: %s)", line, mountPoint)
			} else {
				utils.LogDebugf("排除索引节点行: %s (挂载点: %s)", line, mountPoint)
			}
		} else {
			// 如果格式异常，保留该行
			filteredData.WriteString(line + "\n")
		}
	}

	filteredMsgData := filteredData.String()
	utils.LogDebugf("过滤后的索引节点内容: %s", filteredMsgData)

	r := regexp.MustCompile(`\d+%`)
	matchs := r.FindAllString(filteredMsgData, -1)
	// 转换阈值字符串为整数
	inode_ge1, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Inodeusage.Inode_ge1, "%"))
	inode_ge2, _ := strconv.Atoi(strings.TrimSuffix(rule.Osrule.Inodeusage.Inode_ge2, "%"))
Looop:
	for _, p := range matchs {
		percent, _ := strconv.Atoi(strings.TrimSuffix(p, "%"))
		if percent >= inode_ge1 {
			osshtp.Inodeusage.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,索引节点使用率当前值%d%%超过阈值%d%%,\n建议: 及时清理", osshtp.Hostname.Contents, percent, inode_ge1))
			if percent >= inode_ge2 {
				osshtp.Inodeusage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,索引节点使用率当前值%d%%超过严重阈值%d%%,\n建议: 需尽快清理", osshtp.Hostname.Contents, percent, inode_ge2))
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
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,内存换页si值当前%d超过阈值%d,\n建议: 优化内存使用或增加内存", osshtp.Hostname.Contents, data[6], rule.Osrule.Cpustat.Swap_ge2))
				break Looop
			case data[14] < rule.Osrule.Cpustat.Idle_le2:
				osshtp.Cpustat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,CPU空闲率当前%d%%小于阈值%d%%,\n建议: 优化进程或增加CPU资源", osshtp.Hostname.Contents, data[14], rule.Osrule.Cpustat.Idle_le2))
				break Looop
			case data[6] >= rule.Osrule.Cpustat.Swap_ge1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,内存换页 si 值 %d 超过 %d,\n建议: 优化内存使用", osshtp.Hostname.Contents, data[6], rule.Osrule.Cpustat.Swap_ge1))
			case data[14] < rule.Osrule.Cpustat.Idle_le1:
				osshtp.Cpustat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,CPU 空闲率 %d%% 小于 %d%%,\n建议: 关注负载", osshtp.Hostname.Contents, data[14], rule.Osrule.Cpustat.Idle_le1))
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
		utils.LogWarnf("rule.Osrule.Memstat--->[%v] No data found!!! ", "msgdata")
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
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,内存使用率当前%.1f%%超过90%%,\n建议: 需尽快优化内存使用或增加内存", osshtp.Hostname.Contents, memoryUsagePercent))
	} else if memoryUsagePercent > 80 {
		osshtp.Memstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,内存使用率当前%.1f%%超过80%%,\n建议: 增加内存或优化内存使用", osshtp.Hostname.Contents, memoryUsagePercent))
	}

	// 检查可用内存
	if availableMem < rule.Osrule.Memstat.Available_le2 {
		osshtp.Memstat.Alarm = "R"
		entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,可用内存当前%d MB小于严重阈值%d MB,\n建议: 需尽快增加内存或优化内存使用", osshtp.Hostname.Contents, availableMem, rule.Osrule.Memstat.Available_le2))
	} else if availableMem < rule.Osrule.Memstat.Available_le1 {
		osshtp.Memstat.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,可用内存当前%d MB小于阈值%d MB,\n建议: 增加内存或优化内存使用", osshtp.Hostname.Contents, availableMem, rule.Osrule.Memstat.Available_le1))
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
		utils.LogWarnf("rule.Osrule.Iostat--->[%v] No data found!!! ", "Iostat")
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
				utils.LogErrorf("解析磁盘使用率失败: %v", err)
				continue
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge1 {
				osshtp.Iostat.Alarm = "B"
				entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,磁盘%s使用率当前值%.2f%%超过阈值%.0f%%,\n建议: 优化IO负载", osshtp.Hostname.Contents, msgs[1], data, rule.Osrule.Iostat.Diskutil_ge1))
			}
			if data >= rule.Osrule.Iostat.Diskutil_ge2 {
				osshtp.Iostat.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,磁盘%s使用率当前值%.2f%%超过严重阈值%.0f%%,\n建议: 需尽快优化IO负载", osshtp.Hostname.Contents, msgs[1], data, rule.Osrule.Iostat.Diskutil_ge2))
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
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,透明大页当前配置为%s,\n建议: 关闭透明大页功能以提升系统稳定性", osshtp.Hostname.Contents, strings.TrimSpace(line)))
		}
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}

// Ana_Hugepage 分析大页配置情况
func Ana_Hugepage(rule *utils.RuleInfo, osshtp *structs.OsShts, summaryEntries *structs.SummaryEntries) {
	msgdata := osshtp.Hugepage.Contents
	entry := structs.SummaryEntry{
		Category: "主机系统",
		Nm:       rule.Osrule.Hugepage.Nm,
		Title:    rule.Osrule.Hugepage.Title,
		Desc:     rule.Osrule.Hugepage.Desc,
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
		osshtp.Hugepage.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,未配置大页内存,\n建议: 根据数据库SGA大小配置适当的大页内存", osshtp.Hostname.Contents))
	} else if totalHugePages > 0 {
		// 已配置大页，检查使用率
		if freeHugePages > 0 {
			usagePercent := float64(freeHugePages) / float64(totalHugePages) * 100
			if usagePercent > 70 {
				// 空闲大页超过70%，严重告警
				osshtp.Hugepage.Alarm = "R"
				entry.Severe = append(entry.Severe, fmt.Sprintf("问题: %s主机,大页配置总数为%d,空闲数为%d,使用率仅为%.1f%%,说明大页配置未生效,\n建议: 检查大页相关配置", osshtp.Hostname.Contents, totalHugePages, freeHugePages, 100-usagePercent))
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
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NUMA功能当前已开启,\n建议: 禁用NUMA以提升数据库稳定性", osshtp.Hostname.Contents))
	} else {
		// 其他情况，普通告警
		osshtp.Numa.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NUMA功能状态未知,\n建议: 检查并禁用NUMA以提升数据库稳定性", osshtp.Hostname.Contents))
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
		Title:    rule.Osrule.Ntp.Title,
		Desc:     rule.Osrule.Ntp.Desc,
	}
	if strings.Contains(str, "not running") {
		osshtp.Ntp.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,未配置时钟同步服务,\n建议: 启用NTP或chronyd时钟同步服务", osshtp.Hostname.Contents))
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
		Nm:       rule.Osrule.Selinux.Nm,
		Title:    rule.Osrule.Selinux.Title,
		Desc:     rule.Osrule.Selinux.Desc,
	}

	// 检查SELinux状态，如果为enabled，则为B普通告警
	if strings.Contains(msgdata, "SELinux status:") {
		if strings.Contains(msgdata, "enabled") {
			// SELinux已启用，普通告警
			osshtp.Selinux.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,SELinux安全模块当前已启用,\n建议: 禁用SELinux以提升数据库整体稳定性", osshtp.Hostname.Contents))
		} else if strings.Contains(msgdata, "disabled") {
			// SELinux已禁用，正常状态
			osshtp.Selinux.Alarm = ""
		} else {
			// 其他状态，普通告警
			osshtp.Selinux.Alarm = "B"
			entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,SELinux安全模块状态未知,\n建议: 禁用SELinux以提升数据库整体稳定性", osshtp.Hostname.Contents))
		}
	} else {
		// 未找到SELinux状态信息，普通告警
		osshtp.Selinux.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,未找到SELinux状态信息,\n建议: 禁用SELinux以提升数据库整体稳定性", osshtp.Hostname.Contents))
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
		Nm:       rule.Osrule.Firewall.Nm,
		Title:    rule.Osrule.Firewall.Title,
		Desc:     rule.Osrule.Firewall.Desc,
	}

	// 检查防火墙状态，如果为is running，则为B普通告警
	if strings.Contains(msgdata, "is running") {
		// 防火墙正在运行，普通告警
		osshtp.Firewall.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,防火墙服务当前正在运行,\n建议: 禁用防火墙服务以提升数据库整体稳定性", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "not running") {
		// 防火墙未运行，正常状态
		osshtp.Firewall.Alarm = ""
	} else {
		// 其他状态，普通告警
		osshtp.Firewall.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,防火墙服务状态未知,\n建议: 禁用防火墙服务以提升数据库整体稳定性", osshtp.Hostname.Contents))
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
		Nm:       rule.Osrule.Oslog.Nm,
		Title:    rule.Osrule.Oslog.Title,
		Desc:     rule.Osrule.Oslog.Desc,
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
		Nm:       rule.Osrule.Tmzone.Nm,
		Title:    rule.Osrule.Tmzone.Title,
		Desc:     rule.Osrule.Tmzone.Desc,
	}

	// 检查时区配置，如果不是东八区(+0800)，则为G轻微告警
	if strings.Contains(msgdata, "+0800") {
		// 东八区，正常状态
		osshtp.Tmzone.Alarm = ""
	} else {
		// 非东八区，轻微告警
		osshtp.Tmzone.Alarm = "G"
		entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,当前时区为%s,\n建议: 调整为东八区(+0800)以保持时间一致性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
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
		Nm:       rule.Osrule.Nsswitch.Nm,
		Title:    rule.Osrule.Nsswitch.Title,
		Desc:     rule.Osrule.Nsswitch.Desc,
	}

	// 检查NSSwitch配置，如果发现NIS configuration，则为B普通告警
	if strings.Contains(msgdata, "NIS configuration found") {
		// 发现NIS配置，普通告警
		osshtp.Nsswitch.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NSSwitch配置中发现NIS配置,\n建议: 移除NIS配置以提升系统安全性", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "No NIS configuration found") {
		// 未发现NIS配置，正常状态
		osshtp.Nsswitch.Alarm = ""
	} else {
		// 其他状态，普通告警
		osshtp.Nsswitch.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NSSwitch配置状态未知,\n建议: 检查并移除NIS配置以提升系统安全性", osshtp.Hostname.Contents))
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
		Nm:       rule.Osrule.Lo_mtu.Nm,
		Title:    rule.Osrule.Lo_mtu.Title,
		Desc:     rule.Osrule.Lo_mtu.Desc,
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
						entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,LO网卡MTU值当前为%d,超过16384,\n建议: 调整为16384或更小值以提升网络性能", osshtp.Hostname.Contents, mtu))
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
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,未找到LO网卡MTU配置信息,\n建议: 检查并设置合适的MTU值", osshtp.Hostname.Contents))
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
		Nm:       rule.Osrule.CPU_PERF_MODE.Nm,
		Title:    rule.Osrule.CPU_PERF_MODE.Title,
		Desc:     rule.Osrule.CPU_PERF_MODE.Desc,
	}

	// 检查是否为虚拟化平台，如果是则跳过检查
	if strings.Contains(msgdata, "Virtual machine detected - CPU performance mode check skipped") {
		osshtp.CPU_PERF_MODE.Alarm = ""
		// entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,虚拟化平台检测到,CPU性能模式检查已跳过", osshtp.Hostname.Contents))
		return
	}

	// 当为物理主机时,检查CPU是否在性能模式下运行
	if strings.Contains(msgdata, "CPU is in performance mode") {
		// CPU在性能模式下运行，正常
		osshtp.CPU_PERF_MODE.Alarm = ""
		// entry.Minor = append(entry.Minor, fmt.Sprintf("问题: %s主机,CPU当前在性能模式下运行,状态正常", osshtp.Hostname.Contents))
	} else if strings.Contains(msgdata, "CPU is NOT in performance mode") {
		// CPU不在性能模式下运行，普通告警
		osshtp.CPU_PERF_MODE.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,CPU未在性能模式下运行,\n建议: 检查并启用CPU性能模式以提升性能", osshtp.Hostname.Contents))
	} else {
		// 其他情况，设置为普通告警
		osshtp.CPU_PERF_MODE.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,CPU性能模式状态未知: %s,\n建议: 检查并确保CPU在性能模式下运行", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
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
		Nm:       rule.Osrule.NOZEROCONF.Nm,
		Title:    rule.Osrule.NOZEROCONF.Title,
		Desc:     rule.Osrule.NOZEROCONF.Desc,
	}

	// 检查NOZEROCONF配置
	if strings.Contains(msgdata, "NOZEROCONF=yes is configured") {
		// NOZEROCONF已正确配置，正常状态
		osshtp.NOZEROCONF.Alarm = ""
		// 正常状态不添加到任何告警列表中
	} else if strings.Contains(msgdata, "NOZEROCONF is not") {
		// NOZEROCONF未配置或未明确配置，普通告警
		osshtp.NOZEROCONF.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NOZEROCONF未配置或未明确配置,\n建议: 设置NOZEROCONF=yes以提升网络安全性", osshtp.Hostname.Contents))
	} else {
		// 其他情况，设置为普通告警
		osshtp.NOZEROCONF.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,NOZEROCONF配置状态未知: %s,\n建议: 检查并设置NOZEROCONF=yes以提升网络安全性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
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
		Nm:       rule.Osrule.RPM_PACKAGES.Nm,
		Title:    rule.Osrule.RPM_PACKAGES.Title,
		Desc:     rule.Osrule.RPM_PACKAGES.Desc,
	}

	// 检查是否有安装rpm包，如果msgdata不是"No item detected"则说明有安装存在风险漏洞的rpm包
	if strings.TrimSpace(msgdata) != "" && !strings.Contains(msgdata, "No item detected") {
		// 检测到安装了RPM包，普通告警
		osshtp.RPM_PACKAGES.Alarm = "B"
		entry.Moderate = append(entry.Moderate, fmt.Sprintf("问题: %s主机,检测到安装了存在安全风险的RPM包\n%s,\n建议: 移除存在安全隐患的RPM包以提升系统安全性", osshtp.Hostname.Contents, strings.TrimSpace(msgdata)))
	} else {
		// 未检测到RPM包，正常状态
		osshtp.RPM_PACKAGES.Alarm = ""
	}

	if len(entry.Severe) > 0 || len(entry.Moderate) > 0 || len(entry.Minor) > 0 {
		summaryEntries.Entries = append(summaryEntries.Entries, entry)
	}
}
