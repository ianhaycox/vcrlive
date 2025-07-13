package model

import "github.com/ianhaycox/vcrlive/irsdk/iryaml"

type Driver struct {
	CarIdx        int    `json:"car_idx,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	UserID        int    `json:"user_id,omitempty"`
	CarClassID    int    `json:"car_class_id,omitempty"`
	CarID         int    `json:"car_id,omitempty"`
	ClassPosition int    `json:"class_position,omitempty"`
	IRating       int    `json:"irating,omitempty"`
	ClubID        int    `json:"club_id,omitempty"`
	CarNumberRaw  int    `json:"car_number_raw,omitempty"`
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
	for position, carIdx := range positions {
		if driver, ok := d[carIdx]; ok {
			driver.ClassPosition = position
			d[carIdx] = driver
		}
	}
}
