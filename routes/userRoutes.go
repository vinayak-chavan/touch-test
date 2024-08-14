package routes

import (
    "touch-test/controllers"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
    fileUploadController := &controllers.FileUploadController{DB: db, RDB: rdb}

    router.POST("/users", fileUploadController.UploadFile)
    router.GET("/users", fileUploadController.FetchAllUsers)
		router.GET("/users/:id", fileUploadController.FetchUserByID)
		router.PUT("/users/:id", fileUploadController.UpdateUser)
		router.DELETE("users/:id", fileUploadController.DeleteUserByID)
}
