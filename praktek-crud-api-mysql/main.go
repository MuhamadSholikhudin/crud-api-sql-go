package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {

	//penamaan Camel Cas untuk Import Package supaya bisa di pakai dari luar
	Id   int    `json:"id"` // `` cara membuat penamaan ulang pada golang pada saat di GET
	Name string `json:"name"`
}

func Conn() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_belajar_golang")
	if err != nil {
		return nil, err
	}

	return db, nil

}

var datas []User

func Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	db, err := Conn()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var id = 0
	rows, err := db.Query("select id, name from user where id > ?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		var each = User{}
		var err = rows.Scan(&each.Id, &each.Name)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(resp))

}

func Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id_r, _ := strconv.Atoi(vars["id"])

	var db, err = Conn()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	var result = User{}
	user_id := id_r
	err = db.
		QueryRow("select id, name from user where id = ?", user_id).
		Scan(&result.Id, &result.Name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(resp))

}

func Post(w http.ResponseWriter, r *http.Request) {
	//untuk membuat json pertama kita harus set Header
	w.Header().Set("Content-Type", "application/json")

	// data := r.Body
	//mendecode requset body langsung menjadi json
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(data)

	db, err := Conn()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("insert into user (id, name) values (?, ?)", data.Id, data.Name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("insert success!")

	response := map[string]interface{}{
		"status": "Oke",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Update(w http.ResponseWriter, r *http.Request) {
	//untuk membuat json pertama kita harus set Header
	w.Header().Set("Content-Type", "application/json")

	//mendecode requset body langsung menjadi json
	data := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db, err := Conn()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("update user set name = ? where id = ?", data.Name, data.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Updated success!")

	response := map[string]interface{}{
		"status": "Oke",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {
	//untuk membuat json pertama kita harus set Header
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id_r, _ := strconv.Atoi(vars["id"])

	var db, err = Conn()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	user_id := id_r

	_, err = db.Exec("delete from user where id = ?", user_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Delete success!")

	response := map[string]interface{}{
		"status": "Oke",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/user", Index).Methods("GET")
	r.HandleFunc("/user", Post).Methods("POST")
	r.HandleFunc("/user/{id}", Get).Methods("GET")
	r.HandleFunc("/user/{id}", Update).Methods("PUT")
	r.HandleFunc("/user/{id}", Delete).Methods("DELETE")

	fmt.Println("LIsten on Port 127.0.0.1:8080")
	http.ListenAndServe(":8080", r)

}
