package main

import (
	"testing"
	"time"
)

func Test_parseSchedule(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Weekday
		wantErr bool
	}{
		{name: "test number", args: args{val: "0"}, want: 0, wantErr: false},
		{name: "test number", args: args{val: "14"}, want: -1, wantErr: true},
		{name: "test short name", args: args{val: "Mon"}, want: 1, wantErr: false},
		{name: "test short name lower case", args: args{val: "sat"}, want: 6, wantErr: false},
		{name: "test short name upper case", args: args{val: "FRI"}, want: 5, wantErr: false},
		{name: "test full name", args: args{val: "Monday"}, want: 1, wantErr: false},
		{name: "test full name lower case", args: args{val: "saturday"}, want: 6, wantErr: false},
		{name: "test full name upper case", args: args{val: "FRIDAY"}, want: 5, wantErr: false},
		{name: "test unknown val", args: args{val: "SOMEDAY"}, want: -1, wantErr: true},
		{name: "test number left-padded (one zero)", args: args{val: "01"}, want: 1, wantErr: false},
		{name: "test number left-padded (three zeros)", args: args{val: "0001"}, want: 1, wantErr: false},
		{name: "test invalid number left-padded (three zeros)", args: args{val: "0014"}, want: -1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSchedule(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseSchedule() got = %v, want %v", got, tt.want)
			}
		})
	}
}
