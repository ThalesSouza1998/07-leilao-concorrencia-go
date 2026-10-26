package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	registrador *zap.Logger
)

func init() {
	configuracao := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	registrador, _ = configuracao.Build()
}

func Informacao(mensagem string, campos ...zap.Field) {
	registrador.Info(mensagem, campos...)
	registrador.Sync()
}

func Erro(mensagem string, erro error, campos ...zap.Field) {
	campos = append(campos, zap.NamedError("error", erro))
	registrador.Error(mensagem, campos...)
	registrador.Sync()
}
