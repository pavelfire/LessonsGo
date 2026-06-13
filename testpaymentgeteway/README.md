# Payment Gateway

A Go API that lets merchants process card payments and retrieve payment details. Includes a separate bank simulator that can be swapped for a real acquiring bank in production.

## Requirements

- Go 1.22+
- Docker (optional)

## Project layout

```
cmd/gateway/                 payment gateway HTTP API
cmd/banksim/                 acquiring bank simulator
internal/domain/payment/     payment model, validation, service
internal/bank/               HTTP client for bank simulator
internal/banksim/            bank processing rules
internal/banksimapi/         bank simulator HTTP handlers
internal/api/                gateway HTTP handlers
internal/storage/memory/     in-memory payment repository
```

## Run locally

Start the bank simulator and gateway in separate terminals:

```bash
make run-banksim
make run-gateway
```

Or with `go run`:

```bash
go run ./cmd/banksim -addr :8081
go run ./cmd/gateway -addr :8080 -bank-url http://localhost:8081
```

### Process a payment

```bash
curl -s -X POST http://localhost:8080/payments \
  -H 'Content-Type: application/json' \
  -d '{
    "card_number": "4111111111111111",
    "expiry_month": 12,
    "expiry_year": 2027,
    "cvv": "123",
    "amount": 1500,
    "currency": "USD"
  }' | jq
```

Example response:

```json
{
  "id": "pay_a1b2c3d4e5f6g7h8",
  "status": "succeeded",
  "bank_reference": "bank_...",
  "amount": 1500,
  "currency": "USD",
  "card_masked": "**** **** **** 1111",
  "expiry_month": 12,
  "expiry_year": 2027,
  "created_at": "2026-06-13T12:00:00Z"
}
```

### Retrieve a payment

```bash
curl -s http://localhost:8080/payments/<payment_id> | jq
```

## Bank simulator rules

The simulator declines a payment when:

- amount exceeds `100000` minor units ($1000.00)
- card number ends with `0000`
- card expiry is in the past

Otherwise it returns `succeeded` with a unique `bank_` reference id.

Direct bank API:

```bash
curl -s -X POST http://localhost:8081/process \
  -H 'Content-Type: application/json' \
  -d '{
    "card_number": "4111111111111111",
    "expiry_month": 12,
    "expiry_year": 2027,
    "cvv": "123",
    "amount": 1500,
    "currency": "USD"
  }' | jq
```

## Docker

```bash
make docker-up
```

Gateway: `http://localhost:8080`  
Bank simulator: `http://localhost:8081`

## Test

```bash
make test
```

## Assumptions

- Amounts are expressed in minor currency units (cents/pence).
- Supported currencies: `USD`, `EUR`, `GBP`.
- No merchant authentication is implemented for this exercise.
- Full card numbers and CVV are never persisted; only masked card data is stored.
- Payment storage is in-memory and is lost on restart.

## Areas for improvement

- Persist payments in PostgreSQL with encryption at rest.
- Add merchant API keys and request authentication.
- Support idempotency keys on `POST /payments`.
- Emit webhooks to merchants on payment status changes.
- Add retry with backoff when calling the acquiring bank.
- Tokenize card data via a PCI-compliant vault instead of handling PAN/CVV directly.

## Cloud deployment

For production I would deploy:

- **Gateway**: AWS ECS Fargate or GCP Cloud Run behind an API Gateway / load balancer for auto-scaling and TLS termination.
- **Database**: Amazon RDS PostgreSQL or Cloud SQL for durable payment records and reconciliation.
- **Secrets**: AWS Secrets Manager or GCP Secret Manager for bank API credentials.
- **Observability**: structured logs to CloudWatch / Cloud Logging, metrics via Prometheus, distributed tracing with OpenTelemetry.

This keeps the stateless gateway horizontally scalable while the bank integration remains behind a swappable `BankClient` interface.

## Extra

- Separate bank simulator service with its own HTTP API.
- Graceful shutdown for both services.
- Docker Compose setup for local multi-service runs.
- Structured logging with `log/slog`.
- Unit and HTTP integration tests.
