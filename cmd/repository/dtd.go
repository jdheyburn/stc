package repository

import (
	"github.com/jdheyburn/stc/cmd/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger, _ = zap.NewDevelopment()

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()
}

// DtdRepository provides an abstraction between databases
type DtdRepository interface {
	FindStationByCRS(crs string) (*models.StationData, error)
}