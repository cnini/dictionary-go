package middleware

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Log struct {
	Authorization string
	Timestamp     string
	Method        string
	Url           string
	Duration      time.Duration
}

var logMutex sync.Mutex

func NewLog(
	next http.HandlerFunc,
	w http.ResponseWriter,
	r *http.Request,
	errors chan<- error,
	authorizationStatus string,
) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Start time of the capture
	startTime := time.Now()

	// Call the next handler of the chain :
	// captureLog() captures information about the request,
	// then it passes it to the next() to handle it
	// until another information will be captured by the function
	// and be passed to it
	next.ServeHTTP(w, r)

	// End time of the capture
	endTime := time.Now()

	write(
		Log{
			authorizationStatus,
			startTime.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.String(),
			endTime.Sub(startTime),
		},
		errors,
	)
}

func write(log Log, errors chan<- error) {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		errors <- err
	}

	defer logFile.Close()

	_, err = logFile.WriteString(
		fmt.Sprintf(
			"[%s - %s] %s %s %v\n",
			log.Authorization,
			log.Timestamp,
			log.Method,
			log.Url,
			log.Duration,
		),
	)

	if err != nil {
		errors <- err
	}
}
