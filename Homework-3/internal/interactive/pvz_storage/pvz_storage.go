package pvz_storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/interactive/model"
	"io"
	"os"
	"sync"
)

type pvzStorage struct {
	filePath string
	Wg       sync.WaitGroup
	mx       sync.RWMutex
}

func (s *pvzStorage) AddWG() {
	s.Wg.Add(1)
}

func NewPvzStorage(filePath string) pvzStorage {
	return pvzStorage{filePath: filePath}
}

func (s *pvzStorage) Read(pvzID int) (model.Pvz, error) {
	defer s.Wg.Done()

	all, err := s.readAll()
	if err != nil {
		err = fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
		return model.Pvz{}, err
	}

	for _, pvz := range all {
		if pvz.PvzID == pvzID {
			return pvz, nil
		}
	}

	return model.Pvz{}, nil
}

func (s *pvzStorage) Write(input model.Pvz) error {
	defer s.Wg.Done()

	s.AddWG()
	checkPvz, err := s.Read(input.PvzID)
	if err != nil {
		err = fmt.Errorf("Ошибка при чтении файла для поиска дубликатов: %w", err)
		return err
	}
	if checkPvz.PvzID > 0 {
		err = errors.New("Попытка записать дубликат записи о ПВЗ.")
		return err
	}

	all, err := s.readAll()
	if err != nil {
		err = fmt.Errorf("Не удалось прочитать данные из базы данных: %w", err)
		return err
	}

	all = append(all, input)

	err = s.overwriteBytes(all)
	if err != nil {
		err = fmt.Errorf("Не удалось записать данные в базу данных: %w", err)
		return err
	}

	return nil
}

func (s *pvzStorage) readAll() ([]model.Pvz, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		err = fmt.Errorf("Ошибка при открытии файла на чтение: %w", err)
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		err = fmt.Errorf("Ошибка Seek() при перемещении на начальную позицию в файле: %w", err)
		return nil, err
	}
	reader := bufio.NewReader(file)
	bytesReader, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("Ошибка при вызове ReadAll(): %w", err)
		return nil, err
	}
	pvzs := make([]model.Pvz, 0)
	if len(bytesReader) == 0 {
		return pvzs, nil
	}

	err = json.Unmarshal(bytesReader, &pvzs)
	if err != nil {
		err = fmt.Errorf("Ошибка при распаковке json (json.Unmarshal()): %w", err)
		return nil, err
	}

	return pvzs, nil
}

func (s *pvzStorage) overwriteBytes(pvzs []model.Pvz) error {
	bytesReader, err := json.Marshal(pvzs)
	if err != nil {
		return fmt.Errorf("Ошибка при приведении данных в json (json.Marshal()): %w", err)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	err = os.WriteFile(s.filePath, bytesReader, 0777)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове перезаписи файла: %w", err)
	}

	return nil
}
