package db

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"movie-rating-api/models"
	"strings"
	"time"
)

type PostgresConnection interface {
	GetConnectionString() (string, error)
	ValidateConfig() error
}

type PostgresConfig struct {
	Host         string
	Port         string
	DatabaseName string
	User         string
	Password     string
	SSLMode      string
}

type postgresConnection struct {
	config PostgresConfig
}

func NewPostgresConnection(config PostgresConfig) PostgresConnection {
	return postgresConnection{
		config: config,
	}
}

func (conn postgresConnection) GetConnectionString() (string, error) {

	if err := conn.ValidateConfig(); err != nil {
		return "", err
	}

	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		conn.config.Host,
		conn.config.Port,
		conn.config.DatabaseName,
		conn.config.User,
		conn.config.Password,
		conn.config.SSLMode,
	), nil
}

func (conn postgresConnection) ValidateConfig() error {
	if conn.config.Host == "" {
		return fmt.Errorf("host missing in config")
	}
	if conn.config.Port == "" {
		return fmt.Errorf("port missing in config")
	}
	if conn.config.DatabaseName == "" {
		return fmt.Errorf("databaseName missing in config")
	}
	if conn.config.User == "" {
		return fmt.Errorf("user missing in config")
	}
	if conn.config.Password == "" {
		return fmt.Errorf("password missing in config")
	}
	if conn.config.SSLMode == "" {
		return fmt.Errorf("SSLMode missing in config")
	}

	return nil
}

func getPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:         "postgres",
		Port:         "5432",
		DatabaseName: "postgres",
		User:         "postgres",
		Password:     "docker",
		SSLMode:      "disable",
	}
}

var db *gorm.DB

func InitializeDB() error {
	driver, connString, err := ConnectionInfoFromEnvironment()
	if err != nil {
		return fmt.Errorf("error getting the db driver and connection string: %s", err.Error())
	}

	db = setupDB(driver, connString, true)

	return nil
}

func setupDB(driver string, connString string, autoMigrate bool) *gorm.DB {
	// once the datamodel is settled, then change
	// the 3rd param from true to false which turns off the AutoMigrate
	// and speeds up start time to about 20 seconds
	// The AutoMigrate can take upwards of a few minute or more.
	// to validate the database schema.
	database, err := Connect(driver, connString, autoMigrate)
	if err != nil {
		log.Fatalf("error connecting to database %s", err)
	}
	return database
}

// Get connection to the DB.
// returns the driver name and connection string
func ConnectionInfoFromEnvironment() (string, string, error) {
	connString, err := NewPostgresConnection(getPostgresConfig()).GetConnectionString()
	if err != nil {
		log.Fatalf("unable to get postgres connnection string: %s", err.Error())
	}

	return "postgres", connString, nil
}

func Connect(driver string, connection interface{}, autoMigrate bool) (*gorm.DB, error) {
	dbConnect, err := gorm.Open(driver, connection)
	if err != nil {
		log.Printf("Got error when connecting to database: '%v'\n", err)
		return dbConnect, err
	}

	if autoMigrate {
		dbConnect.LogMode(true)
		log.Println("Running db migration...")

		// created in a database that already has preexisting tables
		dbConnect.CreateTable(&models.Movies{})
		dbConnect.CreateTable(&models.MovieRatings{})
		dbConnect.CreateTable(&models.Ratings{})

		dbConnect.AutoMigrate(
			&models.Movies{},
			&models.MovieRatings{},
			&models.Ratings{},
		)

		dbConnect.Model(&models.Ratings{}).AddForeignKey("movie_ratings_id", "movie_ratings(id)", "RESTRICT", "RESTRICT")

		log.Println("finished db migration")
	}

	// turn this on to see the details of gorm working.
	dbConnect.LogMode(false)

	dbConnect.DB().SetMaxOpenConns(5)
	dbConnect.DB().SetMaxIdleConns(1)
	dbConnect.DB().SetConnMaxLifetime(time.Second * 30)

	return dbConnect, err
}

