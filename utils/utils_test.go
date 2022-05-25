// alibaba-inc.com Inc.
// Copyright (c) 2004-2022 All Rights Reserved.
//
// @Author : huaiyou.cyz
// @Time : 2022/5/25 4:08 PM
// @File : utils_test.go
//

package utils

import "testing"

func TestSplitServer(t *testing.T) {
	type args struct {
		server string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 uint16
	}{
		{
			args: args{
				server: "[1408:4003:10bb:6a01:83b9:6360:c66d:ed3e]:6443",
			},
			want:  "1408:4003:10bb:6a01:83b9:6360:c66d:ed3e",
			want1: 6443,
		},
		{
			args: args{
				server: "1.1.1.1:6443",
			},
			want:  "1.1.1.1",
			want1: 6443,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitServer(tt.args.server)
			if got != tt.want {
				t.Errorf("SplitServer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitServer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
