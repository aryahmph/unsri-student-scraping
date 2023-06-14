package main

import (
	"regexp"
	"strconv"
)

var (
	TrimSpaceRegex = regexp.MustCompile(`\t|\n`)
)

func TextCleaning(arg string) string {
	return TrimSpaceRegex.ReplaceAllString(arg, "")
}

func ListStudentToRows(students []Student) [][]string {
	var rows [][]string
	for _, student := range students {
		row := []string{
			student.Name,
			student.Nim,
			strconv.Itoa(student.BatchYear),
			student.ImageURL,
			student.Gender,
			student.Religion,
			student.SchoolOrigin,
			student.Faculty,
			student.Major,
			student.MajorDetail,
			student.Status,
		}
		rows = append(rows, row)
	}
	return rows
}
