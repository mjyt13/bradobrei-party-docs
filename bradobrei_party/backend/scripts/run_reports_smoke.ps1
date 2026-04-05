param(
    [string]$ApiBaseUrl = "http://localhost:9000/api/v1"
)

$ErrorActionPreference = "Stop"

$backendDir = Split-Path -Parent $PSScriptRoot
$outputDir = Join-Path $backendDir "test_artifacts\report_smoke"
New-Item -ItemType Directory -Force -Path $outputDir | Out-Null

$suffix = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Url,
        [object]$Body = $null,
        [string]$Token = ""
    )

    $headers = @{ Accept = "application/json" }
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }

    if ($null -ne $Body) {
        return Invoke-RestMethod -Method $Method -Uri $Url -Headers $headers -ContentType "application/json" -Body ($Body | ConvertTo-Json -Depth 12)
    }

    return Invoke-RestMethod -Method $Method -Uri $Url -Headers $headers
}

function Save-Json {
    param(
        [string]$Name,
        [object]$Data
    )

    $path = Join-Path $outputDir "$Name.json"
    $Data | ConvertTo-Json -Depth 20 | Set-Content -Path $path -Encoding UTF8
    Write-Host "Saved $path"
}

function Login {
    param(
        [string]$Username,
        [string]$Password
    )

    $resp = Invoke-Api -Method POST -Url "$ApiBaseUrl/auth/login" -Body @{
        username = $Username
        password = $Password
    }

    return $resp.token
}

$adminUsername = "admin_smoke_$suffix"
$workerUsername = "worker_smoke_$suffix"
$clientUsername = "client_smoke_$suffix"
$password = "password123"

Invoke-Api -Method POST -Url "$ApiBaseUrl/auth/register" -Body @{
    username  = $adminUsername
    password  = $password
    full_name = "Smoke Admin"
    phone     = "+79990000001"
    email     = "$adminUsername@example.com"
    role      = "ADMINISTRATOR"
} | Out-Null

$adminToken = Login -Username $adminUsername -Password $password

$salon = Invoke-Api -Method POST -Url "$ApiBaseUrl/salons" -Token $adminToken -Body @{
    name             = "Smoke Salon $suffix"
    address          = "Пермь, Тестовая 1"
    location         = "58.0141, 56.2230"
    working_hours    = '{"mon":"10:00-20:00","tue":"10:00-20:00"}'
    status           = "OPEN"
    max_staff        = 10
    base_hourly_rate = 1500
}

$service = Invoke-Api -Method POST -Url "$ApiBaseUrl/services" -Token $adminToken -Body @{
    name             = "Smoke Service $suffix"
    description      = "Smoke scenario service"
    price            = 1800
    duration_minutes = 75
}

$material = Invoke-Api -Method POST -Url "$ApiBaseUrl/materials" -Token $adminToken -Body @{
    name = "Smoke Material $suffix"
    unit = "ml"
}

Invoke-Api -Method PUT -Url "$ApiBaseUrl/materials/service/$($service.id)" -Token $adminToken -Body @(
    @{
        material_id      = $material.id
        quantity_per_use = 1
    }
) | Out-Null

$employee = Invoke-Api -Method POST -Url "$ApiBaseUrl/employees" -Token $adminToken -Body @{
    username        = $workerUsername
    password        = $password
    full_name       = "Smoke Worker"
    phone           = "+79990000002"
    email           = "$workerUsername@example.com"
    role            = "ADVANCED_MASTER"
    specialization  = "Fade"
    expected_salary = 85000
    work_schedule   = '{"mon":"10:00-19:00"}'
    salon_id        = $salon.id
}

$workerToken = Login -Username $workerUsername -Password $password

Invoke-Api -Method POST -Url "$ApiBaseUrl/auth/register" -Body @{
    username  = $clientUsername
    password  = $password
    full_name = "Smoke Client"
    phone     = "+79990000003"
    email     = "$clientUsername@example.com"
    role      = "CLIENT"
} | Out-Null

$clientToken = Login -Username $clientUsername -Password $password

$now = Get-Date
$confirmedStart = $now.ToUniversalTime().AddDays(1).ToString("o")
$cancelledStart = $now.ToUniversalTime().AddDays(2).ToString("o")
$pendingPastStart = $now.ToUniversalTime().AddDays(-2).ToString("o")

