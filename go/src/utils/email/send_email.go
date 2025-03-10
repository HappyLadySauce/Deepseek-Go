package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"Deepseek-Go/config"
)

// 生成随机验证码
func GenerateVerificationCode() string {
	// 使用随机源初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成6位数字验证码
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", r.Intn(10))
	}
	return code
}

// 发送邮箱验证邮件
func SendVerificationEmail(toEmail, code string) error {
	// 如果邮件配置未初始化
	emailCfg := config.Config.Email
	if emailCfg.Host == "" {
		// 调试模式下，直接返回成功并打印验证码
		log.Printf("邮件服务未配置，验证码: %s 发送到: %s", code, toEmail)
		return nil
	}

	// 记录邮件配置信息
	fromAddr := ExtractEmailAddress(emailCfg.From)
	log.Printf("邮件配置: 发件人=%s (提取地址=%s), 主机=%s:%d, SSL=%v",
		emailCfg.From, fromAddr, emailCfg.Host, emailCfg.Port, emailCfg.EnableSSL)

	subject := "邮箱验证码"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>邮箱验证</h2>
			<p>您的验证码是: <strong>%s</strong></p>
			<p>该验证码3分钟内有效，请不要泄露给他人。</p>
		</body>
		</html>
	`, code)

	return sendEmail(toEmail, subject, body)
}

// 从完整的From字段中提取纯邮箱地址
func ExtractEmailAddress(from string) string {
	// 尝试匹配 "Name <email@example.com>" 格式
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindStringSubmatch(from)
	if len(matches) > 1 {
		return matches[1]
	}

	// 如果不是上述格式，则直接返回原字符串（假设它就是一个邮箱地址）
	return from
}

// 确保From字段格式正确
func EnsureFromFormat(from, username string) string {
	// 如果from已经包含尖括号，说明格式已经是"Name <email@example.com>"
	if strings.Contains(from, "<") && strings.Contains(from, ">") {
		return from
	}

	// 如果from不包含@，说明它可能只是一个名称
	if !strings.Contains(from, "@") {
		return fmt.Sprintf("%s <%s>", from, username)
	}

	// 如果from只是一个邮箱地址
	if from == username {
		return fmt.Sprintf("Verification <%s>", from)
	}

	// 其他情况，构造一个名称 + 邮箱的格式
	return fmt.Sprintf("%s <%s>", from, username)
}

// 发送邮件
func sendEmail(to, subject, body string) error {
	emailCfg := config.Config.Email

	// 确保From字段格式正确
	fromAddr := ExtractEmailAddress(emailCfg.From)
	formattedFrom := EnsureFromFormat(emailCfg.From, emailCfg.Username)

	log.Printf("发送邮件: 收件人=%s, 发件人=%s, 格式化后=%s, 提取地址=%s",
		to, emailCfg.From, formattedFrom, fromAddr)

	// 设置邮件内容
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	header := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n", to, formattedFrom, subject, mime)
	message := header + body

	// 连接地址
	addr := fmt.Sprintf("%s:%d", emailCfg.Host, emailCfg.Port)

	// 根据是否启用SSL/TLS使用不同的发送方式
	if emailCfg.EnableSSL {
		return sendSSLEmail(addr, to, fromAddr, message)
	} else {
		// 普通认证
		auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)
		return smtp.SendMail(addr, auth, fromAddr, []string{to}, []byte(message))
	}
}

// 使用SSL/TLS发送邮件
func sendSSLEmail(addr, to, from, message string) error {
	emailCfg := config.Config.Email

	// 记录TLS连接信息
	log.Printf("开始TLS连接: 地址=%s, 服务器名=%s", addr, emailCfg.ServerName)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         emailCfg.ServerName,
	}

	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接到SMTP服务器失败: %v", err)
	}

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, emailCfg.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 设置认证信息
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 设置发件人 - 使用提取出的纯邮箱地址
	log.Printf("设置发件人: %s", from)
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v (使用地址: %s)", err, from)
	}

	// 设置收件人
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 设置邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取邮件写入器失败: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入器失败: %v", err)
	}

	return client.Quit()
}
