package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

// uniq key we can use to store and get auth status for request context
