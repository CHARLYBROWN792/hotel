package main

import (
    "fmt"
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
	Reservado int
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
var tmpl = template.Must(template.ParseGlob("formhab/*"))

func leernum() int {
    var n int
  
	fmt.Scanf("%d\n", &n)
	fmt.Println()

	return (n)
}

func Inserthabdefecto(n1 int, n2 int) {
    db := dbConn()
    var Id int
    for i := 1; i <= n1 {
        for j := 1; j <= n2; j++{                       	
            Id := Id + 1
            piso := i
		    numero := j
            tipo := ""
            capacidad := 0	
                insForm, err := db.Prepare("INSERT INTO Habitacion(Id, piso, numero, tipo, capacidad) VALUES(?,?,?,?,?)")
                if err != nil {
                    panic(err.Error())
                }
            insForm.Exec(Id, piso, numero, tipo, capacidad)
        }
    }
    defer db.Close()
}

//Funcion para contar cuantos registros hay dentro de la tabla Habitacion
func CountTotalHabitacion() int {
    db := dbConn()
    selDB, err := db.Query("SELECT COUNT(*) FROM Habitacion ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}

    defer db.Close()

    //Creamos una variable para pasar el query a un numero integro
    var count int
    for selDB.Next() {   
        if err := selDB.Scan(&count); err != nil {
            log.Fatal(err)
        }
    }
    return (count)        
}


func DeleteHabitacion(n1 int) {
    db := dbConn()
    for i := 1; i <= n1; i++ {    
        
        delForm, err := db.Prepare("DELETE FROM Habitacion WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        delForm.Exec(i)
        log.Println("DELETE")
              
    }
    defer db.Close()
}

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
		var reservado int
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

func Showhab(w http.ResponseWriter, r *http.Request) {
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
		var reservado int
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

func Newhab(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edithab(w http.ResponseWriter, r *http.Request) {
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
		var reservado int
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

func Inserthab(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
		piso := r.FormValue("piso")
		numero := r.FormValue("numero")
        tipo := r.FormValue("tipo")
        capacidad := r.FormValue("capacidad")
        reservado := r.FormValue("reservado")	
        insForm, err := db.Prepare("INSERT INTO Habitacion(piso, numero, tipo, capacidad, reservado) VALUES(?,?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(piso, numero, tipo, capacidad, reservado)
        log.Println("INSERT: Piso: " + piso + " | Numero: " + numero + "Tipo: " + tipo + " | Capacidad: " + capacidad + " | Reservado: " + reservado )
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Updatehab(w http.ResponseWriter, r *http.Request) {
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

func Deletehab(w http.ResponseWriter, r *http.Request) {
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
    cha:= CountTotalHabitacion()
    fmt.Println ("numero de habitaciones", cha)

    if cha == 0{
        fmt.Println("Por favor ingrese el numero de pisos")
        num1:= leernum()
        fmt.Println("Por favor ingrese el numero de habitaciones por piso")
        num2:= leernum()

        fmt.Println("El total de habitaciones para su hotel es:", num1*num2)
        
        if cha < num1*num2 {
            if cha > 0{
                DeleteHabitacion(cha)
            }       
            Inserthabdefecto(num1, num2)
        } 
        cha = CountTotalHabitacion()
    }
 
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/showhab", Showhab)
    http.HandleFunc("/newhab", Newhab)
    http.HandleFunc("/edithab", Edithab)
    http.HandleFunc("/inserthab", Inserthab)
    http.HandleFunc("/updatehab", Updatehab)
    http.HandleFunc("/deletehab", Deletehab)
    http.ListenAndServe(":8080", nil)
}