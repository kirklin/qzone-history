package usecase

import (
	"context"
)

// ExportUseCase 定义了数据导出相关的用例接口
type ExportUseCase interface {
	// ExportUserDataToJSON 将用户数据导出为JSON格式
	ExportUserDataToJSON(ctx context.Context, userQQ string) error

	// ExportUserDataToExcel 将用户数据导出为Excel格式
	ExportUserDataToExcel(ctx context.Context, userQQ string) error

	// ExportUserDataToHTML 将用户数据导出为HTML格式
	ExportUserDataToHTML(ctx context.Context, userQQ string) error
}
