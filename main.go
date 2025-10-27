package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type Env struct {
	AccessTokenMeli string `json:"acess_token"`
	AuthKey         string `json:"auth_key"`
}

type PaymentRequest struct {
	Value    float64 `json:"value"`
	ItemId   int     `json:"item_id"`
	ItemDesc string  `json:"item_desc"`
	Email    string  `json:"email"`
}

type PaymentResponse struct {
	QRCode string `json:"qrcode"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getPaymentQR", GetPayment).Methods("POST")
	http.Handle("/", r)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // allow all origins — for production, restrict this
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors(r),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	env, err := GetEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Handle preflight requests (OPTIONS method).
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	authBearer := r.Header.Get("Authorization")
	authBearer = strings.TrimPrefix(authBearer, "Bearer ")

	if env.AuthKey != authBearer {
		err := errors.New("unauthorized token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var paymentRequest PaymentRequest
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := config.New(env.AccessTokenMeli)
	if err != nil {
		panic(fmt.Sprintf("Erro ao configurar SDK: %v", err))
	}

	client := payment.NewClient(cfg)

	request := payment.Request{
		TransactionAmount: paymentRequest.Value,    // valor total do pagamento
		Description:       paymentRequest.ItemDesc, // descrição
		Installments:      1,                       // número de parcelas
		PaymentMethodID:   "pix",                   // método de pagamento (ex: visa, master)
		Payer: &payment.PayerRequest{
			Email: paymentRequest.Email, // e-mail do pagador
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		panic(fmt.Sprintf("%s:Erro ao criar pagamento | %v", paymentRequest.Email, err))
	}

	qrCode := resource.PointOfInteraction.TransactionData.QRCode
	json.NewEncoder(w).Encode(&PaymentResponse{
		QRCode: qrCode,
	})
	fmt.Printf("%s:qrcode successfully generated | %s", paymentRequest.Email, qrCode)

	// fmt.Printf("Pagamento criado com sucesso!\nID: %v\nStatus: %v\nDetalhes: %+v\n",
	// 	resource.ID, resource.Status, resource)
}

func GetEnv() (*Env, error) {
	env := Env{
		AccessTokenMeli: os.Getenv("ACCESS_TOKEN"),
		AuthKey:         os.Getenv("AUTH_TOKEN"),
	}
	if env.AccessTokenMeli == "" || env.AuthKey == "" {
		fmt.Println(env.AccessTokenMeli)
		fmt.Println(env.AuthKey)
		return nil, errors.New("Envs not find")
	}
	return &env, nil
}
