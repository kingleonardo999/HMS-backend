package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	LoginId  string `json:"loginId" gorm:"unique;not null"`  // 登录ID，唯一
	Password string `json:"loginPwd" gorm:"not null"`        // 密码，不能为空
	Name     string `json:"name" gorm:"not null"`            // 姓名，不能为空
	Phone    string `json:"phone" gorm:"not null"`           // 电话，不能为空
	Email    string `json:"email"`                           // 邮箱，允许为空
	RoleId   uint   `json:"roleId" gorm:"not null"`          // 角色ID，不能为空
	Role     Role   `gorm:"foreignKey:RoleId;references:ID"` // 关联角色
	ImgId    uint   `json:"imgId"`                           // 头像ID，允许为空
	Img      Img    `gorm:"foreignKey:ImgId;references:ID"`  // 关联头像
}
