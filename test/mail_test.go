package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"go-cloud-disk/core/define"
	"net/smtp"
	"testing"
)
func TestSendMial(t *testing.T) {
	e := email.NewEmail()
	e.From = "xiao dao <xiaodao_yin@163.com>"
	e.To = []string{"1598207459@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>12345</h1>")
	err := e.SendWithTLS("smtp.163.com:465",
		// 注意这里的密码不是邮箱登录密码, 是开启 smtp 服务后获取的一串验证码
		smtp.PlainAuth("", "xiaodao_yin@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
