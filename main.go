package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// Change according to your capacity,
	// make sure to try whether it is affected by the Web Rate Limiter or not
	workerCount := 30

	jobs := make(chan Job, workerCount*15)
	results := make(chan Result, workerCount*15)

	wg := new(sync.WaitGroup)
	wg.Add(workerCount)

	// Start workers
	for i := 1; i <= workerCount; i++ {
		go func(workerId int) {
			defer wg.Done()
			scrapeWorker(workerId, jobs, results)
		}(i)
	}

	// Configure as you need
	minBatch, maxBatch := 2020, 2020
	
	// Start scraping process
	faculties := listFaculties()
	go func() {
		for faculty, majors := range faculties {
			// Restrictions for testing purposes, you can remove these
			if faculty != "Fakultas Ilmu Komputer" {
				continue
			}

			for _, major := range majors {
				currentBatch := minBatch
				for currentBatch <= maxBatch {
					jobs <- Job{
						FacultyName: faculty,
						MajorName:   major.Name,
						MajorDetail: major.Detail,
						MajorCode:   major.Code,
						BatchYear:   currentBatch,
					}
					currentBatch++
				}
			}
		}
		close(jobs)
	}()

	go func() {
		start := time.Now()
		wg.Wait()
		log.Println("Scraping process completed, time elapsed:", time.Since(start))
		close(results)
	}()

	csvWriter := NewCSVWriter("data")
	defer csvWriter.file.Close()
	defer csvWriter.writer.Flush()

	// Process the results
	for result := range results {
		if result.Err != nil {
			log.Println("Error scraping:", result.Err)
			continue
		}

		// Empty data, you can check manually on the web
		if len(result.Students) == 0 {
			log.Println("Empty scraping, major code:", result.MajorCode)
			continue
		}

		err := csvWriter.write(ListStudentToRows(result.Students))
		if err != nil {
			log.Println("Error writing csv:", err)
			continue
		}
		log.Println("Success:", len(result.Students))
	}
}
