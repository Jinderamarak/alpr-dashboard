package data

import (
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupCarRepository(t *testing.T) CarRepository {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal("open database failed:", err)
		return nil
	}

	err = db.AutoMigrate(&model.Car{}, &model.Recognition{})
	if err != nil {
		t.Fatal("migration failed:", err)
		return nil
	}

	repo := NewCarRepository(db)
	return repo
}

func TestGetOrCreateByPlate(t *testing.T) {
	repo := setupCarRepository(t)
	if repo == nil {
		return
	}

	carPlate := "1A2-3456"
	car, err := repo.GetOrCreateByPlate(carPlate)
	if err != nil {
		t.Fatal("Failed to create car with plate:", err)
		return
	}

	if len(car.ID) == 0 || car.ID == (uuid.UUID{}).String() {
		t.Fatal("Uuid was not generated:", car.ID)
		return
	}

	if car.Plate != carPlate {
		t.Fatalf(`Expected car to have plate "%v", but got "%v"`, carPlate, car.Plate)
		return
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatal("Failed to count cars:", err)
		return
	}

	if count != 1 {
		t.Fatalf(`Expected count to be 1, but was %v`, count)
		return
	}

	next, err := repo.GetOrCreateByPlate(carPlate)
	if err != nil {
		t.Fatal("Failed to retreive car with plate:", err)
		return
	}

	if next.ID != car.ID {
		t.Fatalf(`Expected car plates to match, previous was "%v", retreived "%v"`, car.ID, next.ID)
		return
	}
}
