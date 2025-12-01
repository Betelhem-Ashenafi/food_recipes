package models

// User model for authentication and profile
 type User struct {
	ID        int    `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	Name      string `db:"name" json:"name"`
	AvatarURL string `db:"avatar_url" json:"avatar_url"`
}
