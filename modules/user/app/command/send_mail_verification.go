package command

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	templateDomain "tixgo/modules/template/domain"
	"tixgo/modules/user/domain"
	sharedMail "tixgo/shared/events/mail"

	"github.com/duongptryu/gox/messaging"
	"github.com/duongptryu/gox/notification/mail"
	"github.com/duongptryu/gox/syserr"
)

const (
	SlugMailOTP = "mail-verify-mail"
)

type sendOTPVerifyMailHandler struct {
	otpStore         domain.OTPStore
	templateRepo     templateDomain.TemplateRepository
	templateRenderer templateDomain.TemplateRenderer
	eventBus         messaging.EventBus
}

type SendOTPVerifyMailCommand struct {
	Mail string
}

func NewSendOTPVerifyMailHandler(otpStore domain.OTPStore, templateRepo templateDomain.TemplateRepository, templateRenderer templateDomain.TemplateRenderer, eventBus messaging.EventBus) *sendOTPVerifyMailHandler {
	return &sendOTPVerifyMailHandler{
		otpStore:         otpStore,
		templateRepo:     templateRepo,
		templateRenderer: templateRenderer,
		eventBus:         eventBus,
	}
}

func (h *sendOTPVerifyMailHandler) Handle(ctx context.Context, cmd *SendOTPVerifyMailCommand) error {
	otp, err := generateOTP()
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to generate OTP")
	}

	// store otp
	err = h.otpStore.Store(ctx, cmd.Mail, otp)
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to store OTP")
	}

	template, err := h.templateRepo.GetBySlug(ctx, SlugMailOTP)
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to get template")
	}

	// render to html
	rendered, err := h.templateRenderer.Render(ctx, template, map[string]interface{}{
		"otp": otp,
	})
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to render template")
	}

	// send mail
	h.eventBus.PublishEvent(ctx, &sharedMail.EventSendMail{
		ToMail: []mail.EmailAddress{
			{
				Email: cmd.Mail,
				Name:  "",
			},
		},
		Subject:  rendered.Subject,
		HTMLBody: rendered.Content,
		Priority: mail.PriorityHigh,
	})

	return nil
}

// generateOTP generates a 6-digit OTP
func generateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}