$bookingConfirmed = Invoke-Api -Method POST -Url "$ApiBaseUrl/bookings" -Token $clientToken -Body @{
    start_time  = $confirmedStart
    salon_id    = $salon.id
    service_ids = @($service.id)
    notes       = "Подтверждённое бронирование для отчётов"
}

$bookingCancelled = Invoke-Api -Method POST -Url "$ApiBaseUrl/bookings" -Token $clientToken -Body @{
    start_time  = $cancelledStart
    salon_id    = $salon.id
    service_ids = @($service.id)
    notes       = "Клиент отменил запись"
}

$bookingPendingPast = Invoke-Api -Method POST -Url "$ApiBaseUrl/bookings" -Token $clientToken -Body @{
    start_time  = $pendingPastStart
    salon_id    = $salon.id
    service_ids = @($service.id)
    notes       = "Клиент не пришёл"
}

Invoke-Api -Method POST -Url "$ApiBaseUrl/bookings/$($bookingConfirmed.id)/confirm" -Token $workerToken | Out-Null
Invoke-Api -Method POST -Url "$ApiBaseUrl/bookings/$($bookingCancelled.id)/cancel" -Token $clientToken | Out-Null

$paymentSuccess = Invoke-Api -Method POST -Url "$ApiBaseUrl/payments" -Token $adminToken -Body @{
    booking_id               = $bookingConfirmed.id
    amount                   = 1800
    status                   = "SUCCESS"
    external_transaction_id  = "txn_success_$suffix"
}

$paymentRefunded = Invoke-Api -Method POST -Url "$ApiBaseUrl/payments" -Token $adminToken -Body @{
    booking_id               = $bookingCancelled.id
    amount                   = 1800
    status                   = "REFUNDED"
    external_transaction_id  = "txn_refund_$suffix"
}

$paymentPending = Invoke-Api -Method POST -Url "$ApiBaseUrl/payments" -Token $adminToken -Body @{
    booking_id               = $bookingPendingPast.id
    amount                   = 1800
    status                   = "PENDING"
    external_transaction_id  = "txn_pending_$suffix"
}

$review = Invoke-Api -Method POST -Url "$ApiBaseUrl/reviews" -Token $clientToken -Body @{
    text   = "Smoke review for reports"
    rating = 5
}

$from = (Get-Date).AddDays(-7).ToString("yyyy-MM-dd")
$to = (Get-Date).AddDays(7).ToString("yyyy-MM-dd")

$reports = @{
    "2_2_1_employees"          = "$ApiBaseUrl/reports/employees"
    "2_2_2_salon_activity"     = "$ApiBaseUrl/reports/salon-activity?from=$from&to=$to"
    "2_2_3_service_popularity" = "$ApiBaseUrl/reports/service-popularity?from=$from&to=$to"
    "2_2_4_master_activity"    = "$ApiBaseUrl/reports/master-activity?from=$from&to=$to"
    "2_2_5_reviews"            = "$ApiBaseUrl/reports/reviews?from=$from&to=$to"
    "2_2_6_inventory"          = "$ApiBaseUrl/reports/inventory-movement?from=$from&to=$to&salon_id=$($salon.id)"
    "2_2_7_client_loyalty"     = "$ApiBaseUrl/reports/client-loyalty?from=$from&to=$to"
    "2_2_8_cancelled"          = "$ApiBaseUrl/reports/cancelled-bookings?from=$from&to=$to"
    "2_2_9_financial"          = "$ApiBaseUrl/reports/financial-summary?from=$from&to=$to&salon_id=$($salon.id)"
}

foreach ($entry in $reports.GetEnumerator()) {
    $data = Invoke-Api -Method GET -Url $entry.Value -Token $adminToken
    Save-Json -Name $entry.Key -Data $data
}

Save-Json -Name "smoke_entities" -Data @{
    admin            = $adminUsername
    worker_profile   = $employee
    salon            = $salon
    service          = $service
    material         = $material
    payment_success  = $paymentSuccess
    payment_refunded = $paymentRefunded
    payment_pending  = $paymentPending
    review           = $review
}

Write-Host "Smoke scenario completed. Outputs saved to $outputDir"
