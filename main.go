package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Controller struct {
	Database *gorm.DB
}

var (
	dbConnectionString string = "host=localhost user=user password=password dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
)

func main() {
	controller := Controller{}
	controller.initDatabase()

	router := gin.Default()

	//Setup CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/lamps", controller.getLamps)
	router.GET("/lamps/:id", controller.getLampByID)
	router.POST("/lamps", controller.postLamp)
	router.PUT("/lamps/:id", controller.updateLamp)

	router.Run("0.0.0.0:8080")
}

func (c *Controller) initDatabase() {

	db, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	c.Database = db

	c.Database.AutoMigrate(&Lamp{})
}

func (ctrl *Controller) getLamps(c *gin.Context) {
	var lamps []Lamp
	result := ctrl.Database.Find(&Lamp{})
	result.Scan(&lamps)

	c.IndentedJSON(http.StatusOK, lamps)
}

func (ctrl *Controller) getLampByID(c *gin.Context) {
	id := c.Param("id")

	var lamp Lamp
	ctrl.Database.Model(&Lamp{}).First(&lamp, id)

	if lamp.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "lamp not found"})
		return
	}

	c.JSON(http.StatusOK, lamp)
}

func (ctrl *Controller) postLamp(c *gin.Context) {
	var newLamp Lamp
	if err := c.ShouldBindJSON(&newLamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	ctrl.Database.Create(&newLamp)

	c.JSON(http.StatusCreated, newLamp)
}

func (ctrl *Controller) updateLamp(c *gin.Context) {
	id := c.Param("id")
	var updatedLamp Lamp
	if err := c.ShouldBindJSON(&updatedLamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var lamp Lamp
	ctrl.Database.Model(&Lamp{}).First(&lamp, id)

	if lamp.ID != 0 {
		ctrl.Database.Model(&Lamp{}).Where("id = ?", id).Updates(updatedLamp)
		c.JSON(http.StatusOK, updatedLamp)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "lamp not found"})
}
