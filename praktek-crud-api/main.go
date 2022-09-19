package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {

	//penamaan Camel Cas untuk Import Package supaya bisa di pakai dari luar
	Id   int    `json:"id"` // `` cara membuat penamaan ulang pada golang pada saat di GET
	Name string `json:"name"`
}

var datas []User

func Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(resp))

}

func Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	data := User{}
	id, _ := strconv.Atoi(vars["id"])

	for _, item := range datas {
		if item.Id == id {
			data = item
		}
	}

	resp, err := json.Marshal(data)
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
	datas = append(datas, data)
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
	vars := mux.Vars(r)
	// NewUpdate := User{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(data)

	id, _ := strconv.Atoi(vars["id"])
	index := 0
	for i, item := range datas {
		if item.Id == id {
			index = i
		}
	}
	fmt.Println(index)

	datas[index] = datas[len(datas)-1]
	datas[len(datas)-1] = User{}
	datas = datas[:len(datas)-1]

	datas = append(datas, data)

	fmt.Println(datas)

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

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/user", Index).Methods("GET")
	r.HandleFunc("/user/{id}", Get).Methods("GET")
	r.HandleFunc("/user/{id}", Update).Methods("PUT")
	r.HandleFunc("/user", Post).Methods("POST")

	// http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

	// 	//mwnginiliasi method
	// 	switch r.Method {
	// 	case "GET":
	// 		//untuk membuat json pertama kita harus set Header
	// 		w.Header().Set("Content-Type", "application/json")

	// 		//cara parsing struct ke json
	// 		// data := map[string]interface{}{
	// 		// 	"id":   1,
	// 		// 	"name": "mahmud",
	// 		// }

	// 		// json marshal convert json ke object
	// 		resp, err := json.Marshal(datas)
	// 		if err != nil {
	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		}

	// 		w.Write([]byte(resp))

	// 	case "POST":
	// 		//untuk membuat json pertama kita harus set Header
	// 		w.Header().Set("Content-Type", "application/json")

	// 		// data := r.Body
	// 		//mendecode requset body langsung menjadi json
	// 		data := User{}
	// 		err := json.NewDecoder(r.Body).Decode(&data)
	// 		if err != nil {
	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		}
	// 		fmt.Println(data)
	// 		datas = append(datas, data)
	// 		response := map[string]interface{}{
	// 			"status": "Oke",
	// 		}
	// 		err = json.NewEncoder(w).Encode(response)
	// 		if err != nil {
	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}
	// 	}

	// 	// w.Write([]byte("response ini"))
	// })

	fmt.Println("LIsten on Port 127.0.0.1:8080")
	http.ListenAndServe(":8080", r)

}
