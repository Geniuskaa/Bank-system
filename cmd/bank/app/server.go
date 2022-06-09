package app

import (
	"Bank-system/cmd/bank/app/dto"
	"Bank-system/pkg/card"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *mux.Router
}

func NewServer(cardSvc *card.Service, mux *mux.Router) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {

	s.mux.HandleFunc("/getUserCards", s.getUserCards).Methods(http.MethodGet)
	s.mux.HandleFunc("/getCardTransactions", s.getCardTransactions)

	s.mux.HandleFunc("/returnPanic", s.returnPanic).Methods(http.MethodGet)
	s.mux.Use(s.recoverer) // Добавление middleWare который ловит панику и обрабатывает её

	//s.mux.HandleFunc("/addUser", s.addUser)
	//s.mux.HandleFunc("/addCard", s.addCard)
	//s.mux.HandleFunc("/editCard", s.editCard)
	//s.mux.HandleFunc("/removeCard", s.removeCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func marshalDtos(anyObject interface{}) ([]byte, error) {
	respBody, err := json.Marshal(anyObject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return respBody, nil
}

//func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusBadGateway)
//		return
//	}
//
//	var newUserDTO dto.NewUserDTO
//	err = json.Unmarshal(body, &newUserDTO)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	s.cardSvc.AddUser(r.Context(), newUserDTO)
//}

func (s *Server) returnPanic(w http.ResponseWriter, r *http.Request) {
	panic("I will be caught by middleware!")
}

func (s *Server) recoverer(handler http.Handler) http.Handler {
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

func (s *Server) getUserCards(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	log.Println(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.cardSvc.FindUserCardsById(r.Context(), int64(id)) //FindUserCardsById(r.Context(), int64(id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dtos := make([]*dto.CardDTO, len(cards))
	for i, c := range cards {

		if c == nil {
			continue
		}

		dtos[i] = &dto.CardDTO{
			Id:      c.Id,
			Number:  c.Number,
			Issuer:  c.Issuer,
			Type:    c.Type,
			Balance: c.Balance,
		}
	}

	respBody, err := marshalDtos(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	// по умолчанию статус 200 Ok
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getCardTransactions(writer http.ResponseWriter, request *http.Request) {
	cardID, err := strconv.Atoi(request.URL.Query().Get("cardID"))
	if err != nil {
		log.Println(err)
		return
	}

	transactions, err := s.cardSvc.FindCardTransactionsByCardId(request.Context(), int64(cardID))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	respondBody, err := marshalDtos(transactions)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	_, err = writer.Write(respondBody)
	if err != nil {
		log.Println(err)
		return
	}
}

//func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
//	var card dto.NewCardDTO
//
//	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	user, err := s.cardSvc.FindUserById(r.Context(), int64(userID))
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusBadGateway)
//		return
//	}
//	defer func() {
//		err := r.Body.Close()
//		if err != nil {
//			return
//		}
//	}()
//
//	err = json.Unmarshal(body, &card)
//	if err != nil {
//		log.Println(err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	user.AddCard(card)
//
//}
//
//func (s *Server) editCard(w http.ResponseWriter, r *http.Request) {
//	panic("implement me")
//}
//
//func (s *Server) removeCard(w http.ResponseWriter, r *http.Request) {
//	panic("implement me")
//}
