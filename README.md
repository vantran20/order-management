# 🧾 Order Management System

Hệ thống quản lý đơn hàng (Order Management) được xây dựng theo kiến trúc **Clean Architecture**, sử dụng **Golang**, **Gin framework**, **WebSocket** cho real-time notification, và được container hóa bằng **Docker** để dễ dàng triển khai trên nhiều môi trường.

## 🚀 Tính năng chính

### ✅ Authentication
- Đăng ký người dùng (Register)
- Đăng nhập (Login) và xác thực qua middleware
- Cập nhật / xóa người dùng

### 📦 Product Management
- Tạo, cập nhật, xóa sản phẩm
- Quản lý tồn kho

### 🛒 Order Management
- Tạo đơn hàng và validate sản phẩm tồn kho
- Tự động cập nhật tổng tiền đơn hàng
- Cập nhật trạng thái đơn hàng

### 🔔 Real-time Notification
- WebSocket hỗ trợ gửi thông báo **realtime** tới client khi đơn hàng được tạo hoặc cập nhật trạng thái

---

## 🧑‍💻 Cài đặt & chạy project

### Yêu cầu
- Docker + Docker Compose
- Make

### Cài đặt nhanh:

```bash
# Clone dự án
git clone https://github.com/vantran20/order-management.git
cd order-management

# Khởi động local environment [Pull Docker Image + Setup local]
make setup
# Generate ORM
make boilerplate
# Build ứng dụng
make run

Các lệnh hữu ích:
•	maek api-update-vendor: 	Download các thư việc và dọn dẹp các thư viện không dùng đến
•	make teardown:	                Dừng container
•	make api-pg-migrate:    	Chạy migration
•	make test:              	Chạy test



⸻

📡 API Endpoint

Public APIs
	•	POST   /public/users/register – Đăng ký người dùng
	•	POST   /public/users/login – Đăng nhập

Authenticated APIs (require token)
	•	GET    /authenticated/users/profile – Lấy thông tin user
	•	GET    /authenticated/users/:id – Lấy thông tin user by ID
	•	GET    /authenticated/users/list - Lấy thông tin users
	•	POST   /authenticated/products/create – Tạo product
	•	POST   /authenticated/products/update – Cập nhập product
	•	POST   /authenticated/products/delete  – Soft delete product
	•	GET    /authenticated/products/:id  – Lấy thông tin product bằng id
	•	GET    /authenticated/products/list   – Lấy thông tin tất cả products
	•	POST   /authenticated/order/create   – Tạo đơn hàng
	•	POST   /authenticated/order/update/:id   – Cập nhập đơn hàng

WebSocket
	•	ws://localhost:3000/authenticated/order/ws – Lắng nghe thông báo trạng thái đơn hàng
