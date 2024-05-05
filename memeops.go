package main

import (
//"fmt"
"encoding/json"
"log"
"net/http"
"strings"
)

type MemesAPI struct {
    TokenDB *TokenDatabase
}

func NewMemesAPI(tokenDB *TokenDatabase) *MemesAPI {
    return &MemesAPI{
        TokenDB: tokenDB,
    }
}

type Meme struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Image  string `json:"image"`
    Source string `json:"source"`
}

func main() {
	tokenDB := NewTokenDatabase()

    if err := tokenDB.Load("token_balances.json"); err != nil {
        log.Printf("Failed to load token balances: %v", err)
    }

    api := NewMemesAPI(tokenDB)

	http.HandleFunc("/memes", api.handleMemesRequest)
	http.HandleFunc("/update-token-balance", api.handleUpdateTokenBalanceRequest)

	log.Println("Server is up and running on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

func (api *MemesAPI) handleMemesRequest(w http.ResponseWriter, r *http.Request) {

	authToken := r.Header.Get("Authorization")
    if authToken == "" {
        http.Error(w, "Unauthorized: Missing auth token", http.StatusUnauthorized)
        return
    }

	userID := strings.TrimPrefix(authToken, "Bearer ")

    if api.TokenDB.GetBalance(userID) == 0 {
        api.TokenDB.UpdateBalance(userID, 150)
        if err := api.TokenDB.Save("token_balances.json"); err != nil {
            log.Printf("Failed to save token balances: %v", err)
        }
    }

	params := r.URL.Query()
	latitude := params.Get("lat")
	longitude := params.Get("lon")
	query := params.Get("query")

	 var memes []Meme
	 var err error

	 if err = validateCoordinates(latitude, longitude); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	 if latitude != "" && longitude != "" {

		 memes = []Meme{
			 {ID: 1, Title: "Meme at Lat: " + latitude + ", Lon: " + longitude + ", Category: " + query, Image: "https://example.com/meme.jpg"},
		 }
	 } else {
		 memes = []Meme{
			 {ID: 1, Title: "Default Meme 1", Image: "https://example.com/meme1.jpg"},
			 {ID: 2, Title: "Default Meme 2", Image: "https://example.com/meme2.jpg"},
		 }
	 }

	 responseJSON, err := json.Marshal(memes)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseJSON)

}

func (api *MemesAPI) handleUpdateTokenBalanceRequest(w http.ResponseWriter, r *http.Request) {

	var data struct {
        UserID string `json:"user_id"`
        Amount int    `json:"amount"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    api.TokenDB.UpdateBalance(data.UserID, api.TokenDB.GetBalance(data.UserID)+data.Amount)

    if err := api.TokenDB.Save("token_balances.json"); err != nil {
        log.Printf("Failed to save token balances: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

	response := struct {
        Success    bool   `json:"success"`
        Message    string `json:"message"`
        UserID     string `json:"user_id"`
        NewBalance int    `json:"new_balance"`
    }{
        Success:    true,
        Message:    "Token balance updated successfully",
        UserID:     data.UserID,
        NewBalance: api.TokenDB.GetBalance(data.UserID),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }}

