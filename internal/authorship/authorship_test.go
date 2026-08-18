package authorship

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthorshipLedger(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "authorship_test_*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ledger := NewLedger(tempDir)

	// 1. Ghi nhận ý tưởng ban đầu
	rec1, err := ledger.RecordContribution(
		TypePremiseSeed,
		"Nguyễn Văn A",
		"Thiết lập tiền đề vũ trụ Sanbaka 1956 tại ĐHBK Hà Nội",
		"Câu chuyện xoay quanh 3 thế hệ sinh viên kỹ thuật trong bối cảnh Cyberpunk Đông Dương...",
		0,
	)
	if err != nil {
		t.Fatalf("RecordContribution seed failed: %v", err)
	}
	if rec1.ID == "" || rec1.Hash == "" {
		t.Fatalf("expected valid ID and Hash, got %v", rec1)
	}

	// 2. Ghi nhận can thiệp cốt truyện
	_, err = ledger.RecordContribution(
		TypeCreativeSteer,
		"Nguyễn Văn A",
		"Chuyển hướng tình tiết chương 3: Thêm nhân vật kỹ sư trưởng",
		"Cho nhân vật chính gặp cố vấn kỹ thuật ở phòng thí nghiệm tầng hầm C9...",
		3,
	)
	if err != nil {
		t.Fatalf("RecordContribution steer failed: %v", err)
	}

	// 3. Ghi nhận duyệt chương
	_, err = ledger.RecordContribution(
		TypeChapterApproval,
		"Nguyễn Văn A",
		"Nghiệm thu bản thảo chương 3",
		"Đã đọc và đồng ý cho phát hành chương 3.",
		3,
	)
	if err != nil {
		t.Fatalf("RecordContribution approval failed: %v", err)
	}

	// 4. Đọc lại danh sách bản ghi
	records, err := ledger.LoadRecords()
	if err != nil {
		t.Fatalf("LoadRecords failed: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	// 5. Thống kê
	summary, err := ledger.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary failed: %v", err)
	}
	if summary.TotalInterventions != 3 {
		t.Fatalf("expected 3 interventions, got %d", summary.TotalInterventions)
	}
	if summary.PremiseSeeds != 1 || summary.CreativeSteers != 1 || summary.ChapterApprovals != 1 {
		t.Fatalf("unexpected summary breakdown: %+v", summary)
	}

	// 6. Xuất báo cáo bản quyền
	report, err := ledger.GenerateCopyrightReport("Sanbaka Universe", "Nguyễn Văn A")
	if err != nil {
		t.Fatalf("GenerateCopyrightReport failed: %v", err)
	}
	if !strings.Contains(report, "HỒ SƠ CHỨNG MINH QUYỀN TÁC GIẢ") || !strings.Contains(report, "Sanbaka Universe") {
		t.Fatalf("report missing key phrases: %s", report)
	}
}
