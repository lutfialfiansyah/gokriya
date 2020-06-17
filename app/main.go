package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_userHttpDelivery "github.com/bxcodec/go-clean-arch/users/delivery/http"
	_userRepo "github.com/bxcodec/go-clean-arch/users/repository/mysql"
	_userUcase "github.com/bxcodec/go-clean-arch/users/usecase"

	_articleHttpDelivery "github.com/bxcodec/go-clean-arch/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/bxcodec/go-clean-arch/article/delivery/http/middleware"
	_articleRepo "github.com/bxcodec/go-clean-arch/article/repository/mysql"
	_articleUcase "github.com/bxcodec/go-clean-arch/article/usecase"
	_authorRepo "github.com/bxcodec/go-clean-arch/author/repository/mysql"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	//dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	//connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	//dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, 8976, dbUser, dbPass, dbName)

	dbConn, err := sql.Open(`postgres`, psqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	userRepo := _userRepo.NewMysqlUsersRepository(dbConn)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	userUsecase := _userUcase.NewUsersUsecase(userRepo,timeoutContext)
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)
	_userHttpDelivery.NewArticleHandler(e,userUsecase)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
