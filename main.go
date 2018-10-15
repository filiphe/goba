package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Drink holds an object representing a drink recipe
type Drink struct {
	Title        string   `json:"title"`
	Category     string   `json:"category"`
	Instructions string   `json:"instructions"`
	Ingredients  []string `json:"ingredients"`
	Link         string   `json:"source"`
}

// Drinks holds a list of objects of type Drink
type Drinks struct {
	Drinks []Drink `json:"drinks"`
}

var drinks Drinks

func init() {
	raw, err := ioutil.ReadFile("./drinks.json")
	if err != nil {
		return
	}
	json.Unmarshal(raw, &drinks)
}

func queryDrinks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if name, ok := query["name"]; ok {
		drink, err := getDrinkByName(name[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		bytes, err := json.Marshal(drink)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
	var result []Drink
	if ingredients, ok := query["ingredients"]; ok {
		fmt.Println(ingredients)
		drinks, err := getDrinksByIngredient(ingredients)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		for _, drink := range drinks {
			fmt.Println("Before adding by ingredients")
			if !contains(drinks, drink) {
				result = append(result, drink)
			}
		}
	}
	if category, ok := query["category"]; ok {
		fmt.Println(category[0])
		drinks, err := getDrinksByCategory(category[0])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		for _, drink := range drinks {
			fmt.Println("Before adding by category")
			if !contains(drinks, drink) {
				result = append(result, drink)
			}
		}
	}
	fmt.Println(result)
	bytes, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func getDrinkByName(name string) (Drink, error) {
	for _, v := range drinks.Drinks {
		if strings.ToLower(v.Title) == strings.Replace(strings.ToLower(name), "_", " ", -1) {
			return v, nil
		}
	}
	return Drink{}, fmt.Errorf("Drink %s not found in database", name)
}

func getDrinksByCategory(category string) ([]Drink, error) {
	var result []Drink
	for _, v := range drinks.Drinks {
		if strings.ToLower(v.Category) == strings.Replace(strings.ToLower(category), "_", " ", -1) {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return []Drink{}, fmt.Errorf("No drinks of category %s found in database", category)
	}
	return result, nil
}

func getDrinksByIngredient(queryIngredients []string) ([]Drink, error) {
	var result []Drink
	for _, v := range drinks.Drinks {
		for _, ingredient := range v.Ingredients {
			if strings.Contains(ingredient, "cl") {
				is := strings.Split(ingredient, " ")
				i := strings.TrimSpace(strings.Join(is[2:], " "))
				for _, queryIngredient := range queryIngredients {
					if strings.ToLower(i) == strings.ToLower(strings.Replace(queryIngredient, "_", " ", -1)) {
						result = append(result, v)
					}
				}
			}
		}
	}
	if len(result) == 0 {
		return []Drink{}, fmt.Errorf("Drink with ingredients %v not found in databse", queryIngredients)
	}
	return result, nil
}

func contains(drinks []Drink, d Drink) bool {
	for _, dr := range drinks {
		if dr.Equal(d) {
			return true
		}
	}
	return false
}

func (drink Drink) Equal(other Drink) bool {
	value := true
	if drink.Title == other.Title {
		if drink.Category == other.Category {
			if drink.Link == other.Link {
				for i := range drink.Ingredients {
					if drink.Ingredients[i] != other.Ingredients[i] {
						value = false
						break
					}
				}
			}
		}
	}
	return value
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// RESTy routes for "drinks" resource
	r.Route("/api", func(r chi.Router) {
		r.Get("/drink", queryDrinks)
	})

	http.ListenAndServe(":3333", r)
}
