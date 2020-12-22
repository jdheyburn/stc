package repository

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdheyburn/stc/cmd/models"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a record cannot be found
var ErrNotFound = errors.New("not found")

// DBOptions holds the information required to construct a DB connection
type DtdSqlDBOptions struct {
	Host, Port, User, Password, DBName string
}

// DtdRepositorySql is a concrete MySql implementation of a DtdRepository
type DtdRepositorySql struct {
	db *gorm.DB
}

func NewDtdRepositorySql(options *DtdSqlDBOptions) (*DtdRepositorySql, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "creating sql conn")
	}

	fmt.Println("Successfully connected!")

	return &DtdRepositorySql{
		db: db,
	}, nil
}

// FindStationsByCrs returns a location from the given CRS code
func (dtd *DtdRepositorySql) FindStationsByCrs(crs string) ([]*models.StationData, error) {

	var stations []*models.StationData
	err := dtd.db.Unscoped().
		Select("uic", "nlc", "description", "crs", "fare_group", "start_date", "end_date").
		Where("crs = ? AND start_date <= CURDATE() AND end_date > CURDATE()", crs).
		Find(&stations).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for location %s:", crs)
	}

	if len(stations) == 0 {
		return nil, ErrNotFound
	}

	return stations, nil
}

func (dtd *DtdRepositorySql) findFlows(src, dst string, reversed bool) ([]*models.FlowData, error) {

	var directionFilter string
	if reversed {
		directionFilter = " AND direction = 'R'"
	}

	var flows []*models.FlowData
	err := dtd.db.Unscoped().
		Select("flow_id", "origin_code", "destination_code", "route_code", "direction", "start_date", "end_date").
		Where(fmt.Sprintf("origin_code = ? AND destination_code = ? AND start_date <= CURDATE() AND end_date > CURDATE()%s", directionFilter), src, dst).
		Find(&flows).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying flows for source: %s dest: %s", src, dst)
	}

	return flows, err
}

// FindFlowsForStations returns all flows between two NLC codes
func (dtd *DtdRepositorySql) FindFlowsForStations(src, dst string) ([]*models.FlowData, error) {

	reversed := false
	flows, err := dtd.findFlows(src, dst, reversed)
	if err != nil {
		return nil, err
	}

	if len(flows) > 0 {
		return flows, nil
	}

	// Else we didn't get any results, query for the other way round
	reversed = true
	flows, err = dtd.findFlows(dst, src, reversed)
	if err != nil {
		return nil, err
	}

	if len(flows) > 0 {
		return flows, nil
	}

	return nil, ErrNotFound
}

// TODO need to find flows where the second query is not in direction R

func (dtd *DtdRepositorySql) FindAllFlowsForStation(nlc string) ([]*models.FlowData, error) {

	var flows []*models.FlowData
	err := dtd.db.Unscoped().
		Select("flow_id", "origin_code", "destination_code", "route_code", "direction", "start_date", "end_date").
		Where("origin_code = ? OR destination_code = ? AND start_date <= CURDATE() AND end_date > CURDATE()", nlc).
		Find(&flows).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying all flows for NLC %s", nlc)
	}

	if len(flows) == 0 {
		return nil, ErrNotFound
	}

	return flows, nil
}

func (dtd *DtdRepositorySql) FindFaresForFlow(flowId uint) (fares []*models.FareDetail, err error) {

	err = dtd.db.Unscoped().Model(&models.FareData{}).
		Distinct("fare.id", "fare.flow_id", "fare.ticket_code", "fare.fare", "fare.restriction_code", "ticket_type.description", "ticket_type.tkt_class as ticket_class").
		Joins("LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code").
		Where("fare.flow_id IN (?) AND ticket_type.start_date <= CURDATE() AND ticket_type.end_date > CURDATE()", flowId).
		Scan(&fares).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying all fares for flowID %v", flowId)
	}

	if len(fares) == 0 {
		return nil, ErrNotFound
	}

	return fares, nil
}
