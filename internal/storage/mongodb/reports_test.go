package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_Reports(t *testing.T) {

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
		status  []storage.Status
		want    int
		wantErr bool
	}{
		{
			name:    "OK No status",
			status:  nil,
			want:    3,
			wantErr: false,
		},
		{
			name:    "OK Status 1",
			status:  []storage.Status{1},
			want:    3,
			wantErr: false,
		},
		{
			name:    "OK Status 1 & 3",
			status:  []storage.Status{1, 3},
			want:    3,
			wantErr: false,
		},
		{
			name:    "Error Not found status 3",
			status:  []storage.Status{3},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.Reports(context.Background(), tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Reports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Storage.Reports() len = %v, want %v", len(got), tt.want)
			}
		})
	}
}
