# Distributed Commerce System

A production-grade microservices e-commerce platform implementing **distributed transaction patterns**, **event-driven architecture**, and **resilience engineering** principles.  
The system handles order processing, inventory management, payment processing, and notifications across multiple services with eventual consistency guarantees.

---

## 📦 Services

| Service | Port | Responsibility |
|---------|------|----------------|
| **Order Service** | 8081 | Order creation, idempotency, state management |
| **Inventory Service** | 8082 | Stock management, reservation, release |
| **Payment Service** | 8083 | Payment processing, webhook handling, retry logic |
| **Notification Service** | 8084 | Email notifications, DLQ handling |
| **Saga Orchestrator** | - | Distributed transaction coordination |
| **Kafka** | 9092 | Event bus for asynchronous communication |

---

## 🔄 Core Workflow

The following diagram illustrates the end‑to‑end flow of a successful order:
CLIENT
│
│
POST /orders
│
▼
═══════════════════════════════════════════════════════════════════════════════
ORDER SERVICE
═══════════════════════════════════════════════════════════════════════════════

Validate Request
│
▼
Check Idempotency Key
│
┌─────────┴─────────┐
│ │
Already Exists New Request
│ │
▼ ▼
Return Existing Begin Transaction
│
▼
Save Order (CREATED)
│
▼
Save ORDER_CREATED to Outbox
│
▼
Save Idempotency Record
│
▼
COMMIT TX
│
▼
═══════════════════════════════════════════════════════════════════════════════
OUTBOX RELAY
═══════════════════════════════════════════════════════════════════════════════

Poll Pending Events
│
▼
Publish ORDER_CREATED
│
▼
Mark Outbox Published
│
▼
═══════════════════════════════════════════════════════════════════════════════
KAFKA
═══════════════════════════════════════════════════════════════════════════════

order-events
│
▼

═══════════════════════════════════════════════════════════════════════════════
SAGA ORCHESTRATOR
═══════════════════════════════════════════════════════════════════════════════

Receives ORDER_CREATED

│
▼

Create Saga Record

Status = RUNNING

Inventory = PENDING

Payment = PENDING

│
▼

Publish Command

RESERVE_INVENTORY

│
▼

═══════════════════════════════════════════════════════════════════════════════
saga-commands
═══════════════════════════════════════════════════════════════════════════════

RESERVE_INVENTORY
│
▼

═══════════════════════════════════════════════════════════════════════════════
INVENTORY SERVICE
═══════════════════════════════════════════════════════════════════════════════

Receive RESERVE_INVENTORY

│
▼

Begin Transaction

│
▼

Find Inventory

│
▼

Enough Stock ?
│
┌────┴────┐
│ │
YES NO
│ │
│ ▼
│ Build
│ INVENTORY_RESERVATION_FAILED
│ │
│ ▼
│ Save Failed Event
│ into Outbox
│ │
│ Commit
│ │
│ ▼
│ Relay
│ │
│ ▼
│ Publish INVENTORY_RESERVATION_FAILED
│
│
▼

Reserve Stock

↓

Update Inventory

↓

Create Reservation

↓

Save INVENTORY_RESERVED
into Outbox

↓

Commit

↓

Relay publishes

INVENTORY_RESERVED

═══════════════════════════════════════════════════════════════════════════════
inventory-events
═══════════════════════════════════════════════════════════════════════════════

INVENTORY_RESERVED
│
▼

═══════════════════════════════════════════════════════════════════════════════
SAGA ORCHESTRATOR
═══════════════════════════════════════════════════════════════════════════════

Update Saga

Inventory = COMPLETED

│
▼

Publish

PROCESS_PAYMENT

│
▼

═══════════════════════════════════════════════════════════════════════════════
saga-commands
═══════════════════════════════════════════════════════════════════════════════

PROCESS_PAYMENT
│
▼

═══════════════════════════════════════════════════════════════════════════════
PAYMENT SERVICE
═══════════════════════════════════════════════════════════════════════════════

Receive PROCESS_PAYMENT

│
▼

Create Payment

│
▼

Create Provider Intent

│
▼

Simulator Gateway

│
▼

Capture Payment

│
▼

Returns Immediately

│
▼

(10 sec later...)

Webhook arrives

│
▼

Verify Signature

│
▼

Update Payment Status

│
▼

Save Event to Outbox

│
▼

Payment Success ?
│
┌────┴────┐
│ │
YES NO
│ │
▼ ▼

PAYMENT_SUCCEEDED PAYMENT_FAILED

│ │
▼ ▼

Outbox Relay Outbox Relay

