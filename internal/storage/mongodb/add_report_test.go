package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_AddReport(t *testing.T) {

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

	tests := []struct {
		name    string
		rep     storage.Report
		wantErr bool
	}{
		{
			name: "OK",
			rep: storage.Report{
				Number:      1,
				Address:     "Адрес 1",
				Description: "Описание 1",
				Media:       []string{"https://google.com"},
				Geo:         storage.Geo{Coordinates: [2]float64{55.75388130172051, 37.62026781374883}},
			},
			wantErr: false,
		},
		{
			name: "Error Duplicate number",
			rep: storage.Report{
				Number:      1,
				Address:     "Адрес 1",
				Description: "Описание 1",
				Media:       []string{"https://google.com"},
				Geo:         storage.Geo{Coordinates: [2]float64{55.75388130172051, 37.62026781374883}},
			},
			wantErr: true,
		},
		{
			name: "Error Incorrect number",
			rep: storage.Report{
				Number:      -1,
				Address:     "Адрес 1",
				Description: "Описание 1",
				Media:       []string{"https://google.com"},
				Geo:         storage.Geo{Coordinates: [2]float64{55.75388130172051, 37.62026781374883}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.AddReport(context.Background(), tt.rep); (err != nil) != tt.wantErr {
				t.Errorf("Storage.AddReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
