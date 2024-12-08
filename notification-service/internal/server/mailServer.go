package server

import (
	"context"

	"github.com/NeGat1FF/e-commerce/notification-service/internal/email"
	mail "github.com/NeGat1FF/e-commerce/notification-service/proto"
)

type Server struct {
	mail.UnimplementedMailServiceServer

	sender *email.EmailSender
}

func NewServer(sender *email.EmailSender) *Server {
	return &Server{
		sender: sender,
	}
}

func (s *Server) SendMail(ctx context.Context, req *mail.MailRequest) (res *mail.MailResponse, err error) {
	var emailType email.EmailType

	switch int32(req.Type) {
	case 0:
		emailType = email.EmailConfirmation
	case 1:
		emailType = email.PasswordReset
	}

	err = s.sender.SendMail(req.To, emailType, req.Data)

	if err != nil {
		return nil, err
	}

	res = &mail.MailResponse{
		Success: true,
	}

	return
}
