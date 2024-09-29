package mongodb

import (
	"context"
	"os"
	"testing"
)

func TestStorage_CounterInc(t *testing.T) {

	// Создаем пул подключений.
	dbName = testDatabase
	colReport = testCollection
	colCounter = testCounter
	opts := setOpts(path, "admin", os.Getenv("MONGO_DB_PASSWD"))
	st, err := new(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	// Очищаем тестовую коллекцию.
	err = st.trun(colCounter)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    int32
		wantErr bool
	}{
		{
			name:    "OK First",
			want:    1,
			wantErr: false,
		},
		{
			name:    "OK Second",
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.CounterInc(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CounterInc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.CounterInc() = %v, want %v", got, tt.want)
			}
		})
	}
}
