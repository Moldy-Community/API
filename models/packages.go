package models

type Package struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Password    string `json:"password" bson:"password"`
	Version     string `json:"version" bson:"version"`
}
type PackageUpdate struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	Url         string `json:"url" bson:"url"`
	Description string `json:"description" bson:"description"`
	Password    string `json:"password" bson:"password"`
	NewPassword string `json:"newpassword" bson:"newpassword"`
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

type AuthPassword struct {
	Password string `json:"password" bson:"password"`
}

type Packages []*Format
