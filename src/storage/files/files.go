package files

import (
	"encoding/gob"
	"errors"
	"example/main/src/lib/e"
	"example/main/src/storage"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)
type Storage struct {
	basePath string
}
const defaultPerm = 0774
var  ErrNoSavedPages = errors.New("no saved files")

func NewStorage(basePath string) Storage {
	return Storage{basePath: basePath}
}
func (s Storage) Save(page *storage.Page ) (err error) {
	defer func() {err = e.WrapIfErr("can't save page", err)}()
	fPath := filepath.Join(s.basePath, page.UserName)
	if err:= os.MkdirAll(fPath,  defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}
	fPath = filepath.Join(fPath, fName)
	file, err :=os.Create(fPath)
	if err != nil {
		return err
	}
	defer func(){_=file.Close()}()
	if err := gob.NewEncoder(file).Encode(page);err != nil {
		return err
	}
	return nil


}
func (s Storage) PickRandom(userName string) ( page *storage.Page, err error) {
	defer func() {err = e.WrapIfErr("can't pick random page", err)}()
	path := filepath.Join(s.basePath, userName)
	files, err:= os.ReadDir(path)
	if err != nil {
		return nil,err
	}
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}
	// rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n:=rand.Intn(len(files))
	file := files[n]
	return s.decoderPage(filepath.Join(path,file.Name()))
 	
	

}
func (s Storage) Remove(p *storage.Page ) error {
	fileName, err :=fileName(p)
	if err != nil {
		return  e.Wrap("can't remove file", err)
		
	}
	path:=filepath.Join(s.basePath,p.UserName,fileName)
	if err := os.Remove(path); err != nil {
		return  e.Wrap(fmt.Sprintf("can't remove file: %s",path), err)

	}
	return nil

}

func (s Storage) IsExistst (p *storage.Page) (bool, error) {
	fileName, err :=fileName(p)
	if err != nil {
		return  false, e.Wrap("can't find file", err)
		
	}
	path:=filepath.Join(s.basePath,p.UserName,fileName)
	switch _, err := os.Stat(path);  {
	case errors.Is(err,os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("can't find file: %s",path), err)
	
		
	}
	return true, nil

}

func (s Storage) decoderPage(filePath string) (*storage.Page, error) {
	f,err:=os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't open file",err)
	}
	defer func(){_=f.Close()}()
	var p storage.Page
	if err:=gob.NewDecoder(f).Decode(&p); err!=nil {
		return nil, e.Wrap("can't open file",err)
	}
	return &p, nil


}
func fileName(p *storage.Page)(string, error)  {
	return p.Hash()
}