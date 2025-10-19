# 🚀 OrderFlow - Real-time Order Management System

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![NATS](https://img.shields.io/badge/NATS%20Streaming-2.10+-27AE60?style=for-the-badge)](https://docs.nats.io/nats-streaming-concepts/intro)
[![Docker](https://img.shields.io/badge/Docker-28.5.1-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2.40.0-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/compose/)

**Мощная система управления заказами с реальным временем обновления**

[Интерфейс](#-интерфейс) • [Особенности](#-особенности) • [Быстрый старт](#-быстрый-старт)

</div>

## 📊 Интерфейс

![GIF](-) 
*🎥 [MP4](https://drive.google.com/file/d/1yCcb7Cn2Wf2ziM4YrwisDXr_kLhz8Hje/view?usp=sharing#)*

## ✨ Особенности

### 🏗️ Архитектура
- **Microservices-ready** архитектура на Go
- **In-memory кэш** с автоматическим восстановлением из БД
- **Асинхронная обработка** через NATS Streaming
- **RESTful API** с полной документацией

### 🚀 Производительность
- **1.3ms** среднее время ответа API
- **48+ RPS** на получение заказов  
- **Gzip сжатие** для оптимизации трафика
- **HTTP кэширование** на клиентской стороне

### 🔧 Технологический стек
| Компонент | Технология | Назначение |
|-----------|------------|------------|
| **Backend** | Go 1.21+ | Высокопроизводительный API |
| **Database** | PostgreSQL 15 | Надежное хранение данных |
| **Message Broker** | NATS Streaming | Асинхронная коммуникация |
| **Cache** | In-memory + Redis-ready | Быстрый доступ к данным |
| **Frontend** | Vanilla JS + HTML/CSS | Легковесный интерфейс |
| **Containerization** | Docker + Compose | Простое развертывание |

## 🚀 Быстрый старт

### Предварительные требования
- [Docker](https://docs.docker.com/get-docker/) 
- [Docker Compose](https://docs.docker.com/compose/install/)

### Запуск за 60 секунд ⚡

```bash
# Клонируем репозиторий
git clone https://github.com/Fiflocs/order-service.git
cd order-service

# Запускаем всю инфраструктуру одной командой
docker-compose up --build
```
### Система будет доступна по адресу: http://localhost:8080

```bash
# Проверяем health endpoint
curl http://localhost:8080/health

# Ответ:
{
  "status": "healthy",
  "service": "order-service", 
  "cache_size": 120,
  "timestamp": "2025-10-19T14:21:17+07:00"
}
```
