package bootstrap

import "github.com/gorilla/sessions"

// Ganti dengan key rahasia yang aman dari config atau .env
var Store = sessions.NewCookieStore([]byte("a-very-secret-key"))