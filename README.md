# Zapp
A lightweight REST API inspired by X.com, built with GO.

# Technologies 
[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=go)](https://go.dev/)
[![Redis](https://img.shields.io/badge/-Redis-464646?style=flat-square&logo=redis)](https://redis.io/)
[![Kafka](https://img.shields.io/badge/-Kafka-464646?style=flat-square&logo=apache-kafka)](https://kafka.apache.org/)
[![Elasticsearch](https://img.shields.io/badge/-Elasticsearch-464646?style=flat-square&logo=elasticsearch)](https://www.elastic.co/elasticsearch/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![MinIO](https://img.shields.io/badge/-MinIO-464646?style=flat-square&logo=minio)](https://min.io/)
[![golang-migrate](https://img.shields.io/badge/-Migrate-464646?style=flat-square&logo=go)](https://github.com/golang-migrate/migrate)
[![Prometheus](https://img.shields.io/badge/-Prometheus-464646?style=flat-square&logo=prometheus)](https://prometheus.io/)
[![Grafana](https://img.shields.io/badge/-Grafana-464646?style=flat-square&logo=grafana)](https://grafana.com/)
[![Gin](https://img.shields.io/badge/-Gin-464646?style=flat-square&logo=go)](https://gin-gonic.com/)
[![gRPC](https://img.shields.io/badge/-gRPC-464646?style=flat-square&logo=grpc)](https://grpc.io/)
[![Govatar](https://img.shields.io/badge/-Govatar-464646?style=flat-square&logo=go)](https://github.com/alexeyco/govatar)
[![Docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/-Kubernetes-464646?style=flat-square&logo=kubernetes)](https://kubernetes.io/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Tech Stack

### Languages & Frameworks
- **Programming Language**: Go 1.24.4
- **Web Framework**: [Gin](https://gin-gonic.com/) — High-performance HTTP framework for building REST APIs  
- **gRPC**: Implemented services using the gRPC protocol for inter-service communication  
- **Avatar Generation**: [govatar](https://github.com/alexeyco/govatar) — Library for generating random avatars

### Data Storage
- **Relational Database**: [PostgreSQL](https://www.postgresql.org/) — Reliable and scalable RDBMS (recommended version: 16.11+)  
- **Caching / Sessions**: [Redis](https://redis.io/) — Fast key-value store for caching and session management  
- **Search & Analytics**: [Elasticsearch](https://www.elastic.co/elasticsearch/) — Full-text search and data aggregation engine  
- **Object Storage**: [MinIO S3](https://min.io/) — S3-compatible storage for files and media

### Migrations & Schema Management
- **Database Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) — Tool for managing PostgreSQL schema migrations

### Monitoring & Observability
- **Metrics**: [Prometheus](https://prometheus.io/) — Collection and visualization of application metrics  
- **Dashboards**: [Grafana](https://grafana.com/) — Visualization of metrics, logs, and traces

### Infrastructure & Deployment
- **Containerization**: [Docker](https://www.docker.com/) — Building and running containers  
- **Orchestration**: [Kubernetes](https://kubernetes.io/) — Managing microservices in production