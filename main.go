package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/RemiEven/miam/dao"
	"github.com/RemiEven/miam/handler"

	"github.com/gorilla/mux"
)

const port = 8080

var (
	productDao *dao.ProductDao
)

func main() {
	log.Println("Starting") // FIXME this writes to stderr apparently

	productDao, err := dao.NewProductDao()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer productDao.Close()
	handler := handler.NewProductHandler(productDao)

	router := mux.NewRouter()
	router.HandleFunc("/product/{id}", handler.GetProductByID).Methods(http.MethodGet)
	router.HandleFunc("/product", handler.AddProduct).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Printf("Will try to listen on port %d", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wait = time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("Shutting down")
}
