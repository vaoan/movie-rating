package db

import (
	"github.com/jinzhu/gorm"
	"movie-rating-api/models"
	"time"
)

// https://gorm.io/docs/index.html

type DB interface {
	GetMovies() ([]models.Movies, error)
	GetMovieRatings() ([]models.MovieRatings, error)
	CreateMovie(movie models.Movies) error
	CreateMovieRating(rating models.MovieRatings) error
}

type Client interface {
	DB
}

type dbClient struct {
	Gorm *gorm.DB
}

func NewDBCLient(gormDB *gorm.DB) Client {
	if gormDB == nil {
		gormDB = db
	}

	return &dbClient{
		Gorm: gormDB,
	}
}

// GetMovies returns movie information but not ratings
func (d dbClient) GetMovies() ([]models.Movies, error) {
	// this is simulating a slow api call. You can not change this for the purposes of the interview
	time.Sleep(3 * time.Second)

	var result []models.Movies
	err := d.Gorm.Find(&result).Error
	if err != nil {
		return []models.Movies{}, err
	}

	return result, nil
}

// GetMovieRatings returns movie ratings
func (d dbClient) GetMovieRatings() ([]models.MovieRatings, error) {
	// this is simulating a slow api call. You can not change this for the purposes of the interview
	time.Sleep(3 * time.Second)

	var result []models.MovieRatings
	err := d.Gorm.Preload("Ratings").Find(&result).Error
	if err != nil {
		return []models.MovieRatings{}, err
	}

	return result, nil
}

func (d dbClient) CreateMovie(movie models.Movies) error {
	return d.Gorm.Create(&movie).Error
}

func (d dbClient) CreateMovieRating(rating models.MovieRatings) error {
	return d.Gorm.Create(&rating).Error
}
