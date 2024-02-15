package main

import (
	handler "go-nutritioncalculator2/handlers"
	repository "go-nutritioncalculator2/repositories"
	service "go-nutritioncalculator2/services"
	"log"
	"net/http"
	"os"

	_ "go-nutritioncalculator2/docs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Nutrition Calculator API documentation
// @version 1.0.0
// @host go-nutritioncalculatorv2.onrender.com
// @BasePath /
// @description API for record all meal that you have in each day and help you calculate summary nutrition in each meal and you can save favorite menu and favorite meal for track your diet easily and create your own menu

func main() {
	d, err := sqlx.Connect("postgres", "postgres://postgresql_nutritioncalculator_zrbf_user:M86b8FJlY8Y00Ln0XYug0FnpQUUaO6c1@dpg-cn70prdjm4es73bp9aqg-a.singapore-postgres.render.com/postgresql_nutritioncalculator_zrbf")
	if err != nil {
		panic(err)
	}
	userRepo := repository.NewUserRepositoryDB(d)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	menuRepo := repository.NewMenuRositoryDB(d)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	favListRepo := repository.NewFavListRepositoryDB(d)
	favListService := service.NewFavListService(favListRepo)
	favListHandler := handler.NewFavListHandler(favListService)
	recordRepo := repository.NewRecordRepositoryDB(d)
	recordService := service.NewRecordService(recordRepo)
	recordHandler := handler.NewRecordHandler(recordService)
	multiHandler := handler.NewMultiHandler(menuService, userService, favListService)
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	credentialsOk := handlers.AllowCredentials()

	r.HandleFunc("/user/", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/user/{user_id}", userHandler.GetUserDetail).Methods("GET")
	r.HandleFunc("/user/login", userHandler.LogIn).Methods("PUT")
	r.HandleFunc("/user/userdetail", userHandler.UpdateUserDetail).Methods("PUT")

	r.HandleFunc("/menu/", menuHandler.CreateMenu).Methods("POST")
	r.HandleFunc("/menu/{menu_id}", menuHandler.DeleteMenu).Methods("DELETE")
	r.HandleFunc("/menu/", menuHandler.GetAllMenues).Methods("GET")
	r.HandleFunc("/menu/", menuHandler.UpdateMenu).Methods("PUT")

	r.HandleFunc("/favlist/", favListHandler.CreateFavList).Methods("POST")
	r.HandleFunc("/favlist/{favlist_id}", favListHandler.DeleteFavList).Methods("DELETE")
	r.HandleFunc("/favlist/{user_id}", favListHandler.GetFavListsByUserId).Methods("GET")
	r.HandleFunc("/favlist/", favListHandler.UpdateFavList).Methods("PUT")

	r.HandleFunc("/record/", recordHandler.CreateRecord).Methods("POST")
	r.HandleFunc("/record/{record_id}", recordHandler.DeleteRecord).Methods("DELETE")
	r.HandleFunc("/record/{user_id}", recordHandler.GetRecordsByUserId).Methods("GET")
	r.HandleFunc("/record/", recordHandler.UpdateRecord).Methods("PUT")

	r.HandleFunc("/recover/", multiHandler.RecoverDeletedMenu).Methods("PUT")
	r.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(r)))
}
