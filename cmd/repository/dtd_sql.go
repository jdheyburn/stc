package repository

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdheyburn/stc/cmd/models"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
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

	dbLogger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			LogLevel: glogger.Info,
			Colorful: true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		return nil, errors.Wrap(err, "creating sql conn")
	}

	fmt.Println("Successfully connected!")

	return &DtdRepositorySql{
		db: db,
	}, nil
}

// FindStationsByCrs returns locations from the given CRS code
func (dtd *DtdRepositorySql) FindStationsByCrs(crs string) (stations []*models.LocationData, err error) {

	logger.Infof("looking up CRS %v", crs)

	err = dtd.db.Unscoped().
		Select("uic", "nlc", "description", "crs", "fare_group", "start_date", "end_date").
		Where("crs = ?", crs).
		Where("start_date <= CURDATE()").
		Where("end_date > CURDATE()").
		Find(&stations).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for location %s", crs)
	}

	if len(stations) > 0 {
		logger.Infof("found %v from crs %v", stations[0].Description, crs)
		return stations, nil
	}

	return nil, ErrNotFound
}

func (dtd *DtdRepositorySql) FindNLCsRelatedToCrs(crs string) (nlcs []string, err error) {

	logger.Infof("looking up NLCs related to CRS %v", crs)

	err = dtd.db.Raw(nlcs_query, crs, crs).Scan(&nlcs).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for NLCs related to CRS %s", crs)
	}

	return nlcs, nil
}

func (dtd *DtdRepositorySql) FindFaresForNLCs(srcNlcs, dstNlcs []string, season bool, class string) (fares []*models.FareDetailExtreme, err error) {

	logger.Infof("looking up fares related to nlcs")

	err = dtd.db.Raw(fares_query, srcNlcs, dstNlcs, dstNlcs, srcNlcs).Scan(&fares).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for fares related to nlcs")
	}

	// TODO replace with SQL query filter
	if season {
		var filtered []*models.FareDetailExtreme
		for _, fare := range fares {
			if fare.TicketType == "N" {
				filtered = append(filtered, fare)
			}
		}
		return filtered, nil
	}

	return fares, nil
}

func (dtd *DtdRepositorySql) FindFareOverridesForNLCs(srcNlcs, dstNlcs []string) (fares []*models.FareDetailExtreme, err error) {

	logger.Infof("looking up fares overrides related to nlcs")

	err = dtd.db.Raw(nfo_query, srcNlcs, dstNlcs).Scan(&fares).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for fares overrides related to nlcs")
	}

	return fares, nil
}

func (dtd *DtdRepositorySql) FindFlowsForNLCs(srcNlcs []string, dstNlcs []string) (flows []*models.FlowDetail, err error) {

	logger.Infof("looking up flows matching src and dst NLCs")

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
		Where(
			dtd.db.Where("flow.origin_code in ?", srcNlcs).Where("flow.destination_code in ?", dstNlcs),
		).
		Or(dtd.db.Where("flow.origin_code in ?", dstNlcs).Where("flow.destination_code in ?", srcNlcs).Where("flow.direction = 'R'")).
		Where("flow.start_date <= CURDATE()").
		Where("flow.end_date > CURDATE()").
		Where("route.start_date <= CURDATE()").
		Where("route.end_date > CURDATE()").
		Find(&flows).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying for flows")
	}

	return flows, nil
}

func (dtd *DtdRepositorySql) findFlows(src, dst string, reversed bool) (flows []*models.FlowDetail, err error) {

	chain := dtd.db.Unscoped().Model(&models.FlowData{}).
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
		Where("flow.origin_code = ?", src).
		Where("flow.destination_code = ?", dst).
		Where("flow.start_date <= CURDATE()").
		Where("flow.end_date > CURDATE()").
		Where("route.start_date <= CURDATE()").
		Where("route.end_date > CURDATE()")

	if reversed {
		chain = chain.Where("flow.direction = 'R'")
	}

	err = chain.Find(&flows).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying flows for source: %s dest: %s", src, dst)
	}

	return flows, err
}

// FindFlowsForStations returns all flows between two NLC codes
func (dtd *DtdRepositorySql) FindFlowsForStations(src, dst string) (flows []*models.FlowDetail, err error) {

	logger.Infof("searching for all flows between src %v and dst %v", src, dst)

	reversed := false
	flows, err = dtd.findFlows(src, dst, reversed)
	if err != nil {
		return nil, err
	}

	if len(flows) > 0 {
		logger.Infof("returning %v flows between src %v and dst %v", len(flows), src, dst)
		return flows, nil
	}

	logger.Warnf("found no flows for src %v and dst %v - searching flows in the reverse direction", src, dst)
	reversed = true
	flows, err = dtd.findFlows(dst, src, reversed)
	if err != nil {
		return nil, err
	}

	if len(flows) > 0 {
		logger.Infof("returning %v flows between src %v and dst %v", len(flows), src, dst)
		return flows, nil
	}

	return nil, ErrNotFound
}

func (dtd *DtdRepositorySql) FindAllFlowsForStation(nlc string) (flows []*models.FlowDetail, err error) {

	logger.Infof("searching for all flows for nlc %v", nlc)
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
		Where(dtd.db.Where("origin_code = ?", nlc).Or("destination_code = ?", nlc)).
		Where("start_date <= CURDATE()").
		Where("end_date > CURDATE()").
		Find(&flows).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying all flows for NLC %s", nlc)
	}

	if len(flows) == 0 {
		return nil, ErrNotFound
	}

	logger.Infof("returning %v flows for station nlc &v", len(flows), nlc)

	return flows, nil
}

func (dtd *DtdRepositorySql) FindFaresForFlows(flowIds []string) (fares []*models.FareDetail, err error) {

	logger.Infof("finding fares for flowIDs %v", flowIds)

	err = dtd.db.Unscoped().Model(&models.FareData{}).
		Distinct(
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
		// TODO this join causes duplicate records (need a better query - Distinct is a workaround)
		Joins("LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code").
		Where("fare.flow_id IN ?", flowIds).
		Where("ticket_type.start_date <= CURDATE()").
		Where("ticket_type.end_date > CURDATE()").
		Where("ticket_type.tkt_type = 'N'"). // Season tickets only
		Where("ticket_type.tkt_class = 2").  // 2nd class only
		Order("fare ASC").
		Scan(&fares).Error

	if err != nil {
		return nil, errors.Wrapf(err, "querying all fares for flowIDs %v", flowIds)
	}

	if len(fares) == 0 {
		return nil, ErrNotFound
	}

	logger.Infof("returning %v fares for flowIDs &v", len(fares), flowIds)

	return fares, nil
}
