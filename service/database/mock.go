package database

import (
	"github.com/graphql-go/graphql"
	"net/http"
	"io/ioutil"
	"io"
	"log"
	"fmt"
	"encoding/json"
)

var (
	OkComputer      Disc
	TheQueenIsDead  Disc
	BeHereNow       Disc
	AppetiteForDest Disc
	BackToBlack     Disc
	HotelCal        Disc

	DiscData map[int]Disc

	discType *graphql.Object

	MusicSchema graphql.Schema
)

type Disc struct {
	Id     string
	Title  string
	Artist string
	Year   int
}

func init() {
	OkComputer = Disc{
		Title:  "OK Computer",
		Artist: "Radiohead",
		Year:   1997,
		Id:     "1",
	}
	TheQueenIsDead = Disc{
		Title:  "The Queen is dead",
		Artist: "The Smiths",
		Year:   1986,
		Id:     "2",
	}
	BeHereNow = Disc{
		Title:  "Be Here Now",
		Artist: "Oasis",
		Year:   1997,
		Id:     "3",
	}
	AppetiteForDest = Disc{
		Title:  "Appetite for Destruction",
		Artist: "Guns N' Roses ",
		Year:   1987,
		Id:     "4",
	}
	BackToBlack = Disc{
		Title:  "Back To Black",
		Artist: "Amy Winehouse",
		Year:   2006,
		Id:     "5",
	}
	HotelCal = Disc{
		Title:  "Hotel California",
		Artist: "Eagles",
		Year:   1976,
		Id:     "6",
	}

	DiscData = map[int]Disc{
		1: OkComputer,
		2: TheQueenIsDead,
		3: BeHereNow,
		4: AppetiteForDest,
		5: BackToBlack,
		6: HotelCal,
	}

	discType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Disc",
		Description: "A set of songs from one or many artists.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The Identifier of the disc.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if disc, ok := p.Source.(Disc); ok {
						return disc.Id, nil
					}
					return nil, nil
				},
			},
			"title": &graphql.Field{
				Type:        graphql.String,
				Description: "The Title of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if disc, ok := p.Source.(Disc); ok {
						return disc.Title, nil
					}
					return nil, nil
				},
			},
			"artist": &graphql.Field{
				Type:        graphql.String,
				Description: "The Artist of the album.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if disc, ok := p.Source.(Disc); ok {
						return disc.Artist, nil
					}
					return nil, nil
				},
			},
			"year": &graphql.Field{
				Type:        graphql.Int,
				Description: "The release year.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if disc, ok := p.Source.(Disc); ok {
						return disc.Year, nil
					}
					return nil, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"discs": &graphql.Field{
				Type: graphql.NewList(discType), // we return a list of discType, defined above
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetAllDiscs(), nil // <-- every time this query is called, return the result of GetAllDiscs()
				},
			},
		},
	})

	MusicSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		// mutation will be added later
	})
}

//GetAllDiscs return all the dummy data
func GetAllDiscs() []Disc {
	var discs []Disc
	for _, disc := range DiscData {
		discs = append(discs, disc)
	}
	return discs
}

//AddDisc is called every time a createDiscMutation is requested.
// It adds a new disc to the list and returns the list with the newly added object
func AddDisc(newDisc Disc) []Disc {
	DiscData[len(DiscData)+1] = newDisc
	return GetAllDiscs()
}

func GetDiscs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error GetDiscs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error GetDiscs", err)
	}
	var apolloQuery map[string]interface{}
	fmt.Println("Received request GetDiscs")
	if err := json.Unmarshal(body, &apolloQuery); err != nil { // unmarshall body contents as a type query
		fmt.Println(err)
		fmt.Println("Error on Unmarshalling!!!")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error GetDiscs unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	query := apolloQuery["query"]
	variables := apolloQuery["variables"]
	result := graphql.Do(graphql.Params{
		Schema:         MusicSchema,
		RequestString:  query.(string),
		VariableValues: variables.(map[string]interface{}),
	})
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
	return
}
