# Queue System

## Overview
Queue System 是一個商店排隊管理平台，提供商店管理目前正在排隊的顧客以及顧客目前的服務進度。
商家在每次申請帳號後(開店)，可以使用此系統 24 小時，接著系統就會自動將帳號停止（關店）。帳號停止後，商家會收到一份 csv 報表，紀錄當日所有顧客的詳細資訊，方便商家做商業分析。隔天商家要開店並使用此系統，只要再重新申請帳號即可使用。

## Screenshots

## Architecture
![](./images/architecture.png)

### Deployment
Queue System 在 AWS EC2 內使用 MicroK8s 架設 k8s 群集，並使用 `Nginx Ingress Controller` 及 `MetalLB` (Load-Balancer) 讓群集內的服務與外部進行溝通。

k8s 群集內部屬了以下資源：
* Deployments:
  * Backend
  * Frontend
  * gRPC
* Services
  * Backend
  * Frontend
  * gRPC
* CronJob
  * 設定每分鐘執行程式以對 Backend Service 發送 REST API，檢查是否有開店已經 24 小時的商店，將其帳號停止，並寄送 csv 報表給商家，紀錄當日所有顧客的詳細資訊。
* Ingress
  * Nginx Ingress Controller
  * Ingress

### PostgreSQL
Queue System 內所有資料儲存在 PostgreSQL 資料庫內，並且 PostgreSQL 資料庫沒有部屬於 k8s 群集之內。PostgreSQL 內預設會有

### Vault

### Backend

### Frontend

### gRPC

### Client (Store) and Client (Customer)

## TODO
* sse
* snapshots
* structure
* env vars
* code structure
* log format