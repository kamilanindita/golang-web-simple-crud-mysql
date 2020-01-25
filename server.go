package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)


type Buku struct {
    Id    int
    Penulis  string
    Judul string
    Kota string
    Penerbit string
    Tahun int
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "website_crud"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Index.html",
        "views/Footer.html",
    ))
    data:=""
    var error = tmp.ExecuteTemplate(w,"Index",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}

func HandlerBuku(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT id,penulis,judul,kota,penerbit,tahun FROM buku")
    if err != nil {
        panic(err.Error())
    }
    buku := Buku{}
    data := []Buku{}
    for selDB.Next() {
        var id, tahun int
        var penulis, judul, kota, penerbit string
        err = selDB.Scan(&id, &penulis, &judul, &kota, &penerbit, &tahun)
        if err != nil {
            panic(err.Error())
        }
        buku.Id = id
        buku.Penulis = penulis
        buku.Judul = judul
        buku.Kota = kota
        buku.Penerbit = penerbit
        buku.Tahun = tahun
        data = append(data, buku)
    }
    defer db.Close()

    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Buku.html",
        "views/Footer.html",
    ))

    var error = tmp.ExecuteTemplate(w,"Buku", data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
    
}



func HandlerTamabah(w http.ResponseWriter, r *http.Request) {
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Tambah.html",
        "views/Footer.html",
    ))
    data:=""
    var error = tmp.ExecuteTemplate(w,"Tambah",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerSave(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        penulis := r.FormValue("penulis")
        judul := r.FormValue("judul")
        kota := r.FormValue("kota")
        penerbit := r.FormValue("penerbit")
        tahun := r.FormValue("tahun")
        insForm, err := db.Prepare("INSERT INTO buku (penulis,judul,kota,penerbit,tahun) VALUES(?,?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(penulis, judul,  kota, penerbit, tahun)
        log.Println("INSERT: Penulis: " + penulis + " | Judul: " + judul + " | Kota: " + kota+  " | Penerbit: " + penerbit + " | Tahun: " + tahun)
    }
    defer db.Close()
    http.Redirect(w, r, "/buku", 301)
}


func HandlerEdit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id,penulis,judul,kota,penerbit,tahun FROM buku WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    buku := Buku{}
    for selDB.Next() {
        var id, tahun int
        var penulis, judul, kota, penerbit string
        err = selDB.Scan(&id, &penulis, &judul, &kota, &penerbit, &tahun)
        if err != nil {
            panic(err.Error())
        }
        buku.Id = id
        buku.Penulis = penulis
        buku.Judul = judul
        buku.Kota = kota
        buku.Penerbit = penerbit
        buku.Tahun = tahun
    }
    defer db.Close()
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Edit.html",
        "views/Footer.html",
    ))

    var error = tmp.ExecuteTemplate(w,"Edit",buku)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        penulis := r.FormValue("penulis")
        judul := r.FormValue("judul")
        kota := r.FormValue("kota")
        penerbit := r.FormValue("penerbit")
        tahun := r.FormValue("tahun")
        id := r.URL.Query().Get("id")
        insForm, err := db.Prepare("UPDATE buku SET penulis=?, judul=?, kota=?, penerbit=?, tahun=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(penulis, judul,  kota, penerbit, tahun, id)
        log.Println("Update: ID:"+ id +" | Penulis: " + penulis + " | Judul: " + judul + " | Kota: " + kota+  " | Penerbit: " + penerbit + " | Tahun: " + tahun)
    }
    defer db.Close()
    http.Redirect(w, r, "/buku", 301)
}


func HandlerDelete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM buku WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE | ID :" + id)
    defer db.Close()
    http.Redirect(w, r, "/buku", 301)
}


func main() {
    log.Println("Server started on: http://localhost:8000")
    http.HandleFunc("/", HandlerIndex)
    http.HandleFunc("/buku", HandlerBuku)
    http.HandleFunc("/buku/tambah", HandlerTamabah)
    http.HandleFunc("/buku/edit", HandlerEdit)
    http.HandleFunc("/buku/save", HandlerSave)
    http.HandleFunc("/buku/update", HandlerUpdate)
    http.HandleFunc("/buku/delete", HandlerDelete)
    http.ListenAndServe(":8000", nil)
}