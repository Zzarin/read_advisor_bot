package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"read_advisor_bot/internal/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func NewStorage(basePath string) Storage {
	return Storage{basePath: basePath}
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}

func (s Storage) Save(page *storage.Page) (err error) {
	filePath := filepath.Join(s.basePath, page.UserName) // путь до директории куда сохранить файл
	err = os.Mkdir(filePath, defaultPerm)                //создаем директорию с правами доступа
	if err != nil {
		return fmt.Errorf("can't create new directory to save %w", err)
	}

	fileName, err := fileName(page) //формируем имя файла
	if err != nil {
		return fmt.Errorf("can't get file name %w", err)
	}

	filePath = filepath.Join(filePath, fileName) //дописываем имя файла к пути

	file, err := os.Create(filePath) //создаем файл
	if err != nil {
		return fmt.Errorf("can't create file %w", err)
	}
	defer func(f func() error) {
		errClose := f()
		if err == nil {
			err = errClose
		} else if errClose != nil {
			log.Printf("can't close file %s", err.Error())
		}
	}(file.Close)

	err = gob.NewEncoder(file).Encode(page) //записываем страницу в файл в нужном формате
	if err != nil {
		return fmt.Errorf("can't write page to file %w", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	filePath := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't read directory %w", err)
	}

	if len(files) == 0 {
		return nil, errors.New("no saved pages")
	}

	rand.Seed(time.Now().Unix()) //чтобы генератор псевдослучайных чисел всегд возвращал разные числа
	n := rand.Intn(len(files))   // рандомное число от 0 до кол-ва файлов

	file := files[n]                                          //выбрать рандомный файл
	return s.decodePage(filepath.Join(filePath, file.Name())) //декодируем рандомный файл
}

func (s Storage) decodePage(filePath string) (page *storage.Page, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open file %w", err)
	}
	defer func(f func() error) {
		errClose := f()
		if err == nil {
			err = errClose
		} else if errClose != nil {
			log.Printf("can't close file %s", err.Error())
		}
	}(file.Close)

	err = gob.NewDecoder(file).Decode(&page)
	if err != nil {
		return nil, fmt.Errorf("can't open page from file %w", err)
	}
	return page, nil
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can't get file name %w", err)
	}
	pathToFile := filepath.Join(s.basePath, page.UserName, fileName)
	err = os.Remove(pathToFile)
	if err != nil {
		fileName := fmt.Sprintf("can't remove file %s", pathToFile)
		return fmt.Errorf(fileName, err)
	}
	return nil
}

func (s Storage) IsExist(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, fmt.Errorf("can't check if file exists %w", err)
	}
	pathToFile := filepath.Join(s.basePath, page.UserName, fileName)
	switch _, err = os.Stat(pathToFile); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		fileName := fmt.Sprintf("can't check if file %s exists", pathToFile)
		return false, fmt.Errorf(fileName, err)
	}

	return true, nil

}
