package anadata

import (
	"autochk/structs"
	"strings"
)

// format.go 包含格式化 InfoSht 结构体的函数

// Fmt_DbRole 格式化数据库角色信息
func Fmt_DbRole(infstp *structs.InfoSht) {
	msgdata := infstp.DbRole
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.DbRole = value
	}
}

// Fmt_LogMode 格式化日志模式信息
func Fmt_LogMode(infstp *structs.InfoSht) {
	msgdata := infstp.LogMode
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.LogMode = strings.TrimSpace(value)
	}
}

// Fmt_FlashBack 格式化闪回状态信息
func Fmt_FlashBack(infstp *structs.InfoSht) {
	msgdata := infstp.FlashBack
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.FlashBack = strings.TrimSpace(value)
	}
}

// Fmt_DbTotalsize 格式化数据库总大小信息
func Fmt_DbTotalsize(infstp *structs.InfoSht) {
	msgdata := infstp.DbTotalsize
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.DbTotalsize = strings.TrimSpace(value) + " GB"
	}
}

// Fmt_DbFilecount 格式化数据库文件数量信息
func Fmt_DbFilecount(infstp *structs.InfoSht) {
	msgdata := infstp.DbFilecount
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.DbFilecount = strings.TrimSpace(value)
	}
}

// Fmt_DbTblcount 格式化数据库表数量信息
func Fmt_DbTblcount(infstp *structs.InfoSht) {
	msgdata := infstp.DbTblcount
	for index, value := range strings.Split(msgdata, "\n") {
		if index != 2 {
			continue
		}
		infstp.DbTblcount = strings.TrimSpace(value)
	}
}
