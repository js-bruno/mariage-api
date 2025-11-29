package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/js-bruno/mariage-api/internal/adapter"
	"github.com/js-bruno/mariage-api/internal/repository"
	"github.com/js-bruno/mariage-api/internal/utils"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

	r := mux.NewRouter()
	r.HandleFunc("/getPaymentQR", GetPayment).Methods("POST")
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

func GetPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Handle preflight requests (OPTIONS method).
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authBearer := r.Header.Get("Authorization")
	authBearer = strings.TrimPrefix(authBearer, "Bearer ")

	if EnvConfig.ApiAuthToken != authBearer {
		err := errors.New("unauthorized token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paymentRequest PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := config.New(EnvConfig.AccessTokenMeli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	client := payment.NewClient(cfg)
	qrCode, err := adapter.GeneratePIXQRCode(client, paymentRequest.Email, paymentRequest.Value, paymentRequest.ItemDesc)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(&PaymentResponse{
		QRCode: qrCode,
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(response)
}

func GetGift(ctx context.Context, client *mongo.Client, giftId any) (repository.Gift, error) {
	coll := client.Database(repository.DatabaseName).Collection(repository.CollectionName)
	// mongoID, err := primitive.ObjectIDFromHex("0")
	// if err != nil {
	// 	return repository.Gift{}, err
	// }

	// filter := bson.D{{"id", 0}}

	var gift repository.Gift
	err := coll.FindOne(ctx, bson.M{"_id": giftId}).Decode(&gift)
	if err != nil {
		return repository.Gift{}, err
	}
	return gift, nil
}
