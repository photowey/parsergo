/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sets

type Empty struct{}

type String map[string]Empty

func NewString(items ...string) String {
	ss := String{}
	ss.Insert(items...)

	return ss
}

func (str String) Insert(items ...string) String {
	for _, item := range items {
		str[item] = Empty{}
	}
	return str
}

func (str String) Has(item string) bool {
	_, contained := str[item]

	return contained
}
