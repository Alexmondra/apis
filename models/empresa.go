package models

import "gorm.io/gorm"

type Empresa struct {
    gorm.Model
    RUC              string `gorm:"uniqueIndex;size:11"`
    RazonSocial      string
    Estado           string
    Condicion        string
    Direccion        string
    Ubigeo            string
    Departamento     string
    Provincia        string
    Distrito         string
    EsAgenteRetencion bool
    Contactos        []ContactoEmpresa `gorm:"foreignKey:EmpresaID"`
}

type ContactoEmpresa struct {
    gorm.Model
    EmpresaID uint
    Tipo      string // "correo", "celular"
    Valor     string
}