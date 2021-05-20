package repository

import (
	"github.com/jdheyburn/stc/cmd/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var logger, _ = zap.NewDevelopment()
var logger *zap.SugaredLogger

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	slogger, _ := config.Build()
	logger = slogger.Sugar()
}

// DtdRepository provides an abstraction between databases
type DtdRepository interface {
	FindStationsByCrs(crs string) (*models.LocationData, error)
	FindFlowsForStations(src, dst string) ([]*models.FlowData, error)
	FindAllFlowsForStation(nlc string) ([]*models.FlowData, error)
	FindFaresForFlows(flowIds []string) ([]*models.FareDetail, error)
	FindFaresForNLCs(srcNlcs, dstNlcs []string) ([]*models.FareDetailExtreme, error)
	FindFareOverridesForNLCs(srcNlcs, dstNlcs []string) ([]*models.FareDetailExtreme, error)
	FindNLCsRelatedToCrs(crs string) ([]string, error)
	FindFlowsForNLCs(srcNlcs []string, dstNlcs []string) ([]*models.FlowDetail, error)
}
