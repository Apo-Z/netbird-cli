package main

import (
	"regexp"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func red(s string) string   { return colorRed + s + colorReset }
func green(s string) string { return colorGreen + s + colorReset }
func yellow(s string) string { return colorYellow + s + colorReset }
func cyan(s string) string  { return colorCyan + s + colorReset }
func bold(s string) string  { return colorBold + s + colorReset }

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func colorBool(positive bool, labelTrue, labelFalse string) string {
	if labelTrue == "" {
		labelTrue = "yes"
	}
	if labelFalse == "" {
		labelFalse = "no"
	}
	if positive {
		return green(labelTrue)
	}
	return red(labelFalse)
}

func colorStatus(status string) string {
	lower := strings.ToLower(status)
	switch {
	case lower == "active" || lower == "approved" || lower == "connected":
		return green(status)
	case lower == "disabled" || lower == "inactive" || lower == "disconnected" || lower == "error":
		return red(status)
	case lower == "pending" || lower == "in progress":
		return yellow(status)
	}
	return status
}

func colorState(state string) string {
	lower := strings.ToLower(state)
	switch {
	case lower == "up" || lower == "running" || lower == "active" || lower == "ready":
		return green(state)
	case lower == "down" || lower == "stopped" || lower == "error" || lower == "failed":
		return red(state)
	case lower == "provisioning":
		return yellow(state)
	}
	return state
}

func colorDecision(decision string) string {
	switch strings.ToLower(decision) {
	case "allow":
		return green(decision)
	case "deny":
		return red(decision)
	}
	return decision
}
