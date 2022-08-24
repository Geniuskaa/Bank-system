package serviceMiddleware

import (
	"context"
	"github.com/Geniuskaa/Bank-system/pkg/app/contextKey"
	"github.com/Geniuskaa/Bank-system/pkg/app/user"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

// Кастомный serviceMiddleware для ловли паники. Сделан исключительно в целях практики. В проде используй провереренные решения!
func Recoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("new request: %s %s", request.Method, request.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				log.Println("panic occurred:", err)
			}
		}()
		handler.ServeHTTP(writer, request)
	})
}

type ReturningUserRoleFunc func(ctx context.Context, token string) pgx.Row

func Authorization(roleFunc ReturningUserRoleFunc) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			//log.Printf("new request: %s %s", request.Method, request.URL.Path)
			token := request.Header.Get("Authorization")

			row := roleFunc(request.Context(), token)
			if row == nil {
				writer.WriteHeader(http.StatusInternalServerError)
			}
			var userRole user.Role
			err := row.Scan(&userRole)
			if err != nil {
				log.Println(err)
				writer.WriteHeader(http.StatusForbidden)
				handler.ServeHTTP(writer, request)
				return
			}

			ctx := context.WithValue(request.Context(), contextKey.AuthenticationContextKey, userRole)
			request = request.WithContext(ctx)

			handler.ServeHTTP(writer, request)
		})
	}
}
