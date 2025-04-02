package date

import (
	"testing"
	"time"
)

func TestBirthday_ToString(t *testing.T) {
	type fields struct {
		day       int
		monthName string
		Time      time.Time
	}
	type args struct {
		maxMonthLength int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "1", fields: fields{day: 1, monthName: "января"}, args: args{10}, want: "1  января    "},
		{name: "2", fields: fields{day: 13, monthName: "сентября"}, args: args{10}, want: "13 сентября  "},
		{name: "3", fields: fields{day: 25, monthName: "марта"}, args: args{5}, want: "25 марта"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Birthday{
				day:       tt.fields.day,
				monthName: tt.fields.monthName,
				//Time:      tt.fields.Time,
			}
			if got := b.ToStringWithFormatting(tt.args.maxMonthLength); got != tt.want {
				t.Errorf("ToStringWithFormatting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBirthday(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    Birthday
		wantErr bool
	}{
		{name: "1", args: args{"sdfdsfsfsdfsd"}, want: Birthday{}, wantErr: true},
		{name: "2", args: args{"11 января"}, want: Birthday{day: 11, monthName: "января"}, wantErr: false},
		{name: "3", args: args{"  11 января   "}, want: Birthday{day: 11, monthName: "января"}, wantErr: false},
		{name: "4", args: args{"  11 январь   "}, want: Birthday{}, wantErr: true},
		{name: "5", args: args{"  35 января   "}, want: Birthday{}, wantErr: true},
		{name: "6", args: args{"  15    февраля   "}, want: Birthday{day: 15, monthName: "февраля"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBirthday(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBirthday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.day != tt.want.day {
				t.Errorf("ParseBirthday() day = %v, want %v", got.day, tt.want.day)
			}
			if got.monthName != tt.want.monthName {
				t.Errorf("ParseBirthday() monthName = %v, want %v", got.day, tt.want.day)
			}
		})
	}
}
