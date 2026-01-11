package models

import "gorm.io/gorm"

type Persona struct {
    gorm.Model
    DNI       string            `gorm:"uniqueIndex;size:8"`
    Nombres   string
    Paterno   string
    Materno   string
    Contactos []ContactoPersona `gorm:"foreignKey:PersonaID"` 
}

type ContactoPersona struct {
    gorm.Model
    PersonaID uint
    Tipo      string // "correo", "celular"
    Valor     string
}