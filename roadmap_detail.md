🧱 GIAI ĐOẠN 1 – Clean Code & Design Pattern (2–3 tuần)
🎯 Mục tiêu:
Viết code rõ ràng, maintainable, áp dụng design pattern đúng chỗ.

✍️ Việc cần làm:
Khởi tạo project theo Clean Architecture (Folder: cmd, internal/domain, infrastructure, usecase, etc.)

Tạo core domain: Ticket, Event, User, QueueSession

Áp dụng design pattern:

Factory cho tạo vé (VIP, thường)

State cho vòng đời vé (Available, OnHold, Paid, Expired)

Repository pattern để tách DB layer

Strategy cho ticket allocation (FCFS / random)

Observer để update trạng thái qua WebSocket hoặc event bus

✅ Output:
Project scaffold theo Clean Arch

Unit test cho usecase

Logging/error handling chuẩn hóa

🚀 GIAI ĐOẠN 2 – Kiến trúc & Performance (3–4 tuần)
🎯 Mục tiêu:
Thiết kế module có thể scale, xử lý tải lớn, tối ưu DB

✍️ Việc cần làm:
Xây module hàng đợi:

Redis Sorted Set để lưu thứ tự xếp hàng

TTL giữ chỗ, xử lý timeout bằng worker

Lock vé bằng Lua script (chống double booking)

Tối ưu DB:

Indexing theo event_id, user_id

Query profiling & explain analyze

Read replica nếu cần scale

Concurrency Patterns:

Worker pool xử lý thanh toán

Context + channel quản lý timeout

Rate limiting (per user IP/email)

✅ Output:
System benchmark: X xử lý bao nhiêu request/s?

Stress test (1000 user cùng đặt vé)

⚙️ GIAI ĐOẠN 3 – DevOps + CI/CD (2 tuần)
🎯 Mục tiêu:
Build, deploy, monitor như production thực tế

✍️ Việc cần làm:
Docker hóa toàn bộ service (event, ticket, queue, payment)

CI/CD: GitHub Actions cho:

Build

Test

Deploy local/staging (Docker Compose)

Monitoring:

Prometheus + Grafana (số người trong queue, vé còn lại)

Alert khi service down, vé lỗi nhiều

Logging: Gửi logs về Loki hoặc ELK

✅ Output:
Pipeline CI/CD hoạt động đầy đủ

Dashboard performance metrics

Alert rule khi queue delay > threshold

🧠 GIAI ĐOẠN 4 – Software Architecture & Hệ thống lớn (3 tuần)
🎯 Mục tiêu:
Tư duy thiết kế hệ thống hoàn chỉnh, có khả năng mở rộng & bảo trì lâu dài

✍️ Việc cần làm:
Vẽ kiến trúc hệ thống:

Component diagram

Sequence diagram khi người dùng xếp hàng → đặt vé → thanh toán

ADR (Architectural Decision Record):

Redis chosen for queue

Dùng WebSocket hay polling?

Làm gì nếu Redis mất kết nối?

Security:

CSRF, rate limit

OAuth2/JWT cho user auth

HTTPS everywhere (TLS)

Trade-off analysis:

Scale queue bằng Redis vs Kafka?

Hold vé ở memory hay DB?

✅ Output:
Full system design doc

5+ ADR ghi rõ lý do kiến trúc

Security checklist