package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

//Creamos una estructura llamada tipo
type Tipo struct {
    Id    int
    Nombre  string
    Capacidad int
    Descripcion string
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
var tmpl = template.Must(template.ParseGlob("formtipo/*"))

// funcion para mostrar la pagina principal 
func Index(w http.ResponseWriter, r *http.Request) {
	// Aqui hacemos el query a la base de dato
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM Tipo ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	//Creamos estas variables donde se van a introducir los valores que nos de la base de dato para mostrar al usuario despues
    tipo := Tipo{}
    res := []Tipo{}
    for selDB.Next() {
        var id, capacidad int
        var nombre, descripcion string
        err = selDB.Scan(&id, &nombre, &capacidad, &descripcion)
        if err != nil {
            panic(err.Error())
		}
        tipo.Id = id
        tipo.Nombre = nombre
        tipo.Capacidad = capacidad
        tipo.Descripcion = descripcion
        res = append(res, tipo)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Tipo WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    tipo := Tipo{}
    for selDB.Next() {
        var id, capacidad int
        var nombre, descripcion string
        err = selDB.Scan(&id, &nombre, &capacidad, &descripcion)
        if err != nil {
            panic(err.Error())
        }
        tipo.Id = id
        tipo.Nombre = nombre
        tipo.Capacidad = capacidad
        tipo.Descripcion = descripcion
    }
    tmpl.ExecuteTemplate(w, "Show", tipo)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Tipo WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    tipo := Tipo{}
    for selDB.Next() {
        var id, capacidad int
        var nombre, descripcion string
        err = selDB.Scan(&id, &nombre, &capacidad, &descripcion)
        if err != nil {
            panic(err.Error())
        }
        tipo.Id = id
        tipo.Nombre = nombre
        tipo.Capacidad = capacidad
        tipo.Descripcion = descripcion
    }
    tmpl.ExecuteTemplate(w, "Edit", tipo)
    defer db.Close()
}


func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        nombre := r.FormValue("nombre")
        capacidad := r.FormValue("capacidad")
        descripcion := r.FormValue("descripcion")
        insForm, err := db.Prepare("INSERT INTO Tipo(nombre, capacidad, descripcion) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(nombre, capacidad, descripcion)
        log.Println("INSERT: Nombre: " + nombre + "Capacidad: " + capacidad + " | Descripcion: " + descripcion)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        nombre := r.FormValue("nombre")
        capacidad := r.FormValue("capacidad")
        descripcion := r.FormValue("descripcion")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Tipo SET nombre=?, capacidad=?, descripcion=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(nombre, capacidad, descripcion, id)
        log.Println("UPDATE: Nombre: " + nombre + "Capacidad: " + capacidad + " | Descripcion: " + descripcion)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    tipo := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Tipo WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(tipo)
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