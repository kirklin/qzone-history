// export_usecase.go

package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/repository"
	"qzone-history/internal/domain/usecase"
)

type exportUseCase struct {
	momentRepo       repository.MomentRepository
	boardMessageRepo repository.BoardMessageRepository
	friendRepo       repository.FriendRepository
}

func NewExportUseCase(
	momentRepo repository.MomentRepository,
	boardMessageRepo repository.BoardMessageRepository,
	friendRepo repository.FriendRepository,
) usecase.ExportUseCase {
	return &exportUseCase{
		momentRepo:       momentRepo,
		boardMessageRepo: boardMessageRepo,
		friendRepo:       friendRepo,
	}
}

func (u *exportUseCase) ExportUserDataToJSON(ctx context.Context, userQQ string) error {
	// 获取用户数据
	moments, _ := u.momentRepo.FindByUserQQ(ctx, userQQ, -1, 0)
	boardMessages, _ := u.boardMessageRepo.FindByUserQQ(ctx, userQQ, -1, 0)
	friends, _ := u.friendRepo.FindFriendsByUserQQ(ctx, userQQ)

	// 创建导出数据结构
	exportData := struct {
		Moments       []entity.Moment
		BoardMessages []entity.BoardMessage
		Friends       []entity.Friend
	}{
		Moments:       moments,
		BoardMessages: boardMessages,
		Friends:       friends,
	}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// 写入文件
	filename := fmt.Sprintf("%s_export.json", userQQ)
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

func (u *exportUseCase) ExportUserDataToExcel(ctx context.Context, userQQ string) error {
	// TODO 实现Excel导出逻辑
	return fmt.Errorf("ExportUserDataToExcel not implemented")
}

func (u *exportUseCase) ExportUserDataToHTML(ctx context.Context, userQQ string) error {
	// TODO 实现HTML导出逻辑
	return fmt.Errorf("ExportUserDataToHTML not implemented")
}
