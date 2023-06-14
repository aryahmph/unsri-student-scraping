package main

type (
	Faculty map[string][]Major

	Major struct {
		Name   string
		Detail string
		Code   string
	}

	Student struct {
		Name         string
		Nim          string
		BatchYear    int
		ImageURL     string
		Gender       string
		Religion     string
		SchoolOrigin string
		Faculty      string
		Major        string
		MajorDetail  string
		Status       string
	}
)
