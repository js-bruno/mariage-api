curl -X POST http://localhost:8080/getPaymentQR \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TESTETOKEN" \
  -d '{
    "value": 100.50,
    "item_id": 1,
    "item_desc": "Fone de ouvido Bluetooth",
    "email": "cliente@exemplo.com"
  }'
