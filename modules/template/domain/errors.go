package domain

import "github.com/duongptryu/gox/syserr"

// Template domain errors
var (
	ErrTemplateNotFound      = syserr.New(syserr.NotFoundCode, "template not found")
	ErrTemplateAlreadyExists = syserr.New(syserr.ConflictCode, "template already exists")
	ErrInvalidTemplateType   = syserr.New(syserr.InvalidArgumentCode, "invalid template type")
	ErrInvalidTemplateStatus = syserr.New(syserr.InvalidArgumentCode, "invalid template status")
	ErrTemplateInactive      = syserr.New(syserr.ForbiddenCode, "template is inactive")
	ErrTemplateRenderFailed  = syserr.New(syserr.InternalCode, "template rendering failed")
	ErrInvalidTemplateSlug   = syserr.New(syserr.InvalidArgumentCode, "invalid template slug")
	ErrTemplateSyntaxError   = syserr.New(syserr.InvalidArgumentCode, "template syntax error")
)