func InitializeMovies(dbClient Client) error {
	var moviesToCreate []models.Movies
	err := json.Unmarshal([]byte(moviesJsonString), &moviesToCreate)
	if err != nil {
		return err
	}

	for _, movie := range moviesToCreate {
		err = dbClient.CreateMovie(movie)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return err
			}
		}
	}

	var ratingsToCreate []models.MovieRatings
	err = json.Unmarshal([]byte(ratingsJsonString), &ratingsToCreate)
	if err != nil {
		return err
	}

	for _, rating := range ratingsToCreate {
		for i, _ := range rating.Ratings {
			rating.Ratings[i].MovieRatingsID = rating.ID
		}

		err = dbClient.CreateMovieRating(rating)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				return err
			}
		}
	}

	return nil
}

const moviesJsonString = `[
    {
      "Id": 1,
      "Title": "Life of Brian",
      "Year": "1979",
      "Rated": "R",
      "Released": "17 Aug 1979",
      "Runtime": "94 min",
      "Genre": "Comedy",
      "Director": "Terry Jones",
      "Writer": "Graham Chapman, John Cleese, Terry Gilliam",
      "Actors": "Graham Chapman, John Cleese, Michael Palin",
      "Plot": "Born on the original Christmas in the stable next door to Jesus Christ, Brian of Nazareth spends his life being mistaken for a messiah.",
      "Language": "English, Latin",
      "Country": "United Kingdom",
      "Awards": "N/A",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDA1ZWI4ZDItOTRlYi00OTUxLWFlNWQtMzM5NDI0YjA4ZGI2XkEyXkFqcGdeQXVyMjUzOTY1NTc@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "16 Nov 1999",
      "BoxOffice": "$20,206,622",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 2,
      "Title": "Star Wars",
      "Year": "1977",
      "Rated": "PG",
      "Released": "25 May 1977",
      "Runtime": "121 min",
      "Genre": "Action, Adventure, Fantasy",
      "Director": "George Lucas",
      "Writer": "George Lucas",
      "Actors": "Mark Hamill, Harrison Ford, Carrie Fisher",
      "Plot": "Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy from the Empire's world-destroying battle station, while also attempting to rescue Princess Leia from the mysterious Darth Vad",
      "Language": "English",
      "Country": "United States",
      "Awards": "Won 6 Oscars. 63 wins & 29 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNzVlY2MwMjktM2E4OS00Y2Y3LWE3ZjctYzhkZGM3YzA1ZWM2XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "06 Dec 2005",
      "BoxOffice": "$460,998,507",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 3,
      "Title": "The Mitchells vs the Machines",
      "Year": "2021",
      "Rated": "PG",
      "Released": "30 Apr 2021",
      "Runtime": "113 min",
      "Genre": "Animation, Adventure, Comedy",
      "Director": "Michael Rianda, Jeff Rowe",
      "Writer": "Michael Rianda, Jeff Rowe, Peter Szilagyi",
      "Actors": "Abbi Jacobson, Danny McBride, Maya Rudolph",
      "Plot": "A quirky, dysfunctional family's road trip is upended when they find themselves in the middle of the robot apocalypse and suddenly become humanity's unlikeliest last hope.",
      "Language": "English",
      "Country": "United States, Hong Kong",
      "Awards": "Nominated for 1 Oscar. 46 wins & 56 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BOTFjNjc0MTgtYmYwZi00NDcyLTlmMmYtNmJkZTI4MWJjYjM5XkEyXkFqcGdeQXVyMTA5ODEyNTc5._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "30 Apr 2021",
      "BoxOffice": "N/A",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 4,
      "Title": "Coco",
      "Year": "2017",
      "Rated": "PG",
      "Released": "22 Nov 2017",
      "Runtime": "105 min",
      "Genre": "Animation, Adventure, Comedy",
      "Director": "Lee Unkrich, Adrian Molina",
      "Writer": "Lee Unkrich, Jason Katz, Matthew Aldrich",
      "Actors": "Anthony Gonzalez, Gael García Bernal, Benjamin Bratt",
      "Plot": "Aspiring musician Miguel, confronted with his family's ancestral ban on music, enters the Land of the Dead to find his great-great-grandfather, a legendary singer.",
      "Language": "English, Spanish",
      "Country": "United States",
      "Awards": "Won 2 Oscars. 109 wins & 40 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYjQ5NjM0Y2YtNjZkNC00ZDhkLWJjMWItN2QyNzFkMDE3ZjAxXkEyXkFqcGdeQXVyODIxMzk5NjA@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "13 Feb 2018",
      "BoxOffice": "$210,460,015",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 5,
      "Title": "Logan",
      "Year": "2017",
      "Rated": "R",
      "Released": "03 Mar 2017",
      "Runtime": "137 min",
      "Genre": "Action, Drama, Sci-Fi",
      "Director": "James Mangold",
      "Writer": "James Mangold, Scott Frank, Michael Green",
      "Actors": "Hugh Jackman, Patrick Stewart, Dafne Keen",
      "Plot": "In a future where mutants are nearly extinct, an elderly and weary Logan leads a quiet life. But when Laura, a mutant child pursued by scientists, comes to him for help, he must get her to safety.",
      "Language": "English, Spanish",
      "Country": "United States",
      "Awards": "Nominated for 1 Oscar. 28 wins & 80 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYzc5MTU4N2EtYTkyMi00NjdhLTg3NWEtMTY4OTEyMzJhZTAzXkEyXkFqcGdeQXVyNjc1NTYyMjg@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "20 Jun 2017",
      "BoxOffice": "$226,277,068",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 6,
      "Title": "Akira",
      "Year": "1988",
      "Rated": "R",
      "Released": "28 Jun 1991",
      "Runtime": "124 min",
      "Genre": "Animation, Action, Drama",
      "Director": "Katsuhiro Ôtomo",
      "Writer": "Katsuhiro Ôtomo, Izô Hashimoto",
      "Actors": "Mitsuo Iwata, Nozomu Sasaki, Mami Koyama",
      "Plot": "A secret military project endangers Neo-Tokyo when it turns a biker gang member into a rampaging psychic psychopath who can only be stopped by a teenager, his gang of biker friends and a group of psychics.",
      "Language": "Japanese",
      "Country": "Japan",
      "Awards": "1 win",
      "Poster": "https://m.media-amazon.com/images/M/MV5BM2ZiZTk1ODgtMTZkNS00NTYxLWIxZTUtNWExZGYwZTRjODViXkEyXkFqcGdeQXVyMTE2MzA3MDM@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "24 Feb 2009",
      "BoxOffice": "$553,171",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 7,
      "Title": "Forrest Gump",
      "Year": "1994",
      "Rated": "PG-13",
      "Released": "06 Jul 1994",
      "Runtime": "142 min",
      "Genre": "Drama, Romance",
      "Director": "Robert Zemeckis",
      "Writer": "Winston Groom, Eric Roth",
      "Actors": "Tom Hanks, Robin Wright, Gary Sinise",
      "Plot": "The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood swe...",
      "Language": "English",
      "Country": "United States",
      "Awards": "Won 6 Oscars. 50 wins & 75 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNWIwODRlZTUtY2U3ZS00Yzg1LWJhNzYtMmZiYmEyNmU1NjMzXkEyXkFqcGdeQXVyMTQxNzMzNDI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "28 Aug 2001",
      "BoxOffice": "$330,455,270",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 8,
      "Title": "Star Wars: Episode V - The Empire Strikes Back",
      "Year": "1980",
      "Rated": "PG",
      "Released": "20 Jun 1980",
      "Runtime": "124 min",
      "Genre": "Action, Adventure, Fantasy",
      "Director": "Irvin Kershner",
      "Writer": "Leigh Brackett, Lawrence Kasdan, George Lucas",
      "Actors": "Mark Hamill, Harrison Ford, Carrie Fisher",
      "Plot": "After the Rebels are brutally overpowered by the Empire on the ice planet Hoth, Luke Skywalker begins Jedi training with Yoda, while his friends are pursued across the galaxy by Darth Vader and bounty hunter Boba Fett.",
      "Language": "English",
      "Country": "United States",
      "Awards": "Won 1 Oscar. 25 wins & 20 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYmU1NDRjNDgtMzhiMi00NjZmLTg5NGItZDNiZjU5NTU4OTE0XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "21 Sep 2004",
      "BoxOffice": "$292,753,960",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 9,
      "Title": "Rogue One: A Star Wars Story",
      "Year": "2016",
      "Rated": "PG-13",
      "Released": "16 Dec 2016",
      "Runtime": "133 min",
      "Genre": "Action, Adventure, Sci-Fi",
      "Director": "Gareth Edwards",
      "Writer": "Chris Weitz, Tony Gilroy, John Knoll",
      "Actors": "Felicity Jones, Diego Luna, Alan Tudyk",
      "Plot": "In a time of conflict, a group of unlikely heroes band together on a mission to steal the plans to the Death Star, the Empire's ultimate weapon of destruction.",
      "Language": "English",
      "Country": "United States",
      "Awards": "Nominated for 2 Oscars. 24 wins & 85 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjEwMzMxODIzOV5BMl5BanBnXkFtZTgwNzg3OTAzMDI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "04 Apr 2017",
      "BoxOffice": "$532,177,324",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 10,
      "Title": "Terminator 2: Judgment Day",
      "Year": "1991",
      "Rated": "R",
      "Released": "03 Jul 1991",
      "Runtime": "137 min",
      "Genre": "Action, Sci-Fi",
      "Director": "James Cameron",
      "Writer": "James Cameron, William Wisher",
      "Actors": "Arnold Schwarzenegger, Linda Hamilton, Edward Furlong",
      "Plot": "A cyborg, identical to the one who failed to kill Sarah Connor, must now protect her ten-year-old son John from a more advanced and powerful cyborg.",
      "Language": "English, Spanish",
      "Country": "United States",
      "Awards": "Won 4 Oscars. 36 wins & 33 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMGU2NzRmZjUtOGUxYS00ZjdjLWEwZWItY2NlM2JhNjkxNTFmXkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "13 Feb 2007",
      "BoxOffice": "$205,881,154",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 11,
      "Title": "The Lion King",
      "Year": "1994",
      "Rated": "G",
      "Released": "24 Jun 1994",
      "Runtime": "88 min",
      "Genre": "Animation, Adventure, Drama",
      "Director": "Roger Allers, Rob Minkoff",
      "Writer": "Irene Mecchi, Jonathan Roberts, Linda Woolverton",
      "Actors": "Matthew Broderick, Jeremy Irons, James Earl Jones",
      "Plot": "Lion prince Simba and his father are targeted by his bitter uncle, who wants to ascend the throne himself.",
      "Language": "English, Swahili, Xhosa, Zulu",
      "Country": "United States",
      "Awards": "Won 2 Oscars. 39 wins & 35 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BYTYxNGMyZTYtMjE3MS00MzNjLWFjNmYtMDk3N2FmM2JiM2M1XkEyXkFqcGdeQXVyNjY5NDU4NzI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "04 Oct 2011",
      "BoxOffice": "$422,783,777",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 12,
      "Title": "Spider-Man: Into the Spider-Verse",
      "Year": "2018",
      "Rated": "PG",
      "Released": "14 Dec 2018",
      "Runtime": "117 min",
      "Genre": "Animation, Action, Adventure",
      "Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman",
      "Writer": "Phil Lord, Rodney Rothman",
      "Actors": "Shameik Moore, Jake Johnson, Hailee Steinfeld",
      "Plot": "Teen Miles Morales becomes the Spider-Man of his universe, and must join with five spider-powered individuals from other dimensions to stop a threat for all realities.",
      "Language": "English, Spanish",
      "Country": "United States",
      "Awards": "Won 1 Oscar. 82 wins & 57 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMjMwNDkxMTgzOF5BMl5BanBnXkFtZTgwNTkwNTQ3NjM@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "19 Mar 2019",
      "BoxOffice": "$190,241,310",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 13,
      "Title": "Shaun of the Dead",
      "Year": "2004",
      "Rated": "R",
      "Released": "24 Sep 2004",
      "Runtime": "99 min",
      "Genre": "Comedy, Horror",
      "Director": "Edgar Wright",
      "Writer": "Simon Pegg, Edgar Wright",
      "Actors": "Simon Pegg, Nick Frost, Kate Ashfield",
      "Plot": "The uneventful, aimless lives of a London electronics salesman and his layabout roommate are disrupted by the zombie apocalypse.",
      "Language": "English",
      "Country": "United Kingdom, France, United States",
      "Awards": "Nominated for 3 BAFTA 13 wins & 20 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMTg5Mjk2NDMtZTk0Ny00YTQ0LWIzYWEtMWI5MGQ0Mjg1OTNkXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "21 Dec 2004",
      "BoxOffice": "$13,542,874",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 14,
      "Title": "Jumanji: Welcome to the Jungle",
      "Year": "2017",
      "Rated": "PG-13",
      "Released": "20 Dec 2017",
      "Runtime": "119 min",
      "Genre": "Action, Adventure, Comedy",
      "Director": "Jake Kasdan",
      "Writer": "Chris McKenna, Erik Sommers, Scott Rosenberg",
      "Actors": "Dwayne Johnson, Karen Gillan, Kevin Hart",
      "Plot": "Four teenagers are sucked into a magical video game, and the only way they can escape is to work together to finish the game.",
      "Language": "English",
      "Country": "United States",
      "Awards": "5 wins & 15 nominations",
      "Poster": "https://m.media-amazon.com/images/M/MV5BODQ0NDhjYWItYTMxZi00NTk2LWIzNDEtOWZiYWYxZjc2MTgxXkEyXkFqcGdeQXVyMTQxNzMzNDI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "20 Mar 2018",
      "BoxOffice": "$404,540,171",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 15,
      "Title": "Back to the Future",
      "Year": "1985",
      "Rated": "PG",
      "Released": "03 Jul 1985",
      "Runtime": "116 min",
      "Genre": "Adventure, Comedy, Sci-Fi",
      "Director": "Robert Zemeckis",
      "Writer": "Robert Zemeckis, Bob Gale",
      "Actors": "Michael J. Fox, Christopher Lloyd, Lea Thompson",
      "Plot": "Marty McFly, a 17-year-old high school student, is accidentally sent thirty years into the past in a time-traveling DeLorean invented by his close friend, the eccentric scientist Doc Brown.",
      "Language": "English",
      "Country": "United States",
      "Awards": "Won 1 Oscar. 22 wins & 25 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BZmU0M2Y1OGUtZjIxNi00ZjBkLTg1MjgtOWIyNThiZWIwYjRiXkEyXkFqcGdeQXVyMTQxNzMzNDI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "17 Aug 2010",
      "BoxOffice": "$212,836,762",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 16,
      "Title": "Turning Red",
      "Year": "2022",
      "Rated": "PG",
      "Released": "11 Mar 2022",
      "Runtime": "100 min",
      "Genre": "Animation, Adventure, Comedy",
      "Director": "Domee Shi",
      "Writer": "Domee Shi, Julia Cho, Sarah Streicher",
      "Actors": "Rosalie Chiang, Sandra Oh, Ava Morse",
      "Plot": "A 13-year-old girl named Meilin turns into a giant red panda whenever she gets too excited.",
      "Language": "English",
      "Country": "United States, Canada",
      "Awards": "N/A",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNjY0MGEzZmQtZWMxNi00MWVhLWI4NWEtYjQ0MDkyYTJhMDU0XkEyXkFqcGdeQXVyODc0OTEyNDU@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "11 Mar 2022",
      "BoxOffice": "N/A",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 17,
      "Title": "Monty Python and the Holy Grail",
      "Year": "1975",
      "Rated": "PG",
      "Released": "25 May 1975",
      "Runtime": "91 min",
      "Genre": "Adventure, Comedy, Fantasy",
      "Director": "Terry Gilliam, Terry Jones",
      "Writer": "Graham Chapman, John Cleese, Eric Idle",
      "Actors": "Graham Chapman, John Cleese, Eric Idle",
      "Plot": "King Arthur and his Knights of the Round Table embark on a surreal, low-budget search for the Holy Grail, encountering many, very silly obstacles.",
      "Language": "English, French, Latin",
      "Country": "United Kingdom",
      "Awards": "3 wins & 3 nominations",
      "Poster": "https://m.media-amazon.com/images/M/MV5BN2IyNTE4YzUtZWU0Mi00MGIwLTgyMmQtMzQ4YzQxYWNlYWE2XkEyXkFqcGdeQXVyNjU0OTQ0OTY@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "08 Jun 2004",
      "BoxOffice": "$1,827,696",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 18,
      "Title": "Star Wars: Episode III - Revenge of the Sith",
      "Year": "2005",
      "Rated": "PG-13",
      "Released": "19 May 2005",
      "Runtime": "140 min",
      "Genre": "Action, Adventure, Fantasy",
      "Director": "George Lucas",
      "Writer": "George Lucas, John Ostrander, Jan Duursema",
      "Actors": "Hayden Christensen, Natalie Portman, Ewan McGregor",
      "Plot": "Three years into the Clone Wars, the Jedi rescue Palpatine from Count Dooku. As Obi-Wan pursues a new threat, Anakin acts as a double agent between the Jedi Council and Palpatine and is lured into a sinister plan to rule the galaxy.",
      "Language": "English",
      "Country": "United States",
      "Awards": "Nominated for 1 Oscar. 26 wins & 63 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BNTc4MTc3NTQ5OF5BMl5BanBnXkFtZTcwOTg0NjI4NA@@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "01 Nov 2005",
      "BoxOffice": "$380,270,577",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 19,
      "Title": "Pokémon: The First Movie - Mewtwo Strikes Back",
      "Year": "1998",
      "Rated": "G",
      "Released": "10 Nov 1999",
      "Runtime": "96 min",
      "Genre": "Animation, Action, Adventure",
      "Director": "Kunihiko Yuyama, Michael Haigney",
      "Writer": "Satoshi Tajiri, Takeshi Shudo, Norman J. Grossfeld",
      "Actors": "Veronica Taylor, Rachael Lillis, Eric Stuart",
      "Plot": "Scientists genetically create a new Pokémon, Mewtwo, but the results are horrific and disastrous.",
      "Language": "Japanese",
      "Country": "Japan",
      "Awards": "3 wins & 6 nominations",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMTkyNDQxOTg5MF5BMl5BanBnXkFtZTYwODA2MDE3._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "21 Mar 2000",
      "BoxOffice": "$85,744,662",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    },
    {
      "Id": 20,
      "Title": "Toy Story",
      "Year": "1995",
      "Rated": "G",
      "Released": "22 Nov 1995",
      "Runtime": "81 min",
      "Genre": "Animation, Adventure, Comedy",
      "Director": "John Lasseter",
      "Writer": "John Lasseter, Pete Docter, Andrew Stanton",
      "Actors": "Tom Hanks, Tim Allen, Don Rickles",
      "Plot": "A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.",
      "Language": "English",
      "Country": "United States",
      "Awards": "Nominated for 3 Oscars. 27 wins & 23 nominations total",
      "Poster": "https://m.media-amazon.com/images/M/MV5BMDU2ZWJlMjktMTRhMy00ZTA5LWEzNDgtYmNmZTEwZTViZWJkXkEyXkFqcGdeQXVyNDQ2OTk4MzI@._V1_SX300.jpg",
      "Type": "movie",
      "DVD": "23 Mar 2010",
      "BoxOffice": "$223,225,679",
      "Production": "N/A",
      "Website": "N/A",
      "Response": "True"
    }
  ]`
