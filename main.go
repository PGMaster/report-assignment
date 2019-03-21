package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Why not use config file?
const (
	host     = "localhost"
	port     = "5432"
	user     = "william"
	password = "password"
	dbname   = "postgres"
)

// Should be using camel casing instead of snake casing
type Chapter struct {
	CompanyName string    `json:"company"`
	ProjectName string    `json:"project"`
	ChapterName string    `json:"chapter"`
	// Does not follow the json tag names based on the coding task
	// Should have been 'versions' instead of 'version'
	Version     []Version `json:"versions"`
}
type Version struct {
	CreatedBy        string `json:"created_by"`
	ChapterVersionId int    `json:"chapter_version_id"`
	VersionNumber    int    `json:"version_number"`
	Created          string `json:"created"`
	AppVersion       string `json:"appversion"`
}

// Why is this declared in the global scope?
var Chapters []Chapter
var Versions []Version

func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/chapter_versions/", handler)
	http.ListenAndServe(":80", nil)
}

func initDb() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	keys, ok := r.URL.Query()["chapter_id"]
	// fmt.Println(r.URL.EscapedPath())
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'chapter_id' is missing")
		return
	}

	chapterId := keys[0]
	getResponse(chapterId)
	getVersionInfo(chapterId)

	w.Header().Set("Content-Type", "application/json")

	// This logic is completely wrong here
	// The results versions are added as a new chapter, meaning the result would be
	/*
		[
	 		{
				"company": "CtrlPrint",
				"project": "Testnot",
				"chapter": "Test Chapter  1",
				"versions": null
 	 		},
 	  		{
				"company": "",
				"project": "",
				"chapter": "",
				"versions": [
						{
							"created_by": "ctrl-romain",
							"chapter_version_id": 262,
							"version_number": 4,
							"created": "2015-08-26T10:49:53.059514Z",
							"appversion": "CC 2015"
						},
						{
							"created_by": "ctrl-jens",
							"chapter_version_id": 261,
							"version_number": 3,
							"created": "2015-08-26T10:48:41.795795Z",
							"appversion": "CC 2015"
						}
            		],
					...
			},
 	 	]
	*/
	Chapters = append(Chapters, Chapter{Version: Versions})

	json.NewEncoder(w).Encode(Chapters)
}

func getResponse(chapterId string) {
	rows, err := db.Query(`select chapter_name, project_name, company_name 
									from chapter 
									inner join project on chapter_id = $1 AND project_id = chapter_project_id 
									inner join company on company_id = project_company_id`,
									chapterId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var chapterName string
		var projectName string
		var companyName string
		if err := rows.Scan(&chapterName, &projectName, &companyName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(projectName)
		Chapters = append(Chapters, Chapter{ChapterName: chapterName, ProjectName: projectName, CompanyName: companyName})
	}
}
func getVersionInfo(chapterId string) {
	rows, err := db.Query(`SELECT person_username, chapter_version_id, chapter_version_number, chapter_version_create_date, chapter_version_appversion 
									from chapter_version 
									  inner join person on person_id = chapter_version_person_id 
									where chapter_version_chapter_id=$1`, chapterId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var personUsername string
		var chapterVersionId int
		var chapterVersionNumber int
		var chapterVersionCreateDate string
		var chapterVersionAppversion string
		if err := rows.Scan(&personUsername, &chapterVersionId, &chapterVersionNumber, &chapterVersionCreateDate, &chapterVersionAppversion); err != nil {
			log.Fatal(err)
		}
		// why not make this a function?
		switch chapterVersionAppversion {
		case "11.0":
			chapter_version_appversion := "CC 2015"
			Versions = append(Versions, Version{CreatedBy: personUsername, ChapterVersionId: chapterVersionId, VersionNumber: chapterVersionNumber, Created: chapterVersionCreateDate, AppVersion: chapter_version_appversion})

		case "12.0":
			chapter_version_appversion := "CC 2017"
			Versions = append(Versions, Version{CreatedBy: personUsername, ChapterVersionId: chapterVersionId, VersionNumber: chapterVersionNumber, Created: chapterVersionCreateDate, AppVersion: chapter_version_appversion})

		default:
			chapter_version_appversion := "CC 2015"
			Versions = append(Versions, Version{CreatedBy: personUsername, ChapterVersionId: chapterVersionId, VersionNumber: chapterVersionNumber, Created: chapterVersionCreateDate, AppVersion: chapter_version_appversion})
		}

	}
}
