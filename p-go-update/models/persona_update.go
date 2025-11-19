package models

type PersonaUpdate struct {
	Nombre   string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido string `json:"apellido,omitempty" bson:"apellido,omitempty"`
	Edad     int    `json:"edad,omitempty" bson:"edad,omitempty"`
}
