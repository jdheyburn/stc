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
	findStationsByCrsQuery = "SELECT `uic`,`nlc`,`description`,`crs`,`fare_group`,`start_date`,`end_date` FROM `location` WHERE crs = ? AND start_date <= CURDATE() AND end_date > CURDATE()"
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

func TestDtdRepositorySql_FindStationByCrs_Single(t *testing.T) {

	opts := &DtdSqlDBOptions{
		User:     "root",
		Password: "password123",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "fares",
	}
	repo, err := NewDtdRepositorySql(opts)
	if err != nil {
		panic(err)
	}

	expected := []*models.StationData{
		{
			UIC:         "7054330",
			StartDate:   newDateField(2020, 9, 9),
			EndDate:     infiniteTime,
			NLC:         "5433",
			Description: "SANDERSTEAD",
			CRS:         "SNR",
			FareGroup:   "5433",
		},
	}

	// rows := sqlmock.NewRows([]string{"uic", "start_date", "end_date", "nlc", "description", "crs", "fare_group"}).
	// 	AddRow(expected.UIC, expected.StartDate, expected.EndDate, expected.NLC, expected.Description, expected.CRS, expected.FareGroup)

	// mock.ExpectQuery(fmt.Sprintf(findStationByCrsQueryTest, "SNR")).WillReturnRows(rows)

	actual, err := repo.FindStationsByCrs("SNR")

	assert.NotNil(t, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

// func TestDtdRepositorySql_FindStationByCrs_NoneError(t *testing.T) {

// 	db, mock := newMock()
// 	repo := &DtdRepositorySql{db}

// 	defer func() {
// 		repo.db.Close()
// 	}()

// 	rows := sqlmock.NewRows([]string{"uic", "start_date", "end_date", "nlc", "description", "crs", "fare_group"})

// 	mock.ExpectQuery(fmt.Sprintf(findStationByCrsQueryTest, "SNR")).WillReturnRows(rows)

// 	actual, err := repo.FindStationByCrs("SNR")

// 	assert.Empty(t, actual)
// 	assert.EqualError(t, err, ErrNotFound.Error())
// }

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
		want    []*models.StationData
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
				rows := sqlmock.NewRows([]string{"uic", "nlc", "description", "crs", "fare_group", "start_date", "end_date"}).
					AddRow("7054330", "5433", "SANDERSTEAD", "SNR", "5433", newDateField(2020, 9, 9), infiniteTime)
				mock.ExpectQuery(regexp.QuoteMeta(findStationsByCrsQuery)).WithArgs("SNR").WillReturnRows(rows)
			},
			want: []*models.StationData{
				{
					UIC:         "7054330",
					StartDate:   newDateField(2020, 9, 9),
					EndDate:     infiniteTime,
					NLC:         "5433",
					Description: "SANDERSTEAD",
					CRS:         "SNR",
					FareGroup:   "5433",
				},
			},
		},
		{
			name: "should return not found error given no records found",
			fields: fields{
				db: db,
			},
			args: args{
				crs: "SNR",
			},
			setUp: func(a args) {
				rows := sqlmock.NewRows([]string{"uic", "nlc", "description", "crs", "fare_group", "start_date", "end_date"})
				mock.ExpectQuery(regexp.QuoteMeta(findStationsByCrsQuery)).WithArgs("SNR").WillReturnRows(rows)
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
		})
	}
}

// func TestDtdRepositorySql_FindFlowsForStations(t *testing.T) {

//     db, mock := newMock()

//     type fields struct {
//         db *sql.DB
//     }
//     type args struct {
//         src string
//         dst string
//     }
//     tests := []struct {
//         name    string
//         fields  fields
//         args    args
//         setUp   func(args)
//         want    []*models.FlowData
//         wantErr error
//     }{
//         {
//             name: "should return not found error given no records found",
//             fields: fields{
//                 db: db,
//             },
//             args: args{
//                 src: ""
//                 crs: "SNR",
//             },
//             setUp: func(a args) {
//                 rows := sqlmock.NewRows([]string{"uic", "startsdate", "end_date", "nlc", "description", "crs", "fare_group"})
//                 mock.ExpectQuery(fmt.Sprintf(findStationByCrsQueryTest, a.crs)).WillReturnRows(rows)
//             },
//             wantErr: ErrNotFound,
//         },
//     }
//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             dtd := &DtdRepositorySql{
//                 db: tt.fields.db,
//             }
//             // mock.ExpectQuery(fmt.Sprintf(findStationByCrsQueryTest, tt.args.src, tt.args.dst)).WillReturnRows(tt.returnRows)
//             got, err := dtd.FindFlowsForStations(tt.args.src, tt.args.dst)
//             if err != nil && tt.wantErr == nil {
//                 assert.Fail(t, fmt.Sprintf(
//                     "Error not expected but got one:\n"+
//                         "error: %q", err),
//                 )
//                 return
//             }
//             if tt.wantErr != nil {
//                 assert.EqualError(t, err, tt.wantErr.Error())
//                 return
//             }
//             assert.Equal(t, tt.want, got)
//         })
//     }
// }
