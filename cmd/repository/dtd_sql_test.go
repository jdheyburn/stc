package repository

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jdheyburn/stc/cmd/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	findStationsByCrsQuery             = "SELECT `uic`,`nlc`,`description`,`crs`,`fare_group`,`start_date`,`end_date` FROM `location` WHERE crs = ? AND start_date <= CURDATE() AND end_date > CURDATE()"
	findFlowsForStationsQuery          = "SELECT flow.flow_id,flow.origin_code,flow.destination_code,flow.direction,flow.start_date,flow.end_date,flow.route_code,route.description as route_desc FROM `flow` LEFT JOIN route on flow.route_code = route.route_code WHERE (flow.origin_code = ?) AND flow.destination_code = ? AND flow.start_date <= CURDATE() AND flow.end_date > CURDATE() AND route.start_date <= CURDATE() AND route.end_date > CURDATE()"
	findFlowsForStationsDirectionQuery = "SELECT flow.flow_id,flow.origin_code,flow.destination_code,flow.direction,flow.start_date,flow.end_date,flow.route_code,route.description as route_desc FROM `flow` LEFT JOIN route on flow.route_code = route.route_code WHERE (flow.origin_code = ?) AND flow.destination_code = ? AND flow.start_date <= CURDATE() AND flow.end_date > CURDATE() AND route.start_date <= CURDATE() AND route.end_date > CURDATE() AND flow.direction = 'R'"
	findAllFlowsForStationQuery        = "SELECT flow.flow_id,flow.origin_code,flow.destination_code,flow.direction,flow.start_date,flow.end_date,flow.route_code,route.description as route_desc FROM `flow` LEFT JOIN route on flow.route_code = route.route_code WHERE ((origin_code = ?) OR destination_code = ?) AND start_date <= CURDATE() AND end_date > CURDATE()"
	findFaresForFlowQuery              = "SELECT fare.id,fare.flow_id,fare.ticket_code,fare.fare,fare.restriction_code,ticket_type.description as ticket_description,ticket_type.tkt_class as ticket_class,ticket_type.tkt_type as ticket_type,restriction_header.description as restriction_desc,restriction_header.desc_out as restriction_desc_out,restriction_header.desc_ret as restriction_desc_rtn FROM `fare` LEFT JOIN ticket_type on fare.ticket_code = ticket_type.ticket_code LEFT JOIN restriction_header on fare.restriction_code = restriction_header.restriction_code WHERE fare.flow_id IN (?) AND ticket_type.start_date <= CURDATE() AND ticket_type.end_date > CURDATE()"

	findStationsByCrsQueryNew = "with grouped_locations as ( select lgm.member_uic_code , lgm.member_crs_code , lgm.group_uic_code, lg.description from location_group_member lgm left join location_group lg on lgm.group_uic_code = lg.group_uic_code where lgm.end_date > CURDATE() and lg.start_date <= CURDATE() AND lg.end_date > CURDATE() ) select location.uic , location.nlc , location.crs , location.description ,location.fare_group , location.start_date , location.end_date , grouped_locations.group_uic_code , grouped_locations.description as group_description from location left join grouped_locations on location.uic = grouped_locations.member_uic_code where location.crs = '?' and location.start_date <= CURDATE() and location.end_date > CURDATE();"
)

func newDateField(year int, month time.Month, day int) *time.Time {
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &d
}

var infiniteTime = newDateField(2999, 12, 31)

func newMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.22"))

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return gdb, mock
}

