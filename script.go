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
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
)

type Scripting struct {
	vm        *otto.Otto
	listeners map[string][]otto.Value
}

func NewScripting() (s *Scripting, err error) {
	s = &Scripting{
		vm:        otto.New(),
		listeners: map[string][]otto.Value{},
	}
	return
}

func (s *Scripting) setAPI() {
	var err = errors.New("Usage: addEventListener(string, func)")
	s.vm.Set("addEventListener", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			panic(err)
		}
		if !call.Argument(0).IsString() {
			panic(err)
		}
		if !call.Argument(1).IsFunction() {
			panic(err)
		}
		var (
			callbacks []otto.Value
			present   bool
			eventName = call.Argument(0).String()
		)
		if callbacks, present = s.listeners[eventName]; !present {
			s.listeners[eventName] = callbacks
		}
		s.listeners[eventName] = append(s.listeners[eventName], call.Argument(1))
		return otto.Value{}
	})
}

func (s *Scripting) LoadScript(path string) (err error) {
	var (
		script *otto.Script
	)
	if script, err = s.vm.Compile(path, nil); err != nil {
		return
	}
	s.setAPI()
	if _, err = s.vm.Run(script); err != nil {
		return
	}
	return
}

func (s *Scripting) convertArguments(args []interface{}) (converted []interface{}, err error) {
	var (
		i   int
		arg interface{}
	)
	converted = make([]interface{}, len(args))
	for i, arg = range args {
		if converted[i], err = s.vm.ToValue(arg); err != nil {
			return
		}
	}
	return
}

func (s *Scripting) TriggerEvent(eventName string, rawArguments ...interface{}) (err error) {
	var (
		present   bool
		callbacks []otto.Value
		callback  otto.Value
		response  otto.Value
		arguments []interface{}
	)
	if callbacks, present = s.listeners[eventName]; !present {
		return // Not an error to have no listeners.
	}
	if arguments, err = s.convertArguments(rawArguments); err != nil {
		return
	}
	for _, callback = range callbacks {
		if response, err = callback.Call(callback, arguments...); err != nil {
			return
		}
		if response.IsString() {
			fmt.Printf("Response from callback: %s\n", response.String())
		}
	}
	return
}
