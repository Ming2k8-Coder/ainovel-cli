# Bước 2 RFC: Engine Chạy Trực Tiếp Worker (7 Câu Hỏi Bắt Buộc)

> Trạng thái: Định bản (2026-07-12). Xem chi tiết tại `docs/engine-arbiter.md`.

---

## 1. Mạch Thực Thi Worker: Gọi Chương Trình `subagent.Runner`

`Runner.Run(agent, task)` mỗi lần khởi động một `agentcore.AgentLoop` hoàn chỉnh. Engine gọi trực tiếp, toàn bộ cấu hình trong `build.go` có hiệu lực nguyên vẹn: Model vai trò + Failover, Prompt Cache Key, ThinkingLevel, UsageRecorder, SessionLogger, ContextManagerFactory, RestorePack, StopGuardFactory, StopAfterTools.

---

## 2. Vòng Đời Engine

Vòng lặp chuỗi đơn goroutine; `ctx` cancel = Tạm dừng / Hủy (lan truyền vào vòng lặp Worker, checkpoint đảm bảo không mất mát dữ liệu); Resume/Continue = Bắt đầu vòng lặp mới.

---

## 3. Phân Loại Lỗi (Định Tính Trước)

- Lỗi có thể thử lại (mạng / giới hạn tốc độ / stream-idle): Được xử lý nội bộ tại subagent (MaxRetries=7).
- Lỗi Worker trả về (escalate/hard_stop/lỗi công cụ): Engine thử lại cùng một chỉ thị 1 lần → Vẫn thất bại → Hỏi Arbiter `worker_failure` (retry / reroute / abort).
- Lỗi thông số / Agent không xác định: Tạm dừng trực tiếp + thông báo.

---

## 4. Giao Thức Bế Tắc (Deadlock Protocol)

Mỗi vòng lặp ghi lại khóa chỉ thị `Agent+Task`. Nếu sau khi thực thi, `Route` vẫn tạo ra cùng một khóa, nghĩa là điều kiện sau của nhiệm vụ chưa được đáp ứng, `repeat++`. 
- `repeat == 3` ➔ Hỏi Arbiter `deadlock`.
- `repeat == 5` ➔ Tạm dừng cứng + thông báo.
