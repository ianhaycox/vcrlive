package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

type Driver struct {
	CarIdx        int    `json:"car_idx"`
	UserName      string `json:"user_name"`
	UserID        int    `json:"user_id"`
	CarClassID    int    `json:"car_class_id"`
	CarID         int    `json:"car_id"`
	ClassPosition int    `json:"class_position"`
	LapsCompleted int    `json:"laps_completed"`
	IRating       int    `json:"irating"`
	ClubID        int    `json:"club_id"`
	CarNumberRaw  int    `json:"car_number_raw"`
}

type Drivers map[int]Driver

func NewDrivers(drivers []iryaml.Driver) Drivers {
	d := make(Drivers, len(drivers))

	for _, driver := range drivers {
		if driver.IsPaceCar() || driver.Spectating() || driver.IsAI() {
			continue
		}

		d[driver.CarIdx] = Driver{
			CarIdx:       driver.CarIdx,
			UserName:     driver.UserName,
			UserID:       driver.UserID,
			CarClassID:   driver.CarClassID,
			CarID:        driver.CarID,
			IRating:      driver.IRating,
			ClubID:       driver.ClubID,
			CarNumberRaw: driver.CarNumberRaw,
		}
	}

	return d
}

func (d Drivers) SetPositions(positions []int) {
	for carIdx, position := range positions {
		if driver, ok := d[carIdx]; ok {
			driver.ClassPosition = position
			d[carIdx] = driver
		}
	}
}

func (d Drivers) SetLaps(laps []int) {
	for carIdx, lap := range laps {
		if driver, ok := d[carIdx]; ok {
			driver.LapsCompleted = lap
			d[carIdx] = driver
		}
	}
}
