package doctor

import (
	"time"
)

/*

CREATE TABLE `favordoctor` (
`favor_id` int(11) NOT NULL,
`user_id` int(11) NULL,
`doctor_id` int(11) NULL,
PRIMARY KEY (`favor_id`)
);
*/

type FavorDoctor struct {
	ID       uint
	CreateAt time.Time
	UpdateAt time.Time
	UserId   int
	DoctorId int
}
