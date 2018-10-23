package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var URLArray []string
var igcMap = make(map[int]igc.Track)
var timeStarted = time.Now()
var id int

type URLForm struct {
	URL string `jason:"URL"`
}

type trackDB struct {
	Uid          string
	Pilot        string
	H_date       string
	Glider       string
	Glider_ID    string
	Track_length string
	Url          string
}

func connectDB() *mongo.Collection {
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("igc").Collection("igcTracks")

	return collection
}
func insertToDB(collection *mongo.Collection, trackFile trackDB) {

	res, err := collection.InsertOne(context.Background(), trackFile)

	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Print(id)
	if id == nil {
		fmt.Print("ID nil!")
	}

}

func validateURL(collection *mongo.Collection, url string) int64 {

	filter := bson.NewDocument(bson.EC.String("url", ""+url+""))

	cur, err := collection.Count(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	return cur
}
func getIndex(x []string, y string) int {
	for i, j := range x {
		if j == y {
			//found
			return i
		}
	}
	return -1
}

func elapsedTime(start time.Time) string {

	duration := time.Since(start)

	days := int(duration.Hours() / 24)
	years := int(days / 365)
	ddays := days % 365
	realMonths := days / 30
	realDays := ddays % 30
	hours := int(duration.Hours()) % 24
	min := int(duration.Minutes()) % 60
	sec := int(duration.Seconds()) % 60

	return fmt.Sprintf("P%dY%dM%dD%dH%dm%ds", years, realMonths, realDays, hours, min, sec)
}

func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

//Handler is the first handler function which will be responsible for the first page
func Handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/igcinfo/api" {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}
	// Set response content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	pathVars := mux.Vars(r)
	if len(pathVars) != 0 {
		http.Error(w, "400 - Bad Request, too many URL arguments.", 400)
		return
	}
	resp := "{\n"
	resp += "\"uptime\": \" " + elapsedTime(timeStarted) + "\" ,\n"
	resp += "\"info\": \"Service for Paragliding tracks.\",\n"
	resp += "\"version\": \"v1\",\n"
	resp += "\"no\": \"" + fmt.Sprintf("%d", len(URLArray)) + "\" \n"
	resp += "\n}"

	fmt.Fprint(w, resp)
}

func postHANDLER1(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	pathVars := mux.Vars(r)
	if len(pathVars) != 0 {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}

	URL := &URLForm{}

	err := json.NewDecoder(r.Body).Decode(&URL) //obtaining the URL recived from the request's body

	if err != nil {
		http.Error(w, err.Error(), 400) //checking for errors in the process and returning bad request if so
		return
	}

	track, err1 := igc.ParseLocation(URL.URL) //Used for parsing the obtained URL
	if err1 != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest) //checking for errors in the process and returning bad request if so
		return
	}
	collection := connectDB()
	var trackFile trackDB
	if validateURL(collection, URL.URL) == 0 {

		URLArray = append(URLArray, URL.URL)
		track, _ = igc.ParseLocation(URL.URL)
		igcMap[len(URLArray)-1] = track

		uID, error := collection.Count(context.Background(), nil)
		if error != nil {
			fmt.Print("Err count")
		}
		track.UniqueID = fmt.Sprintf("%d", uID) //I decided to use the array index as UniqueID

		//client := connectDB() //connecting to DB
		trackFile = trackDB{track.UniqueID, track.Pilot, track.Date.String(), track.GliderType, track.GliderID, fmt.Sprintf("%f", trackLength(track)), URL.URL}
		insertToDB(collection, trackFile)
	}
	//result := bson.NewDocument()

	trackDBobj := trackDB{}
	filter := bson.NewDocument(bson.EC.String("url", ""+URL.URL+""))
	error := collection.FindOne(context.Background(), filter).Decode(&trackDBobj)
	if error != nil {
		log.Fatal(err)
	}

	//uID := len(URLArray)

	track.UniqueID = trackDBobj.Uid

	resp := "{\n\"id\": " + "\"" + track.UniqueID + "\"\n}" //formating the response in json format

	fmt.Fprint(w, resp)

}

