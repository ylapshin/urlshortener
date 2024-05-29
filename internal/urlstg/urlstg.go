package urlstg

import (
	"github.com/sethvargo/go-password/password"
	"github.com/ylapshin/urlshortener/internal/minidb"
)

const KeyLen = 6

type UrlStg struct {
	data *minidb.MiniDb
}

func New() *UrlStg {
	return &UrlStg{
		data: minidb.New(),
	}
}

func (stg *UrlStg) Reg(url string) (id string, err error) {
	// если url уже существует - отдаем его идентификатор
	if id, err = stg.data.Resolve(url); err == nil {
		return id, nil
	}

	// если url'а нет, то генерируем для него уникальный id и записываем в минибд
	for {
		id, _ = password.Generate(KeyLen, 0, 0, false, true)
		if _, err = stg.data.Get(id); err != nil {
			break
		}
	}
	err = stg.data.Set(id, url)

	return id, err
}

func (stg *UrlStg) Get(id string) (url string, err error) {
	return stg.data.Get(id)
}
