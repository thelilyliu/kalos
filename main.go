package main

import (
	// "encoding/binary"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/husobee/vestigo"
)

type Page struct {
	Viewer string
}

type Entry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	errorStatusCode = 555
	serverName      = "GWS"
	user            = "Lily"
)

func main() {
	router := vestigo.NewRouter()

	// set up router global CORS policy
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*"},
		AllowCredentials: false,
		MaxAge:           3600 * time.Second,
	})

	fileServerKalos := http.FileServer(http.Dir("kalos"))
	router.Get("/kalos/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", serverName)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/kalos")
		fileServerKalos.ServeHTTP(w, r)
	})

	// poll
	router.Get("/loadPolls", loadPolls)
	router.Get("/loadPoll/:pollID", loadPoll)
	router.Get("/insertPoll", insertPoll)
	router.Post("/updatePoll/:pollID", updatePoll)
	router.Delete("/deletePoll/:pollID", deletePoll)

	// other
	router.Post("/submitCode", submitCode)
	router.Post("/submitResponse/:pollID", submitResponse)
	router.Get("/getResults/:pollID", getResults)
	router.Get("/generateResponses/:pollID", generateResponses)

	// lambda
	router.Post("/createPoll", createPoll)
	router.Post("/viewResults", viewResults)

	// pages
	router.Get("/edit", viewAdmin)
	router.Get("/manage", viewAdmin)
	router.Get("/poll", viewClient)
	router.Get("/", viewClient)

	log.Println("Listening...")
	if err := http.ListenAndServe(":2323", router); err != nil {
		log.Println(err)
	}
}

/*
  ========================================
  Lambda
  ========================================
*/

func createPoll(w http.ResponseWriter, r *http.Request) {
	log.Println("=== createPoll ===")

	returnCode := 0

	var err error
	var entries []Entry
	var byteSlice []byte
	poll := new(Poll)

	// https://golang.org/pkg/io/ioutil/#ReadAll
	if byteSlice, err = ioutil.ReadAll(r.Body); err != nil {
		returnCode = 1
		log.Println("createPoll err 1", err)
	}

	if returnCode == 0 {
		log.Println(string(byteSlice))

		// https://golang.org/pkg/encoding/json/#Unmarshal
		if err := json.Unmarshal(byteSlice, &entries); err != nil {
			returnCode = 2
			log.Println("createPoll err 2", err)
		}
	}

	if returnCode == 0 {
		if err := insertPollDB(poll); err != nil {
			returnCode = 3
			log.Println("createPoll err 3", err)
		}
	}

	if returnCode == 0 {
		// initialize poll
		poll.Question = entries[0].Value
		poll.Options = make([]string, len(entries)-1)

		for i := 1; i < len(entries); i++ {
			poll.Options[i-1] = entries[i].Value
		}

		log.Println(*poll)

		if err := updatePollDB(poll); err != nil {
			returnCode = 4
			log.Println("createPoll err 4", err)
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(poll.Code); err != nil {
			returnCode = 5
			log.Println("createPoll err 5", err)
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Poll could not be created.", w)
	}
}

func viewResults(w http.ResponseWriter, r *http.Request) {
	log.Println("=== viewResults ===")

	returnCode := 0

	var err error
	entry := new(Entry)
	poll := new(Poll)

	if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
		returnCode = 1
		log.Println("viewResults err 1", err)
	}

	if returnCode == 0 {
		if poll.Code, err = strconv.Atoi(entry.Value); err != nil {
			returnCode = 2
			log.Println("viewResults err 2", err)
		}
	}

	if returnCode == 0 {
		if err := submitCodeDB(poll); err != nil {
			returnCode = 3
			log.Println("viewResults err 3", err)
		}
	}

	if returnCode == 0 {
		generateResponsesDB(poll)
		poll.Results = calculateResultsDB(poll.Options, poll.Responses, poll.Results)
		log.Println(poll.Results)

		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 4
			log.Println("viewResults err 4", err)
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Results could not be gotten.", w)
	}
}

/*
  ========================================
  Pages
  ========================================
*/

func viewClient(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	setHeader(w)
	var page Page

	layout := path.Join("kalos/html", "client.html")
	content := path.Join("kalos/html", "content.html")

	t, err := template.ParseFiles(layout, content)
	if err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := t.ExecuteTemplate(w, "my-template", page); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Client page could not be viewed.", w)
	}
}

