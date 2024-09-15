package app

import (
	"context"
	"log"
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
	log.Println("开始检查本地登录状态...")
	user, loggedIn, err := a.authUseCase.CheckLocalLoginStatus(ctx)
	if err != nil {
		_ = logError("检查本地登录状态失败", err)
	}

	if !loggedIn {
		log.Println("本地未登录，获取登录二维码...")
		qrCode, qrsig, err := a.authUseCase.GetLoginQRCode(ctx)
		if err != nil {
			return logError("获取登录二维码失败", err)
		}

		log.Println("保存并显示二维码...")
		qrPath, err := qrcode.SaveQRCode(qrCode)
		if err != nil {
			return logError("保存和显示二维码失败", err)
		}
		// 打印二维码路径
		log.Printf("二维码已保存至 %s，请使用手机QQ扫描登录\n", qrPath)
		log.Println("等待用户扫描二维码...")

		loginSuccess := false
		for !loginSuccess {
			status, res, err := a.authUseCase.CheckQRCodeLoginStatus(ctx, qrsig)
			if err != nil {
				return logError("检查二维码登录状态失败", err)
			}

			switch status {
			case entity.LoginStatusSuccess:
				log.Println("二维码扫描成功，完成登录过程...")
				user, err = a.authUseCase.CompleteLogin(ctx, res)
				if err != nil {
					return logError("完成登录失败", err)
				}
				loginSuccess = true
			case entity.LoginStatusWaiting:
				// log.Println("等待扫描二维码")
			case entity.LoginStatusScanning:
				log.Println("二维码认证中")
			case entity.LoginStatusExpired:
				return logError("二维码已失效", nil)
			default:
				log.Println("未知二维码状态")
			}

			if !loginSuccess {
				time.Sleep(2 * time.Second)
			}
		}
	}

	log.Println("登录成功，获取活动记录...")
	_, err = a.activityUseCase.FetchActivities(ctx, *user)
	if err != nil {
		return logError("获取活动记录失败", err)
	}

	log.Println("开始数据重建过程...")
	err = a.reconstructionUseCase.ReconstructMomentsFromActivities(ctx, user.QQ)
	if err != nil {
		return logError("重建 Moments 失败", err)
	}

	err = a.reconstructionUseCase.ReconstructBoardMessagesFromActivities(ctx, user.QQ)
	if err != nil {
		return logError("重建留言板失败", err)
	}

	log.Println("导出用户数据到 JSON 格式...")
	err = a.exportUseCase.ExportUserDataToJSON(ctx, user.QQ)
	if err != nil {
		return logError("导出用户数据到 JSON 失败", err)
	}

	log.Println("应用程序成功完成")
	return nil
}

func logError(message string, err error) error {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
	return err
}
