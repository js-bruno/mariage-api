curl -X POST "http://localhost:8080/getPaymentQR" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TESTe" \
  -d '{
    "value": 100.50,
    "itemDesc": "Produto Exemplo",
    "email": "cliente@example.com"
  }'

curl -X POST "http://localhost:8080/health-check" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TESTe" \
  -d '{
    "value": 100.50,
    "itemDesc": "Produto Exemplo",
    "email": "cliente@example.com"
  }'
