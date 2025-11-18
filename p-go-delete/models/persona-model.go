package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Persona struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Documento string             `bson:"documento" json:"documento"`
	Nombre    string             `bson:"nombre" json:"nombre"`
	Apellido  string             `bson:"apellido" json:"apellido"`
	Edad      int                `bson:"edad" json:"edad"`
	Correo    string             `bson:"correo" json:"correo"`
	Telefono  string             `bson:"telefono" json:"telefono"`
	Direccion string             `bson:"direccion" json:"direccion"`
}
