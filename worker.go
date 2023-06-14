package main

type Job struct {
	MajorName   string
	FacultyName string
	MajorCode   string
	MajorDetail string
	BatchYear   int
}

type Result struct {
	Students  []Student
	MajorCode string
	Err       error
}

func scrapeWorker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		var err error
		var students []Student

		retries := 5
		for retries > 0 {
			students, err = listStudentsByBatch(job.MajorName, job.FacultyName, job.MajorDetail, job.MajorCode, job.BatchYear)
			if err == nil {
				break
			}
			retries--
		}

		results <- Result{
			Students:  students,
			MajorCode: job.MajorCode,
			Err:       err,
		}
	}
}