func getHANDLER1(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/igcinfo/api/igc" {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	pathVars := mux.Vars(r)
	if len(pathVars) != 0 {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}

	collection := connectDB()
	trackFile := trackDB{}
	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	lenght, err := collection.Count(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp := "["
	var i int64 = 0
	for cur.Next(context.Background()) {
		err := cur.Decode(&trackFile)
		if err != nil {
			log.Fatal(err)
		}
		resp += trackFile.Uid
		if i+1 == lenght {
			break
		}
		resp += ","
		i++
	}
	resp += "]"

	// resp := "["
	// for i := range URLArray {

	// 	resp += strconv.Itoa(i)
	// 	if i+1 == len(URLArray) {
	// 		break
	// 	}
	// 	resp += ","
	// }
	// resp += "]"

	fmt.Fprint(w, resp)

}

//Handler2 is the handler which will be responsible for requests that contain an IDs
func Handler2(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	pathVars := mux.Vars(r)
	if len(pathVars) != 1 {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}

	// validation
	if pathVars["id"] == "" {

		http.Error(w, "400 - Bad Request, you entered an empty ID.", http.StatusBadRequest)
		return

	}
	id, err := strconv.Atoi(pathVars["id"])
	if err != nil {

		http.Error(w, "400 - Bad Request, you entered an ID which is not numeric!", 400)
		return

	}
	if id > len(igcMap)-1 {

		http.Error(w, "404 - Not found, you entered an ID which is not in our system!", 404)
		return

	}
	//end of validation

	resp := "{\n"
	resp += "  \"H_date\": " + "\"" + igcMap[id].Date.String() + "\",\n"
	resp += "  \"pilot\": " + "\"" + igcMap[id].Pilot + "\",\n"
	resp += "  \"glider\": " + "\"" + igcMap[id].GliderType + "\",\n"
	resp += "  \"glider_id\": " + "\"" + igcMap[id].GliderID + "\",\n"
	resp += "  \"track_lenght\": " + "\"" + fmt.Sprintf("%f", trackLength(igcMap[id])) + "\"\n"
	resp += "}"

	fmt.Fprint(w, resp)

}

//Handler3 is the handler that's responsible for requests that will contain an ID and a field
func Handler3(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}

	pathVars := mux.Vars(r)
	if len(pathVars) != 2 {
		http.Error(w, "400 - Bad Request, too many URL arguments.", http.StatusBadRequest)
		return
	}
	// validation
	if pathVars["id"] == "" {

		http.Error(w, "400 - Bad Request, you entered an empty ID.", http.StatusBadRequest)
		return

	}
	id, err := strconv.Atoi(pathVars["id"])
	if err != nil {

		http.Error(w, "400 - Bad Request, you entered an ID which is not numeric!", 400)
		return

	}
	if id > len(igcMap)-1 {

		http.Error(w, "404 - Not found, you entered an ID which is not in our system!", 404)
		return

	}
	if pathVars["field"] == "" {

		http.Error(w, "400 - Bad Request, you entered an empty Field.", http.StatusBadRequest)
		return

	}

	switch pathVars["field"] {

	case "pilot":
		fmt.Fprintf(w, "%s", igcMap[id].Pilot)

	case "glider":
		fmt.Fprintf(w, "%s", igcMap[id].GliderType)

	case "glider_id":
		fmt.Fprintf(w, "%s", igcMap[id].GliderID)

	case "track_length":
		fmt.Fprintf(w, "%f", trackLength(igcMap[id]))

	case "H_date":
		fmt.Fprintf(w, "%s", igcMap[id].Date.String())

	default:
		http.Error(w, "", http.StatusNotFound)

	}

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/igcinfo/api", Handler).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc", getHANDLER1).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc", postHANDLER1).Methods("POST")
	r.HandleFunc("/igcinfo/api/igc/{id}", Handler2).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc/{id}/{field}", Handler3).Methods("GET")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// collection := client.Database("igc").Collection("igcTracks")

// res, err := collection.InsertOne(context.Background(), map[string]string{
// 	"uID":          "" + track.UniqueID + "",
// 	"pilot":        "" + track.Pilot + "",
// 	"h_date":       "" + track.Date.String() + "",
// 	"glider":       "" + track.GliderType + "",
// 	"glider_ID":    "" + track.GliderID + "",
// 	"track_length": "" + fmt.Sprintf("%f", trackLength(igcMap[id])) + "",
// 	"url":          "" + URL.URL + ""})

// if err != nil {
// 	log.Fatal(err)
// }
// id := res.InsertedID
// if id == nil {
// 	fmt.Print("ID nil!")
// }