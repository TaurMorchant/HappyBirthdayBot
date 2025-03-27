package handlers

import "testing"

func Test_getDaysWord(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{1}, want: "день"},
		{name: "2", args: args{2}, want: "дня"},
		{name: "5", args: args{5}, want: "дней"},
		{name: "10", args: args{10}, want: "дней"},
		{name: "11", args: args{11}, want: "дней"},
		{name: "12", args: args{12}, want: "дней"},
		{name: "15", args: args{15}, want: "дней"},
		{name: "20", args: args{20}, want: "дней"},
		{name: "31", args: args{31}, want: "день"},
		{name: "33", args: args{33}, want: "дня"},
		{name: "37", args: args{37}, want: "дней"},
		{name: "111", args: args{111}, want: "дней"},
		{name: "211", args: args{211}, want: "дней"},
		{name: "254", args: args{254}, want: "дня"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDaysWord(tt.args.n); got != tt.want {
				t.Errorf("getDaysWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
