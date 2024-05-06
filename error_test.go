package httphelper

import (
	"reflect"
	"testing"

	"github.com/solsw/generichelper"
)

func Test_objMsg_NoType(t *testing.T) {
	type args struct {
		herr    *Error[generichelper.NoType]
		bb      []byte
		options ErrorOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *Error[generichelper.NoType]
		wantErr bool
	}{
		{name: "withObject",
			args: args{
				herr:    &Error[generichelper.NoType]{},
				bb:      []byte("qwerty"),
				options: ErrorOptions{withObject: true},
			},
			want: &Error[generichelper.NoType]{},
		},
		{name: "withMessage",
			args: args{
				herr:    &Error[generichelper.NoType]{},
				bb:      []byte("qwerty"),
				options: ErrorOptions{withMessage: true},
			},
			want: &Error[generichelper.NoType]{Message: "qwerty"},
		},
		{name: "withObjectwithMessage",
			args: args{
				herr:    &Error[generichelper.NoType]{},
				bb:      []byte("qwerty"),
				options: ErrorOptions{withObject: true, withMessage: true},
			},
			want: &Error[generichelper.NoType]{Message: "qwerty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objMsg(tt.args.herr, tt.args.bb, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("objMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

type E struct {
	I int
	S string
}

func Test_objMsg_E(t *testing.T) {
	type args struct {
		herr    *Error[E]
		bb      []byte
		options ErrorOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *Error[E]
		wantErr bool
	}{
		// {name: "object",
		// 	args: args{
		// 		herr:    &Error[E]{},
		// 		bb:      []byte(`{"I":1,"S":"one"}`),
		// 		options: ErrorOptions{withObject: true, withMessage: true},
		// 	},
		// 	want: &Error[E]{Object: E{I: 1, S: "one"}},
		// },
		// {name: "not JSON",
		// 	args: args{
		// 		herr:    &Error[E]{},
		// 		bb:      []byte("qwerty"),
		// 		options: ErrorOptions{withObject: true, withMessage: true},
		// 	},
		// 	want: &Error[E]{Message: "qwerty"},
		// },
		// {name: "cannot unmarshal string",
		// 	args: args{
		// 		herr:    &Error[E]{},
		// 		bb:      []byte(`"qwerty"`),
		// 		options: ErrorOptions{withObject: true, withMessage: true},
		// 	},
		// 	want: &Error[E]{Message: `"qwerty"`},
		// },
		{name: "cannot unmarshal",
			args: args{
				herr:    &Error[E]{},
				bb:      []byte(`{"X":1,"Y":"one"}`),
				options: ErrorOptions{withObject: true, withMessage: true},
			},
			want: &Error[E]{Message: `{"X":1,"Y":"one"}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objMsg(tt.args.herr, tt.args.bb, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("objMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
