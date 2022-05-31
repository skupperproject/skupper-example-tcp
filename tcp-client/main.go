//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package main

import (
        "bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "HOST", "PORT")
		return
	}

	host := os.Args[1]
	port := os.Args[2]

	conn, err := net.Dial("tcp", host + ":" + port)

        if err != nil {
                fmt.Println(err)
                return
        }

        send_text, _ := bufio.NewReader(os.Stdin).ReadString('\n')

        fmt.Fprintf(conn, send_text + "\n")

        receive_text, _ := bufio.NewReader(conn).ReadString('\n')

        fmt.Print(receive_text)
}
