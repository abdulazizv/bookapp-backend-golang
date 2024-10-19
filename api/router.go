package api

import (
	h "net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "gitlab.com/bookapp/api/docs" // docs
	v1 "gitlab.com/bookapp/api/handler/v1"
	"gitlab.com/bookapp/api/middleware"
	"gitlab.com/bookapp/api/tokens"
	"gitlab.com/bookapp/config"
	"gitlab.com/bookapp/pkg/logger"
	"gitlab.com/bookapp/storage"
	"gitlab.com/bookapp/storage/repo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Options struct {
	Cfg            config.Config
	Storage        storage.StorageI
	Log            logger.Logger
	CasbinEnforcer *casbin.Enforcer
	Redis          repo.InMemoryStorageI
}

// @title           Book store API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(opt *Options) *gin.Engine {
	router := gin.New()
	basicAuth := middleware.BasicAuth()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corConfig := cors.DefaultConfig()
	corConfig.AllowAllOrigins = true
	corConfig.AllowCredentials = true
	corConfig.AllowHeaders = []string{"*"}
	corConfig.AllowBrowserExtensions = true
	corConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corConfig))

	jwtHandler := tokens.JWTHandler{
		SigninKey: opt.Cfg.SigningKey,
		Log:       opt.Log,
	}
	handlerV1 := v1.New(&v1.HandlerV1Option{
		Cfg:        &opt.Cfg,
		Storage:    opt.Storage,
		Log:        opt.Log,
		JwtHandler: jwtHandler,
		Redis:      opt.Redis,
	})
	router.Use(middleware.NewAuth(opt.CasbinEnforcer, jwtHandler, config.Load()))
	router.Use(basicAuth.Middleware)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(h.StatusOK, gin.H{
			"message": "Server is running!!!",
		})
	})
	router.MaxMultipartMemory = 8 << 20 // 8 Mib

	api := router.Group("/v1")
	// media
	media := api.Group("/store")
	media.GET("/:filename", handlerV1.GetFile)
	media.POST("/upload", handlerV1.Upload)
	//category
	category := api.Group("/category")
	category.POST("/", handlerV1.CreateCategory)
	category.GET("/:id", handlerV1.GetByIdCategory)
	category.GET("/list", handlerV1.GetListCategory)
	category.PUT("/", handlerV1.UpdateCategory)
	category.DELETE("/:id", handlerV1.DeleteCategory)
	category.GET("/books", handlerV1.GetCategoryId)

	// subcategory
	subcategory := api.Group("/subcategory")
	subcategory.POST("/", handlerV1.CreateSubCategory)
	subcategory.GET("/:id", handlerV1.GetSubCategoryById)
	subcategory.PUT("/", handlerV1.UpdateSubCategory)
	subcategory.DELETE("/:id", handlerV1.DeleteSubCategory)

	// books
	book := api.Group("/book")
	book.POST("/", handlerV1.CreateBook)
	book.GET("/search", handlerV1.GetListBooks)
	book.GET("/top", handlerV1.GetBookTop)
	book.GET("/:id", handlerV1.GetByIdBook)
	book.PUT("/", handlerV1.UpdateBook)
	book.DELETE("/:id", handlerV1.DeleteBook)
	book.POST("/like", handlerV1.CreateBookLike)
	book.DELETE("/like", handlerV1.DeleteBookLike)
	book.GET("/filter", handlerV1.GetBooksFilter)
	book.GET("/mostread", handlerV1.GetBookReadALot)
	book.GET("/audios", handlerV1.GetBookAudios)

	// comment
	comment := api.Group("/comment")
	comment.POST("/post", handlerV1.CreateComment)
	comment.PUT("/put", handlerV1.UpdateComment)
	comment.DELETE("/:id", handlerV1.DeleteComment)

	// user
	user := api.Group("/client")
	user.POST("/register", handlerV1.RegisterUser)
	user.POST("/login", handlerV1.LoginUser)
	user.GET("/:id", handlerV1.GetUser)
	user.PUT("/update", handlerV1.UpdateUser)
	user.DELETE("/:id", handlerV1.DeleteUser)

	// author
	author := api.Group("/author")
	author.POST("/", handlerV1.CreateAuthor)
	author.GET("/:id", handlerV1.GetAuthor)
	author.GET("/list", handlerV1.GetAuthorList)
	author.PUT("/", handlerV1.UpdateAuthor)
	author.DELETE("/:id", handlerV1.DeleteAuthor)

	//admin
	admin := api.Group("/admin")
	admin.POST("/", handlerV1.AddAdmin)
	admin.POST("/login", handlerV1.LoginAdmin)
	admin.DELETE("/:id", handlerV1.DeleteAdmin)
	admin.GET("/", handlerV1.GetAllAdmin)

	//superadmin
	superadmin := api.Group("/superadmin")
	superadmin.POST("/login", handlerV1.LoginSuperAdmin)
	superadmin.POST("/", handlerV1.AddSuperAdmin)

	//statistic
	statistic := api.Group("/statistic")
	statistic.GET("", handlerV1.GetStatistic)
	statistic.GET("/category/bookcount", handlerV1.GetCategoryBookCount)
	statistic.GET("/week/bookcount", handlerV1.GetAddedWeekBook)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
