package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"reflect"
	"testing"
)

func TestStorage_Statistic(t *testing.T) {

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
		want    storage.Statistic
		wantErr bool
	}{
		{
			name:    "OK",
			want:    storage.Statistic{Total: 3, Unverified: 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.Statistic(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Statistic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Statistic() = %v, want %v", got, tt.want)
			}
		})
	}
}
