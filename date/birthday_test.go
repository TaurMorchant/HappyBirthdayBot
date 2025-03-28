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
				Time:      tt.fields.Time,
			}
			if got := b.ToString(tt.args.maxMonthLength); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
