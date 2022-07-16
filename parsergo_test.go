/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parsergo

import (
	"testing"
)

func Test_scanner_Scan(t *testing.T) {
	type fields struct {
		Paths []string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "test scanner#Scan()",
			fields: fields{
				Paths: []string{"./tests/structx"},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scr := &scanner{
				Paths: []string{"./tests/structx"},
			}
			got := scr.Scan()
			if len(got) != tt.want {
				t.Errorf("scan the path:%s error: got %v, want %v", scr.Paths, len(got), tt.want)
			}

			// TODO assert struct
			// TODO assert methods
			// TODO assert annotations
		})
	}
}
