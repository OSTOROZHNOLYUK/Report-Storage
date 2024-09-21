package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_ReportsByPoly(t *testing.T) {

	// Создаем пул подключений.
	dbName = "goUnitTestDB"
	colName = "goUnitTestCollection"
	opts := setOpts(path, "admin", os.Getenv("MONGO_DB_PASSWD"))
	st, err := new(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	// Очищаем тестовую коллекцию.
	err = st.trun()
	if err != nil {
		t.Fatal(err)
	}

	// Заполняем коллекцию тестовыми заявками.
	for _, v := range reports {
		_, err := st.addOne(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	type args struct {
		poly   [][2]float64
		status []storage.Status
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Ok Two reports",
			args: args{
				poly: [][2]float64{
					{55.77581155422086, 37.56283858534519},
					{55.77747634327944, 37.68767076147129},
					{55.73347907589756, 37.68093411091992},
					{55.73474339729524, 37.57243320885774},
				},
				status: nil,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Ok One report",
			args: args{
				poly: [][2]float64{
					{59.96147407967736, 30.251580071744584},
					{59.957866124720105, 30.362775713726272},
					{59.92556446632214, 30.365209288382434},
					{59.92978562806268, 30.25710241423357},
				},
				status: nil,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Ok One report with status",
			args: args{
				poly: [][2]float64{
					{59.96147407967736, 30.251580071744584},
					{59.957866124720105, 30.362775713726272},
					{59.92556446632214, 30.365209288382434},
					{59.92978562806268, 30.25710241423357},
				},
				status: []storage.Status{1},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error Not found by status",
			args: args{
				poly: [][2]float64{
					{59.96147407967736, 30.251580071744584},
					{59.957866124720105, 30.362775713726272},
					{59.92556446632214, 30.365209288382434},
					{59.92978562806268, 30.25710241423357},
				},
				status: []storage.Status{2},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error Not found by polygon",
			args: args{
				poly: [][2]float64{
					{59.961346513327165, 30.347008519215102},
					{59.96115909604834, 30.446785080117863},
					{59.91624211554908, 30.465130489064336},
					{59.919995340061696, 30.34569813286178},
				},
				status: nil,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.ReportsByPoly(context.Background(), tt.args.poly, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReportsByPoly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Storage.ReportsByPoly() len = %d, want %d", len(got), tt.want)
			}
		})
	}
}
