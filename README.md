🧾 Order Management System

The Order Management System is built using Clean Architecture, developed in Golang with the Gin framework, supports WebSocket for real-time notifications, and is containerized using Docker for easy deployment across environments.

🚀 Key Features

✅ Authentication

•	User registration

•	Login with middleware-based authentication

•	Update / delete user

📦 Product Management

•	Create, update, delete products

•	Manage inventory


🛒 Order Management

•	Create orders with inventory validation

•	Automatically calculate total order cost

•	Update order status


🔔 Real-time Notification

•	WebSocket support for real-time notifications when an order is created or its status changes

⸻

🧑‍💻 Setup & Run Project

Requirements

•	Docker + Docker Compose

•	Make CMD


# Quick Setup:

## Clone the project
1. `git clone https://github.com/vantran20/order-management.git`

2. `cd order-management`

## Start local environment [Pull Docker Image + Setup local]
`make setup`
## Generate ORM
`make boilerplate`
## Build and run the app
`make run`

Useful commands:

•	`make api-update-vendor`:   Download and clean up unused dependencies

•	`make teardown`:            Stop containers

•	`make api-pg-migrate`:      Run database migrations

•	`make test`:                Run tests

⸻

📡 API Endpoints


## Public APIs:

•	POST   /public/users/register – Register user

•	POST   /public/users/login – Login

## Authenticated APIs (require token):

•	GET    /authenticated/users/profile – Get user profile

•	GET    /authenticated/users/:id – Get user by ID

•	GET    /authenticated/users/list – Get user list

•	POST   /authenticated/products/create – Create product

•	POST   /authenticated/products/update – Update product

•	POST   /authenticated/products/delete – Soft delete product

•	GET    /authenticated/products/:id – Get product by ID

•	GET    /authenticated/products/list – Get all products

•	POST   /authenticated/order/create – Create order

•	POST   /authenticated/order/update/:id – Update order


## WebSocket:
•	ws://localhost:3000/authenticated/order/ws – Listen to order status updates in real time