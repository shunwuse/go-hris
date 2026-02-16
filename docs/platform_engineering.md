# Go HRIS Platform & SRE Roadmap (The "Arrows")

This roadmap focuses on the "Arrows" between services—how the system is deployed, monitored, scaled, and secured as a whole, moving beyond simple API development.

---

## Phase 1: Automation & Software Supply Chain (自動化與供應鏈安全)
**Goal:** Shift from "Manual Startup" to "Automated Delivery."

### 1. Advanced CI/CD Pipeline [Not Started]
*   **Target:** GitHub Actions.
*   **Task:** Not just `go test`, but a pipeline that:
    1. Runs Linters (golangci-lint).
    2. Performs Security Scanning (Trivy) on Docker images.
    3. Automatically tags and pushes to a Container Registry (GHCR/DockerHub).
*   **Architectural Value:** Understanding the lifecycle of an artifact and ensuring code quality at the gate.

### 2. Infrastructure as Code (IaC) - Local/Cloud [Not Started]
*   **Target:** Terraform or Pulumi.
*   **Task:** Define your Redis, MySQL, and NATS setup in code instead of `docker-compose.yaml` (or use the Terraform Docker provider).
*   **Architectural Value:** Learning how to manage environment drift and reproducible infrastructure.

---

## Phase 2: Data Excellence & Reliability (資料卓越與可靠性)
**Goal:** Understanding how data lives and dies in a distributed system.

### 3. Database High Availability (HA) & Scaling
*   **Task:**
    1. Set up **MySQL Master-Slave Replication** in Docker.
    2. Implement **Read/Write Splitting** in your Go service (Writer to Master, Reader to Slave).
    3. Simulate a Master failure and practice manual/automatic promotion.
*   **Architectural Value:** Learning about replication lag, consistency trade-offs, and failover strategies.

### 4. Database Performance Observability
*   **Task:** Integrate Percona Monitoring or similar tools to track "Slow Queries" and "Index Usage" from a system level.
*   **Architectural Value:** Transitioning from "writing SQL" to "managing query performance at scale."

---

## Phase 3: Communication & Distributed Patterns (通訊與分散式模式)
**Goal:** Moving away from synchronous-only communication.

### 5. Event-Driven Architecture (EDA) [Not Started]
*   **Target:** NATS JetStream or RabbitMQ.
*   **Task:** Refactor a business flow (e.g., User Creation) to be asynchronous. The API saves to DB and publishes an event; a separate Worker service handles the welcome email.
*   **Architectural Value:** Decoupling services, handling backpressure, and understanding "Eventual Consistency."

### 6. Service Mesh & Network Traffic Control
*   **Target:** Nginx/Envoy or Istio (if moving to K8s).
*   **Task:** Implement **Global Rate Limiting** and **mTLS** at the network edge, not in the Go code.
*   **Architectural Value:** Understanding how to protect your system from DDoS and how to secure internal traffic.

---

## Phase 4: Full-Stack Observability (全鏈路觀測)
**Goal:** Seeing what is happening in the dark.

### 7. The Three Pillars (The Proper Way)
*   **Metrics:** Prometheus + Grafana dashboards tracking RED (Rate, Errors, Duration) metrics.
*   **Logging:** Centralized log aggregation using **Loki** or **ELK Stack**.
*   **Tracing:** Distributed tracing using **Tempo** or **Jaeger** (OpenTelemetry).
*   **Architectural Value:** Learning to diagnose "it's slow" without guessing.

### 8. Alerting & SLO Management
*   **Task:** Define **SLIs (Service Level Indicators)**. Setup alerts that trigger when the 99th percentile latency > 500ms for more than 5 minutes.
*   **Architectural Value:** Moving from "Is it running?" to "Is it meeting the business commitment?"

---

## Phase 5: Cloud Native & Resilience (雲原生與彈性)
**Goal:** Operating in a professional, orchestrated environment.

### 9. Kubernetes Migration (The Final Boss)
*   **Target:** Kind, Minikube, or K3s.
*   **Task:** Deploy Go-HRIS as a set of K8s Deployments, Services, and ConfigMaps. Implement HPA (Horizontal Pod Autoscaler) based on CPU/Memory.
*   **Architectural Value:** Mastering the industry-standard "Operating System" of the cloud.

### 10. Chaos Engineering (Basic)
*   **Task:** Randomly kill the Redis container or add 2 seconds of latency to the DB network using `pumba` or `toxiproxy`. See how your API server survives.
*   **Architectural Value:** Validating that your "Resilience Wrappers" (Circuit Breakers) actually work in a disaster.