func TestDtdRepositorySql_FindStationsByCrs(t *testing.T) {

	db, mock := newMock()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		crs string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setUp   func(args)
		want    *models.LocationWithGroups
		wantErr error
	}{
		{
			name: "should return one station location given existing CRS",
			fields: fields{
				db: db,
			},
			args: args{
				crs: "SNR",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"uic", "nlc", "crs", "description", "fare_group", "start_date", "end_date", "group_uic_code", "group_description"}).
					AddRow("7054330", "5433", "SNR", "SANDERSTEAD", "5433", newDateField(2020, 9, 9), infiniteTime, nil, nil)
				mock.ExpectQuery(regexp.QuoteMeta(findStationsByCrsQueryNew)).WithArgs("SNR").WillReturnRows(rows)
			},
			want: &models.LocationWithGroups{
				UIC:         "7054330",
				StartDate:   newDateField(2020, 9, 9),
				EndDate:     infiniteTime,
				NLC:         "5433",
				Description: "SANDERSTEAD",
				CRS:         "SNR",
				FareGroup:   "5433",
				Groups:      []*models.LocationGroup{},
			},
		},
		{
			name: "should return not found error given no records found",
			fields: fields{
				db: db,
			},
			args: args{
				crs: "NOPE",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"uic", "nlc", "crs", "description", "fare_group", "start_date", "end_date", "group_uic_code", "group_description"})
				mock.ExpectQuery(regexp.QuoteMeta(findStationsByCrsQueryNew)).WithArgs("NOPE").WillReturnRows(rows)
			},
			wantErr: ErrNotFound,
		},
		{
			name: "should return grouped station location given CRS for grouped station",
			fields: fields{
				db: db,
			},
			args: args{
				crs: "MCV",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"uic", "nlc", "crs", "description", "fare_group", "start_date", "end_date", "group_uic_code", "group_description"}).
					AddRow("7029700", "2970", "MCV", "MANCHESTER VIC", "0438", newDateField(2020, 9, 9), infiniteTime, "7004380", "MANCHESTER STNS").
					AddRow("7029700", "2970", "MCV", "MANCHESTER VIC", "0438", newDateField(2020, 9, 9), infiniteTime, "70L0050", "GM METROLNK Z1-4").
					AddRow("7029700", "2970", "MCV", "MANCHESTER VIC", "0438", newDateField(2020, 9, 9), infiniteTime, "70L0090", "METROLINK Z1-2")
				mock.ExpectQuery(regexp.QuoteMeta(findStationsByCrsQueryNew)).WithArgs("MCV").WillReturnRows(rows)
			},
			want: &models.LocationWithGroups{
				UIC:         "7029700",
				StartDate:   newDateField(2020, 9, 9),
				EndDate:     infiniteTime,
				NLC:         "2970",
				Description: "MANCHESTER VIC",
				CRS:         "MCV",
				FareGroup:   "0438",
				Groups: []*models.LocationGroup{
					{
						UIC:         "7004380",
						Description: "MANCHESTER STNS",
					},
					{
						UIC:         "70L0050",
						Description: "GM METROLNK Z1-4",
					},
					{
						UIC:         "70L0090",
						Description: "METROLINK Z1-2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dtd := &DtdRepositorySql{
				db: tt.fields.db,
			}
			tt.setUp(tt.args)
			got, err := dtd.FindStationsByCrs(tt.args.crs)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
				return
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Equal(t, tt.want, got)
			if err := mock.ExpectationsWereMet(); err != nil {
				assert.Fail(t, "Not all mocks hit", err)
			}
		})
	}
}

func TestDtdRepositorySql_FindFlowsForStations(t *testing.T) {

	db, mock := newMock()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setUp   func(args)
		want    []*models.FlowDetail
		wantErr error
	}{
		{
			name: "should return flows for src and dst",
			fields: fields{
				db: db,
			},
			args: args{
				src: "5432",
				dst: "5433",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "direction", "start_date", "end_date", "route_code", "route_desc"}).
					AddRow("136210", "5432", "5433", "R", newDateField(2020, 1, 3), infiniteTime, "01000", ".")
				mock.ExpectQuery(regexp.QuoteMeta(findFlowsForStationsQuery)).WithArgs("5432", "5433").WillReturnRows(rows)
			},
			want: []*models.FlowDetail{
				{
					FlowID:          "136210",
					OriginCode:      "5432",
					DestinationCode: "5433",
					RouteCode:       "01000",
					Direction:       "R",
					StartDate:       newDateField(2020, 1, 3),
					EndDate:         infiniteTime,
					RouteDesc:       ".",
				},
			},
		},
		{
			name: "should query reverse of route given empty first result",
			fields: fields{
				db: db,
			},
			args: args{
				src: "5433",
				dst: "5432",
			},
			setUp: func(a args) {
				firstQuery := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "direction", "start_date", "end_date", "route_code", "route_desc"})
				mock.ExpectQuery(regexp.QuoteMeta(findFlowsForStationsQuery)).WithArgs("5433", "5432").WillReturnRows(firstQuery)
				secondQuery := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "direction", "start_date", "end_date", "route_code", "route_desc"}).
					AddRow("136210", "5432", "5433", "R", newDateField(2020, 1, 3), infiniteTime, "01000", ".")
				mock.ExpectQuery(regexp.QuoteMeta(findFlowsForStationsDirectionQuery)).WithArgs("5432", "5433").WillReturnRows(secondQuery)
			},
			want: []*models.FlowDetail{
				{
					FlowID:          "136210",
					OriginCode:      "5432",
					DestinationCode: "5433",
					RouteCode:       "01000",
					Direction:       "R",
					StartDate:       newDateField(2020, 1, 3),
					EndDate:         infiniteTime,
					RouteDesc:       ".",
				},
			},
		},
		{
			name: "should return not found error if no records found",
			fields: fields{
				db: db,
			},
			args: args{
				src: "5433",
				dst: "5432",
			},
			setUp: func(a args) {
				firstQuery := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "direction", "start_date", "end_date", "route_code", "route_desc"})
				mock.ExpectQuery(regexp.QuoteMeta(findFlowsForStationsQuery)).WithArgs("5433", "5432").WillReturnRows(firstQuery)
				mock.ExpectQuery(regexp.QuoteMeta(findFlowsForStationsDirectionQuery)).WithArgs("5432", "5433").WillReturnRows(firstQuery)
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dtd := &DtdRepositorySql{
				db: tt.fields.db,
			}
			tt.setUp(tt.args)
			got, err := dtd.FindFlowsForStations(tt.args.src, tt.args.dst)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
				return
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Equal(t, tt.want, got)
			if err := mock.ExpectationsWereMet(); err != nil {
				assert.Fail(t, "Not all mocks hit", err)
			}
		})
	}
}

