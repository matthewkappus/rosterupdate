package store

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/matthewkappus/syncup/src/types"
)

func testDB() *sql.DB {
	return nil
}

func TestRosters_GetRosters(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRosters types.Rosters
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &Rosters{
				DB: tt.fields.DB,
			}
			gotRosters, err := rs.GetRosters(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rosters.GetRosters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRosters, tt.wantRosters) {
				t.Errorf("Rosters.GetRosters() = %v, want %v", gotRosters, tt.wantRosters)
			}
		})
	}
}
