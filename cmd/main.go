package main

import (
	"fmt"
	"log"
	"meal_prep/internal/db"
	"meal_prep/internal/recipes"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DBPath = "/tmp/meal_prep.db"
)

func main() {
	mealDB, err := db.Open(DBPath)
	defer mealDB.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = db.Init(mealDB)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/v1")
	{
		api.GET("/recipes", func(c *gin.Context) { recipes.ListRecipesHandler(c, mealDB) })
		api.GET("/recipes/:id", func(c *gin.Context) { recipes.GetRecipeHandler(c, mealDB) })
		api.POST("/recipes", func(c *gin.Context) { recipes.CreateRecipeHandler(c, mealDB) })
	}

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("exiting...")
}
