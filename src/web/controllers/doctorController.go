package controllers

import (
	"DrFinder/src/utils"
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"time"

	_ "github.com/kataras/iris/sessions/sessiondb/boltdb"
)

var validate *validator.Validate

type DoctorController struct {
	Ctx     iris.Context
	Service service.DoctorService
	GeoService service.GeoService
	AffiliationService service.AffiliationService
	AwardService service.AwardService
	CerService service.CertificationService
	ClinicService service.ClinicalService
	EduService service.EducationService
	LangService service.LangService
	MemberService service.MembershipService
	CollectionService service.CollectionService
}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	validate = validator.New()
	b.Handle(iris.MethodPost, utils.AddDoctor, "AddDoctor")
	b.Handle(iris.MethodPost, utils.GetDoctorById, "GetDoctorById")
	b.Handle(iris.MethodPost, utils.SearchDoctor, "SearchDoctor")
	b.Handle(iris.MethodPost, utils.UpdateDoctorById, "UpdateDoctorById")
	b.Handle(iris.MethodPost, utils.DeleteDoctorById, "DeleteDoctorById")
	b.Handle(iris.MethodPost, utils.SearchDoctorByPage, "SearchDoctorByPage")
	b.Handle(iris.MethodPost, utils.GetHotSearchDoctors, "GetHotSearchDoctors")
	b.Handle(iris.MethodPost, utils.GetDoctorInfoWithNpi, "GetDoctorInfoWithNpi")
	b.Handle(iris.MethodPost, utils.GetRelatedDoctors, "GetRelatedDoctors")
	b.Handle(iris.MethodPost, utils.GetCollections, "GetCollections")
	b.Handle(iris.MethodPost, utils.GetCollectionStatus, "GetCollectionStatus")
	b.Handle(iris.MethodPost, utils.DeleteCollection, "DeleteCollection")
}

func (c *DoctorController) AddDoctor() {
	type Param struct {
		Npi          int  `validate:"gt=0,numeric"`
		LastName     string `validate:"gt=0"`
		FirstName    string `validate:"gt=0"`
		MiddleName   string
		FullName     string `validate:"gt=0"`
		NamePrefix   string
		Credential   string `validate:"gt=0"`
		Gender       string `validate:"len=1"`
		Address      string `validate:"gt=0"`
		AddressSuit  string
		City         string `validate:"gt=0"`
		State        string `validate:"gt=0"`
		Zip          string `validate:"gt=0"`
		Phone        string
		Fax 		 string
		Specialty    string `validate:"gt=0"`
		SubSpecialty string
		JobTitle     string
		Summary      string
	}

	var param Param

	err:= utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	newDoctor:= &models.Doctor{
		Npi: param.Npi,
		LastName: param.LastName,
		FirstName: param.FirstName,
		MiddleName: param.MiddleName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FullName: param.FullName,
		NamePrefix: param.NamePrefix,
		Gender: param.Gender,
		Address: param.Address,
		City: param.City,
		State: param.State,
		Zip: param.Zip,
		Phone: param.Phone,
		Fax: param.Fax,
		Specialty: param.Specialty,
		SubSpecialty: param.SubSpecialty,
		JobTitle: param.JobTitle,
		Summary: param.Summary,
	}

	ok := c.Service.Add(newDoctor)

	if ok == true {
		response.Success(c.Ctx, response.Successful, nil)
	}else {
		response.Fail(c.Ctx, response.Err, "", nil)
	}
}

func (c *DoctorController)GetDoctorById() {
	type Param struct {
		DoctorId int `validate:"gt=0,numeric"`
	}

	var param Param

	err:= utils.ValidateParam(c.Ctx, validate, &param)
	if err != nil {
		return
	}

	doctor := c.Service.GetDoctorById(param.DoctorId)

	if  doctor != nil {
		response.Success(c.Ctx, response.Successful, doctor)
	}else {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}
}

func (c *DoctorController) GetDoctorBySpecialty(specialty string)  {
	type Param struct {
		Specialty string `validate:"gt=0"`
	}

	var param Param
	err:= utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	doctor:= c.Service.GetDoctorBySpecialty(param.Specialty)

	if doctor != nil {
		response.Success(c.Ctx, response.Successful, doctor)
	}else {
		response.Fail(c.Ctx, response.Err, response.NotFound, nil)
	}
}

func (c *DoctorController) SearchDoctor()  {
	type Param struct {
		FirstName string
		LastName string
		Specialty string
		Gender string
		City string
	}

	var param Param
	err:= utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var doctor models.Doctor
	doctor.FirstName = param.FirstName
	doctor.LastName = param.LastName
	doctor.Specialty = param.Specialty
	doctor.Gender = param.Gender
	doctor.City = param.City

	doctors:= c.Service.SearchDoctor(&doctor)

	response.Success(c.Ctx, response.Successful, doctors)
}

