package main

import (
	"database/sql"

	"github.com/GabrielMSosa/menu-frontend-go/cmd/server/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {
	//abrimos la conexion con la DB
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bootdb")
	if err != nil {
		panic(err)
	}

	//inicializamos el server en gin
	eng := gin.New()
	//definimos los middleware manualmente porque vamos a usar uno custom
	//si no podemos usar simplemente el gin.Default() en vez del New y ya me levanta por defecto
	eng.Use(gin.Recovery(), gin.Logger())

	//swagger doc
	HOST := "localhost:8080"
	docs.SwaggerInfo.Host = HOST
	eng.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := routes.NewRouter(eng, db)
	router.MapRoutes()
	if err := eng.Run("localhost:8080"); err != nil {
		panic(err)
	}

}
