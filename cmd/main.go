package main

import (
	"context"
	"github.com/crazybolillo/ldapb/internal"
	"github.com/glauth/ldap"
	"os"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	err := serve(ctx)
	if err != nil {
		return 1
	}

	return 0
}

func serve(ctx context.Context) error {
	server := ldap.NewServer()
	handler := internal.Handler{}

	server.BindFunc("", &handler)
	server.UnbindFunc("", &handler)
	server.SearchFunc("", &handler)
	server.AbandonFunc("", &handler)

	return server.ListenAndServe(":1389")
}
