package main

import (
	"io"

	"gopkg.in/yaml.v3"

	"github.com/committeddb/db-connection-parser/pkg/parser"
)

type Routes struct {
	Routes []*RouteParser `yaml:"routes"`
}

type RouteParser struct {
	TableName     string `yaml:"table"`
	DatastoreName string `yaml:"datastore"`
	Query         string `yaml:"query"`
	JWTVerifyURL  string `yaml:"jwtVerifyURL"`
}

type Route struct {
	Datastore    parser.Datastore
	TableName    string
	Query        string
	JWTVerifyURL string
}

func Parse(r io.Reader) ([]*Route, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var rs Routes
	err = yaml.Unmarshal(bs, &rs)
	if err != nil {
		return nil, err
	}

	dm, err := parser.Parse(r)
	if err != nil {
		return nil, err
	}

	var routes []*Route
	for _, r := range rs.Routes {
		routes = append(
			routes,
			&Route{
				Datastore:    dm[r.DatastoreName],
				TableName:    r.TableName,
				Query:        r.Query,
				JWTVerifyURL: r.JWTVerifyURL,
			},
		)
	}

	return routes, nil
}
