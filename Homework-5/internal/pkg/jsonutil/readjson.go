package jsonutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadJSON(v interface{}, path string) (interface{}, error) {

	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		err = fmt.Errorf("Ошибка при открытии файла на чтение: %w", err)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bytesReader, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("Ошибка при вызове ReadAll(): %w", err)
		return nil, err
	}

	err = json.Unmarshal(bytesReader, &v)
	if err != nil {
		err = fmt.Errorf("Ошибка при десериализации json: %w", err)
		return nil, err
	}

	return v, nil
}
