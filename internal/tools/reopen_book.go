package tools

import (
	"fmt"
	"slices"

	"github.com/voocel/ainovel-cli/internal/domain"
	"github.com/voocel/ainovel-cli/internal/errs"
	"github.com/voocel/ainovel-cli/internal/store"
)

// ReopenBook mở lại cuốn tiểu thuyết đã hoàn thành để chuyển sang trạng thái làm lại/viết lại, được Engine gọi ở ranh giới can thiệp của người dùng.
// Sau khi hoàn thành tác phẩm, completePhaseGate sẽ chặn mọi hoạt động điều phối subagent, khiến người dùng không thể viết lại các chương đã hoàn thành.
// Hàm này không thông qua subagent, có thể gọi được trong giai đoạn hoàn thành: chuyển đổi phase về writing một cách nguyên tử, đưa các chương mục tiêu vào
// PendingRewrites, đặt flow=rewriting. Sau đó Flow Router sẽ điều phối writer viết lại từng chương theo hàng đợi,
// khi hàng đợi hoàn tất sẽ tự động đóng sổ tác phẩm qua commit_chapter. Toàn bộ logic lõi của Gate / Router / edit / commit không cần phải sửa đổi.
func ReopenBook(s *store.Store, chapters []int, reason string) error {
	if len(chapters) == 0 {
		return fmt.Errorf("chapters không được để trống, cần chỉ rõ chương muốn làm lại: %w", errs.ErrToolArgs)
	}

	progress, err := s.Progress.Load()
	if err != nil {
		return fmt.Errorf("load progress: %w: %w", errs.ErrStoreRead, err)
	}
	if progress == nil {
		return fmt.Errorf("progress chưa được khởi tạo: %w", errs.ErrToolPrecondition)
	}
	// Chỉ cho phép làm lại các chương đã hoàn thành; chương không nằm trong danh sách hoàn thành bị từ chối rõ ràng và hướng dẫn người dùng điều chỉnh độ dài tác phẩm.
	var invalid []int
	for _, ch := range chapters {
		if !slices.Contains(progress.CompletedChapters, ch) {
			invalid = append(invalid, ch)
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("chương %v chưa hoàn thành, reopen chỉ có thể làm lại các chương đã viết xong (nếu muốn thêm cốt truyện vui lòng dùng điều chỉnh độ dài): %w", invalid, errs.ErrToolPrecondition)
	}

	// Kiểm tra điều kiện tiên quyết của phase được đảm bảo bên trong store.Reopen (chỉ có thể gọi khi phase=complete).
	if err := s.Progress.Reopen(chapters, reason); err != nil {
		return fmt.Errorf("reopen: %w: %w", errs.ErrStoreWrite, err)
	}

	// checkpoint: đối xứng với complete_book (GlobalScope + meta/progress.json).
	if _, err := s.Checkpoints.AppendArtifact(domain.GlobalScope(), "reopen", "meta/progress.json"); err != nil {
		return fmt.Errorf("checkpoint reopen: %w: %w", errs.ErrStoreWrite, err)
	}
	return nil
}
