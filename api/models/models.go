package models

type MoviesReturnObject struct {
	Title         string    `json:"title"`
	Ratings       []Ratings `json:"ratings"`
	AverageRating int       `json:"average_rating"`
	Plot          string
	Genre         string
}
type Movies struct {
	ID    int    `gorm:"primary_key"`
	Title string `gorm:"unique;not null"`
	Plot  string
	Genre string
	Year  string
	Rated string
}

type MovieRatings struct {
	ID      int       `gorm:"primary_key"`
	Title   string    `gorm:"unique;not null"`
	Ratings []Ratings `json:"ratings"`
}

type Ratings struct {
	MovieRatingsID int           `json:"movie_ratings_id,omitempty"`
	MovieRatings   *MovieRatings `json:"movie_ratings,omitempty" gorm:"foreignKey:MovieRatingsID"`
	Source         string        `json:"source"`
	Value          int           `json:"value"`
}
