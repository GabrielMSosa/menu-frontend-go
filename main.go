package main

import (
	"database/sql"

	"github.com/GabrielMSosa/menu-frontend-go/cmd/server/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	eureka "github.com/xuanbo/eureka-client"
)

func main() {
	//abrimos la conexion con la DB
	db, err := sql.Open("mysql", "root:root@tcp(my-db:3306)/bootdb")
	if err != nil {
		panic(err)
	}
	// create eureka client
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://host.docker.internal:8761/eureka/",
		App:                   "menu-frontend-go",
		Port:                  8080,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})
	// start client, register、heartbeat、refresh
	client.Start()
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
	if err := eng.Run(); err != nil {
		panic(err)
	}

}
