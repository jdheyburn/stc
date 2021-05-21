package cmd

import (
	"math"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jdheyburn/stc/cmd/models"
	"github.com/jdheyburn/stc/cmd/repository"
)

var logger, _ = zap.NewDevelopment()
var fromStation, toStation string
var season bool

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()
	rootCmd.AddCommand(calcCmd)
	calcCmd.Flags().StringVarP(&fromStation, "from", "f", "", "Origin station CRS code")
	calcCmd.Flags().StringVarP(&toStation, "to", "t", "", "Destination station CRS code")
	calcCmd.Flags().BoolVarP(&season, "season", "s", false, "Whether to lookup season tickets only")
	calcCmd.MarkFlagRequired("from")
	calcCmd.MarkFlagRequired("to")
}

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Calculate a season ticket",
	Long:  `TBC`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("Arguments", zap.String("from", fromStation), zap.String("to", toStation), zap.Bool("season", season))
		if err := calc(fromStation, toStation, season); err != nil {
			logger.Error("error running calc", zap.Error(err))
			os.Exit(1)
		}
	},
}

type Fares struct {
	WeeklyStd       float64
	MonthlyStd      float64
	ThreeMonthlyStd float64
	SixMonthlyStd   float64
	AnnualStd       float64
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func calculateFares(weeklyFare uint) *Fares {
	unit := 0.1
	// i, err := strconv.ParseFloat(weeklyFare, 64)
	// if err != nil {
	// 	panic(err)
	// }
	week := float64(weeklyFare) / 100
	monthly := Round(week*3.84, unit)
	threeMonthly := Round(week*3.84*3, unit)
	sixMonthly := Round(week*3.84*6, unit)
	annual := Round(week*40, unit)
	return &Fares{
		WeeklyStd:       week,
		MonthlyStd:      monthly,
		ThreeMonthlyStd: threeMonthly,
		SixMonthlyStd:   sixMonthly,
		AnnualStd:       annual,
	}
}

type GetFaresConfig struct {
	Repo        *repository.DtdRepositorySql
	FromStation string
	ToStation   string
	Season      bool
	Class       string
}

func GetFares(cfg *GetFaresConfig) ([]*models.FareDetailExtreme, error) {

	logger.Info("searching for fares with config", zap.Any("cfg", cfg))

	src, err := cfg.Repo.FindStationsByCrs(cfg.FromStation)

	if err != nil {
		return nil, errors.Wrapf(err, "finding stations for source crs")
	}

	logger.Debug("found station for crs", zap.String("crs", cfg.FromStation), zap.Any("station", src))

	dst, err := cfg.Repo.FindStationsByCrs(cfg.ToStation)

	if err != nil {
		return nil, errors.Wrapf(err, "finding stations for destination crs")
	}

	logger.Debug("found station for crs", zap.String("crs", cfg.ToStation), zap.Any("station", src))

	srcNlcs, err := cfg.Repo.FindNLCsRelatedToCrs(src[0].CRS)

	if err != nil {
		return nil, errors.Wrapf(err, "finding NLCs related to source CRS")
	}

	logger.Debug("found NLCs related to crs", zap.String("crs", cfg.FromStation), zap.Any("nlcs", srcNlcs))

	dstNlcs, err := cfg.Repo.FindNLCsRelatedToCrs(dst[0].CRS)

	if err != nil {
		return nil, errors.Wrapf(err, "finding NLCs related to destination CRS")
	}

	logger.Debug("found NLCs related to crs", zap.String("crs", cfg.ToStation), zap.Any("nlcs", dstNlcs))

	fares, err := cfg.Repo.FindFaresForNLCs(srcNlcs, dstNlcs, cfg.Season, cfg.Class)

	if err != nil {
		return nil, errors.Wrapf(err, "finding fares for src and dst NLCs")
	}

	logger.Info("found fares for src and dst NLCs", zap.Int("numFares", len(fares)))

	if !cfg.Season {
		logger.Info("season ticket not specified, retrieving fare overrides")
		overrides, err := cfg.Repo.FindFareOverridesForNLCs(srcNlcs, dstNlcs)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieving fare overrides")
		}

		logger.Info("found fare overrides", zap.Int("numFares", len(overrides)))

		for _, fare := range overrides {
			fares = append(fares, fare)
		}
	}

	return fares, nil
}

// Kinda using this just for testing locally atm
func calc(fromStation, toStation string, season bool) error {

	opts := &repository.DtdSqlDBOptions{
		User:     "root",
		Password: "password123",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "fares",
	}
	repo, err := repository.NewDtdRepositorySql(opts)
	if err != nil {
		panic(err)
	}

	cfg := &GetFaresConfig{
		Repo:        repo,
		FromStation: fromStation,
		ToStation:   toStation,
		Season:      season,
		Class:       "2",
	}

	fares, err := GetFares(cfg)

	printer := tableprinter.New(os.Stdout)
	printer.Print(fares)

	return nil
}
