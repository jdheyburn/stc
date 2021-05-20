package cmd

import (
	"fmt"
	"math"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jdheyburn/stc/cmd/repository"
)

var logger, _ = zap.NewDevelopment()
var fromStation, toStation string

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()
	rootCmd.AddCommand(calcCmd)
	calcCmd.Flags().StringVarP(&fromStation, "from", "f", "", "Where the season ticket should start from")
	calcCmd.Flags().StringVarP(&toStation, "to", "t", "", "Where the season ticket should end at")
	calcCmd.MarkFlagRequired("from")
	calcCmd.MarkFlagRequired("to")
}

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Calculate a season ticket",
	Long:  `TBC`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("Arguments", zap.String("from", fromStation), zap.String("to", toStation))
		if err := calc(fromStation, toStation); err != nil {
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

func calc(fromStation, toStation string) error {

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

	// 1. Get the NLCs from CRS
	// TODO don't assume the user has input CRS

	src, err := repo.FindStationsByCrs(fromStation)

	if err != nil {
		panic(err)
	}

	dst, err := repo.FindStationsByCrs(toStation)

	if err != nil {
		panic(err)
	}

	srcNlcs, err := repo.FindNLCsRelatedToCrs(src[0].CRS)

	if err != nil {
		panic(err)
	}

	dstNlcs, err := repo.FindNLCsRelatedToCrs(dst[0].CRS)

	if err != nil {
		panic(err)
	}
	logger.Info("src")
	logger.Info(fmt.Sprint(srcNlcs))
	logger.Info("dst")
	logger.Info(fmt.Sprint(dstNlcs))

	// flows, err := repo.FindFlowsForNLCs(srcNlcs, dstNlcs)

	// if err != nil {
	// 	panic(err)
	// }

	// logger.Info("flows")
	// logger.Info(fmt.Sprint(flows))

	// var flowIds []string
	// for _, flow := range flows {
	// 	flowIds = append(flowIds, flow.FlowID)
	// }

	// fares, err := repo.FindFaresForFlows(flowIds)

	fares, err := repo.FindFaresForNLCs(srcNlcs, dstNlcs)

	if err != nil {
		panic(err)
	}

	logger.Info("fares")
	var fareIds []uint
	for _, fare := range fares {
		fareIds = append(fareIds, fare.FlowID)
	}

	logger.Info(fmt.Sprint(len(fareIds)))
	logger.Info(fmt.Sprint(fareIds))
	printer := tableprinter.New(os.Stdout)

	printer.Print(fares)

	/* Below commented out while testing NLCs (old flow)
	// 2. Get Flows for stations

	flows, err := repo.FindFlowsForStations(src[0].NLC, dst[0].NLC)

	if err != nil {
		panic(err)
	}

	if len(flows) > 1 {
		panic("More than one flow found")
	}

	// 3. Get fares for flows

	flow := flows[0]
	i, err := strconv.ParseUint(flow.FlowID, 10, 64)
	if err != nil {
		panic(err)
	}
	fares, err := repo.FindFaresForFlow(i)

	// 4. Get the 7DS ticket_code

	// filter down to get just 7DS ticket_code
	// TODO maybe have the DB do the filter?

	found := []*models.FareDetail{}
	for _, fare := range fares {
		if fare.TicketCode == "7DS" {
			found = append(found, fare)
		}
	}

	if len(found) == 0 {
		panic("no 7DS found")
	}

	if len(found) > 1 {
		panic("more than one 7DS found")
	}

	logger.Info(fmt.Sprint(found[0].ID))

	// 5. Now use the fare to calculate the prices
	seasonTickets := calculateFares(found[0].Fare)
	logger.Info(fmt.Sprint(seasonTickets))

	*/

	return nil
}