const ratingsJsonString = `[
  {
    "Id": 1,
    "Title": "Life of Brian",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 80
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 95
      },
      {
        "Source": "Metacritic",
        "Value": 77
      }
    ]
  },
  {
    "Id": 2,
    "Title": "Star Wars",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 86
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 92
      },
      {
        "Source": "Metacritic",
        "Value": 90
      }
    ]
  },
  {
    "Id": 3,
    "Title": "The Mitchells vs the Machines",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 77
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 97
      },
      {
        "Source": "Metacritic",
        "Value": 81
      }
    ]
  },
  {
    "Id": 4,
    "Title": "Coco",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 84
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 97
      },
      {
        "Source": "Metacritic",
        "Value": 81
      }
    ]
  },
  {
    "Id": 5,
    "Title": "Logan",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 81
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 94
      },
      {
        "Source": "Metacritic",
        "Value": 77
      }
    ]
  },
  {
    "Id": 6,
    "Title": "Akira",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 80
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 90
      },
      {
        "Source": "Metacritic",
        "Value": 67
      }
    ]
  },
  {
    "Id": 7,
    "Title": "Forrest Gump",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 88
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 70
      },
      {
        "Source": "Metacritic",
        "Value": 82
      }
    ]
  },
  {
    "Id": 8,
    "Title": "Star Wars: Episode V - The Empire Strikes Back",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 87
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 94
      },
      {
        "Source": "Metacritic",
        "Value": 82
      }
    ]
  },
  {
    "Id": 9,
    "Title": "Rogue One: A Star Wars Story",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 78
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 84
      },
      {
        "Source": "Metacritic",
        "Value": 65
      }
    ]
  },
  {
    "Id": 10,
    "Title": "Terminator 2: Judgment Day",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 86
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 93
      },
      {
        "Source": "Metacritic",
        "Value": 75
      }
    ]
  },
  {
    "Id": 11,
    "Title": "The Lion King",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 85
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 93
      },
      {
        "Source": "Metacritic",
        "Value": 88
      }
    ]
  },
  {
    "Id": 12,
    "Title": "Spider-Man: Into the Spider-Verse",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 84
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 97
      },
      {
        "Source": "Metacritic",
        "Value": 87
      }
    ]
  },
  {
    "Id": 13,
    "Title": "Shaun of the Dead",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 79
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 92
      },
      {
        "Source": "Metacritic",
        "Value": 76
      }
    ]
  },
  {
    "Id": 14,
    "Title": "Jumanji: Welcome to the Jungle",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 69
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 76
      },
      {
        "Source": "Metacritic",
        "Value": 58
      }
    ]
  },
  {
    "Id": 15,
    "Title": "Back to the Future",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 86
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 96
      },
      {
        "Source": "Metacritic",
        "Value": 87
      }
    ]
  },
  {
    "Id": 16,
    "Title": "Turning Red",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 70
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 94
      },
      {
        "Source": "Metacritic",
        "Value": 83
      }
    ]
  },
  {
    "Id": 17,
    "Title": "Monty Python and the Holy Grail",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 82
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 97
      },
      {
        "Source": "Metacritic",
        "Value": 91
      }
    ]
  },
  {
    "Id": 18,
    "Title": "Star Wars: Episode III - Revenge of the Sith",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 76
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 80
      },
      {
        "Source": "Metacritic",
        "Value": 68
      }
    ]
  },
  {
    "Id": 19,
    "Title": "Pokémon: The First Movie - Mewtwo Strikes Back",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 62
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 16
      },
      {
        "Source": "Metacritic",
        "Value": 35
      }
    ]
  },
  {
    "Id": 20,
    "Title": "Toy Story",
    "Ratings": [
      {
        "Source": "Internet Movie Database",
        "Value": 83
      },
      {
        "Source": "Rotten Tomatoes",
        "Value": 100
      },
      {
        "Source": "Metacritic",
        "Value": 95
      }
    ]
  }
]`
