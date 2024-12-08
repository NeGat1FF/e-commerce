# E-Commerce Project

This is an e-commerce project developed for learning purposes, following the guidelines from [roadmap.sh](https://roadmap.sh). The project consists of several microservices, each responsible for a specific functionality within the e-commerce platform.

## Services

1. **Product Service**: Manages product information.
2. **Notification Service**: Handles notifications to users.
3. **Order Service**: Manages customer orders.
4. **Payment Service**: Processes payments.
5. **Search Service**: Provides search functionality using Elasticsearch.
6. **Shopping Cart Service**: Manages user shopping carts.
7. **User Service**: Manages user information and authentication.

## Tech Stack

- MongoDB
- PostgreSQL
- Elasticsearch
- Redis
- RabbitMQ
- gRPC
- Gin
- Docker

## Docker

Each service has its own Dockerfile, allowing for containerized deployment and easy scaling.

## Continuous Integration

GitHub Actions is used for continuous integration (CI) to ensure code quality and automate the build and deployment process.

## Getting Started

To get started with the project, clone the repository and follow the instructions for each service to build and run the containers.

```bash
git clone https://github.com/NeGat1FF/e-commerce.git
cd e-commerce
```

[roadmap.sh](https://roadmap.sh/projects/scalable-ecommerce-platform)

