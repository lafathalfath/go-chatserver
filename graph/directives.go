package graph

import (
	"context"
	"errors"
	"github.com/lafathalfath/go-chatserver/helpers"

	"github.com/99designs/gqlgen/graphql"
)

func AuthDirective(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	if _, ok := helpers.GetUserId(ctx); !ok {
		return nil, errors.New("Unauthorized")
	}
	return next(ctx)
}