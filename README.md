# MetaReviews

MetaReviews is a learning project focused on building microservices in Go. It explores various technologies often used in production-grade microservice environments, including:

- **gRPC**  
- **Protocol Buffers (protoBuf)**  
- **Consul**  
- **MySQL**  
- **Kafka**  
- **Docker**  
- **Kubernetes**  

## Overview

MetaReviews consists of multiple independent services that communicate primarily via gRPC. The current services are:

1. **Metadata**  
   Handles the metadata of movies (title, director, release date, etc.).

2. **Rating**  
   Manages the rating information for different movies.

3. **Movie**  
   Aggregates movie information by communicating with the Metadata and Rating services.

**Planned**: Additional services will be introduced as the project expands.

## Key Features

- **Service Discovery**: Uses **Consul** to locate and register services.  
- **Communication**: Uses **gRPC** and **protoBuf** for efficient inter-service communication.  
- **Messaging**: Integrates **Kafka** for event streaming and asynchronous communication between services.  
- **Data Persistence**: Uses **MySQL** for relational data storage.  
- **Containerization & Orchestration**: Runs the entire setup in **Docker** containers, with plans to manage them in **Kubernetes**.

## Project Structure

- **movie-service**: Aggregates data from `metadata-service` and `rating-service`.  
- **metadata-service**: Stores and retrieves movie metadata.  
- **rating-service**: Processes and provides movie ratings.  
- **deployments**: Dockerfiles and Kubernetes configurations.

## Getting Started

1. **Clone the Repository**  
   ```bash
   git clone https://github.com/your-username/meta-reviews.git
   cd meta-reviews