│ │
▼ ▼

payment-events payment-events

═══════════════════════════════════════════════════════════════════════════════
PAYMENT_SUCCEEDED PATH
═══════════════════════════════════════════════════════════════════════════════

Saga receives

PAYMENT_SUCCEEDED

│
▼

Update Saga

Payment = COMPLETED

Saga = COMPLETED

│
▼

Publish

COMPLETE_ORDER

│
▼

═══════════════════════════════════════════════════════════════════════════════
saga-commands
═══════════════════════════════════════════════════════════════════════════════

COMPLETE_ORDER
│
▼

═══════════════════════════════════════════════════════════════════════════════
ORDER SERVICE
═══════════════════════════════════════════════════════════════════════════════

Receive COMPLETE_ORDER

│
▼

Find Order

│
▼

CREATED

↓

PAID

│
▼

DONE

═══════════════════════════════════════════════════════════════════════════════
INVENTORY FAILURE COMPENSATION
═══════════════════════════════════════════════════════════════════════════════

Inventory

↓

INVENTORY_RESERVATION_FAILED

↓

Saga

↓

Inventory = FAILED

Saga = FAILED

↓

Publish

CANCEL_ORDER

↓

Order Service

↓

Order

CREATED

↓

CANCELLED

↓

END

═══════════════════════════════════════════════════════════════════════════════
PAYMENT FAILURE COMPENSATION
═══════════════════════════════════════════════════════════════════════════════

Payment

↓

PAYMENT_FAILED

↓

Saga

↓

Payment = FAILED

Saga = FAILED

↓

Publish

RELEASE_INVENTORY

↓

Inventory Service

↓

Find Reservation

↓

Release Reserved Quantity

↓

Reservation

RELEASED

↓

Save

INVENTORY_RELEASED

↓

Outbox Relay

↓

Publish INVENTORY_RELEASED

↓

Saga

↓

Inventory = COMPENSATED

↓

Publish

CANCEL_ORDER

↓

Order Service

↓

CREATED

↓

CANCELLED

↓

END

text

---

## 🎯 Distributed Patterns Implemented

### 1. SAGA Orchestration Pattern

![Saga Orchestrator](Screenshot%202026-07-14%20at%205.13.06%20PM.png)

The **Saga Orchestrator** coordinates distributed transactions across multiple services:

- **Centralized Orchestration** – Order Service acts as the orchestrator.
- **Sequence Control** – Manages step‑by‑step execution (Reserve Inventory → Process Payment → Complete Order).
- **Compensation Logic** – Automatic rollback via compensating transactions:
  - Inventory failure → Cancel Order.
  - Payment failure → Release Inventory → Cancel Order.
- **State Management** – Tracks each step's status (`PENDING` → `COMPLETED` / `FAILED`).

---

### 2. Outbox Pattern

![Outbox Pattern](Screenshot%202026-07-10%20at%2.54.45%20PM.png)

**Atomic dual‑write problem** solved with the Outbox Pattern:

- **Transaction Atomicity** – Events are saved within the same database transaction as business data.
- **Background Relay** – Polls pending events every 2 seconds using `SKIP LOCKED` for concurrent safety.
- **At‑Least‑Once Delivery** – Events are marked as published only after Kafka acknowledgment.
- **Idempotency Key** – Prevents duplicate event processing.

