package mongodb

import (
	"context"
	"os"
	"testing"
)

func TestStorage_ReportByNum(t *testing.T) {

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
		num     int
		wantNum int64
		wantErr bool
	}{
		{
			name:    "OK Number2",
			num:     2,
			wantNum: 2,
			wantErr: false,
		},
		{
			name:    "Error Incorrect num",
			num:     -1,
			wantNum: 0,
			wantErr: true,
		},
		{
			name:    "Error Report not found",
			num:     5,
			wantNum: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.ReportByNum(context.Background(), tt.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReportByNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Number != tt.wantNum {
				t.Errorf("Storage.ReportByNum() = %d, want %d", got.Number, tt.wantNum)
			}
		})
	}
}
