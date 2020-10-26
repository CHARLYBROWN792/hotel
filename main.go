package main

import(
	"fmt"
	"database/sql"
	"text/template"
	"net/http"
	"log"

	_ "github.go-sql-driver/mysql"
)

// We create a struct called habitacion
type habitacion struct{
	Id int
	Piso int
	Numero int
	Tipo string
	Capacidad int
	Status bool
}

//we create the conexion with the database
func dbConn() (db *sqlDb){
	dbDriver = "mysql"
	dbUser = "root"
	dbPass = "password"
	dbName = "hotel"
	db, sql = sql.Open(dbDriver, dbuser":+"dbPass"@/"dbName)
	if err != nil{
		panic (err.error())
	}

	return db
}

//This is the root for the project in HTML
var tmpl = template.Must(template.parseglob("form/*"))

//Function to show the Index "/"
func index = http.ResponseWriter, r *http.Request {
	// we make the query to our basedate and our table
	db := dbConn()
	selDb, err := db.Query("SELECT * FROM Habitacion id Desc")
	if err != nil{
		panic err =(err.error())
	}
	
	//now we create the var habitacion where is going to save the value that send the basedata
	habitacion := Habitacion{}
	res := []Habitacion
	//We create the variables, and We make a for to introduce all the data coming for database to those variables
	for selDb.next() {
		var id, piso, numero, Capacidad int
		var tipo string
		var reservado bool
		err = selDb.scan(&id, &piso, &numero, &tipo, &capacidad, &reservado)
		if err != nil{
			panic err= (err.error())
		}
		// Now we pass all the values for thos variables to the struct habitacion
		habitacion.Id = id
		habitacion.Piso = piso
		habitacion.Numero = numero
		habitacion.Tipo = tipo 
		habitacion.Capacidad = capacidad
		habitacion.reservado = reservado
		res = append(res, habitacion)
		}
		tmpl.ExecuteTemplate(w, "index", res)
		defer db.close()
	}
}


