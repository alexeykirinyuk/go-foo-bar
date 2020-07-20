package sigma

import (
	"fmt"
	"github.com/alexeykirinyuk/tech-task-go/data"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type sigmaStorage struct {
	db *gorm.DB
}

func newStorage(dbProvider data.IDatabaseProvider) sigmaStorage {
	db, err := dbProvider.CreateConnection()
	if err != nil {
		panic(err)
	}

	return sigmaStorage{db: db}
}

func (s sigmaStorage) add(Sigma sigma) error {
	if err := s.db.Create(Sigma).Error; err != nil {
		return fmt.Errorf("error when trying crete new sigma: %s", err)
	}

	return nil
}

func (f sigmaStorage) getAll() (items []sigma, err error) {
	err = f.db.Find(&items).Error
	if err != nil {
		err = fmt.Errorf("error when trying get all sigma: %s", err)
		return
	}

	return
}

func (s sigmaStorage) getById(id uuid.UUID) (Sigma sigma, err error) {
	err = s.db.First(&Sigma, id).Error
	if err != nil {
		err = fmt.Errorf("error when trying get sigma by Id: %s", err)
		return
	}

	return
}

func (s sigmaStorage) update(Sigma sigma) error {
	if err := s.db.Update(Sigma).Error; err != nil {
		return fmt.Errorf("error when trying update sigma with ID %s: %s", Sigma.Id, err)
	}

	return nil
}

func (s sigmaStorage) delete(id uuid.UUID) error {
	if err := s.db.Where("id = ?", id).Delete(&sigma{}); err != nil {
		return fmt.Errorf("error when trying delete sigma with ID %s: %s", id, err)
	}

	return nil
}
