package mongodb

import (
	"context"
	"os"
	"testing"
)

func TestStorage_DeleteByNum(t *testing.T) {

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

	// Вставляем тестовую заявку.
	_, err = st.addOne(reports[0])
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		num     int
		wantErr bool
	}{
		{
			name:    "OK",
			num:     1,
			wantErr: false,
		},
		{
			name:    "Error Incorrect number",
			num:     -1,
			wantErr: true,
		},
		{
			name:    "Error Not found",
			num:     5,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.DeleteByNum(context.Background(), tt.num); (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteByNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
