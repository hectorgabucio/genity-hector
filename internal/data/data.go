package data

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

var once sync.Once

var (
	instance *gorm.DB
)

type Data struct {
	Title     string `gorm:"primary_key"`
	UUID      string
	CreatedAt time.Time
}

type DataRepository interface {
	CloseConn()
	Get(dataWhere *Data) (*Data, error)
	Add(data *Data) error
}

type DataRepositoryImpl struct {
	DB *gorm.DB
}

func (u *DataRepositoryImpl) CloseConn() {
	u.DB.Close()
}

func (u *DataRepositoryImpl) Get(dataWhere *Data) (*Data, error) {
	var data Data
	err := u.DB.First(&data, dataWhere).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &data, err
}

func (u *DataRepositoryImpl) Add(data *Data) error {
	data.UUID = uuid.NewV4().String()
	return u.DB.Create(data).Error
}

func NewDataRepository() DataRepository {
	return &DataRepositoryImpl{DB: initConnection()}
}

func initConnection() *gorm.DB {
	once.Do(func() { // <-- atomic, does not allow repeating
		addr := fmt.Sprintf("postgresql://root@%s:%s/postgres?sslmode=disable", os.Getenv("DB_SERVICE_HOST"), os.Getenv("DB_SERVICE_PORT"))
		db, err := gorm.Open("postgres", addr)
		if err != nil {
			log.Fatal(err)
		}
		db.LogMode(true)
		db.AutoMigrate(&Data{})
		instance = db
	})
	return instance
}
