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

package structx

import (
	"fmt"
)

type HelloService interface {
	SayHello(name string) string
	MultiLineFunc(name string, age *int) (string, error)
}

// HelloServiceImpl An implementation of the HelloService interface
// @Service("helloService")
// @ComponentScan({"path":"github.com/photowey/parsergo/tests","excludes":["github.com/photowey/parsergo/tests/structx"]})
type HelloServiceImpl struct{}

func (s HelloServiceImpl) SayHello(name string) string {
	return "Hello " + name
}

func (s *HelloServiceImpl) MultiLineFunc(
	name string,
	age *int,
) (string, error) {
	return fmt.Sprintf("Hello %s, age %d", name, age), nil
}
