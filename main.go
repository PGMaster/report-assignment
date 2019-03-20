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

const (
	host     = "localhost"
	port     = "5432"
	user     = "william"
	password = "password"
	dbname   = "postgres"
)

type Chapter struct {
	Company_name string    `json:"company"`
	Project_name string    `json:"project"`
	Chapter_name string    `json:"chapter"`
	Version      []Version `json:"version"`
}
type Version struct {
	Created_By         string `json:"created_by"`
	Chapter_Version_Id int    `json:"chapter_version_id"`
	Version_Number     int    `json:"version_number"`
	Created            string `json:"created"`
	Appversion         string `json:"appversion"`
}

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
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port,
		user, password, dbname)

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

	chapter_id := keys[0]
	getResponse(chapter_id)
	getVersionInfo(chapter_id)
	w.Header().Set("Content-Type", "application/json")

	Chapters = append(Chapters, Chapter{Version: Versions})

	json.NewEncoder(w).Encode(Chapters)
}

func getResponse(chapter_id string) {
	rows, err := db.Query("select chapter_name, project_name, company_name from chapter inner join project on chapter_id = $1 AND project_id = chapter_project_id inner join company on company_id = project_company_id", chapter_id)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var chapter_name string
		var project_name string
		var company_name string
		if err := rows.Scan(&chapter_name, &project_name, &company_name); err != nil {
			log.Fatal(err)
		}
		fmt.Println(project_name)
		Chapters = append(Chapters, Chapter{Chapter_name: chapter_name, Project_name: project_name, Company_name: company_name})

	}
}
func getVersionInfo(chapter_id string) {
	rows, err := db.Query("SELECT person_username, chapter_version_id, chapter_version_number, chapter_version_create_date, chapter_version_appversion from chapter_version inner join person on person_id = chapter_version_person_id where chapter_version_chapter_id=$1", chapter_id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var person_username string
		var chapter_version_id int
		var chapter_version_number int
		var chapter_version_create_date string
		var chapter_version_appversion string
		if err := rows.Scan(&person_username, &chapter_version_id, &chapter_version_number, &chapter_version_create_date, &chapter_version_appversion); err != nil {
			log.Fatal(err)
		}
		switch chapter_version_appversion {
		case "11.0":
			chapter_version_appversion := "CC 2015"
			Versions = append(Versions, Version{Created_By: person_username, Chapter_Version_Id: chapter_version_id, Version_Number: chapter_version_number, Created: chapter_version_create_date, Appversion: chapter_version_appversion})

		case "12.0":
			chapter_version_appversion := "CC 2017"
			Versions = append(Versions, Version{Created_By: person_username, Chapter_Version_Id: chapter_version_id, Version_Number: chapter_version_number, Created: chapter_version_create_date, Appversion: chapter_version_appversion})

		default:
			chapter_version_appversion := "CC 2015"
			Versions = append(Versions, Version{Created_By: person_username, Chapter_Version_Id: chapter_version_id, Version_Number: chapter_version_number, Created: chapter_version_create_date, Appversion: chapter_version_appversion})
		}

	}
}
