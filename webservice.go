package main

import (
	"database/sql"								//used for unifrom database access
	"fmt"									//print statements
	"github.com/go-martini/martini"						//extra frame work build on net/http
	_ "github.com/lib/pq"							//go sql driver 
	"net/http"								//framework
	"strings"				
	"os"									//used to pull env vars
)

func SetupDB() *sql.DB {							//Change host to fit use 10.254.76.103 
										//Current setup is for open shift							
	db, err := sql.Open("postgres", "host=postgresql user=docker password=docker dbname=postgres sslmode=disable") 		//my only lib/pq usage? login into postgres database
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ShowDB(db *sql.DB, r *http.Request, rw http.ResponseWriter) { //localhost:3000/?search=name
		search := "%" + r.URL.Query().Get("search") + "%"
		rows, err := db.Query(`SELECT fname, LName, cost, city 
                           FROM custauth 
                           WHERE FName ILIKE $1
                           OR LName ILIKE $1
                           OR city ILIKE $1`, search)
		PanicIf(err)
		defer rows.Close()

		var FirstName, LastName, cost, city string
		for rows.Next() {
			err := rows.Scan(&FirstName, &LastName, &cost, &city)
			PanicIf(err)
			fmt.Fprintf(rw, "First Name: %s\nLast Name: %s\nCost: %s\nCity: %s\n\n", FirstName, LastName, cost, city)
		}
	}

func InsertPur(r *http.Request, db *sql.DB){
//	_, err := db.Query("INSERT INTO custauth (custid, fname, lname, city,state, country,email,cost,errorflag) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)",
	_, err := db.Exec("INSERT INTO custauth (custid, fname, lname, city,state, country,email,cost,errorflag) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)",
	r.FormValue("custid"),
		r.FormValue("fname"),
		r.FormValue("lname"),
		r.FormValue("city"),
		r.FormValue("state"),
		r.FormValue("country"),
		r.FormValue("email"),
		r.FormValue("cost"),
		r.FormValue("errorflag"))

	PanicIf(err)
}
	
	

func main() {
	m := martini.Classic()
	m.Map(SetupDB())
	m.Get("/", func() string {return  "Welcome to GoSQL database"})
  	m.Get("/var", func() string {
	for _, e := range os.Environ() {
        pair := strings.Split(e, "=")
		 fmt.Println(pair[0])}	
		 return "hello"
	})
	m.Get("/show", ShowDB)
	m.Get("/print1",func() string {return os.Getenv("POSTGRESDB_SERVICE_HOST")})
	m.Get("/print4",func() string {return os.Getenv("OPENSHIFT_POSTGRESQL_PASSWORD")})
	m.Get("/print5",func() string {return os.Getenv("OPENSHIFT_POSTGRESQL_USER")})
	m.Get("/print6",func() string {return os.Getenv("POSTGRESQL_USER")})
	m.Get("/var", func() string {
	for _, e := range os.Environ() {
        pair := strings.Split(e, "=")
		 fmt.Println(pair[0])}	
		 return "yo"
	})
	m.Post("/add", InsertPur)
	m.RunOnAddr(":8080")
}
