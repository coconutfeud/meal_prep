package recipes

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Servings    *int       `json:"servings,omitempty"`
	PrepTime    *int       `json:"prep_time,omitempty"`
	CookTime    *int       `json:"cook_time,omitempty"`
	IsPublic    bool       `json:"is_public"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateRecipeRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Servings    *int    `json:"servings"`
	PrepTime    *int    `json:"prep_time"`
	CookTime    *int    `json:"cook_time"`
	IsPublic    *bool   `json:"is_public"`
}

func ListRecipesHandler(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
SELECT id, title, description, servings, prep_time, cook_time, is_public, created_at, updated_at
FROM recipes
ORDER BY created_at DESC
LIMIT 100
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query recipes"})
		return
	}

	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		var isPublicInt int
		if err := rows.Scan(
			&r.ID, &r.Title, &r.Description, &r.Servings,
			&r.PrepTime, &r.CookTime, &isPublicInt, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}

		r.IsPublic = isPublicInt == 1
		recipes = append(recipes, r)
	}

	c.JSON(http.StatusOK, recipes)
}

func GetRecipeHandler(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var r Recipe
	var isPublicInt int
	err = db.QueryRow(`
SELECT id, title, description, servings, prep_time, cook_time, is_public, created_at, updated_at
FROM recipes
WHERE id = ?
	`, id).Scan(
		&r.ID, &r.Title, &r.Description, &r.Servings,
		&r.PrepTime, &r.CookTime, &isPublicInt, &r.CreatedAt, &r.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	r.IsPublic = isPublicInt == 1
	c.JSON(http.StatusOK, r)
}

func CreateRecipeHandler(c *gin.Context, db *sql.DB) {
	var req CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	isPublic := 0
	if req.IsPublic != nil && *req.IsPublic {
		isPublic = 1
	}

	res, err := db.Exec(`
		INSERT INTO recipes (title, description, servings, prep_time, cook_time, is_public)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Title, req.Description, req.Servings, req.PrepTime, req.CookTime, isPublic)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert"})
		return
	}

	id64, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get id"})
		return
	}
	id := int(id64)

	// Return the created recipe (simple fetch)
	var r Recipe
	var isPublicInt int
	err = db.QueryRow(`
		SELECT id, title, description, servings, prep_time, cook_time,
		       is_public, created_at, updated_at
		FROM recipes
	`, id).Scan(
		&r.ID, &r.Title, &r.Description, &r.Servings,
		&r.PrepTime, &r.CookTime, &isPublicInt, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "created but failed to reload"})
		return
	}
	r.IsPublic = isPublicInt == 1

	c.JSON(http.StatusCreated, r)
}
