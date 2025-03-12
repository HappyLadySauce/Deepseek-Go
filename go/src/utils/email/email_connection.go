package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"Deepseek-Go/config"
)

// TestConnection 测试邮件服务器连接
// 这个函数可以用来验证当前配置是否能正确连接到SMTP服务器
func TestConnection() error {
	emailCfg := config.Config.Email

	if emailCfg.Host == "" {
		log.Println("邮件服务未配置")
		return fmt.Errorf("邮件服务未配置")
	}

	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)

	// 提取纯邮箱地址
	fromAddr := ExtractEmailAddress(emailCfg.From)
	formattedFrom := EnsureFromFormat(emailCfg.From, emailCfg.Username)

	log.Printf("测试邮件配置: 主机=%s, 端口=%d, 用户名=%s, SSL=%v",
		emailCfg.Host, emailCfg.Port, emailCfg.Username, emailCfg.EnableSSL)
	log.Printf("From字段处理: 原始值=%s, 格式化后=%s, 提取地址=%s",
		emailCfg.From, formattedFrom, fromAddr)

	// 根据是否使用SSL/TLS选择不同的连接方式
	if emailCfg.EnableSSL {
		return testSSLConnection(addr, fromAddr)
	} else {
		return testPlainConnection(addr, fromAddr)
	}
}

// 测试普通连接
func testPlainConnection(addr, from string) error {
	emailCfg := config.Config.Email

	log.Printf("开始普通连接测试: 地址=%s, 发件人=%s", addr, from)

	// 创建SMTP客户端
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("连接SMTP服务器失败: %v", err)
		return fmt.Errorf("无法连接到SMTP服务器: %v", err)
	}
	defer client.Close()

	// 验证身份
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("SMTP认证失败: %v", err)
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 尝试设置发件人
	if err = client.Mail(from); err != nil {
		log.Printf("设置发件人失败: %v (使用地址: %s)", err, from)
		return fmt.Errorf("设置发件人失败: %v (使用地址: %s)", err, from)
	}

	log.Printf("连接测试成功")
	return nil
}

// 测试SSL/TLS连接
func testSSLConnection(addr, from string) error {
	emailCfg := config.Config.Email

	log.Printf("开始SSL/TLS连接测试: 地址=%s, 发件人=%s, 服务器名=%s",
		addr, from, emailCfg.ServerName)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         emailCfg.ServerName,
	}

	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		log.Printf("TLS连接失败: %v", err)
		return fmt.Errorf("连接到SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, emailCfg.Host)
	if err != nil {
		log.Printf("创建SMTP客户端失败: %v", err)
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 验证身份
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("SMTP认证失败: %v", err)
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 尝试设置发件人
	if err = client.Mail(from); err != nil {
		log.Printf("设置发件人失败: %v (使用地址: %s)", err, from)
		return fmt.Errorf("设置发件人失败: %v (使用地址: %s)", err, from)
	}

	log.Printf("SSL/TLS连接测试成功")
	return nil
}
