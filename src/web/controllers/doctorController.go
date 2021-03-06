package controllers

import (
	"DrFinder/src/models"
	"DrFinder/src/response"
	"DrFinder/src/service"
	"DrFinder/src/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"time"

	_ "github.com/kataras/iris/sessions/sessiondb/boltdb"
)

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
	UserTrackService service.UserTrackService
	DoctorElasticService service.DoctorElasticService

}

func (c *DoctorController) BeforeActivation(b mvc.BeforeActivation)  {
	b.Handle(iris.MethodPost, utils.AddDoctor, "AddDoctor")
	b.Handle(iris.MethodPost, utils.GetDoctorById, "GetDoctorById")
	b.Handle(iris.MethodPost, utils.SearchDoctor, "SearchDoctor")
	b.Handle(iris.MethodPost, utils.UpdateDoctorById, "UpdateDoctorById")
	b.Handle(iris.MethodPost, utils.DeleteDoctorById, "DeleteDoctorById")
	b.Handle(iris.MethodPost, utils.SearchDoctorByPage, "SearchDoctorByPage")
	b.Handle(iris.MethodPost, utils.GetHotSearchDoctors, "GetHotSearchDoctors")
	b.Handle(iris.MethodPost, utils.GetDoctorInfoWithNpi, "GetDoctorInfoWithNpi")
	b.Handle(iris.MethodPost, utils.GetRelatedDoctors, "GetRelatedDoctors")
	b.Handle(iris.MethodPost, utils.GetCollections, "GetCollections", j.Serve)
	b.Handle(iris.MethodPost, utils.GetCollectionStatus, "GetCollectionStatus", j.Serve)
	b.Handle(iris.MethodPost, utils.DeleteCollection, "DeleteCollection", j.Serve)
	b.Handle(iris.MethodPost, utils.SearchDoctorES, "SearchDoctorES")
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

func (c *DoctorController) SearchDoctorES()  {
	type Param struct {
		Name string
		Specialty string
		Gender int
		Address string
		City string
		State string
		ZipCode int
		Page int
		PageSize int
		Lat float64
		Lng float64
		Platform string
		UserID int
	}

	var param Param
	err:= utils.ValidateParam(c.Ctx, validate, &param)

	if err != nil {
		return
	}

	npiList := make([]int, 0)

	npiList = c.DoctorElasticService.QueryDoctor(param.Name,
			param.Specialty,
			param.Gender,
			param.State,
			param.City,
			param.Address,
			param.ZipCode,
			param.Page,
			param.PageSize,
		)

	doctors := make([]*models.Doctor, 0)
	if len(npiList) > 0 {
		doctors = c.Service.GetDoctorByNpiList(npiList)
	}

	response.Success(c.Ctx, response.Successful, doctors)


	//add Track records
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	gender := "M"
	if param.Gender == 2 {
		gender = "F"
	}

	record := &models.UserSearchDrRecord{
		Name: param.Name,
		Specialty: param.Specialty,
		Gender: gender,
		City: param.City,
		State: param.State,
		Lat: param.Lat,
		Lng: param.Lng,
		Page: param.Page,
		PageSize: param.PageSize,
		Platform: param.Platform,
		UserID: param.UserID,
		CreatedDate: now,
	}
	_ = c.UserTrackService.AddSearchDrsRecord(record)

	var records []models.DrSearchResultRecord
	for i := 0; i < len(doctors); i ++ {
		doctor := doctors[i]
		records = append(records, models.DrSearchResultRecord{
			Npi:         doctor.Npi,
			Lat:         param.Lat, //可以检索这个医生被哪里的病人查询过
			Lng:         param.Lng,
			Platform:    param.Platform,
			UserID:      param.UserID,
			CreatedDate: now,
		})
	}

	if len(records) > 0 {
		_ = c.UserTrackService.AddSearchDrResultRecords(records)
	}
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
		Name string
		Specialty string
		Gender string
		City string
		State string
		Lat float64
		Lng float64
		Page int `validate:"gt=0"`
		PageSize int `validate:"gt=0"`
		Platform string `validate:"-"`
		UserID int `validate:"-"`
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
		LastName: param.Name,
		FirstName: param.Name,
		Specialty: param.Specialty,
		Gender: param.Gender,
		City: param.City,
		State: param.State,
		}, param.Lat, param.Lng, param.Page, param.PageSize)

	response.Success(c.Ctx, response.Successful, doctors)

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	record := &models.UserSearchDrRecord{
		Name: param.Name,
		Specialty: param.Specialty,
		Gender: param.Gender,
		City: param.City,
		State: param.State,
		Lat: param.Lat,
		Lng: param.Lng,
		Page: param.Page,
		PageSize: param.PageSize,
		Platform: param.Platform,
		UserID: param.UserID,
		CreatedDate: now,
	}
	_ = c.UserTrackService.AddSearchDrsRecord(record)

	var records []models.DrSearchResultRecord
	for i := 0; i < len(doctors); i ++ {
		doctor := doctors[i]
		records = append(records, models.DrSearchResultRecord{
			Npi:         doctor.Npi,
			Lat:         doctor.Lat,
			Lng:         doctor.Lng,
			Platform:    param.Platform,
			UserID:      param.UserID,
			CreatedDate: now,
		})
	}

	if len(records) > 0 {
		_ = c.UserTrackService.AddSearchDrResultRecords(records)
	}
}

func (c *DoctorController) GetHotSearchDoctors()  {
	doctors := c.Service.GetHotSearchDoctors()
	response.Success(c.Ctx, response.Successful, doctors)
}

func (c *DoctorController) GetRelatedDoctors()  {
	var doc models.Doctor
	err := utils.ValidateParam(c.Ctx, validate, &doc)
	if err != nil {
		return
	}

	doctors := c.Service.GetRelatedDoctors(&doc)
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
		Affiliation []*models.Affiliation
		Award []*models.Award
		Certification []*models.Certification
		Clinic []*models.Clinical
		Education []*models.Education
		Geo models.Geo
		Lang models.Lang
		MemberShip []*models.Membership
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