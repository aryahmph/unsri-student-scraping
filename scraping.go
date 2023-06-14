package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	ListFacultiesURL      = "http://old.unsri.ac.id/?act=daftar_mahasiswa"
	ListFacultiesSelector = ".mainContent-news-element > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > ul:nth-child(1)"

	ListStudentByBatchURLPattern = "http://old.unsri.ac.id/?act=daftar_mahasiswa&fak_prodi=%s&angkatan=%d"
	ListStudentByBatchSelector   = ".mainContent-news-element > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(6) > tbody:nth-child(1)"
	StudentImageURLPattern       = "https://akademik.unsri.ac.id/images/foto_mhs/%d/%s.jpg"

	StudentDetailURL      = "http://old.unsri.ac.id/?act=detil_mahasiswa&mhs=%s-%s&akt=%d"
	StudentDetailSelector = ".mainContent-news-element > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1)"
)

var (
	TitleRegexPattern            = regexp.MustCompile(`^([A-Za-z\s]+) - ([A-Za-z0-9\s]+(?: \([A-Za-z\s]+\))?)( \([A-Za-z0-9\s]+\))$`)
	FacultyMajorCodeRegexPattern = regexp.MustCompile(`fak_prodi=(\d+-\d+-\d+)`)
)

func listFaculties() Faculty {
	// Send an HTTP GET request to a URL
	response, err := http.Get(ListFacultiesURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Load the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	// Extract information from the document using selectors
	list := make(Faculty)
	document.Find(ListFacultiesSelector).Children().Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("b").Text()
		if strings.ContainsRune(title, '(') {
			/*
				Matching example:
				"Fakultas Ekonomi - Ekonomi Pembangunan (S1 Kampus Indralaya)"

				Output, []string:
				0: Fakultas Ekonomi - Ekonomi Pembangunan (S1 Kampus Indralaya)
				1: Fakultas Ekonomi
				2: Ekonomi Pembangunan
				3: (S1 Kampus Indralaya)
			*/
			matches := TitleRegexPattern.FindStringSubmatch(title)
			if len(matches) < 2 {
				log.Println("Anomalous title, unsuccessful match :", title)
				return
			}

			faculty := TextCleaning(matches[1])
			major := TextCleaning(matches[2])
			detail := ""
			if len(matches) == 4 {
				detail = TextCleaning(matches[3])
			}

			// Faculty major code, example: https://old.unsri.ac.id/?act=daftar_mahasiswa&fak_prodi=1-10001-12&angkatan=2022
			// Code: 1-10001-12

			url, exist := selection.Find("a:nth-child(17)").Attr("href")
			if !exist {
				log.Println("Anomalous major, student batch not found :", title)
				return
			}

			matches = FacultyMajorCodeRegexPattern.FindStringSubmatch(url)
			if len(matches) != 2 {
				log.Println("Anomalous major batch url :", title)
				return
			}

			list[faculty] = append(list[faculty], Major{
				Name:   major,
				Detail: detail,
				Code:   matches[1],
			})
		}
	})

	return list
}

func listStudentsByBatch(majorName, facultyName, majorDetail, majorCode string, batchYear int) ([]Student, error) {
	// Send an HTTP GET request to a URL
	response, err := http.Get(fmt.Sprintf(ListStudentByBatchURLPattern, majorCode, batchYear))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Load the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	// Extract information from the document using selectors
	list := make([]Student, 0)
	document.Find(ListStudentByBatchSelector).Children().Each(func(i int, selection *goquery.Selection) {
		// Skip table head
		if i <= 1 {
			return
		}

		nim := ""
		selection.Children().Each(func(i int, selection *goquery.Selection) {
			if i == 3 {
				nim = TextCleaning(selection.Text())
				if nim == "" {
					err = errors.New("empty page")
					return
				}
			}
		})

		var student Student
		student, err = getStudent(nim, majorCode, batchYear)
		student.Faculty = facultyName
		student.Major = majorName
		student.MajorDetail = majorDetail
		list = append(list, student)
	})

	return list, err
}

func getStudent(nim, code string, batchYear int) (Student, error) {
	// Send an HTTP GET request to a URL
	response, err := http.Get(fmt.Sprintf(StudentDetailURL, nim, code, batchYear))
	if err != nil {
		return Student{}, err
	}
	defer response.Body.Close()

	// Load the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return Student{}, err
	}

	// Extract information from the document using selectors
	student := Student{BatchYear: batchYear}
	document.Find(StudentDetailSelector).Children().Each(func(i int, selection *goquery.Selection) {
		text := TextCleaning(selection.Find("td").Text())
		switch i {
		case 1:
			student.Nim = text
		case 2:
			student.Name = text
		case 4:
			student.Gender = text
		case 5:
			student.Religion = text
		case 6:
			student.SchoolOrigin = text
		case 7:
			student.Status = text
		}
	})
	student.ImageURL = fmt.Sprintf(StudentImageURLPattern, batchYear, student.Nim)
	return student, nil
}
