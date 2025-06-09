package model

type Character struct {
	Name                string `json:"name"`
	Origin              string `json:"origin"`
	Species             string `json:"species"`
	AdditionalAttribute string `json:"additional_attribute"`
}
