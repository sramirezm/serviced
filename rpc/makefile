# Copyright 2014 The Serviced Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: build test clean

default: build

build:
	cd agent && go build
	cd master && go build

test: build
	cd agent && go test $(GOTEST_FLAGS)
	cd master && go test $(GOTEST_FLAGS)

clean:
	cd agent && go clean
	cd master && go clean
