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

// FindStationsByCrs returns locations from the given CRS code
// TODO how to handle multiple locations
func (dtd *DtdRepositorySql) FindStationsByCrs(crs string) ([]*models.LocationData, error) {

	var stations []*models.LocationData
	err := dtd.db.Unscoped().
		Select("uic", "nlc", "description", "crs", "fare_group", "start_date", "end_date").
		Where("crs = ? AND start_date <= CURDATE() AND end_date > CURDATE()", crs).
		Find(&stations).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for location %s:", crs)
	}

	if len(stations) > 0 {
		return stations, nil
	}

	// Check to see if this is a grouped station (i.e. LBG -> London Terminals)
	// Query is in dbeaver
	// err := dtd.db.Unscoped().
	// Select()

	return nil, ErrNotFound

}

func (dtd *DtdRepositorySql) findFlows(src, dst string, reversed bool) (flows []*models.FlowDetail, err error) {

	var directionFilter string
	if reversed {
		directionFilter = " AND flow.direction = 'R'"
	}

	err = dtd.db.Unscoped().Model(&models.FlowData{}).
		Select(
			"flow.flow_id",
			"flow.origin_code",
			"flow.destination_code",
			"flow.direction",
			"flow.start_date",
			"flow.end_date",
			"flow.route_code",
			"route.description as route_desc",
		).
		Joins("LEFT JOIN route on flow.route_code = route.route_code").
		Where(fmt.Sprintf("flow.origin_code = ? AND flow.destination_code = ? AND flow.start_date <= CURDATE() AND flow.end_date > CURDATE() AND route.start_date <= CURDATE() AND route.end_date > CURDATE()%s", directionFilter), src, dst).
		Find(&flows).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying flows for source: %s dest: %s", src, dst)
	}

	return flows, err
}

// FindFlowsForStations returns all flows between two NLC codes
// TODO need to find flows where the second query is not in direction R
func (dtd *DtdRepositorySql) FindFlowsForStations(src, dst string) (flows []*models.FlowDetail, err error) {

	reversed := false
	flows, err = dtd.findFlows(src, dst, reversed)
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

func (dtd *DtdRepositorySql) FindAllFlowsForStation(nlc string) (flows []*models.FlowDetail, err error) {

	err = dtd.db.Unscoped().Model(&models.FlowData{}).
		Select(
			"flow.flow_id",
			"flow.origin_code",
			"flow.destination_code",
			"flow.direction",
			"flow.start_date",
			"flow.end_date",
			"flow.route_code",
			"route.description as route_desc",
		).
		Joins("LEFT JOIN route on flow.route_code = route.route_code").
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
		Select(
			"fare.id",
			"fare.flow_id",
			"fare.ticket_code",
			"fare.fare",
			"fare.restriction_code",
			"ticket_type.description as ticket_description",
			"ticket_type.tkt_class as ticket_class",
			"ticket_type.tkt_type as ticket_type",
			"restriction_header.description as restriction_desc",
			"restriction_header.desc_out as restriction_desc_out",
			"restriction_header.desc_ret as restriction_desc_rtn",
		).
		Joins("LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code").
		Joins("LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code").
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
