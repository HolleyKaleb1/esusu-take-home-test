package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "sync"
)

type TokenDatabase struct {
    sync.RWMutex
    Tokens map[string]int 

func NewTokenDatabase() *TokenDatabase {
    return &TokenDatabase{
        Tokens: make(map[string]int),
    }
}

func (db *TokenDatabase) Load(filename string) error {
    db.Lock()
    defer db.Unlock()

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // File does not exist, no need to load
        }
        return err
    }

    if err := json.Unmarshal(data, &db.Tokens); err != nil {
        return err
    }

    return nil
}

func (db *TokenDatabase) Save(filename string) error {
    db.RLock()
    defer db.RUnlock()

    data, err := json.Marshal(db.Tokens)
    if err != nil {
        return err
    }

    if err := ioutil.WriteFile(filename, data, 0644); err != nil {
        return err
    }

    return nil
}

func (db *TokenDatabase) GetBalance(userID string) int {
    db.RLock()
    defer db.RUnlock()

    return db.Tokens[userID]
}

func (db *TokenDatabase) UpdateBalance(userID string, balance int) {
    db.Lock()
    defer db.Unlock()

    db.Tokens[userID] = balance
}
