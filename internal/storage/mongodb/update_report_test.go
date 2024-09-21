package mongodb

import (
	"Report-Storage/internal/storage"
	"context"
	"os"
	"reflect"
	"testing"
)

func TestStorage_UpdateReport(t *testing.T) {

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

	type args struct {
		number int
		desc   string
		media  []string
		status storage.Status
	}
	tests := []struct {
		name      string
		reportNum int
		args      args
		wantErr   bool
	}{
		{
			name:      "OK Number 1",
			reportNum: 0,
			args:      args{number: 1, desc: "Новое описание заявки 1", media: []string{"https://bing.com", "https://ya.ru"}, status: 3},
			wantErr:   false,
		},
		{
			name:      "Error Incorrect number",
			reportNum: 1,
			args:      args{number: -1, desc: "Новое описание заявки 2", media: []string{}, status: 1},
			wantErr:   true,
		},
		{
			name:      "Error Not found",
			reportNum: 2,
			args:      args{number: 4, desc: "Новое описание заявки 3", media: []string{}, status: 1},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Вставляем в коллекцию тестовую заявку и получаем ее.
			id, err := st.addOne(reports[tt.reportNum])
			if err != nil {
				t.Fatal(err)
			}
			want, err := st.getOne(id)
			if err != nil {
				t.Fatal(err)
			}

			new := want
			new.Number = int64(tt.args.number)
			new.Description = tt.args.desc
			new.Media = tt.args.media
			new.Status = tt.args.status

			// Выполняем изменение.
			got, err := st.UpdateReport(context.Background(), new)
			if err != nil {
				if tt.wantErr {
					t.Skip()
				}
				t.Errorf("Storage.UpdateReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Storage.UpdateReport() = %v, want %v", got, want)
			}
		})
	}
}
