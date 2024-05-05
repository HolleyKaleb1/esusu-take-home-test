package main

import "strconv"

func validateCoordinates(latitude, longitude string) error {
    if latitude == "" && longitude == "" {
        return nil 
    }

    if latitude != "" {
        lat, err := strconv.ParseFloat(latitude, 64)
        if err != nil || lat < -90 || lat > 90 {
            return NewValidationError("Invalid latitude")
        }
    }

    if longitude != "" {
        lon, err := strconv.ParseFloat(longitude, 64)
        if err != nil || lon < -180 || lon > 180 {
            return NewValidationError("Invalid longitude")
        }
    }

    return nil
}

func NewValidationError(message string) error {
    return &ValidationError{message}
}

type ValidationError struct {
    message string
}

func (e *ValidationError) Error() string {
    return e.message
}