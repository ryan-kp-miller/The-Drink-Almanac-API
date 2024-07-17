package model

import (
	"testing"
)

func TestDefaultEnv(t *testing.T) {
	type args struct {
		envVarName   string
		defaultValue string
	}
	tests := []struct {
		name         string
		args         args
		want         string
		shouldSetEnv bool
	}{
		{
			name: "Env variable found",
			args: args{
				envVarName:   "AWS_ENDPOINT",
				defaultValue: "default",
			},
			want:         "test_aws_endpoint",
			shouldSetEnv: true,
		},
		{
			name: "Env variable not found",
			args: args{
				envVarName:   "AWS_ENDPOINT",
				defaultValue: "default",
			},
			want:         "default",
			shouldSetEnv: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSetEnv {
				t.Setenv(tt.args.envVarName, tt.want)
			}
			if got := DefaultEnv(tt.args.envVarName, tt.args.defaultValue); got != tt.want {
				t.Errorf("DefaultEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
