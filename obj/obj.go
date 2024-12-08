package obj

type UserObj struct {
	TitleId      string `json:"TitleId" validate:"required"`
	UserName     string `json:"UserName" validate:"required,min=8"`
	NickName     string `json:"NickName"`
	FirstName    string `json:"FirstName" validate:"required"`
	LastName     string `json:"LastName" validate:"required"`
	Gender       string `json:"Gender"`
	Location     string `json:"Location"`
	AgeRange     string `json:"AgeRange"`
	Email        string `json:"email" validate:"required,email"`
	MobileNumber string `json:"MobileNumber" validate:"required,min=8"`
	Status       string `json:"Status" validate:"required,min=1"`
	Password     string `json:"Password" validate:"required,min=8"`
}

type UserResponse struct {
	Id           uint
	TitleId      string
	UserName     string
	NickName     string
	FirstName    string
	LastName     string
	Email        string
	Gender       string
	Location     string
	AgeRange     string
	MobileNumber string
	Status       string
	CreatedAt    string
}

type CartObj struct {
	UserId      int    `gorm:"column:UserId"`
	CartTypeId  int    `gorm:"column:CartTypeId"`
	CartName    string `gorm:"column:CartName"`
	Description string `gorm:"column:Description"`
	GroupId     int    `gorm:"column:GroupId"`
	CreatedById int    `gorm:"column:CreatedById"`
}

type EmailObj struct {
	ToEmail  string `gorm:"column:ToEmail"`
	Subject  string `gorm:"column:Subject"`
	MailBody string `gorm:"column:MailBody"`
}

type CartUserObj struct {
	RingMasterEmail string `gorm:"column:RingMasterEmail"`
	MemberEmail     string `gorm:"column:MemberEmail"`
	CartId          int    `gorm:"column:CartId"`
	RingStatus      int    `gorm:"column:RingStatus"`
}

type ResponseMessage struct {
	ResponseCode    string
	ResponseMessage string
}

type RemoveUserFromCartObj struct {
	CartId          int    `gorm:"column:CartId"`
	RingMasterEmail string `gorm:"column:RingMasterEmail"`
	MemberEmail     string `gorm:"column:MemberEmail"`
}

type CloseCartObj struct {
	CartId          int    `gorm:"column:CartId"`
	RingMasterEmail string `gorm:"column:RingMasterEmail"`
	Reason          string `gorm:"column:Reason"`
}

type ConfigStruct struct {
	CreateTable          bool
	IsDropExistingTables bool
}

type TokenResp struct {
	Token string
}
