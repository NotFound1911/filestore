package zap

import (
	"github.com/NotFound1911/filestore/internal/logger/di"
	z "github.com/NotFound1911/filestore/pkg/zap"
	"go.uber.org/zap"
)

var _ di.Logger = &Log{}

type Log struct {
	service *z.Service
}

func (z *Log) Debug(format string, a ...di.Field) {
	z.service.Logger.Debug(format, z.toArgs(a)...)
}

func (z *Log) Info(format string, a ...di.Field) {
	z.service.Logger.Info(format, z.toArgs(a)...)
}

func (z *Log) Warn(format string, a ...di.Field) {
	z.service.Logger.Warn(format, z.toArgs(a)...)
}

func (z *Log) Error(format string, a ...di.Field) {
	z.service.Logger.Error(format, z.toArgs(a)...)
}

func (z *Log) toArgs(args []di.Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Val))
	}
	return res
}
func NewLogger(s *z.Service) di.Logger {
	return &Log{service: s}
}
