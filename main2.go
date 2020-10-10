package main

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
)

type Input struct {
	Panjang     int    `json:"panjang"`
	Lebar       int    `json:"lebar"`
	Tinggi       int    `json:"tinggi"`
}

type Result struct {
	JenisBangun string `json:"jenis_bangun"`
	Volume        int    `json:"volume"`
}

func main() {

    router := mux.NewRouter()
	router.HandleFunc("/api/hitung-volume", HitungVolume)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HitungVolume(w http.ResponseWriter, r *http.Request) {
	var hasilHitung []Result
	var volume []Input

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close() // function di bawahnya tetap dijalankan, walaupun ada error
	if err != nil {
		WrapAPIError2(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &volume) // untuk decode data json hasil inputan user
	if err != nil { // jika error, jalankan fungsi error handling WrapAPIError
		WrapAPIError2(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range volume {
		hasilHitung = append(hasilHitung, Result{
			JenisBangun: v.JenisBangun2(),
			Volume:        v.RumusVolume2(),
		})
	}

	WrapAPIData2(w, r, hasilHitung, http.StatusOK, "success")

	fmt.Println(hasilHitung)
}

func (v *Input) JenisBangun2() string {  // method pada Volume Jenis Bangunan
	var jenis string
	if (v.Panjang == v.Lebar){
		if (v.Lebar == v.Tinggi) {
			jenis = "Kubus"
		}
	}else{
		jenis = "Balok"
	}
	return jenis
}

func (v *Input) RumusVolume2() int { //method yang dimiliki oleh struck
	return v.Panjang * v.Lebar * v.Tinggi
}

func WrapAPIError2(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":          code,
		"error_type":    http.StatusText(code),
		"error_details": message,
	})
	if err == nil {
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API error : %s", err))
	}
}

func WrapAPISuccess2(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API success : %s", err))
	}
}

func WrapAPIData2(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
		"data":   data,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API data : %s", err))
	}
}
