package mongodb

import (
	"context"
	"os"
	"testing"
)

func TestStorage_ReportByID(t *testing.T) {

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
	var ids []string
	for _, v := range reports {
		id, err := st.addOne(v)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}

	tests := []struct {
		name    string
		id      string
		wantID  string
		wantErr bool
	}{
		{
			name:    "OK ID 2",
			id:      ids[1],
			wantID:  ids[1],
			wantErr: false,
		},
		{
			name:    "Error Empty ID",
			id:      "",
			wantID:  "000000000000000000000000",
			wantErr: true,
		},
		{
			name:    "Error Incorrect ID",
			id:      "asdf",
			wantID:  "000000000000000000000000",
			wantErr: true,
		},
		{
			name:    "Error ID not found",
			id:      "66d03711a91e800d8c53b3d8",
			wantID:  "000000000000000000000000",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.ReportByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReportByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID.Hex() != tt.wantID {
				t.Errorf("Storage.ReportByID() ID = %v, want %v", got.ID.Hex(), tt.wantID)
			}
		})
	}
}
