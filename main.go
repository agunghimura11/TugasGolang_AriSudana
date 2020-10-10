package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Volume struct { // deklarasi variable struct Volumen untuk menyimpan input dari user
	Panjang     int    `json:"panjang"`
	Lebar       int    `json:"lebar"`
	Tinggi      int    `json:"tinggi"`
}

type Hasil struct { // deklarasi variable struct Hasil untuk menyimpan hasil perhitungan
	JenisBangun string `json:"jenis_bangun"`
	Volume      int    `json:"volume"`
}

func main() {
	router := mux.NewRouter() //mengimplementasikan permintaan router dan dispatcher untuk mencocokkan permintaan yang masuk ke handlernya masing-masing.
	router.HandleFunc("/api/hitung-volume", HitungVolume)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HitungVolume(w http.ResponseWriter, r *http.Request) { // fungsi untuk menghitung volume

	var hasilHitung []Hasil // deklarasi variable dengan tipe Hasil
	var volume []Volume // deklarasi variable dengan tipe Volume
	if r.Method != "POST" {
		ErrorHandler(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body) //io utility function untuk read seluruh data dari io.Reader
	defer r.Body.Close() // function di bawahnya tetap dijalankan, walaupun ada error
	if err != nil { //jika error jalankan errorHandler
		ErrorHandler(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &volume) // convert respose body ke struct volume
	if err != nil { // jika error, jalankan fungsi error handling ErrorHandler
		ErrorHandler(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError) //mengirimkan pesan error
		return
	}

	for _, v := range volume { //memasukan data ke 
		hasilHitung = append(hasilHitung, Hasil{
			JenisBangun: v.JenisBangun(), // memangil method dari struct untuk menghitung Jenius Bangun
			Volume:        v.RumusVolume(), // memanggil method Rumus Volume
		})
	}

	DataWrapper(w, r, hasilHitung, http.StatusOK, "success") // hasil dari perhitungan diteruskan ke di DataWrapper
}

func (v *Volume) JenisBangun() string {  // method pada Volume Jenis Bangunan
	jenis:= "Kubus"; // inisiasi nilai dari jenis
	hasil:= &jenis // variable hasil menyimpan pointer dari jenis
	if (v.Panjang == v.Lebar){ // jika panjang lebar dan tinggi sama maka merupakan bangun kubus
		if (v.Lebar == v.Tinggi) {
			*hasil = "Kubus"
		}else {
			*hasil = "Balok"
		}
	}
	return jenis
}

func ErrorMessage(message string) {
	log.Println("Error")
}

func (v *Volume) RumusVolume() (int, bool) { //method yang dimiliki oleh struct untuk menghitung volume
	return v.Panjang * v.Lebar * v.Tinggi
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, message string, code int) {
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

func WrapAPISuccess(w http.ResponseWriter, r *http.Request, message string, code int) {
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

func DataWrapper(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
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