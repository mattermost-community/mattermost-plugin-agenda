package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/undefinedlabs/go-mpatch"
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

func Test_nextWeekdayDateInWeekSkippingDay(t *testing.T) {
	var patch, err = mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2021, 11, 15, 00, 00, 00, 0, time.UTC)
	})
	if patch == nil || err != nil {
		t.Errorf("error creating patch")
	}
	type args struct {
		meetingDays []time.Weekday
		nextWeek    bool
		dayToSkip   time.Weekday
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{

		{name: "test skip tuesday in week, today", args: args{[]time.Weekday{1, 2, 3, 4, 5, 6, 7}, false, time.Weekday(2)}, want: time.Now().AddDate(0, 0, 0), wantErr: false},
		{name: "test skip tuesday in few days", args: args{[]time.Weekday{2, 3, 4, 5, 6, 7}, false, time.Weekday(2)}, want: time.Now().AddDate(0, 0, 2), wantErr: false},
		{name: "test skip monday with nextWeek true", args: args{[]time.Weekday{1, 2, 3, 4}, true, time.Weekday(1)}, want: time.Now().AddDate(0, 0, 8), wantErr: false},
		{name: "test only meeting day is skipped", args: args{[]time.Weekday{3}, false, time.Weekday(3)}, want: time.Now().AddDate(0, 0, 2), wantErr: false},
		{name: "test only meeting day is skipped with nextWeek true", args: args{[]time.Weekday{3}, true, time.Weekday(3)}, want: time.Now().AddDate(0, 0, 9), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextWeekdayDateInWeekSkippingDay(tt.args.meetingDays, tt.args.nextWeek, tt.args.dayToSkip)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextWeekdayDateInWeekSkippingDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("nextWeekdayDateInWeekSkippingDay() got = %v, want %v", got, tt.want)
			}
		})
	}
}
