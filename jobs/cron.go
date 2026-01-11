package jobs

import (
	"api-reniec/database"
	"api-reniec/models"
	"log"
	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New()

	_, err := c.AddFunc("0 0 1 * *", func() {
		result := database.DB.Model(&models.Client{}).Where("status = ?", true).Update("usage", 0)
		
		if result.Error != nil {
			log.Printf("[CRON ERROR]: %v", result.Error)
		} else {
			log.Printf("[CRON SUCCESS]: Cuotas reiniciadas. Clientes afectados: %d", result.RowsAffected)
		}
	})

	if err != nil {
		log.Fatal("Error al configurar cron mensual:", err)
	}
	c.Start()
}