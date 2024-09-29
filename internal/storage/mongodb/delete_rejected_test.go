package mongodb

import (
	"context"
	"os"
	"testing"
)

func TestStorage_DeleteRejected(t *testing.T) {

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

	// Вставляем тестовую заявку и устанавливаем статус Rejected.
	_, err = st.addOne(reports[0])
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.UpdateStatus(context.Background(), 1, 5)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "OK",
			want:    1,
			wantErr: false,
		},
		{
			name:    "Error Not found",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.DeleteRejected(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteRejected() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.DeleteRejected() = %d, want %d", got, tt.want)
			}
		})
	}
}
