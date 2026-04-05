#!/usr/bin/env bash

set -euo pipefail

API_BASE_URL="${API_BASE_URL:-${1:-http://localhost:9000/api/v1}}"
ADMIN_USERNAME="${ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-password}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUTPUT_DIR="${BACKEND_DIR}/test_artifacts/report_smoke"
mkdir -p "${OUTPUT_DIR}"

require_bin() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

require_bin curl
require_bin jq

RESPONSE_STATUS=""
RESPONSE_BODY=""

request() {
  local method="$1"
  local url="$2"
  local token="${3:-}"
  local body="${4:-}"

  local tmp
  tmp="$(mktemp)"
  local status

  if [[ -n "${body}" ]]; then
    if [[ -n "${token}" ]]; then
      status="$(curl -sS -o "${tmp}" -w "%{http_code}" \
        -X "${method}" "${url}" \
        -H "Accept: application/json" \
        -H "Authorization: Bearer ${token}" \
        -H "Content-Type: application/json" \
        --data "${body}")"
    else
      status="$(curl -sS -o "${tmp}" -w "%{http_code}" \
        -X "${method}" "${url}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        --data "${body}")"
    fi
  else
    if [[ -n "${token}" ]]; then
      status="$(curl -sS -o "${tmp}" -w "%{http_code}" \
        -X "${method}" "${url}" \
        -H "Accept: application/json" \
        -H "Authorization: Bearer ${token}")"
    else
      status="$(curl -sS -o "${tmp}" -w "%{http_code}" \
        -X "${method}" "${url}" \
        -H "Accept: application/json")"
    fi
  fi

  RESPONSE_STATUS="${status}"
  RESPONSE_BODY="$(cat "${tmp}")"
  rm -f "${tmp}"
}

assert_status() {
  local expected="$1"
  local context="$2"
  if [[ "${RESPONSE_STATUS}" != "${expected}" ]]; then
    echo "Unexpected status for ${context}: got ${RESPONSE_STATUS}, expected ${expected}" >&2
    echo "${RESPONSE_BODY}" >&2
    exit 1
  fi
}

save_json() {
  local name="$1"
  local body="$2"
  echo "${body}" | jq '.' > "${OUTPUT_DIR}/${name}.json"
  echo "Saved ${OUTPUT_DIR}/${name}.json"
}

extract_json() {
  local expr="$1"
  echo "${RESPONSE_BODY}" | jq -r "${expr}"
}

echo "Logging in as ${ADMIN_USERNAME} via ${API_BASE_URL}"
request POST "${API_BASE_URL}/auth/login" "" "$(jq -n \
  --arg username "${ADMIN_USERNAME}" \
  --arg password "${ADMIN_PASSWORD}" \
  '{username:$username,password:$password}')"
assert_status 200 "admin login"
ADMIN_TOKEN="$(extract_json '.token')"

SUFFIX="$(date +%s)"
WORKER_USERNAME="worker_smoke_${SUFFIX}"
CLIENT_USERNAME="client_smoke_${SUFFIX}"
PASSWORD="password123"

request POST "${API_BASE_URL}/salons" "${ADMIN_TOKEN}" "$(jq -n \
  --arg name "Smoke Salon ${SUFFIX}" \
  --arg address "Пермь, Тестовая 1" \
  --arg location "58.0141, 56.2230" \
  --arg working_hours '{"mon":"10:00-20:00","tue":"10:00-20:00"}' \
  '{name:$name,address:$address,location:$location,working_hours:$working_hours,status:"OPEN",max_staff:10,base_hourly_rate:1500}')"
assert_status 201 "salon create"
SALON_ID="$(extract_json '.id')"
SALON_BODY="${RESPONSE_BODY}"

request POST "${API_BASE_URL}/services" "${ADMIN_TOKEN}" "$(jq -n \
  --arg name "Smoke Service ${SUFFIX}" \
  '{name:$name,description:"Smoke scenario service",price:1800,duration_minutes:75}')"
assert_status 201 "service create"
SERVICE_ID="$(extract_json '.id')"
SERVICE_BODY="${RESPONSE_BODY}"

request POST "${API_BASE_URL}/materials" "${ADMIN_TOKEN}" "$(jq -n \
  --arg name "Smoke Material ${SUFFIX}" \
  '{name:$name,unit:"ml"}')"
assert_status 201 "material create"
MATERIAL_ID="$(extract_json '.id')"
MATERIAL_BODY="${RESPONSE_BODY}"

