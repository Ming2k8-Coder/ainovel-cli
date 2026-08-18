# Tầng Văn Phong & Tiêu Chuẩn Sáng Tác (Voice Layer)

Tầng văn phong định nghĩa các tiêu chuẩn chất lượng và giọng văn dành cho Writer và Editor.

## 1. Ba Tầng Ghi Đè Văn Phong (3-Tier Voice Overrides)

Hệ thống áp dụng văn phong theo 3 tầng từ chung đến riêng:

1. **Mặc định Hệ thống (`assets/voice.md`)**: Tiêu chuẩn văn phong cốt lõi tích hợp sẵn.
2. **Thể loại Văn học (`assets/styles/<style>.md`)**: Phong cách chuyên biệt theo thể loại (Kỳ ảo, Ngôn tình, Trinh thám...).
3. **Cấp Sách Cụ Thể (`<ThưMụcSách>/style/voice.md`)**: Văn phong riêng do người dùng tùy chỉnh cho duy nhất bộ sách này.

## 2. Loại Bỏ "Văn Phong AI" (De-AI Tone)

Tầng văn phong áp dụng các quy tắc tự động nhằm loại bỏ các dấu hiệu viết văn đặc trưng của AI:
- Tránh sử dụng các cụm từ đệm công thức ("không phải... mà là...", "trong một khoảnh khắc...").
- Đa dạng hóa cấu trúc câu và nhịp điệu chương.
- Không lặp lại tóm tắt sự kiện chương trước ở đầu chương mới.
