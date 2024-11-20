package models

type App struct {
	ID            int32  `db:"id"`
	Name          string `db:"name"`
	Secret        string `db:"secret"`
	RefreshSecret string `db:"refresh_secret"`
}
