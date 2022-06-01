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
	"fmt"
	"io"
	"net"
	"os"
	"strings"
        "strconv"
)

// Handle a connection. Read messages from it and reply with the
// same message uppercased. Terminate the connection when we reach
// EOF, or an error.
func conn_handler(conn_number int, host string, conn net.Conn) {
	buffer := make([]byte, 512)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err != io.EOF {
                                fmt.Println("Connection read error:", err.Error())
			}
			break
		}

		// If we have received a message, uppercase it and send it
		// back as the reply.
		if n > 0 {
			message := string(buffer[0:n])
                        message = strings.TrimSuffix(message, "\n")

                        fmt.Println("Received from connection " + strconv.Itoa(conn_number) + ": " + message)

			reply := host + ": " + strings.ToUpper(message)

                        conn.Write([]byte(reply + "\n"))
		}
	}

	conn.Close()

        fmt.Println("Closed connection", conn_number)
}

func main() {
        host := os.Getenv("HOSTNAME")
	port := "9090"
	conn_count := 0

	listener, err := net.Listen("tcp", ":" + port)

        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Println("Listening on port", port);

	// Handle each new connection in its own goroutine. Tell each
	// one what number it is so the user can see which handler is
	// printing out each message.
	for {
		conn, _ := listener.Accept()
		conn_count++

		fmt.Println("Accepted connection", conn_count)

		go conn_handler(conn_count, host, conn)
	}
}
