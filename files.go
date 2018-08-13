package main

import (
	"path/filepath"
	"strings"
)

type FileContent struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	ChangeI int    `json:"changeIndex"`
}

func filepathMatches(p string, matchSet []string) bool {
	ext := filepath.Ext(p)

	for _, s := range matchSet {
		if ext == s {
			return true
		}
		if strings.HasSuffix(p, s) {
			return true
		}
		if strings.Contains(p, string(filepath.Separator)+s) {
			return true
		}
	}
	return false
}
