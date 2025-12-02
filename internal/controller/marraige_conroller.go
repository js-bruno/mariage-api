package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/js-bruno/mariage-api/internal/adapter"
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

type APIController struct {
	Env       utils.Env
	SqlClient *adapter.SqliteClient
}

func (a APIController) GetPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Handle preflight requests (OPTIONS method).
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authBearer := r.Header.Get("Authorization")
	err := a.authCheck(authBearer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	var paymentRequest PaymentRequest
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg, err := config.New(a.Env.AccessTokenMeli)
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
func (a APIController) GetPaymentFromDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Handle preflight requests (OPTIONS method).
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authBearer := r.Header.Get("Authorization")
	err := a.authCheck(authBearer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	var paymentRequest PaymentRequest
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gift, err := a.SqlClient.GetGiftByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Succefully returned message: %d, %s", gift.ID, gift.Name)
	json.NewEncoder(w).Encode(&PaymentResponse{
		QRCode: gift.QRCode,
	})
}

func (a APIController) authCheck(authBearer string) error {
	authBearer = strings.TrimPrefix(authBearer, "Bearer ")
	if a.Env.ApiAuthToken != authBearer {
		err := errors.New("unauthorized token")
		return err
	}
	return nil
}
