package models

import "time"

type Question struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Statement string    `gorm:"type:text" json:"statement"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    SubjectID int    `json:"-"`
    Subject   Subject `json:"subject" gorm:"foreignKey:SubjectID"`
    
    UserID int     `json:"-"`
    User   User     `json:"user" gorm:"foreignKey:UserID"`
    
    Topics  []Topic  `gorm:"-" json:"topics"`  // Manually loaded
    Sources []Source `gorm:"-" json:"sources"` // Manually loaded
    
    Answers  []Answer  `json:"answers"`
    Comments []Comment `json:"comments"`
}

func (Question) TableName() string {
	return "question"
}

/*type Question struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Statement string `gorm:"not null"`
	SubjectID int `gorm:"not null"`
	UserID    int `gorm:"not null"`
}

func (Question) TableName() string {
	return "question"
}*/
