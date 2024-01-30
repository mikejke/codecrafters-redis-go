package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleClient(conn)
	}
}

const (
	Ping = "PING"
	Echo = "ECHO"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		commands, err := parseCommands(reader)

		if err == io.EOF {
			return
		}

		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			writer.WriteString("Error reading!\n")
			writer.Flush()
			return
		}

		// Handle commands
		var response *RESPType
		switch strings.ToUpper(commands[0]) {
		case Ping:
			response, err = handlePing(commands)
		case Echo:
			response, err = handleEcho(commands)
		default:
			fmt.Println("Unknown command: ", commands[0])
			writer.WriteString("Unknown command\n")
			writer.Flush()
			return
		}

		if err != nil || response == nil {
			fmt.Println("Error handling command: ", err.Error())
			writer.WriteString("Error handling command!\n")
			writer.Flush()
			return
		}

		writer.Write((*response).Prepare())
	}
}

func parseCommands(reader *bufio.Reader) ([]string, error) {
	rawMessage, err := parseRESPType(reader)
	if err != nil {
		return []string{}, err
	}

	if commandsArray, ok := (*rawMessage).(*RESPArray); ok {
		fmt.Println("Received: ", commandsArray)
		commands := make([]string, len(commandsArray.Value))

		for i, v := range commandsArray.Value {
			if bulkString, ok := v.(*RESPBulkString); ok {
				commands[i] = bulkString.Value
			} else if simpleString, ok := v.(*RESPSimpleString); ok {
				commands[i] = simpleString.Value
			} else {
				fmt.Println("Unknown type: ", v)
			}
		}
		return commands, nil
	}

	return []string{}, fmt.Errorf("expected RESPArray, got %v", rawMessage)
}

const (
	CLRF         = "\r\n"
	SimpleString = '+'
	BulkString   = '$'
	Array        = '*'
)

type RESPType interface {
	Prepare() []byte
}

type RESPSimpleString struct {
	Value string
}

func (t *RESPSimpleString) Prepare() []byte {
	bytes := []byte(t.Value)

	result := []byte(string(SimpleString))
	result = append(result, bytes...)
	result = append(result, []byte(CLRF)...)

	return result
}

type RESPBulkString struct {
	Value string
}

func (t *RESPBulkString) Prepare() []byte {
	bytes := []byte(t.Value)
	l := len(bytes)

	result := []byte(string(BulkString) + strconv.Itoa(l) + CLRF)
	result = append(result, bytes...)
	result = append(result, []byte(CLRF)...)

	return result
}

type RESPArray struct {
	Value []RESPType
}

func (t *RESPArray) Prepare() []byte {
	length := len(t.Value)
	result := []byte(string(Array) + strconv.Itoa(length) + CLRF)

	for _, value := range t.Value {
		result = append(result, value.Prepare()...)
	}

	return result
}

func parseRESPType(reader *bufio.Reader) (*RESPType, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	if len(line) == 0 {
		return nil, io.EOF
	}

	fmt.Println("Received: ", string(line))

	firstChar := line[0]

	switch firstChar {
	case SimpleString:
		simpleString := &RESPSimpleString{string(line[1:])}

		fmt.Println("Unmarshalling simple string of length: ", len(simpleString.Value))

		value := RESPType(simpleString)

		return &value, nil
	case BulkString:
		length, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}

		fmt.Println("Unmarshalling bulk string of length: ", length)

		bulkString, err := parseBulkString(reader, length)
		if err != nil {
			return nil, err
		}

		value := RESPType(bulkString)

		return &value, nil
	case Array:
		length, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}

		fmt.Println("Unmarshalling array of length: ", length)

		arr, err := parseArray(reader, length)
		if err != nil {
			return nil, err
		}

		value := RESPType(arr)

		return &value, nil
	default:
		return nil, fmt.Errorf("unknown type: %v", firstChar)

	}
}

func parseBulkString(r *bufio.Reader, length int) (*RESPBulkString, error) {
	value := &RESPBulkString{}
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	if length < 1 {
		// Empty string
		return value, nil
	}

	stringData := line[:length]
	value.Value = string(stringData)

	return value, nil
}

func parseArray(r *bufio.Reader, length int) (*RESPArray, error) {
	value := &RESPArray{}
	for i := 0; i < length; i++ {
		respType, err := parseRESPType(r)
		if err != nil {
			return nil, err
		}
		value.Value = append(value.Value, *respType)
	}
	return value, nil
}

func handlePing(commands []string) (*RESPType, error) {
	response := &RESPSimpleString{"PONG"}
	wrappedResponse := RESPType(response)
	return &wrappedResponse, nil
}

func handleEcho(commands []string) (*RESPType, error) {
	response := &RESPBulkString{commands[1]}
	wrappedResponse := RESPType(response)
	return &wrappedResponse, nil
}
