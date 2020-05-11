package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
	"html/template"
	"log"
	"net"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {

	var templates = template.Must(template.New("locateip").ParseFiles("asklocation.html"))

	err := templates.ExecuteTemplate(w, "asklocation.html", nil)

	if err != nil {
		panic(err)
	}

}

func returnLatLong(w http.ResponseWriter, r *http.Request) {

	// not recommended to open the file each time per request
	// we put it here for tutorial sake.

	db, err := geoip2.Open("GeoLite2-City.mmdb")

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	if r.Method == "POST" {
		ipAddress := r.FormValue("ajax_post_data")

		// If you are using strings that may be invalid, check that IP is not nil
		// and a valid IP address -- see https://www.socketloop.com/tutorials/golang-validate-ip-address

		if ipAddress != "" {
			ip := net.ParseIP(ipAddress)
			record, err := db.City(ip)
			if err != nil {
				fmt.Println(err)
			}
			if len(record.Subdivisions) > 0 {
				fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
			}
			fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
			fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
			fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
			w.Write([]byte(fmt.Sprintf("{\"Country\":\"%v\",\"Lat\": \"%v\",\"Long\":\"%v\"}", record.Country.Names["en"], record.Location.Latitude, record.Location.Longitude)))

		}
	}
	db, err = geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	if r.Method == "POST" {
		ipAddress := r.FormValue("ajax_post_data")

		// If you are using strings that may be invalid, check that IP is not nil
		// and a valid IP address -- see https://www.socketloop.com/tutorials/golang-validate-ip-address

		if ipAddress != "" {
			ip := net.ParseIP(ipAddress)
			record, err := db.ASN(ip)
			if err != nil {
				fmt.Println(err)
			}
			number := record.AutonomousSystemNumber
			log.Print(number)
			organization := record.AutonomousSystemOrganization
			log.Println(organization)
		}
	}
}
func main() {

	mux := mux.NewRouter()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/returncoord", returnLatLong)

	panic(http.ListenAndServe(":8080", mux))
}
