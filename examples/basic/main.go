// Copyright 2012 Arne Roomann-Kurrik
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

package main

import (
	"github.com/kurrik/twodee"
	"fmt"
	"os"
)

func PrintError(err error) {
	fmt.Printf("[error]: %v\n", err)
}

func main() {
	var (
		system *twodee.System
		window *twodee.Window
		err error
		run bool = true
	)
	if system, err = twodee.Init(); err != nil {
		PrintError(err)
		os.Exit(1)
	}
	defer system.Terminate()
	window = &twodee.Window{Width: 640, Height: 480}
	system.Open(window)
	for run {
		system.Paint()
		run = system.Key(twodee.KeyEsc) == 0 && window.Opened()
	}
}