```go
// Outbox flow
BEGIN TX
  Save Business Data
  Save Outbox Event (PENDING)
  Save Idempotency Key
COMMIT TX

// Background Relay
SELECT * FROM outbox WHERE status='PENDING' FOR UPDATE SKIP LOCKED
  Publish to Kafka
  UPDATE status='PUBLISHED'
3. Flow Control & Resilience
https://Screenshot%25202026-06-28%2520at%25203.10.04%2520PM.png

Multi‑layer defensive architecture:

Rate Limiting – Token bucket algorithm prevents request flooding.

Worker Queue – Bounded queue with backpressure mechanism.

Load Shedding – Priority‑based admission control when overloaded:

Reject low‑priority requests.

Critical requests continue.

Circuit Breakers – Protect downstream dependencies.

4. Circuit Breaker Pattern
https://Screenshot%25202026-06-28%2520at%25203.56.08%2520PM.png

Three‑state circuit breaker implementation:

text
CLOSED → OPEN → HALF-OPEN → CLOSED
CLOSED – Normal operation, tracking failures.

OPEN – Rejects requests immediately (fast‑fail) after threshold exceeded.

HALF‑OPEN – Allows test requests to check recovery.

Real implementation logs:

log
2026/07/06 07:45:44 [CircuitBreaker] CLOSED -> OPEN
2026/07/06 07:45:44 [CircuitBreaker] Transition CLOSED -> OPEN
5. Bulkhead Pattern
https://Screenshot%25202026-06-28%2520at%25203.56.08%2520PM.png

Resource isolation for service protection:

text
Payment Worker Pool (100 workers)
Inventory Worker Pool (100 workers)
Notification Worker Pool (100 workers)
Failure Isolation – Payment exhaustion does not affect Inventory/Notification.

Thread Pool Limits – Configurable per service.

Cached Responses – Degraded but useful responses during failures.

6. Idempotency Pattern
https://Screenshot%25202026-06-28%2520at%25203.56.08%2520PM.png

Duplicate request prevention:

Idempotency Key – Client‑provided UUID for request deduplication.

Inbox Pattern – Tracks processed events to prevent duplicate processing.

Record Existing Returns – Already processed requests return cached results.

7. Distributed Concurrency Control
https://Screenshot%25202026-06-30%2520at%252010.44.59%2520AM.png

Strategy	Use Case	Implementation
Optimistic Locking	High‑read, low‑write	Version numbers, timestamps
Pessimistic Locking	High‑contention	SELECT FOR UPDATE
Distributed Lock	Cross‑service coordination	Redis locks with fencing tokens
Fencing Token	Prevent stale operations	Monotonically increasing tokens
🔄 Event‑Driven Architecture
Kafka Event Flow
https://Screenshot%25202026-07-09%2520at%252011.52.38%E2%80%AFPM.png

Event Structure:

json
{
  "metadata": {
    "event_id": "uuid",
    "event_type": "ORDER_CREATED",
    "aggregate_id": "order-123",
    "source": "order-service",
    "timestamp": "2026-07-09T23:52:38Z"
  },
  "payload": {
    "order_id": "123",
    "customer_id": "456",
    "items": [...],
    "total": 99.99
  }
}
Consumer Group Setup:

Notification Service – Listens to order-events topic.

Saga Orchestrator – Subscribes to multiple event types.

Idempotent Consumers – Inbox pattern prevents duplicates.

Retry & Dead Letter Queue (DLQ)
https://Screenshot%25202026-07-07%2520at%25207.19.06%E2%80%AFPM.png

Exponential Backoff with Circuit Breaker:

log
2026/07/07 16:59:55 retry attempt 1 failed: temporary email provider failure sleeping for 1.437498355s
2026/07/07 16:59:56 retry attempt 2 failed: temporary email provider failure sleeping for 5.946782779s
2026/07/07 17:00:01 retry attempt 3 failed: temporary email provider failure sleeping for 10.259526923s
2026/07/07 17:00:11 retry attempt 4 failed: temporary email provider failure sleeping for 11.946052988s
2026/07/07 17:00:26 [CircuitBreaker] CLOSED -> OPEN
2026/07/07 17:00:26 POST /notifications 201 29.831144291s
DLQ Flow:

Max Retries – Configurable attempts with exponential backoff.

DLQ Storage – Failed messages stored for manual investigation.

Replay Capability – Fix issues and replay from DLQ.

Inbox/Outbox Integration
https://Screenshot%25202026-06-28%2520at%25207.06.11%E2%80%AFPM.png

End‑to‑end exactly‑once processing:

text
Business Operation
    ↓
Local Database Transaction
    ↓
Business Data + Outbox Record
    ↓
Outbox Publisher → Kafka
    ↓
Consumer Receives Event
    ↓
Inbox Check (Duplicate?)
    ├── Duplicate → Ignore
    └── New → Idempotent Logic
         ↓
    Processing Successful?
    ├── Yes → Commit
    └── No → Retry → DLQ
🚦 Resiliency Patterns
Graceful Degradation
When dependencies fail:

Cached Responses – Return stale but valid data.

Default Values – Sensible defaults for non‑critical fields.

Backup Services – Failover to secondary providers.

Fallback Logic – Execute alternative business flows.

Backpressure Implementation
Reactive backpressure:

Queue Monitoring – Track worker queue depth.

Producer Feedback – Slow down upstream producers when queues fill.

Reject Policy – Configurable rejection strategies:

Reject Low Priority

Reject Randomly

Admission Control

📊 Performance Metrics
https://Screenshot%25202026-07-04%2520at%25208.43.31%E2%80%AFAM.png

Load Test Results (100 concurrent requests):

text
Summary:
  Total: 0.1839 secs
  Slowest: 0.1033 secs
  Fastest: 0.0020 secs
  Average: 0.0333 secs
  Requests/sec: 543.6803

Status code distribution:
  [201] 34 responses
  [400] 66 responses
🏥 Health Checks & Observability
Service Startup Logs
Payment Service:

https://Screenshot%25202026-07-05%2520at%25208.48.51%E2%80%AFPM.png

log
rudransh@Rudranshs-MacBook-Air payment % go run ./cmd/main.go
2026/07/05 20:47:37 Payment Database Connected
[payment] 2026/07/05 20:47:37 app.go:97: Payment Service Started

Fiber v2.52.13
http://127.0.0.1:8083
(bound on host 0.0.0.0 and port 8083)

Handlers ...... 7 Processes ...... 1
Prefork ...... Disabled PID ...... 14405
Notification Service:

https://Screenshot%25202026-07-07%2520at%25205.01.49%E2%80%AFPM.png

log
rudransh@Rudranshs-MacBook-Air notification % go run cmd/main.go
2026/07/07 16:59:45 Database Connected
[notification] 2026/07/07 16:59:45 app.go:100: Notification Service Started

Handlers ...... 12 Processes ...... 1
Prefork ...... Disabled PID ...... 38748
Webhook Processing Example
https://Screenshot%25202026-07-06%2520at%25207.46.35%E2%80%AFAM.png

Complete webhook flow:

log
2026/07/06 07:46:00 1. Webhook request received
2026/07/06 07:46:00 3. JSON parsed
2026/07/06 07:46:00 5. Event mapped
2026/07/06 07:46:00 7. Payload generated
2026/07/06 07:46:00 9. Signature verified
2026/07/06 07:46:00 ===== WEBHOOK RECEIVED ======
2026/07/06 07:46:00 EventID: 3005d028-3e19-47ba-9780-02dc4cc8955a
2026/07/06 07:46:00 1. Duplicate check
2026/07/06 07:46:00 2. Saving webhook event
2026/07/06 07:46:00 3. Finding payment
2026/07/06 07:46:00 4. Transition payment
2026/07/06 07:46:00 5. Create attempt
2026/07/06 07:46:00 Attempt 1 succeeded
2026/07/06 07:46:00 11. Webhook completed
🛠️ Technology Stack
Component	Technology	Purpose
Language	Go 1.21+	High‑performance microservices
Web Framework	Fiber v2.52.13	Fast HTTP routing
Database	PostgreSQL	ACID transactions, outbox storage
Message Queue	Apache Kafka	Event streaming, decoupling
Event Processing	Sarama	Kafka client for Go
Logging	Logrus	Structured logging
Configuration	Viper	Environment‑based config
Circuit Breaker	Custom implementation	Resilience patterns
Distributed Lock	Redis	Concurrency control
🚀 Getting Started
Prerequisites
Go 1.21+

Docker & Docker Compose

PostgreSQL 14+

Kafka 3.0+

Running Services
bash
# Start infrastructure
docker-compose up -d postgres kafka redis

# Run services in separate terminals
cd services/order && go run cmd/main.go
cd services/inventory && go run cmd/main.go
cd services/payment && go run cmd/main.go
cd services/notification && go run cmd/main.go
cd services/saga && go run cmd/main.go
Testing the Flow
bash
# Create an order
curl -X POST http://localhost:8081/orders \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "req-123",
    "customer_id": "cust-456",
    "items": [{"product_id": "prod-789", "quantity": 2}]
  }'

# Check order status
curl http://localhost:8081/orders/{order-id}

# View Kafka events
docker exec -it kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic order-events --from-beginning
🔍 Monitoring & Debugging
Circuit Breaker States
go
// Check breaker status
GET /admin/circuit-breaker/{service}
Retry Statistics
go
// View retry metrics
GET /metrics/retries
DLQ Management
go
// Replay failed messages
POST /admin/dlq/replay
🎯 Key Architectural Decisions
Saga Orchestration over Choreography – Centralized control for complex workflows.

Outbox Pattern over 2PC – Avoids distributed transaction coordinator.

Eventual Consistency – Trades immediate consistency for availability.

Circuit Breakers Everywhere – Protects against cascading failures.

Bulkhead Isolation – Prevents resource exhaustion in one service affecting others.

Idempotency First – Every operation designed to be retry‑safe.

🎓 Lessons Learned
Circuit Breakers are essential – Detected and isolated failures before they cascade.

Outbox pattern with SKIP LOCKED – Maintains performance under high concurrency.

Webhooks need idempotency – External systems retry aggressively.

Compensation is critical – Always design for failure scenarios.

Monitoring logs tell the story – Structured logging is crucial for debugging distributed systems.

