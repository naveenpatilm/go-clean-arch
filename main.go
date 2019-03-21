package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	_articleHttpDeliver "github.com/naveenpatilm/go-clean-arch/article/delivery/http"
	_articleRepo "github.com/naveenpatilm/go-clean-arch/article/repository"
	_articleUcase "github.com/naveenpatilm/go-clean-arch/article/usecase"
	_authorRepo "github.com/naveenpatilm/go-clean-arch/author/repository"
	"github.com/naveenpatilm/go-clean-arch/middleware"
	"github.com/naveenpatilm/go-clean-arch/models"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

}

func main() {

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	sslEnable := viper.GetString(`ssl.mode`)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPass, sslEnable)

	fmt.Print(dsn)
	dbConn, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Print("Helloe error her")
		log.Fatal(err)
		os.Exit(1)
	}

	defer dbConn.Close()

	dbConn.AutoMigrate(&models.Article{})

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDeliver.NewArticleHttpHandler(e, au)

	e.Start(viper.GetString("server.address"))
}