request PUT "${API_BASE_URL}/materials/service/${SERVICE_ID}" "${ADMIN_TOKEN}" "$(jq -n \
  --argjson material_id "${MATERIAL_ID}" \
  '[{material_id:$material_id,quantity_per_use:1}]')"
assert_status 200 "service materials update"

request POST "${API_BASE_URL}/material-expenses" "${ADMIN_TOKEN}" "$(jq -n \
  --argjson material_id "${MATERIAL_ID}" \
  --argjson salon_id "${SALON_ID}" \
  '{material_id:$material_id,salon_id:$salon_id,purchase_price:350,quantity:30}')"
assert_status 201 "material expense create"
MATERIAL_EXPENSE_ID="$(extract_json '.id')"
MATERIAL_EXPENSE_BODY="${RESPONSE_BODY}"

request POST "${API_BASE_URL}/employees" "${ADMIN_TOKEN}" "$(jq -n \
  --arg username "${WORKER_USERNAME}" \
  --arg password "${PASSWORD}" \
  --arg full_name "Smoke Worker" \
  --arg phone "+79990000002" \
  --arg email "${WORKER_USERNAME}@example.com" \
  --argjson salon_id "${SALON_ID}" \
  '{username:$username,password:$password,full_name:$full_name,phone:$phone,email:$email,role:"ADVANCED_MASTER",specialization:"Fade",expected_salary:85000,work_schedule:"{\"mon\":\"10:00-19:00\"}",salon_id:$salon_id}')"
assert_status 201 "employee hire"
WORKER_BODY="${RESPONSE_BODY}"

request POST "${API_BASE_URL}/auth/login" "" "$(jq -n \
  --arg username "${WORKER_USERNAME}" \
  --arg password "${PASSWORD}" \
  '{username:$username,password:$password}')"
assert_status 200 "worker login"
WORKER_TOKEN="$(extract_json '.token')"

request POST "${API_BASE_URL}/auth/register" "" "$(jq -n \
  --arg username "${CLIENT_USERNAME}" \
  --arg password "${PASSWORD}" \
  --arg full_name "Smoke Client" \
  --arg phone "+79990000003" \
  --arg email "${CLIENT_USERNAME}@example.com" \
  '{username:$username,password:$password,full_name:$full_name,phone:$phone,email:$email,role:"CLIENT"}')"
assert_status 201 "client register"

request POST "${API_BASE_URL}/auth/login" "" "$(jq -n \
  --arg username "${CLIENT_USERNAME}" \
  --arg password "${PASSWORD}" \
  '{username:$username,password:$password}')"
assert_status 200 "client login"
CLIENT_TOKEN="$(extract_json '.token')"

CONFIRMED_START="$(date -u -d '+1 day' '+%Y-%m-%dT%H:%M:%SZ')"
CANCELLED_START="$(date -u -d '+2 day' '+%Y-%m-%dT%H:%M:%SZ')"
PENDING_PAST_START="$(date -u -d '-2 day' '+%Y-%m-%dT%H:%M:%SZ')"

create_booking() {
  local start_time="$1"
  local notes="$2"
  request POST "${API_BASE_URL}/bookings" "${CLIENT_TOKEN}" "$(jq -n \
    --arg start_time "${start_time}" \
    --arg notes "${notes}" \
    --argjson salon_id "${SALON_ID}" \
    --argjson service_id "${SERVICE_ID}" \
    '{start_time:$start_time,salon_id:$salon_id,service_ids:[$service_id],notes:$notes}')"
  assert_status 201 "booking create"
}

create_booking "${CONFIRMED_START}" "Подтверждённое бронирование для отчетов"
BOOKING_CONFIRMED_ID="$(extract_json '.id')"

create_booking "${CANCELLED_START}" "Клиент отменил запись"
BOOKING_CANCELLED_ID="$(extract_json '.id')"

create_booking "${PENDING_PAST_START}" "Клиент не пришел"
BOOKING_PENDING_ID="$(extract_json '.id')"

request POST "${API_BASE_URL}/bookings/${BOOKING_CONFIRMED_ID}/confirm" "${WORKER_TOKEN}"
assert_status 200 "booking confirm"

request POST "${API_BASE_URL}/bookings/${BOOKING_CANCELLED_ID}/cancel" "${CLIENT_TOKEN}"
assert_status 200 "booking cancel"