func TestDtdRepositorySql_FindAllFlowsForStation(t *testing.T) {

	db, mock := newMock()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		nlc string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setUp   func(args)
		want    []*models.FlowDetail
		wantErr error
	}{
		{
			name: "should return multiple flows for stations",
			fields: fields{
				db: db,
			},
			args: args{
				nlc: "5433",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "direction", "start_date", "end_date", "route_code", "route_desc"}).
					AddRow("6145", "0258", "5433", "R", newDateField(2020, 1, 2), infiniteTime, "00700", "NOT VIA LONDON").
					AddRow("44089", "1402", "5433", "S", newDateField(2020, 5, 18), infiniteTime, "00000", "ANY PERMITTED").
					AddRow("135925", "5433", "5417", "S", newDateField(2020, 1, 2), infiniteTime, "01000", ".").
					AddRow("132215", "5433", "5611", "R", newDateField(2020, 1, 2), infiniteTime, "00700", "NOT VIA LONDON")
				mock.ExpectQuery(regexp.QuoteMeta(findAllFlowsForStationQuery)).WithArgs("5433", "5433").WillReturnRows(rows)
			},
			want: []*models.FlowDetail{
				{
					FlowID:          "6145",
					OriginCode:      "0258",
					DestinationCode: "5433",
					RouteCode:       "00700",
					Direction:       "R",
					StartDate:       newDateField(2020, 1, 2),
					EndDate:         infiniteTime,
					RouteDesc:       "NOT VIA LONDON",
				},
				{
					FlowID:          "44089",
					OriginCode:      "1402",
					DestinationCode: "5433",
					RouteCode:       "00000",
					Direction:       "S",
					StartDate:       newDateField(2020, 5, 18),
					EndDate:         infiniteTime,
					RouteDesc:       "ANY PERMITTED",
				},
				{
					FlowID:          "135925",
					OriginCode:      "5433",
					DestinationCode: "5417",
					RouteCode:       "01000",
					Direction:       "S",
					StartDate:       newDateField(2020, 1, 2),
					EndDate:         infiniteTime,
					RouteDesc:       ".",
				},
				{
					FlowID:          "132215",
					OriginCode:      "5433",
					DestinationCode: "5611",
					RouteCode:       "00700",
					Direction:       "R",
					StartDate:       newDateField(2020, 1, 2),
					EndDate:         infiniteTime,
					RouteDesc:       "NOT VIA LONDON",
				},
			},
		},
		{
			name: "should return not found error given no flows",
			fields: fields{
				db: db,
			},
			args: args{
				nlc: "5433",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"flow_id", "origin_code", "destination_code", "route_code", "direction", "start_date", "end_date"})
				mock.ExpectQuery(regexp.QuoteMeta(findAllFlowsForStationQuery)).WithArgs("5433", "5433").WillReturnRows(rows)
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dtd := &DtdRepositorySql{
				db: tt.fields.db,
			}
			tt.setUp(tt.args)
			got, err := dtd.FindAllFlowsForStation(tt.args.nlc)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
				return
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Equal(t, tt.want, got)
			if err := mock.ExpectationsWereMet(); err != nil {
				assert.Fail(t, "Not all mocks hit", err)
			}
		})
	}
}

