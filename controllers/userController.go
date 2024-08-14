package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"touch-test/middlewares"
	"touch-test/models"
	"touch-test/services"
	"touch-test/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type FileUploadController struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// Method to upload and process the file
func (c *FileUploadController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	savedFilePath, err := middlewares.SaveFile(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save the file"})
		return
	}

	if err := services.ProcessXlsxFile(c.DB, savedFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded and processed successfully"})
}

// Common function to fetch users, cache them, and return the users
func (c *FileUploadController) CacheAndFetchUsers() ([]models.User, error) {
	const cacheKey = "users"
	var users []models.User

	// Try to fetch data from Redis
	cachedData, err := c.RDB.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		if err := c.DB.Find(&users).Error; err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			return nil, err
		}

		if err := utils.AddToRedis(c.RDB, cacheKey, jsonData, 5*time.Minute); err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal([]byte(cachedData), &users); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// Common function to clear cache
func (c *FileUploadController) ClearUsersCache() error {
	const cacheKey = "users"
	return utils.ClearRedis(c.RDB, cacheKey)
}

// Common method to fetch a user by ID
func (c *FileUploadController) getUserByID(id string) (*models.User, error) {
	var user models.User

	if err := c.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Method to fetch all users from the database or Redis
func (c *FileUploadController) FetchAllUsers(ctx *gin.Context) {
	users, err := c.CacheAndFetchUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *FileUploadController) FetchUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.getUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user from database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *FileUploadController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.getUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user from database"})
		return
	}

	var input map[string]interface{}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := c.DB.Model(user).Updates(input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user in database"})
		return
	}

	if err := c.ClearUsersCache(); err != nil && err != redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not clear cache in Redis"})
		return
	}

	if _, err := c.CacheAndFetchUsers(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cache in Redis"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func (c *FileUploadController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.getUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user from database"})
		return
	}

	if err := c.DB.Delete(user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user from database"})
		return
	}

	if err := c.ClearUsersCache(); err != nil && err != redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete cache in Redis"})
		return
	}

	if _, err := c.CacheAndFetchUsers(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cache in Redis"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
