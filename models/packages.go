package models

type Package struct {
	Name    string  `json:"name" bson:"name"`
	Author  string  `json:"author" bson:"author"`
	Version float32 `json:"version" bson:"version"`
	Url     string  `json:"url" bson:"url"`
}

type Packages []*Package
