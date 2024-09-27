package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_ReportsByRadius(t *testing.T) {

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
		radius int
		point  storage.Geo
		status []storage.Status
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "OK Two reports",
			args: args{
				radius: 3000,
				point: storage.Geo{
					Type:        "Point",
					Coordinates: [2]float64{55.75583793441133, 37.620437229927795},
				},
				status: nil,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "OK One report",
			args: args{
				radius: 1000,
				point: storage.Geo{
					Type:        "Point",
					Coordinates: [2]float64{59.93820021009978, 30.3172182590169},
				},
				status: []storage.Status{},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error Not found by status",
			args: args{
				radius: 3000,
				point: storage.Geo{
					Type:        "Point",
					Coordinates: [2]float64{55.75583793441133, 37.620437229927795},
				},
				status: []storage.Status{2},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error Not found by radius",
			args: args{
				radius: 3000,
				point: storage.Geo{
					Type:        "Point",
					Coordinates: [2]float64{56.866556467082695, 35.91217524893142},
				},
				status: nil,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error Incorrect radius",
			args: args{
				radius: -1000,
				point: storage.Geo{
					Type:        "Point",
					Coordinates: [2]float64{55.75583793441133, 37.620437229927795},
				},
				status: nil,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.ReportsByRadius(context.Background(), tt.args.radius, tt.args.point, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReportsByRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Storage.ReportsByRadius() len = %v, want %v", len(got), tt.want)
			}
		})
	}
}
