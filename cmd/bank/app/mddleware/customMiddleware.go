package mddleware

import (
	"log"
	"net/http"
)

//Кастомный middleware для ловли паники. Сделан исключительно в целях практики. В проде используй провереренные решения!
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

//func (s *Server) Authorization(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//		log.Printf("new request: %s %s", request.Method, request.URL.Path)
//
//		token := request.Header.Get("Authorization")
//
//		row := s.pool.QueryRow(request.Context(), `SELECT role from user_to_token join users ON
//   	user_to_token.user_id = users.id where token = $1`, token)
//		if row == nil {
//			writer.WriteHeader(http.StatusInternalServerError)
//		}
//		var userRole user.Role
//		err := row.Scan(&userRole)
//		if err != nil {
//			log.Println(err)
//			writer.WriteHeader(http.StatusInternalServerError)
//			handler.ServeHTTP(writer, request)
//			return
//		}
//
//		ctx := context.WithValue(request.Context(), authenticationContextKey, userRole)
//		request = request.WithContext(ctx)
//
//		handler.ServeHTTP(writer, request)
//	})
//}
