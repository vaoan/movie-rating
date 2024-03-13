package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"movie-rating-api/app"
	"net/http"
)

func ConfigureRouter(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", Health).Methods("GET")
	//api.HandleFunc("/movies", GetMovies).Methods("GET")

}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"health": "OK"})
	if err != nil {
		log.Printf("failed to encode health endpoint json response with err: %s\n", err.Error())
	}
	log.Println("returning healthy 200 response")
	return
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.GetMovies(r.URL.Query())
	if err != nil {
		err = writeJSONResponse(w, err, http.StatusInternalServerError)
		if err != nil {
			fmt.Println("failed to write err body:", err.Error())
			return
		}
		return
	}

	err = writeJSONResponse(w, movies, http.StatusOK)
	if err != nil {
		fmt.Println("failed to write movies body:", err.Error())

		err = writeJSONResponse(w, err, http.StatusInternalServerError)
		if err != nil {
			fmt.Println("failed to write err body:", err.Error())
			return
		}
		return
	}

	return
}

func writeJSONResponse(w http.ResponseWriter, responseBody interface{}, httpStatusCode int) error {
	// marshal json bytes
	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	// write content type header
	w.Header().Add("Content-Type", "application/json")

	// write status code
	w.WriteHeader(httpStatusCode)

	// write body
	_, err = w.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}
