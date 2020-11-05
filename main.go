package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

//Creamos una estructura llamada task
type Habitacion struct {
    Id    int
    Piso  int
	Numero int
	Tipo string
	Capacidad int
	Reservado int64
}

//Funcion para crear la conexion al servidor mysql
func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "password"
    dbName := "hotel"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

//Variable donde esta el path de los las template
var tmpl = template.Must(template.ParseGlob("form/*"))

// funcion para mostrar la pagina principal 
func Index(w http.ResponseWriter, r *http.Request) {
	// Aqui hacemos el query a la base de dato
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM Habitacion ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	//Creamos estas variables donde se van a introducir los valores que nos de la base de dato para mostrar al usuario despues
    habitacion := Habitacion{}
    res := []Habitacion{}
    for selDB.Next() {
        var id, piso, numero, capacidad int
		var tipo string
		var reservado int64
        err = selDB.Scan(&id, &piso, &numero, &tipo, &capacidad, &reservado)
        if err != nil {
            panic(err.Error())
		}
        habitacion.Id = id
        habitacion.Piso = piso
		habitacion.Numero = numero
		habitacion.Tipo = tipo
		habitacion.Capacidad = capacidad
		habitacion.Reservado = reservado
        res = append(res, habitacion)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Habitacion WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    habitacion := Habitacion{}
    for selDB.Next() {
        var id, piso, numero, capacidad int
		var tipo string
		var reservado int64
        err = selDB.Scan(&id, &piso, &numero, &tipo, &capacidad, &reservado)
        if err != nil {
            panic(err.Error())
		}
        habitacion.Id = id
        habitacion.Piso = piso
		habitacion.Numero = numero
		habitacion.Tipo = tipo
		habitacion.Capacidad = capacidad
		habitacion.Reservado = reservado
    }
    tmpl.ExecuteTemplate(w, "Show", habitacion)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Task WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    habitacion := Habitacion{}
    for selDB.Next() {
        var id, piso, numero, capacidad int
		var tipo string
		var reservado int64
        err = selDB.Scan(&id, &piso, &numero, &tipo, &capacidad, &reservado)
        if err != nil {
            panic(err.Error())
		}
        habitacion.Id = id
        habitacion.Piso = piso
		habitacion.Numero = numero
		habitacion.Tipo = tipo
		habitacion.Capacidad = capacidad
		habitacion.Reservado = reservado
    }
    tmpl.ExecuteTemplate(w, "Edit", habitacion)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
		piso := r.FormValue("piso")
		numero := r.FormValue("numero")
        tipo := r.FormValue("tipo")
        capacidad := r.FormValue("capacidad")
        reservado := r.FormValue("reservado")	
        insForm, err := db.Prepare("INSERT INTO Habitacion(piso, numero, tipo, capacidad, reservado) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(piso, numero, tipo, capacidad, reservado)
        log.Println("INSERT: Piso: " + piso + " | Numero: " + numero + "Tipo: " + tipo + " | Capacidad: " + capacidad + " | Reservado: " + reservado )
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        piso := r.FormValue("piso")
		numero := r.FormValue("numero")
        tipo := r.FormValue("tipo")
        capacidad := r.FormValue("capacidad")
        reservado := r.FormValue("reservado")		
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Habitacion SET piso=?, numero=?, tipo=?, capacidad=?, reservado=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(piso, numero, tipo, capacidad, reservado, id)
        log.Println("UPDATE: Piso: " + piso + " | Numero: " + numero + "Tipo: " + tipo + " | Capacidad: " + capacidad + " | Reservado: " + reservado )
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    habitacion := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Habitacion WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(habitacion)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}