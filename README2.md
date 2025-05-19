# 💼 Product API - GoLang | AWS ECS Fargate | CI/CD

### 🔗 Live Endpoint

**[http://product-api-alb-449406175.us-east-1.elb.amazonaws.com/api/products](http://product-api-alb-449406175.us-east-1.elb.amazonaws.com/api/products)**

---

## 📁 Project Summary

A fully containerized, scalable **GoLang-based Product API**, deployed using **AWS ECS Fargate**, with CI/CD pipelines managed by **AWS CodePipeline & CodeBuild**, and exposed to the internet via a stable **Application Load Balancer (ALB)**.

---

## ⚙️ Stack Used

| Component      | Technology                   |
| -------------- | ---------------------------- |
| Language       | Go (Golang)                  |
| Web Framework  | Gin                          |
| Container      | Docker                       |
| CI/CD          | AWS CodePipeline + CodeBuild |
| Hosting        | AWS ECS Fargate              |
| Image Registry | Amazon ECR                   |
| Load Balancer  | Application Load Balancer    |
| Monitoring     | AWS CloudWatch Logs          |

---

## 🧱 Step-by-Step Setup

### ✅ 1. Code Structure

```
.
├── cmd/main.go              # Entry point
├── internal/
│   ├── handlers/            # Route handlers
│   ├── models/              # Data models
│   └── utils/               # Utility functions
├── products.json            # Sample dataset
└── main_test.go             # Unit tests
```

---

### ✅ 2. Dockerization

#### 📄 Dockerfile

```Dockerfile
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]
```

---

### ✅ 3. CodeBuild Configuration

#### 📄 buildspec.yml

```yaml
version: 0.2
env:
  variables:
    IMAGE_TAG: "dev"

phases:
  pre_build:
    commands:
      - go test ./...
      - aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_URI
  build:
    commands:
      - docker build -t $ECR_URI:$IMAGE_TAG .
      - docker tag $ECR_URI:$IMAGE_TAG $ECR_URI:latest
      - docker push $ECR_URI:$IMAGE_TAG
      - docker push $ECR_URI:latest
  post_build:
    commands:
      - aws ecs update-service --cluster $ECS_CLUSTER --service $ECS_SERVICE --force-new-deployment
```

---

### ✅ 4. ECR Repository

* Created ECR repository: `products-api-dev`
* Docker images pushed by CodeBuild

---

### ✅ 5. CodePipeline

Pipeline stages:

* **Source**: GitHub
* **Build**: CodeBuild (uses `buildspec.yml`)
* **Deploy**: ECS Fargate (auto-updates ECS Service)

---

### ✅ 6. Task Definition (ECS)

* **Launch type**: `FARGATE`
* **CPU**: `1024`
* **Memory**: `3072`
* **Container name**: `products-api-container`
* **Port Mapping**: `8080`
* **Image**: Pulled from ECR
* **Logs**: `/ecs/products-api` in CloudWatch

---

### ✅ 7. ECS Cluster & Service

* **Cluster**: `products-cluster-dev`
* **Service**: `products-api-service-with-alb`
* **Tasks**: 1 running on Fargate
* **Uses**: latest task definition

---

### ✅ 8. Application Load Balancer

* **Name**: `products-api-alb`
* **Scheme**: Internet-facing
* **Listener**: HTTP on port 80
* **Target Group**: `products-api-tg`
* **Type**: IP
* **Target Port**: 8080
* **Health Check Path**: `/api/health` or `/api/products`

#### 🔐 Security Groups:

* **ALB** allows inbound HTTP (port 80) from `0.0.0.0/0`
* **ECS Task** allows inbound from ALB SG on port `8080`

---

## ✅ Sample API Endpoints

| Method | Endpoint            | Description           |
| ------ | ------------------- | --------------------- |
| GET    | `/api/products`     | Returns all products  |
| GET    | `/api/products/:id` | Returns product by ID |
| GET    | `/api/health`       | Health check for ALB  |

---

## 🧪 Testing

Run tests locally with:

```bash
go test ./...
```

**Output**:
Basic tests validate `/api/products` response structure and availability.

# 💼 Product API - GoLang | AWS ECS Fargate | CI/CD

### 🔗 Live Endpoint

**[http://product-api-alb-449406175.us-east-1.elb.amazonaws.com/api/products](http://product-api-alb-449406175.us-east-1.elb.amazonaws.com/api/products)**

---

## 📁 Project Summary

A fully containerized, scalable **GoLang-based Product API**, deployed using **AWS ECS Fargate**, with CI/CD pipelines managed by **AWS CodePipeline & CodeBuild**, and exposed to the internet via a stable **Application Load Balancer (ALB)**.

---

## ⚙️ Stack Used

| Component      | Technology                   |
| -------------- | ---------------------------- |
| Language       | Go (Golang)                  |
| Web Framework  | Gin                          |
| Container      | Docker                       |
| CI/CD          | AWS CodePipeline + CodeBuild |
| Hosting        | AWS ECS Fargate              |
| Image Registry | Amazon ECR                   |
| Load Balancer  | Application Load Balancer    |
| Monitoring     | AWS CloudWatch Logs          |

---

## 🧱 Step-by-Step Setup

### ✅ 1. Code Structure

```
.
├── cmd/main.go              # Entry point
├── internal/
│   ├── handlers/            # Route handlers
│   ├── models/              # Data models
│   └── utils/               # Utility functions
├── products.json            # Sample dataset
└── main_test.go             # Unit tests
```

---

### ✅ 2. Dockerization

#### 📄 Dockerfile

```Dockerfile
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]
```

---

### ✅ 3. CodeBuild Configuration

#### 📄 buildspec.yml

```yaml
version: 0.2
env:
  variables:
    IMAGE_TAG: "dev"

phases:
  pre_build:
    commands:
      - go test ./...
      - aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_URI
  build:
    commands:
      - docker build -t $ECR_URI:$IMAGE_TAG .
      - docker tag $ECR_URI:$IMAGE_TAG $ECR_URI:latest
      - docker push $ECR_URI:$IMAGE_TAG
      - docker push $ECR_URI:latest
  post_build:
    commands:
      - aws ecs update-service --cluster $ECS_CLUSTER --service $ECS_SERVICE --force-new-deployment
```

---

### ✅ 4. ECR Repository

* Created ECR repository: `products-api-dev`
* Docker images pushed by CodeBuild

---

### ✅ 5. CodePipeline

Pipeline stages:

* **Source**: GitHub
* **Build**: CodeBuild (uses `buildspec.yml`)
* **Deploy**: ECS Fargate (auto-updates ECS Service)

---

### ✅ 6. Task Definition (ECS)

* **Launch type**: `FARGATE`
* **CPU**: `1024`
* **Memory**: `3072`
* **Container name**: `products-api-container`
* **Port Mapping**: `8080`
* **Image**: Pulled from ECR
* **Logs**: `/ecs/products-api` in CloudWatch

---

### ✅ 7. ECS Cluster & Service

* **Cluster**: `products-cluster-dev`
* **Service**: `products-api-service-with-alb`
* **Tasks**: 1 running on Fargate
* **Uses**: latest task definition

---

### ✅ 8. Application Load Balancer

* **Name**: `products-api-alb`
* **Scheme**: Internet-facing
* **Listener**: HTTP on port 80
* **Target Group**: `products-api-tg`
* **Type**: IP
* **Target Port**: 8080
* **Health Check Path**: `/api/health` or `/api/products`

#### 🔐 Security Groups:

* **ALB** allows inbound HTTP (port 80) from `0.0.0.0/0`
* **ECS Task** allows inbound from ALB SG on port `8080`

---

## ✅ Sample API Endpoints

| Method | Endpoint            | Description           |
| ------ | ------------------- | --------------------- |
| GET    | `/api/products`     | Returns all products  |
| GET    | `/api/products/:id` | Returns product by ID |
| GET    | `/api/health`       | Health check for ALB  |

---

## 🧪 Testing

Run tests locally with:

```bash
go test ./...
```

**Output**:
Basic tests validate `/api/products` response structure and availability.

---

## 🧽 Step-by-Step Workflow & Dependencies

### 1. Codebase Setup

Structure your code into proper folders (`cmd/`, `internal/`, etc.) to follow Go best practices.

This allows for:

* Clear separation of concerns
* Easier testing and maintenance
* Sample dataset (`products.json`) is used for demo purposes.

✅ **Why it matters:** This is your foundation. Everything else builds on this structure.

---

### 2. Dockerization (Packaging)

Create a `Dockerfile` to containerize the app.

Docker ensures the app runs identically everywhere (local, build, production).

```bash
docker build -t product-api .
docker run -p 8080:8080 product-api
```

✅ **Why it matters:** Without Docker, deployments would depend on system configurations — containers eliminate that issue.

---

### 3. Amazon ECR Setup (Image Hosting)

Create an ECR repository (`products-api-dev`) to host your Docker images.

ECR is your private Docker registry within AWS.

✅ **Why it matters:** ECS pulls the container image from ECR during deployment.

---

### 4. Unit Tests Before Build

Use `go test ./...` to ensure your app works before building or deploying.

Tests are placed in the `pre_build` phase of CodeBuild.

✅ **Why it matters:** Never deploy broken code. Automated tests catch issues early.

---

### 5. CI/CD Pipeline (CodePipeline + CodeBuild)

* **Source**: GitHub repo watches for commits.
* **Build**: CodeBuild uses `buildspec.yml` to:

  * Run tests
  * Build Docker image
  * Push to ECR
* **Deploy**: CodePipeline triggers ECS deployment.

✅ **Why it matters:** Automates your deployments end-to-end — from commit to deployment — with zero manual intervention.

---

### 6. Task Definition & ECS Service

* ECS needs a **Task Definition**: blueprint for container config (CPU, memory, image, ports).
* ECS Service uses this Task Definition to launch the app on **Fargate** (serverless containers).

✅ **Why it matters:** Task Definitions ensure ECS knows how to run your containerized app.

---

### 7. Application Load Balancer (ALB)

* ALB exposes your app to the internet.
* Forwards HTTP requests on port 80 to your ECS service.
* Health check: regularly pings `/api/health` to restart failing tasks.

✅ **Why it matters:** Without ALB, the app wouldn’t be publicly accessible. It also keeps your app healthy.

---

### 8. Logs & Monitoring (CloudWatch)

* ECS task logs go to CloudWatch.
* Helps in debugging and monitoring real-time application behavior.

✅ **Why it matters:** Crucial for visibility into what your service is doing in production.

---

## 🔄 End-to-End Flow

```
Developer commits code ➝ GitHub triggers CodePipeline ➝
CodeBuild runs tests & builds image ➝ Image pushed to ECR ➝
CodePipeline triggers ECS update ➝
ECS pulls image & restarts task ➝
ALB routes traffic to ECS ➝
Logs sent to CloudWatch ➝
Health checks keep app stable
```

---

## ✅ TL;DR: Setup Dependencies (Order Matters)

1. Write and test your Go app
2. Create Dockerfile
3. Push code to GitHub
4. Set up ECR to host images
5. Create ECS Cluster + Task Definition
6. Create CodeBuild project
7. Write buildspec.yml (linking to ECR, tests, and ECS updates)
8. Set up CodePipeline (Source ➝ Build ➝ Deploy)
9. Create ALB + security groups + health checks
10. Deploy and monitor via CloudWatch

WHAT HAPPENS WHEN YOU USE THIS PIPELINE
🧱 Pipeline Structure:
java
Copy
Edit
GitHub → CodePipeline
         ├── Source Stage (GitHub)
         ├── Build Stage (CodeBuild)
         └── Deploy Stage (ECS)
✅ Step-by-Step Breakdown
🟩 1. Source Stage (GitHub)
AWS CodePipeline watches your GitHub repo (dev branch in your case).

When a new commit is pushed:

The pipeline is triggered automatically.

The code is zipped and passed as an artifact named SourceOutput.

🟦 2. Build Stage (CodeBuild)
CodeBuild is triggered with:

The zipped source code from GitHub

Your buildspec.yml file which defines the steps to run

🔧 Inside buildspec.yml, CodeBuild:
Logs into Docker Hub (if configured, to avoid rate limits)

Logs into Amazon ECR

Builds your Docker image using your Dockerfile

Tags the image:
Example:
products-api:latest → 654654500855.dkr.ecr.us-east-1.amazonaws.com/products-api:latest

Pushes the image to your ECR repo

Generates an imagedefinitions.json file like this:

json
Copy
Edit
[
  {
    "name": "productsapi-container",
    "imageUri": "654654500855.dkr.ecr.us-east-1.amazonaws.com/products-api:latest"
  }
]
This JSON tells ECS what container to update and with which image.

Passes this file (imagedefinitions.json) as the output artifact named BuildOutput to the deploy stage.

🟨 3. Deploy Stage (Amazon ECS)
CodePipeline reads the imagedefinitions.json from the BuildOutput

It finds:

The ECS container name to update (productsapi-container)

The new Docker image URI

It then:

Creates a new ECS Task Definition revision using the new image

Updates the ECS Service (e.g. productsapi-service)

Triggers a new deployment

ECS will replace the old tasks with new ones that use the latest Docker image

📦 Summary
Stage	What it does
Source	Pulls latest GitHub code on push
Build	Builds Docker image, pushes to ECR, writes imagedefinitions.json
Deploy	Reads imagedefinitions.json, updates ECS task/service with new image

🔁 From Now On:
Just push to GitHub, and the following will happen automatically:

✅ Image built
✅ Pushed to ECR
✅ ECS service updated
✅ Live with the latest code!

Let me know if you want:

Automatic tagging (e.g. v1.0.3)?

Blue/green or canary deployments?

Manual approval before ECS deployment?

You're now in full CI/CD territory — well done!