# SSO Service

A Go-based Single Sign-On (SSO) service providing authentication, authorization, and user management using gRPC. This project supports JWT-based authentication, OAuth2 integration, and includes monitoring with Prometheus and Grafana. It is containerized with Docker and orchestrated using Docker Compose.

---

## **Features**
- **Authentication:** Login and register users with secure password hashing (bcrypt).
- **Authorization:** Role-based checks, including `admin`, `common`, `moderator`, etc.
- **OAuth2 Integration:** Login via popular providers (Google, GitHub, etc.).
- **Token Management:** Supports access and refresh tokens.
- **Microservices Architecture:** Auth, Permissions, and Info services.
- **Database Support:** PostgreSQL as the primary database.
- **Caching:** Redis for caching user-related data.
- **Monitoring:** Prometheus and Grafana for observability.

---

## **Technologies Used**
- **Go:** Backend services.
- **PostgreSQL:** Data storage.
- **Redis:** Caching.
- **Prometheus & Grafana:** Monitoring and visualization.
- **Docker & Docker Compose:** Containerization and orchestration.
- **gRPC:** Inter-service communication.
  
Protobuf contract: https://github.com/nikitauty/protos

---

## **Services**

### **1. Auth Service**
Handles user authentication and token generation.
- Endpoints:
    - `Login(email, password)`
    - `RegisterNewUser(email, password)`
    - `OAuthLogin(provider, code)`

### **2. Permissions Service**
Manages user roles and permissions.

### **3. Info Service**
Stores and retrieves user-related metadata.

---

## **Setup**

### **Requirements**
- Docker and Docker Compose installed.

### **Steps**
1. Clone the repository:
   ```bash
   git clone https://github.com/nktgv/sso.git
   cd sso
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. Access the following services:
    - **Auth Service:** [http://localhost:8081](http://localhost:8081)
    - **Permissions Service:** [http://localhost:8082](http://localhost:8082)
    - **Info Service:** [http://localhost:8083](http://localhost:8083)
    - **Prometheus:** [http://localhost:9090](http://localhost:9090)
    - **Grafana:** [http://localhost:3000](http://localhost:3000) (Login: `admin` | Password: `admin`)

---

## **Environment Variables**

| Variable         | Description                       | Default Value   |
|-------------------|-----------------------------------|-----------------|
| `DB_HOST`         | PostgreSQL host                  | `database`      |
| `DB_PORT`         | PostgreSQL port                  | `5432`          |
| `DB_USER`         | PostgreSQL username              | `postgres`      |
| `DB_PASSWORD`     | PostgreSQL password              | `secret`        |
| `DB_NAME`         | Database name                    | `sso`           |
| `JWT_SECRET`      | Secret key for JWT tokens        | `your_jwt_secret` |
| `REDIS_HOST`      | Redis host                       | `redis`         |
| `REDIS_PORT`      | Redis port                       | `6379`          |

---

## **Monitoring**
Prometheus scrapes metrics from all services. Grafana provides a dashboard for visualization.

### **Adding Custom Metrics**
To add custom metrics:
1. Use the `prometheus/client_golang` library in your Go code.
2. Register metrics in the respective service.
3. Update the `prometheus.yml` configuration file.

---

## **API Documentation**
gRPC API documentation is available in the `proto` directory.

### **Compiling Protobuf Files**
1. Install `protoc` and `protoc-gen-go`:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```
2. Generate `.pb.go` files:
   ```bash
   protoc -I=proto --go_out=. --go-grpc_out=. proto/*.proto
   ```

---

## **Future Enhancements**
- Add Kubernetes deployment with Helm charts.
- Add full worked Permissions and User info services.
- Implement rate limiting for enhanced security.
- Improve OAuth provider support.
- Enable horizontal scaling of services.

---

## **Contributing**
We welcome contributions! Please follow these steps:
1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Open a pull request.

---

## **License**
This project is licensed under the MIT License.