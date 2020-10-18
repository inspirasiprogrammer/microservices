package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Infogempa struct {
	UserID  uint64   `xml:"user_id"`
	XMLName xml.Name `xml:"Infogempa"`
	Text    string   `xml:",chardata"`
	Gempa   struct {
		Text    string `xml:",chardata"`
		Tanggal string `xml:"Tanggal"`
		Jam     string `xml:"Jam"`
		Point   struct {
			Text        string `xml:",chardata"`
			Coordinates string `xml:"coordinates"`
		} `xml:"point"`
		Lintang   string `xml:"Lintang"`
		Bujur     string `xml:"Bujur"`
		Magnitude string `xml:"Magnitude"`
		Kedalaman string `xml:"Kedalaman"`
		Symbol    string `xml:"_symbol"`
		Wilayah1  string `xml:"Wilayah1"`
		Wilayah2  string `xml:"Wilayah2"`
		Wilayah3  string `xml:"Wilayah3"`
		Wilayah4  string `xml:"Wilayah4"`
		Wilayah5  string `xml:"Wilayah5"`
		Potensi   string `xml:"Potensi"`
	} `xml:"gempa"`
}

//Response struct
/**
* Seperti lazimnya web API , perlu adanya format standard untuk membangun struktur data API
* Disini saya memecah menjadi 3 bagian yakni :
* - status (berisi code status , misal 1 : success, 0: failed,dst)
* - message (penjelasan mengenai status)
* - data (isi data yang akan di sampaikan , dalam hal ini data produk dalam bentuk slice)
 */
type ResponseInfo struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Infogempa
}

// funtion untuk memparsing data MySQL ke JSON
func returnBMKG(w http.ResponseWriter, r *http.Request) {
	var result Infogempa
	if xmlBytes, err := getXML("https://data.bmkg.go.id/autogempa.xml"); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}

	fmt.Print(result)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

	// responseProd. Status = 1
	//mengisi valus status = 1 , dengan asumsi pasti success
	// responseProd.Message = "Success"
	// responseProd. Data = arr_products
	// mengisi komponen Data dengan data slice arr_products
	// //mengubah data sstruct menjadi JSON
	// json.NewEncoder (w). Encode(responseProd)
}

// tweaked from: https://stackoverflow.com/a/42718113/1170664
func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

//fungsi main diisi dengan routing dan fungsi http GET untuk menerima response dari request HTTP
// jika program ini dijalankan maka anda bisa megkases via browser/postman dengan URL : http://localhost:1234/bmkg
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/bmkg", returnBMKG).Methods("GET") // menjalurkan URL untuk dapat mengkases data JSON API product ke /bmkg
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":4321", router))

}
