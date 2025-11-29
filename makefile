sent:
	curl -x post https://mariage-api-gtg5.onrender.com/getpaymentqr \
		-h "content-type: application/json" \
		-h "authorization: bearer testetoken" \
		-d '{
			"value": 100.50,
			"item_id": 1,
			"item_desc": "fone de ouvido bluetooth",
			"email": "cliente@exemplo.com"
		}'
