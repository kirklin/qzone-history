package qrcode

import (
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"os"
	"path/filepath"
)

// SaveQRCode 用于将二维码保存为图片文件
func SaveQRCode(qrData []byte) (string, error) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 创建文件名
	filename := filepath.Join(dir, "qrcode.png")

	// 写入文件
	err = os.WriteFile(filename, qrData, 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// Display 用于在终端上显示二维码
func Display(qrData []byte) {
	qrterminal.GenerateHalfBlock(string(qrData), qrterminal.M, os.Stdout)
	fmt.Println("请使用手机QQ扫描上方二维码登录")
}

// SaveAndDisplayQRCode 整合保存和显示二维码的功能
func SaveAndDisplayQRCode(qrData []byte) error {
	// 保存二维码
	qrPath, err := SaveQRCode(qrData)
	if err != nil {
		return fmt.Errorf("保存二维码失败: %w", err)
	}

	// 打印二维码路径
	fmt.Printf("二维码已保存至 %s，请使用手机QQ扫描登录\n", qrPath)

	// 在终端显示二维码
	Display(qrData)

	return nil
}
