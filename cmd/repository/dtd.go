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
	FindStationsByCrs(crs string) (*models.LocationData, error)
	FindFlowsForStations(src, dst string) ([]*models.FlowData, error)
	FindAllFlowsForStation(nlc string) ([]*models.FlowData, error)
	FindFaresForFlow(flowId uint) ([]*models.FareDetail, error)
}
