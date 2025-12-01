package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/js-bruno/mariage-api/internal/adapter"
	"github.com/js-bruno/mariage-api/internal/controller"
	"github.com/js-bruno/mariage-api/internal/utils"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type PaymentRequest struct {
	Value    float64 `json:"value"`
	ItemId   int     `json:"item_id"`
	ItemDesc string  `json:"item_desc"`
	Email    string  `json:"email"`
}

type PaymentResponse struct {
	QRCode string `json:"qrcode"`
}

var EnvConfig utils.Env

func main() {
	var err error
	utils.SetStructuredLogging()
	EnvConfig, err = utils.GetEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	client, close := adapter.NewSqliteClient("gift.db")
	defer close()

	StartBackgroundDBUpdater(client)
	controller := controller.APIController{Env: EnvConfig, SqlClient: client}

	r := mux.NewRouter()
	r.HandleFunc("/getPaymentQR", controller.GetPayment).Methods("POST")
	r.HandleFunc("/getPaymentDB/{id}", controller.GetPaymentFromDatabase).Methods("POST")
	r.HandleFunc("/health-check", HealthHandler).Methods("POST")
	http.Handle("/", r)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // allow all origins — for production, restrict this
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{
		Addr:         EnvConfig.URL,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors(r),
	}

	log.Printf("WebService started at %s", EnvConfig.URL)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(response)
}

func StartBackgroundDBUpdater(sqlClient *adapter.SqliteClient) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			<-ticker.C               // espera 10 minutos
			UpdateQRcodes(sqlClient) // executa a função de update
		}
	}()
}

func UpdateQRcodes(sqlClient *adapter.SqliteClient) error {
	log.Println("starting qr update process..")
	gifts, err := sqlClient.ListGifts()
	if err != nil {
		return err
	}
	cfg, err := config.New(EnvConfig.AccessTokenMeli)
	if err != nil {
		return err
	}
	for _, gift := range gifts {
		client := payment.NewClient(cfg)
		log.Println(gift.Price)
		qrCode, err := adapter.GeneratePIXQRCode(client, "brunocebrsilva@gmail.com", gift.Price, gift.Name)
		if err != nil {
			return err
		}

		gift.QRCode = qrCode
		err = sqlClient.UpdateGiftQRCode(gift.ID, gift.QRCode)
		if err != nil {
			return err
		}
	}
	log.Printf("Update total of %d gifts successfully", len(gifts))
	return nil
}