create_payment() {
  local booking_id="$1"
  local status="$2"
  local ext_id="$3"
  request POST "${API_BASE_URL}/payments" "${ADMIN_TOKEN}" "$(jq -n \
    --argjson booking_id "${booking_id}" \
    --arg status "${status}" \
    --arg external_transaction_id "${ext_id}" \
    '{booking_id:$booking_id,amount:1800,status:$status,external_transaction_id:$external_transaction_id}')"
  assert_status 201 "payment create ${status}"
  echo "${RESPONSE_BODY}"
}

PAYMENT_SUCCESS_BODY="$(create_payment "${BOOKING_CONFIRMED_ID}" "SUCCESS" "txn_success_${SUFFIX}")"
PAYMENT_REFUNDED_BODY="$(create_payment "${BOOKING_CANCELLED_ID}" "REFUNDED" "txn_refunded_${SUFFIX}")"
PAYMENT_PENDING_BODY="$(create_payment "${BOOKING_PENDING_ID}" "PENDING" "txn_pending_${SUFFIX}")"

request POST "${API_BASE_URL}/reviews" "${CLIENT_TOKEN}" "$(jq -n \
  '{text:"Smoke review for reports",rating:5}')"
assert_status 201 "review create"
REVIEW_BODY="${RESPONSE_BODY}"

FROM_DATE="$(date -u -d '-7 day' '+%Y-%m-%d')"
TO_DATE="$(date -u -d '+7 day' '+%Y-%m-%d')"

declare -A REPORTS=(
  ["2_2_1_employees"]="${API_BASE_URL}/reports/employees"
  ["2_2_2_salon_activity"]="${API_BASE_URL}/reports/salon-activity?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_3_service_popularity"]="${API_BASE_URL}/reports/service-popularity?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_4_master_activity"]="${API_BASE_URL}/reports/master-activity?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_5_reviews"]="${API_BASE_URL}/reports/reviews?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_6_inventory"]="${API_BASE_URL}/reports/inventory-movement?from=${FROM_DATE}&to=${TO_DATE}&salon_id=${SALON_ID}"
  ["2_2_7_client_loyalty"]="${API_BASE_URL}/reports/client-loyalty?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_8_cancelled"]="${API_BASE_URL}/reports/cancelled-bookings?from=${FROM_DATE}&to=${TO_DATE}"
  ["2_2_9_financial"]="${API_BASE_URL}/reports/financial-summary?from=${FROM_DATE}&to=${TO_DATE}&salon_id=${SALON_ID}"
)

for name in "${!REPORTS[@]}"; do
  request GET "${REPORTS[$name]}" "${ADMIN_TOKEN}"
  assert_status 200 "report ${name}"
  save_json "${name}" "${RESPONSE_BODY}"
done

request POST "${API_BASE_URL}/services/${SERVICE_ID}/use" "${ADMIN_TOKEN}" "$(jq -n \
  --argjson salon_id "${SALON_ID}" \
  '{salon_id:$salon_id,quantity:1}')"
assert_status 200 "service use"
SERVICE_USE_BODY="${RESPONSE_BODY}"

save_json "smoke_entities" "$(jq -n \
  --arg admin_username "${ADMIN_USERNAME}" \
  --arg worker_username "${WORKER_USERNAME}" \
  --arg client_username "${CLIENT_USERNAME}" \
  --argjson salon "${SALON_BODY}" \
  --argjson service "${SERVICE_BODY}" \
  --argjson material "${MATERIAL_BODY}" \
  --argjson material_expense "${MATERIAL_EXPENSE_BODY}" \
  --argjson worker_profile "${WORKER_BODY}" \
  --argjson payment_success "${PAYMENT_SUCCESS_BODY}" \
  --argjson payment_refunded "${PAYMENT_REFUNDED_BODY}" \
  --argjson payment_pending "${PAYMENT_PENDING_BODY}" \
  --argjson review "${REVIEW_BODY}" \
  --argjson service_use "${SERVICE_USE_BODY}" \
  '{admin_username:$admin_username,worker_username:$worker_username,client_username:$client_username,salon:$salon,service:$service,material:$material,material_expense:$material_expense,worker_profile:$worker_profile,payment_success:$payment_success,payment_refunded:$payment_refunded,payment_pending:$payment_pending,review:$review,service_use:$service_use}')"

echo "Smoke scenario completed. Outputs saved to ${OUTPUT_DIR}"
