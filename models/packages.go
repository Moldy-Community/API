package models

type Package struct {
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Version     string `json:"version" bson:"version"`
}

type Packages []*Package
