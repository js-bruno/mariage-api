package adapter

import (
	"context"
	"log"

	"github.com/mercadopago/sdk-go/pkg/payment"
)

func GeneratePIXQRCode(
	client payment.Client,
	email string,
	value float64,
	desc string,
) (string, error) {
	request := payment.Request{
		TransactionAmount: value, // valor total do pagamento
		Description:       desc,  // descrição
		Installments:      1,     // número de parcelas
		PaymentMethodID:   "pix", // método de pagamento (ex: visa, master)
		Payer: &payment.PayerRequest{
			Email: email, // e-mail do pagador
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	qrCode := resource.PointOfInteraction.TransactionData.QRCode
	log.Printf("%s:qrcode successfully generated | %s \n", email, qrCode)

	return qrCode, nil
}
