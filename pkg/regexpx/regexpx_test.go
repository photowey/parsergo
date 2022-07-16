package regexpx

import (
	"testing"
)

func TestRegexpExtract(t *testing.T) {
	type args struct {
		regex string
		src   string
		temp  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test extract annotation name",
			args: args{
				regex: `^@(?P<annotation>[\S]+)\(.*\)`,
				src:   "@Service(\"helloService\")",
				temp:  "$annotation",
			},
			want: "Service",
		},
		{
			name: "Test extract annotation string value",
			args: args{
				regex: `^@.*\((?P<value>[\S]+)\)`,
				src:   "@Service(\"helloService\")",
				temp:  "$value",
			},
			want: "\"helloService\"",
		},
		{
			name: "Test extract annotation json value",
			args: args{
				regex: `^@.*\((?P<value>[\S]+)\)`,
				src:   "@ComponentScan({\"path\":\"github.com/photowey/parsergo/tests\",\"excludes\":[\"github.com/photowey/parsergo/tests/structx\"]})",
				temp:  "$value",
			},
			want: "{\"path\":\"github.com/photowey/parsergo/tests\",\"excludes\":[\"github.com/photowey/parsergo/tests/structx\"]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RegexpExtract(tt.args.regex, tt.args.src, tt.args.temp)
			if got != tt.want {
				t.Errorf("RegexpExtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