func viewAdmin(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	setHeader(w)
	var page Page

	layout := path.Join("kalos/html", "admin.html")
	content := path.Join("kalos/html", "content.html")

	t, err := template.ParseFiles(layout, content)
	if err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := t.ExecuteTemplate(w, "my-template", page); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Admin page could not be viewed.", w)
	}
}

/*
  ========================================
  Load Polls
  ========================================
*/

func loadPolls(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	var polls []Poll

	if returnCode == 0 {
		if err := loadPollsDB(&polls); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(polls); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Polls could not be loaded.", w)
	}
}

/*
  ========================================
  Load Poll
  ========================================
*/

func loadPoll(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)
	poll.ID = vestigo.Param(r, "pollID")

	if err := loadPollDB(poll); err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Poll could not be loaded.", w)
	}
}

/*
  ========================================
  Insert Poll
  ========================================
*/

func insertPoll(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)

	if returnCode == 0 {
		if err := insertPollDB(poll); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Poll could not be inserted.", w)
	}
}

/*
  ========================================
  Update Poll
  ========================================
*/

func updatePoll(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)
	poll.ID = vestigo.Param(r, "pollID")

	if err := json.NewDecoder(r.Body).Decode(poll); err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := updatePollDB(poll); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 3
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Poll could not be updated.", w)
	}
}

/*
  ========================================
  Delete Poll
  ========================================
*/

func deletePoll(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	pollID := vestigo.Param(r, "pollID")

	if err := deletePollDB(pollID); err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(true); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Poll could not be deleted.", w)
	}
}

/*
  ========================================
  Submit Code
  ========================================
*/

func submitCode(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)

	if err := json.NewDecoder(r.Body).Decode(poll); err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := submitCodeDB(poll); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 3
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Code could not be submitted.", w)
	}
}

/*
  ========================================
  Submit Response
  ========================================
*/

func submitResponse(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	response := new(Response)
	pollID := vestigo.Param(r, "pollID")

	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		returnCode = 1
	}

	if returnCode == 0 {
		if err := submitResponseDB(pollID, response); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(response); err != nil {
			returnCode = 3
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Response could not be submitted.", w)
	}
}

/*
  ========================================
  Get Results
  ========================================
*/

func getResults(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)
	poll.ID = vestigo.Param(r, "pollID")

	if returnCode == 0 {
		if err := loadPollDB(poll); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		poll.Results = calculateResultsDB(poll.Options, poll.Responses, poll.Results)

		if err := json.NewEncoder(w).Encode(true); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Results could not be gotten.", w)
	}
}

/*
  ========================================
  Generate Responses
  ========================================
*/

func generateResponses(w http.ResponseWriter, r *http.Request) {
	returnCode := 0

	poll := new(Poll)
	poll.ID = vestigo.Param(r, "pollID")

	if returnCode == 0 {
		if err := loadPollDB(poll); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		generateResponsesDB(poll)

		if err := json.NewEncoder(w).Encode(poll); err != nil {
			returnCode = 2
		}
	}

	// error handling
	if returnCode != 0 {
		handleError(returnCode, errorStatusCode, "Results could not be gotten.", w)
	}
}

/*
  ========================================
  Basics
  ========================================
*/

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Server", serverName)
}

func handleError(returnCode, statusCode int, message string, w http.ResponseWriter) {
	error := new(Error)
	error.Code = returnCode
	error.Message = message

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		log.Println(err)
	}
}
