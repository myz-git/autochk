package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Rule struct {
	ID    int
	Name  string
	Title string
	Type  string
	Level int
	Desc  string
}

func GetCHKDBRuleInfo() ([]Rule, error) {
	// 连接到 SQLite 数据库
	db, err := sql.Open("sqlite3", "../chk.db")
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 查询 rules 表中的所有数据
	rows, err := db.Query("SELECT id, nm, title, tp, level, desc FROM rules")
	if err != nil {
		return nil, fmt.Errorf("查询 rules 表失败: %v", err)
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.Title, &rule.Type, &rule.Level, &rule.Desc); err != nil {
			log.Printf("读取规则失败: %v", err)
			continue
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历 rules 表失败: %v", err)
	}

	return rules, nil
}
