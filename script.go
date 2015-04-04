// Copyright 2014 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package twodee

import (
	"github.com/robertkrimen/otto"
)

type Scripting struct {
	vm *otto.Otto
}

func NewScripting() (s *Scripting, err error) {
	s = &Scripting{
		vm: otto.New(),
	}
	return
}

func (s *Scripting) LoadScript(path string) (err error) {
	var (
		script *otto.Script
	)
	if script, err = s.vm.Compile(path, nil); err != nil {
		return
	}
	if _, err = s.vm.Run(script); err != nil {
		return
	}
	return
}
