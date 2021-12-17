package report

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		r    Status
		want string
	}{
		{
			name: "StatusPending",
			r:    StatusPending,
			want: "Pending",
		},
		{
			name: "StatusComplete",
			r:    StatusComplete,
			want: "Complete",
		},
		{
			name: "StatusError",
			r:    StatusError,
			want: "Error",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			require.Equal(tt.want, tt.r.String())
		})
	}
}

func TestStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       Status
		want    []byte
		wantErr string
	}{
		{
			name: "StatusPending",
			r:    StatusPending,
			want: []byte(`"Pending"`),
		},
		{
			name: "StatusComplete",
			r:    StatusComplete,
			want: []byte(`"Complete"`),
		},
		{
			name: "StatusError",
			r:    StatusError,
			want: []byte(`"Error"`),
		},
		{
			name:    "StatusFake",
			r:       Status(-1),
			wantErr: "-1 is an invalid enum for Status",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			got, err := tt.r.MarshalJSON()

			if tt.wantErr == "" {
				require.NoError(err)
			} else {
				require.EqualError(err, tt.wantErr)
			}
			require.Equal(tt.want, got)
		})
	}
}

func TestStatus_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		r       *Status
		args    args
		wantErr string
	}{
		{
			name: "StatusPending",
			r:    statusToStatusPointer(StatusPending),
			args: args{
				data: []byte(`"Pending"`),
			},
		},
		{
			name: "StatusComplete",
			r:    statusToStatusPointer(StatusComplete),
			args: args{
				data: []byte(`"Complete"`),
			},
		},
		{
			name: "StatusError",
			r:    statusToStatusPointer(StatusError),
			args: args{
				data: []byte(`"Error"`),
			},
		},
		{
			name: "StatusFake",
			r:    statusToStatusPointer(Status(-1)),
			args: args{
				data: []byte(`"fake"`),
			},
			wantErr: `"fake" is an invalid enum for Status`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			err := tt.r.UnmarshalJSON(tt.args.data)

			if tt.wantErr == "" {
				require.NoError(err)
			} else {
				require.EqualError(err, tt.wantErr)
			}
		})
	}
}

func statusToStatusPointer(r Status) *Status {
	return &r
}
