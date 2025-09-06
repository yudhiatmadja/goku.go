package middleware

import (
    "net/http"
    "goku-framework/bootstrap"
)

func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, _ := bootstrap.Store.Get(r, "goku-session")
        
        auth, ok := session.Values["authenticated"].(bool)
        if !ok || !auth {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}