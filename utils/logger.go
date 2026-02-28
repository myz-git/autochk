package utils

import (
    "log"
    "os"
    "strings"
)

type LogLevel int

const (
    LevelError LogLevel = iota
    LevelWarn
    LevelInfo
    LevelDebug
)

var currentLevel = LevelError

func init() {
    setLevelFromEnv()
}

func setLevelFromEnv() {
    lvl := strings.ToLower(strings.TrimSpace(os.Getenv("AUTOCHK_LOG_LEVEL")))
    if lvl == "" {
        currentLevel = LevelError
        return
    }
    switch lvl {
    case "debug":
        currentLevel = LevelDebug
    case "info":
        currentLevel = LevelInfo
    case "warn":
        currentLevel = LevelWarn
    default:
        currentLevel = LevelError
    }
}

func ShouldLog(level LogLevel) bool {
    return level <= currentLevel
}

func LogErrorf(format string, v ...any) {
    if ShouldLog(LevelError) {
        log.Printf(format, v...)
    }
}

func LogWarnf(format string, v ...any) {
    if ShouldLog(LevelWarn) {
        log.Printf(format, v...)
    }
}

func LogInfof(format string, v ...any) {
    if ShouldLog(LevelInfo) {
        log.Printf(format, v...)
    }
}

func LogDebugf(format string, v ...any) {
    if ShouldLog(LevelDebug) {
        log.Printf(format, v...)
    }
}


