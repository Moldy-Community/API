package models

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UpdateUser struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	OldPassword string `json:"oldpassword" bson:"oldpassword"`
	NewPassword string `json:"newpassword" bson:"newpassword"`
}
