package routes

import (
	controllers "api-golang-compro/controllers"
	"api-golang-compro/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)


func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
    
    // Convert c.Handler to a gin.HandlerFunc
    corsMiddleware := func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
    
    // Use the converted middleware in gin
    r.Use(corsMiddleware)
    
    // Buat objek limiter
    var limiter = rate.NewLimiter(rate.Limit(10), 1) // Contoh: 10 permintaan per detik

    // Tambahkan middleware rate limiting ke router utama
    r.Use(func(c *gin.Context) {
        // Gunakan limiter untuk memeriksa rate limiting
        if limiter.Allow() == false {
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    })
    
    
    // set db to gin context
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
    })
    
    r.Static("/public","./public")
    // r.POST("/file", controllers.HandleUploadFile)
    r.DELETE("/file/:name", controllers.HandleRemoveFile)
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)
	r.GET("/career", controllers.GetAllCareer)
    r.GET("/career/:id", controllers.GetCareerById)
	r.GET("/category", controllers.GetAllCategory)
    r.GET("/category/:id", controllers.GetCategoryById)
	r.GET("/position", controllers.GetAllPosition)
    r.GET("/position/:id", controllers.GetPositionById)
	r.GET("/category_home", controllers.GetAllCategoryHome)
    r.GET("/category_home/:id", controllers.GetCategoryHomeById)
	r.GET("/form", controllers.GetAllForm)
    r.GET("/form/:id", controllers.GetFormById)
	r.GET("/home", controllers.GetAllHome)
    r.GET("/home/:id", controllers.GetHomeByid)
	r.GET("/service", controllers.GetAllService)
    r.GET("/service/:id", controllers.GetServiceByid)
	r.GET("/portfolio", controllers.GetAllPortfolio)
    r.GET("/portfolio/:id", controllers.GetPortfolioById)
	r.GET("/navbar", controllers.GetAllNavbar)
    r.GET("/navbar/:id", controllers.GetNavbarById)
    r.GET("/portfolio_hepytech", controllers.GetAllPortfolioHepytech)
    r.GET("/portfolio_hepytech/:id", controllers.GetPortfolioHepytechById)
    r.GET("/portfolio_hepytech/slug/:slug", controllers.GetPortfolioHepytechBySlug)
    r.GET("/logactivity/", controllers.GetAllLogActivity)
    r.GET("/logo/", controllers.GetAllLogo)
    r.GET("/logo/:id", controllers.GetLogoByid)
    r.GET("/impact", controllers.GetAllImpact)
    r.GET("/impact/:id", controllers.GetImpactByid)
    
	Career := r.Group("/career")
    Career.Use(middlewares.JwtAuthMiddleware()) //use jwt
    Career.POST("/", controllers.CreateCareer)
    Career.PATCH("/:id", controllers.UpdateCareer)
    Career.DELETE("/:id", controllers.DeleteCareer)


	Category := r.Group("/category")
	Category.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Category.POST("/", controllers.CreateCategory)
    Category.PATCH("/:id", controllers.UpdateCategory)
    Category.DELETE("/:id", controllers.DeleteCategory)

	Position := r.Group("/position")
	Position.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Position.POST("/", controllers.CreatePosition)
    Position.PATCH("/:id", controllers.UpdatePosition)
    Position.DELETE("/:id", controllers.DeletePosition)

	Category_home := r.Group("/category_home")
	Category_home.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Category_home.POST("/", controllers.CreateCategoryHome)
    Category_home.PATCH("/:id", controllers.UpdateCategoryHome)
    Category_home.DELETE("/:id", controllers.DeleteCategoryHome)

    Form := r.Group("/form")
	Form.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Form.POST("/", controllers.CreateForm)
    Form.PATCH("/:id", controllers.UpdateForm)
    Form.DELETE("/:id", controllers.DeleteForm)


    Home := r.Group("/home")
	Home.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Home.POST("/", controllers.CreateHome)
    Home.PATCH("/:id", controllers.UpdateHome)
    Home.DELETE("/:id", controllers.DeleteHome)

    Navbar := r.Group("/navbar")
	Navbar.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Navbar.POST("/", controllers.CreateNavbar)
    Navbar.PATCH("/:id", controllers.UpdateNavbar)
    Navbar.DELETE("/:id", controllers.DeleteNavbar)

    Portfolio := r.Group("/portfolio")
    Portfolio.Use(middlewares.JwtAuthMiddleware()) //use jwt
    Portfolio.POST("/", controllers.CreatePortfolio)
    Portfolio.PATCH("/:id", controllers.UpdatePortfolio)
    Portfolio.DELETE("/:id", controllers.DeletePortfolio)



    Service := r.Group("/service")
	Service.Use(middlewares.JwtAuthMiddleware())  //use jwt
    Service.POST("/", controllers.CreateService)
    Service.PATCH("/:id", controllers.UpdateService)
    Service.DELETE("/:id", controllers.DeleteService)

    PortfolioHepytech := r.Group("/portfolio_hepytech")
    PortfolioHepytech.Use(middlewares.JwtAuthMiddleware())  //use jwt
    PortfolioHepytech.POST("/", controllers.CreatePortfolioHepytech)
    PortfolioHepytech.PATCH("/:id", controllers.UpdatePortfolioHepytech)
    PortfolioHepytech.DELETE("/:id", controllers.DeletePortfolioHepytech)
    
    logo := r.Group("/logo")
    logo.Use(middlewares.JwtAuthMiddleware())  //use jwt
    logo.POST("/", controllers.CreateLogo)
    logo.PATCH("/:id", controllers.UpdateLogo)
    logo.DELETE("/:id", controllers.DeleteLogo)

    impact := r.Group("/impact")
    impact.Use(middlewares.JwtAuthMiddleware())  //use jwt
    impact.POST("/", controllers.CreateImpact)
    impact.PATCH("/:id", controllers.UpdateImpact)
    impact.DELETE("/:id", controllers.DeleteImpact)

    return r
}
