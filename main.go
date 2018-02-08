package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// our main function
func main() {

	snowmen = append(snowmen, Snowman{ID: "Mega Snowman", Weight: 100})
	snowmen = append(snowmen, Snowman{ID: "Tall Snowman", Weight: 40})
	snowmen = append(snowmen, Snowman{ID: "Medium Snowman", Weight: 30})
	snowmen = append(snowmen, Snowman{ID: "Baby Snowman", Weight: 10})

	router := mux.NewRouter()
	router.HandleFunc("/snowmen", GetSnowmen).Methods("GET")
	router.HandleFunc("/snowman/build", BuildSnowman).Methods("POST")
	router.HandleFunc("/flamethrower", FlameSnowmen).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetSnowmen returns a list of all snowmen.
func GetSnowmen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snowmen)
}

// BuildSnowman creates a snowman.
func BuildSnowman(w http.ResponseWriter, r *http.Request) {
	var snowman Snowman
	_ = json.NewDecoder(r.Body).Decode(&snowman)
	snowmen = append(snowmen, snowman)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snowman)
}

// FlameSnowmen damadges and/or destroys snowmen.
func FlameSnowmen(w http.ResponseWriter, r *http.Request) {
	var fuel int
	_ = json.NewDecoder(r.Body).Decode(&fuel)

	// Damage caused: 1 litre of fuel causes the loss of 2kg of snow.
	const fuelToSnowLoss = 2
	var totalDamadge = fuelToSnowLoss * fuel

	// How many get hit.
	var numberOfDamagedSnowmen = rand.Intn(len(snowmen))

	// Shuffle snowmen.
	var disorder = snowmen
	for i := range disorder {
		j := rand.Intn(i + 1)
		disorder[i], disorder[j] = disorder[j], disorder[i]
	}

	// Burn Baby Burn.
	var i = 0
	var wastedFuel = 0
	for i <= numberOfDamagedSnowmen {
		var weightLeft = 0
		var damadge = rand.Intn(totalDamadge)
		totalDamadge = totalDamadge - damadge
		weightLeft = disorder[i].CauseDamadge(damadge)
		wastedFuel = wastedFuel + weightLeft/fuelToSnowLoss
		if totalDamadge <= 0 {
			i = numberOfDamagedSnowmen
		} else {
			i = i + 1
		}
	}

	fmt.Println("totalDamadge ", totalDamadge)
	// Report.
	var damadgeReport = DamadgeReport{Fuel: fuel, Survivors: disorder, Wasted: wastedFuel}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(damadgeReport)
}

// Snowman definition.
type Snowman struct {
	ID     string `json:"id,omitempty"`
	Weight int    `json:"weight,omitempty"`
}

// CauseDamadge definition.
func (s Snowman) CauseDamadge(damadge int) int {
	var left = s.Weight - damadge
	if left < 0 {
		s.Weight = 0
		left = 0
	} else {
		s.Weight = s.Weight - damadge
	}

	return left
}

// DamadgeReport definition.
type DamadgeReport struct {
	Fuel      int       `json:"fuel, omitempty"`
	Survivors []Snowman `json:"survivors"`
	Wasted    int       `json:"wasted_fuel, omitempty"`
}

var snowmen []Snowman
