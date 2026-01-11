package database

import (
	"api-reniec/models"
	"log"
	"os"
	"time" 
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DB_DSN") 
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/reniec?charset=utf8mb4&parseTime=True&loc=Local"
	}

	var err error
	// Intentar conectar hasta 10 veces (esperando 2 segundos entre cada una)
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Conexión exitosa a MariaDB")
			break
		}
		log.Printf("⏳ Base de datos no lista (intento %d/10)... esperando 2s", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ No se pudo conectar a la DB después de varios intentos:", err)
	}

	// Migraciones automáticas
	DB.AutoMigrate(&models.Persona{}, &models.ContactoPersona{}, &models.Empresa{}, &models.ContactoEmpresa{},&models.Client{})
}