func (c *DoctorController) UpdateDoctorById() {
	type Param struct {
		ID        int
		FirstName string
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	var doctor models.Doctor
	doctor.FirstName = param.FirstName
	doctor.ID = param.ID

	err = c.Service.UpdateDoctorById(&doctor)

	if err != nil {
		response.Fail(c.Ctx, response.Err, "update failed", nil)
	} else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}

func (c *DoctorController) DeleteDoctorById()  {
	type Param struct {
		ID int `validate:"gt=0"`
	}

	var param Param

	err:= utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	ok:= c.Service.DeleteDoctorById(param.ID)

	if ok {
		response.Success(c.Ctx, response.Successful, nil )
	}else {
		response.Fail(c.Ctx, response.Err, "delete failed", nil)
	}
}

func (c *DoctorController) SearchDoctorByPage()  {
	type Param struct {
		FirstName string
		LastName string
		Specialty string
		Gender string
		City string
		State string
		Lat float64
		Lng float64
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	type DoctorGeo struct {
		Doctor models.Doctor
		Geo models.GeoDistance
	}

	doctors := c.Service.FindDoctorByPage(&models.Doctor{
		FirstName: param.FirstName,
		LastName: param.LastName,
		Specialty: param.Specialty,
		Gender: param.Gender,
		City: param.City,
		State: param.State,
		}, param.Lat, param.Lng, param.Page, param.PageSize)

	response.Success(c.Ctx, response.Successful, doctors)

	//if len(param.FirstName) == 0 && len(param.LastName) == 0 && len(param.Specialty) == 0 { // 没有任何的搜索条件
	//	geos := c.GeoService.GetNearByDoctorGeoInfo(param.Lat, param.Lng, param.Page, param.PageSize)
	//
	//	var npiList []int
	//	nMap := make(map[int]models.GeoDistance)
	//	for i:= 0; i < len(geos); i ++ {
	//		npiList = append(npiList, geos[i].Npi)
	//		nMap[geos[i].Npi] = geos[i]
	//	}
	//
	//	doctors := c.Service.SearchDoctorsByNpiList(npiList)
	//
	//	var dgList []*DoctorGeo
	//	for i := 0; i < len(doctors); i ++ {
	//		dr := doctors[i]
	//		dgList = append(dgList, &DoctorGeo{Doctor: dr, Geo: nMap[dr.Npi]})
	//	}
	//
	//	response.Success(c.Ctx, response.Successful, dgList)
	//}else {
	//	var doctor models.Doctor
	//	doctor.FirstName = param.FirstName
	//	doctor.LastName = param.LastName
	//	doctor.Specialty = param.Specialty
	//	doctor.Gender = param.Gender
	//	doctor.City = param.City
	//	doctor.State = param.State
	//
	//	doctors := c.Service.SearchDoctorByPage(&doctor, param.Page, param.PageSize)
	//
	//	var npiList []int
	//	nMap := make(map[int] models.Doctor)
	//
	//	for i := 0; i < len(doctors); i ++ {
	//		dr := doctors[i]
	//		npiList = append(npiList, dr.Npi)
	//		nMap[dr.Npi] = dr
	//	}
	//
	//	geos := c.GeoService.GetDoctorGeoInfoByNpiList(param.Lat, param.Lng, npiList)
	//
	//	var dgList []*DoctorGeo
	//	for i := 0; i < len(geos); i ++ {
	//		geo := geos[i]
	//		dgList = append(dgList, &DoctorGeo{Doctor:nMap[geo.Npi], Geo:geo})
	//	}
	//
	//	response.Success(c.Ctx, response.Successful, dgList)
	//}
}

func (c *DoctorController) GetHotSearchDoctors()  {
	doctors := c.Service.GetHotSearchDoctors()

	response.Success(c.Ctx, response.Successful, doctors)
}

func (c *DoctorController) GetRelatedDoctors()  {
	doctors := c.Service.GetRelatedDoctors(nil)
	response.Success(c.Ctx, response.Successful, doctors)
}

func (c *DoctorController) GetDoctorInfoWithNpi()  {
	type Param struct {
		Npi int `validate:"gt=0"`
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	aff := c.AffiliationService.GetAffiliationByNpi(param.Npi)
	award := c.AwardService.GetAwardByNpi(param.Npi)
	cer := c.CerService.GetCertificationByNpi(param.Npi)
	clinic := c.ClinicService.GetClinicalByNpi(param.Npi)
	edu := c.EduService.GetEducationByNpi(param.Npi)
	geo := c.GeoService.GetGeoInfoByNpi(param.Npi)
	lang := c.LangService.GetLangByNpi(param.Npi)
	member := c.MemberService.GetMemberShipByNpi(param.Npi)

	type DoctorInfo struct {
		Npi int
		Affiliation []models.Affiliation
		Award []models.Award
		Certification []models.Certification
		Clinic []models.Clinical
		Education []models.Education
		Geo models.Geo
		Lang models.Lang
		MemberShip []models.Membership
	}

	var info = &DoctorInfo{
		Npi:param.Npi,
		Affiliation:aff,
		Award:award,
		Certification:cer,
		Clinic:clinic,
		Education:edu,
		MemberShip:member,
	}

	if geo != nil {
		info.Geo = *geo
	}

	if lang != nil {
		info.Lang = *lang
	}

	response.Success(c.Ctx, response.Successful, info)
}

func (c *DoctorController) GetCollections()  {
	type Param struct {
		UserID int
		ObjectType int
	}

	var param Param

	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	collections := c.CollectionService.GetCollections(param.UserID, param.ObjectType)

	response.Success(c.Ctx, response.Successful, collections)
}

func (c *DoctorController) GetCollectionStatus()  {
	type Param struct {
		UserID int
		ObjectID int
		ObjectType int
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.CollectionService.GetIsHasCollected(param.UserID, param.ObjectID, param.ObjectType)

	isExist := true
	if err != nil {
		isExist = false
	}

	response.Success(c.Ctx, response.Successful, isExist)
}

func (c *DoctorController) DeleteCollection()  {
	type Param struct {
		UserID int
		ObjectID int
		ObjectType int
	}

	var param Param
	err := utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	err = c.CollectionService.Delete(param.UserID, param.ObjectID, param.ObjectType)

	if err != nil {
		response.Fail(c.Ctx, response.Err, response.UnknownErr, nil)
	}else {
		response.Success(c.Ctx, response.Successful, nil)
	}
}