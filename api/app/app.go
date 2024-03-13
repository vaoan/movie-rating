package app

import (
	"movie-rating-api/db"
	"movie-rating-api/models"
	"net/url"
)

type App interface {
	GetMovies(query url.Values) ([]models.MoviesReturnObject, error)
}

func GetMovies(query url.Values) ([]models.MoviesReturnObject, error) {
	// dbClient represents a slow microservice that brings back data
	dbClient := db.NewDBCLient(nil)

	movies, err := dbClient.GetMovies()
	if err != nil {
		return []models.MoviesReturnObject{}, err
	}

	ratings, err := dbClient.GetMovieRatings()
	if err != nil {
		return []models.MoviesReturnObject{}, err
	}

	var movieReturn []models.MoviesReturnObject
	for _, movie := range movies {
		for _, rating := range ratings {
			if movie.Title == rating.Title {
				var total int
				for _, r := range rating.Ratings {
					total += r.Value
				}

				var avg int
				if len(rating.Ratings) != 0 {
					avg = total / len(rating.Ratings)
				}

				movieReturn = append(movieReturn, models.MoviesReturnObject{
					Title:         movie.Title,
					Genre:         movie.Genre,
					Ratings:       rating.Ratings,
					AverageRating: avg,
					Plot:          movie.Plot,
				})
			}
		}
	}

	return movieReturn, nil
}