func TestDtdRepositorySql_FindFaresForFlow(t *testing.T) {

	db, mock := newMock()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		flowId uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setUp     func(args)
		wantFares []*models.FareDetail
		wantErr   error
	}{
		{
			name: "should return fares for flow",
			fields: fields{
				db: db,
			},
			args: args{
				flowId: 136210,
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"id", "flow_id", "ticket_code", "fare", "restriction_code", "ticket_description", "ticket_class", "ticket_type", "restriction_desc", "restriction_desc_out", "restriction_desc_rtn"}).
					AddRow(1051573, 136210, "0AE", 450, nil, "SMART SDR", 2, "R", nil, nil, nil).
					AddRow(1051574, 136210, "0AF", 270, nil, "SMART SDS", 2, "S", nil, nil, nil).
					AddRow(1051578, 136210, "PAP", 240, "PF", "PAYG PEAK INFO", 2, "S", "PAYG PEAK INFO", "PAY AS YOU GO PEAK - OYSTER CARD REQUIRED", "PAY AS YOU GO PEAK - OYSTER CARD REQUIRED").
					AddRow(1051579, 136210, "POP", 220, "PG", "PAYG OFFPK INFO", 2, "S", "PAYG OFF-PEAK INFO", "PAY AS YOU GO OFF-PEAK - OYSTER CARD REQUIRED", "PAY AS YOU GO OFF-PEAK - OYSTER CARD REQUIRED")
				mock.ExpectQuery(regexp.QuoteMeta(findFaresForFlowQuery)).WithArgs(136210).WillReturnRows(rows)
			},
			wantFares: []*models.FareDetail{
				{
					Model:             gorm.Model{ID: 1051573},
					FlowID:            136210,
					TicketCode:        "0AE",
					Fare:              450,
					RestrictionCode:   "",
					TicketDescription: "SMART SDR",
					TicketClass:       2,
					TicketType:        "R",
				},
				{
					Model:             gorm.Model{ID: 1051574},
					FlowID:            136210,
					TicketCode:        "0AF",
					Fare:              270,
					RestrictionCode:   "",
					TicketDescription: "SMART SDS",
					TicketClass:       2,
					TicketType:        "S",
				},
				{
					Model:              gorm.Model{ID: 1051578},
					FlowID:             136210,
					TicketCode:         "PAP",
					Fare:               240,
					RestrictionCode:    "PF",
					TicketDescription:  "PAYG PEAK INFO",
					TicketClass:        2,
					TicketType:         "S",
					RestrictionDesc:    "PAYG PEAK INFO",
					RestrictionDescOut: "PAY AS YOU GO PEAK - OYSTER CARD REQUIRED",
					RestrictionDescRtn: "PAY AS YOU GO PEAK - OYSTER CARD REQUIRED",
				},
				{
					Model:              gorm.Model{ID: 1051579},
					FlowID:             136210,
					TicketCode:         "POP",
					Fare:               220,
					RestrictionCode:    "PG",
					TicketDescription:  "PAYG OFFPK INFO",
					TicketClass:        2,
					TicketType:         "S",
					RestrictionDesc:    "PAYG OFF-PEAK INFO",
					RestrictionDescOut: "PAY AS YOU GO OFF-PEAK - OYSTER CARD REQUIRED",
					RestrictionDescRtn: "PAY AS YOU GO OFF-PEAK - OYSTER CARD REQUIRED",
				},
			},
		},
		{
			name: "should return error given no rows",
			fields: fields{
				db: db,
			},
			args: args{
				flowId: 136210,
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"id", "flow_id", "ticket_code", "fare", "restriction_code", "description", "ticket_class"})
				mock.ExpectQuery(regexp.QuoteMeta(findFaresForFlowQuery)).WithArgs(136210).WillReturnRows(rows)
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dtd := &DtdRepositorySql{
				db: tt.fields.db,
			}
			tt.setUp(tt.args)
			gotFares, err := dtd.FindFaresForFlow(tt.args.flowId)
			if err != nil && tt.wantErr == nil {
				assert.Fail(t, fmt.Sprintf(
					"Error not expected but got one:\n"+
						"error: %q", err),
				)
				return
			}
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			assert.Equal(t, tt.wantFares, gotFares)
			if err := mock.ExpectationsWereMet(); err != nil {
				assert.Fail(t, "Not all mocks hit", err)
			}
		})
	}
}
