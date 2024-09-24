package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"testing"
)

func TestStorage_UpdateStatus(t *testing.T) {

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

	// Вставляем в коллекцию тестовую заявку.
	_, err = st.addOne(reports[0])
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		num        int
		status     storage.Status
		wantStatus storage.Status
		wantErr    bool
	}{
		{
			name:       "OK",
			num:        1,
			status:     2,
			wantStatus: 2,
			wantErr:    false,
		},
		{
			name:       "Error Incorrect number",
			num:        -1,
			status:     2,
			wantStatus: 1,
			wantErr:    true,
		},
		{
			name:       "Error Incorrect status",
			num:        1,
			status:     6,
			wantStatus: 1,
			wantErr:    true,
		},
		{
			name:       "Error Not found",
			num:        5,
			status:     1,
			wantStatus: 1,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.UpdateStatus(context.Background(), tt.num, tt.status)
			if err != nil {
				if tt.wantErr {
					t.Skip()
				}
				t.Errorf("Storage.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Status != tt.wantStatus {
				t.Errorf("Storage.UpdateStatus() status = %d, want %d", got.Status, tt.wantStatus)
			}
		})
	}
}
