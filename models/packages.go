package models

type Package struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Version     string `json:"version" bson:"version"`
}
type PackageUpdate struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Version     string `json:"version" bson:"version"`
}

type Format struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Version     string `json:"version" bson:"version"`
}

type Packages []*Format
