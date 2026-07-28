#!/bin/bash

# Configuration
PORT=8080
API_URL="http://localhost:${PORT}/api"

echo "======================================================"
echo "Finduo AI API Verification Script"
echo "======================================================"

# Health Check
echo -e "\n1. Checking Server Health..."
curl -s -X GET "http://localhost:${PORT}/health"
echo ""

# Create Users
echo -e "\n2. Creating Users..."
USER_1_RESP=$(curl -s -X POST "${API_URL}/users" \
  -H "Content-Type: application/json" \
  -d '{"name": "Rafael", "salary": 6000.00}')
echo "User 1 (Rafael): ${USER_1_RESP}"
USER_1_ID=$(echo "${USER_1_RESP}" | grep -o '"id":"[^"]*' | head -n1 | cut -d'"' -f4)

USER_2_RESP=$(curl -s -X POST "${API_URL}/users" \
  -H "Content-Type: application/json" \
  -d '{"name": "Partner", "salary": 4000.00}')
echo "User 2 (Partner): ${USER_2_RESP}"
USER_2_ID=$(echo "${USER_2_RESP}" | grep -o '"id":"[^"]*' | head -n1 | cut -d'"' -f4)

if [ -z "${USER_1_ID}" ] || [ -z "${USER_2_ID}" ]; then
  echo "Error: Failed to create users and parse IDs. Make sure server is running and database is active."
  exit 1
fi

echo "User IDs: Rafael=${USER_1_ID}, Partner=${USER_2_ID}"

# List Users
echo -e "\n3. Listing Users..."
curl -s -X GET "${API_URL}/users"
echo ""

# Create Expenses
echo -e "\n4. Adding Expenses..."

EXP_1_RESP=$(curl -s -X POST "${API_URL}/expenses" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Supermarket\",
    \"amount\": 500.00,
    \"date\": \"2026-06-10\",
    \"category\": \"Mercado\",
    \"payer_id\": \"${USER_1_ID}\",
    \"is_shared\": true
  }")
echo "Expense 1 (Supermarket): ${EXP_1_RESP}"
EXP_1_ID=$(echo "${EXP_1_RESP}" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

EXP_2_RESP=$(curl -s -X POST "${API_URL}/expenses" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Dog Food\",
    \"amount\": 200.00,
    \"date\": \"2026-06-12\",
    \"category\": \"Dog\",
    \"payer_id\": \"${USER_1_ID}\",
    \"is_shared\": true
  }")
echo "Expense 2 (Dog Food): ${EXP_2_RESP}"

EXP_3_RESP=$(curl -s -X POST "${API_URL}/expenses" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Rent\",
    \"amount\": 300.00,
    \"date\": \"2026-06-15\",
    \"category\": \"Casa\",
    \"payer_id\": \"${USER_2_ID}\",
    \"is_shared\": true
  }")
echo "Expense 3 (Rent): ${EXP_3_RESP}"

EXP_4_RESP=$(curl -s -X POST "${API_URL}/expenses" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Personal Gym Clothes\",
    \"amount\": 150.00,
    \"date\": \"2026-06-15\",
    \"category\": \"Lazer\",
    \"payer_id\": \"${USER_1_ID}\",
    \"is_shared\": false
  }")
echo "Expense 4 (Individual Gym Clothes): ${EXP_4_RESP}"

# Fetch Monthly Summary
echo -e "\n5. Fetching Monthly Summary for 2026-06..."
curl -s -X GET "${API_URL}/summary?month=2026-06"
echo ""

# Update an Expense (e.g. update Rent from 300 to 350)
if [ ! -z "${EXP_1_ID}" ]; then
  echo -e "\n6. Updating Supermarket Expense (Amount 500 -> 550)..."
  curl -s -X PUT "${API_URL}/expenses/${EXP_1_ID}" \
    -H "Content-Type: application/json" \
    -d "{\"amount\": 550.00}"
  echo ""
fi

# Fetch Monthly Summary after update
echo -e "\n7. Fetching Updated Monthly Summary..."
curl -s -X GET "${API_URL}/summary?month=2026-06"
echo ""

# Settle the Month
echo -e "\n8. Settling the Month..."
curl -s -X POST "${API_URL}/summary/settle" \
  -H "Content-Type: application/json" \
  -d "{
    \"year\": 2026,
    \"month\": 6,
    \"is_settled\": true,
    \"settled_by_id\": \"${USER_1_ID}\"
  }"
echo ""

# Fetch Monthly Summary after settlement
echo -e "\n9. Fetching Final Settled Monthly Summary..."
curl -s -X GET "${API_URL}/summary?month=2026-06"
echo ""

echo "======================================================"
echo "API verification sequence completed."
echo "======================================================"
