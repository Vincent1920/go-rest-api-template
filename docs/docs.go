import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "todo-api/docs"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	{
		// endpoint
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}