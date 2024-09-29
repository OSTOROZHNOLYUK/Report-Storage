package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_ReportsWithFilter(t *testing.T) {

	// Создаем пул подключений.
	dbName = testDatabase
	colReport = testCollection
	opts := setOpts(path, "admin", os.Getenv("MONGO_DB_PASSWD"))
	st, err := new(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	// Очищаем тестовую коллекцию.
	err = st.trun(colReport)
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

	tests := []struct {
		name    string
		filter  storage.Filter
		want    int
		wantErr bool
	}{
		{
			name:    "OK Without filter",
			filter:  storage.Filter{},
			want:    3,
			wantErr: false,
		},
		{
			name:    "OK With filter",
			filter:  storage.Filter{Count: 2, Sort: 1, Status: []storage.Status{1}},
			want:    2,
			wantErr: false,
		},
		{
			name:    "OK Incorrect filter",
			filter:  storage.Filter{Count: -1, Sort: 10, Status: nil},
			want:    3,
			wantErr: false,
		},
		{
			name:    "Error Not found by status",
			filter:  storage.Filter{Count: 2, Status: []storage.Status{2}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.ReportsWithFilter(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReportsWithFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Storage.ReportsWithFilter() len = %d, want %d", len(got), tt.want)
			}
		})
	}
}
