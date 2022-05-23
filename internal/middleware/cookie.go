package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/url-shortener/pkg/hash"
)

const cookieName = "user_id"

func Cookie(keyManager hash.KeyManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			cookie, err := request.Cookie(cookieName)

			if err != nil {
				if errors.Is(http.ErrNoCookie, err) {
					userId := uuid.New().String()
					cookie = &http.Cookie{
						Name:  cookieName,
						Value: keyManager.Encode(userId),
						Path:  "/",
					}
					http.SetCookie(writer, cookie)
				} else {
					log.Println(err)
					return
				}
			}

			userId, err := keyManager.Decode(cookie.Value)

			if err != nil {
				log.Println(err)
				return
			}

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), cookieName, userId)))
		})
	}
}
