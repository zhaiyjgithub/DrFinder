package dao
import (
"DrFinder/src/models"
"github.com/jinzhu/gorm"
)

type Membership struct {
	engine *gorm.DB
}

func NewMembership(engine *gorm.DB) *Membership {
	return &Membership{engine:engine}
}

func (d *Membership) Add(clinic *models.Membership) error {
	db := d.engine.Create(clinic)
	return db.Error
}

func (d *Membership) GetMemberShipByNpi(npi int) *models.Membership {
	var MemberShip models.Membership

	db := d.engine.Where("npi = ? ", npi).Find(&MemberShip)

	if db.Error != nil {
		return nil
	}

	return &MemberShip
}
