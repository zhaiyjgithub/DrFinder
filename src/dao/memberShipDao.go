package dao
import (
"DrFinder/src/models"
"github.com/jinzhu/gorm"
)

type MembershipDao struct {
	engine *gorm.DB
}

func NewMembershipDao(engine *gorm.DB) *MembershipDao {
	return &MembershipDao{engine:engine}
}

func (d *MembershipDao) Add(clinic *models.Membership) error {
	db := d.engine.Create(clinic)
	return db.Error
}

func (d *MembershipDao) GetMemberShipByNpi(npi int) []models.Membership {
	var MemberShip []models.Membership

	db := d.engine.Where("npi = ? ", npi).Find(&MemberShip)

	if db.Error != nil {
		return nil
	}

	return MemberShip
}
