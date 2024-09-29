package tools

import (
	"encoding/json"
	"errors"
	"strings"
	"wechat-back/internals/models"
)

var ErrorIncorrectLength = errors.New("payload not formatted correctly")

// ReadBinaryWebsocketMessage reads the binary message and destructures
func ReadBinaryWebsocketMessage(data []byte, target any) ([][]byte, error) {

	var result [][]byte

	dataDestructured := strings.Split(string(data), models.WEBSOCKET_BINARY_SEPARATOR)
	if len(dataDestructured) != 2 {
		return result, ErrorIncorrectLength
	}

	err := json.Unmarshal([]byte(dataDestructured[0]), target)
	if err != nil {
		return result, err
	}

	filesDestructured := strings.Split(dataDestructured[1], models.WEBSOCKET_FILE_SEPARATOR)

	for _, file := range filesDestructured {
		result = append(result, []byte(file))
	}

	return result, nil
}
