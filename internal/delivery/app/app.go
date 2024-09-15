package app

import (
	"context"
	"fmt"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/domain/usecase"
	"qzone-history/pkg/qrcode"
	"time"
)

type App struct {
	authUseCase           usecase.AuthUseCase
	momentUseCase         usecase.MomentUseCase
	boardMessageUseCase   usecase.BoardMessageUseCase
	friendUseCase         usecase.FriendUseCase
	exportUseCase         usecase.ExportUseCase
	activityUseCase       usecase.ActivityUseCase
	reconstructionUseCase usecase.ReconstructionUseCase
}

func NewApp(
	authUseCase usecase.AuthUseCase,
	momentUseCase usecase.MomentUseCase,
	boardMessageUseCase usecase.BoardMessageUseCase,
	friendUseCase usecase.FriendUseCase,
	exportUseCase usecase.ExportUseCase,
	activityUseCase usecase.ActivityUseCase,
	reconstructionUseCase usecase.ReconstructionUseCase,
) *App {
	return &App{
		authUseCase:           authUseCase,
		momentUseCase:         momentUseCase,
		boardMessageUseCase:   boardMessageUseCase,
		friendUseCase:         friendUseCase,
		exportUseCase:         exportUseCase,
		activityUseCase:       activityUseCase,
		reconstructionUseCase: reconstructionUseCase,
	}
}

func (a *App) Run(ctx context.Context) error {
	// 检查本地登录状态
	user, loggedIn, err := a.authUseCase.CheckLocalLoginStatus(ctx)
	//if err != nil {
	//	return fmt.Errorf("failed to check local login status: %w", err)
	//}
	var responseText string
	if !loggedIn {
		// 获取登录二维码
		qrCode, qrsig, err := a.authUseCase.GetLoginQRCode(ctx)
		if err != nil {
			return fmt.Errorf("failed to get login QR code: %w", err)
		}

		// 显示二维码
		err = qrcode.SaveAndDisplayQRCode(qrCode)
		if err != nil {
			return err
		}

		// 等待用户扫描二维码
		for {
			status, res, err := a.authUseCase.CheckQRCodeLoginStatus(ctx, qrsig)
			if err != nil {
				return fmt.Errorf("failed to check QR code login status: %w", err)
			}

			if status == entity.LoginStatusSuccess {
				responseText = res
				break
			} else {
				switch status {
				case entity.LoginStatusWaiting:
					fmt.Println("等待扫描二维码")
				case entity.LoginStatusScanning:
					fmt.Println("二维码认证中")
				case entity.LoginStatusExpired:
					return fmt.Errorf("二维码已失效")
				default:
				}
			}

			// 等待一段时间后再次检查
			time.Sleep(2 * time.Second)
		}

		// 完成登录过程
		user, err = a.authUseCase.CompleteLogin(ctx, responseText)
		if err != nil {
			return fmt.Errorf("failed to complete login: %w", err)
		}
	}
	// 获取活动记录
	_, err = a.activityUseCase.FetchActivities(ctx, *user)
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}
	// 开始数据重建过程
	err = a.reconstructionUseCase.ReconstructMomentsFromActivities(ctx, user.QQ)
	if err != nil {
		return fmt.Errorf("failed to reconstruct moments: %w", err)
	}

	err = a.reconstructionUseCase.ReconstructBoardMessagesFromActivities(ctx, user.QQ)
	if err != nil {
		return fmt.Errorf("failed to reconstruct board messages: %w", err)
	}

	// 导出用户数据
	err = a.exportUseCase.ExportUserDataToJSON(ctx, user.QQ)
	if err != nil {
		return fmt.Errorf("failed to export user data to JSON: %w", err)
	}

	fmt.Println("Application completed successfully")
	return nil
}
