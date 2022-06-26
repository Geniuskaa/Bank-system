package handlr

import (
	"Bank-system/cmd/bank/app/contextKey"
	"Bank-system/cmd/bank/app/dto"
	"Bank-system/cmd/bank/app/serviceMiddleware/cache"
	"time"

	//"Bank-system/cmd/bank/app/server"
	"Bank-system/pkg/additionalservice/cinema"
	"Bank-system/pkg/additionalservice/cinema/model"
	"Bank-system/pkg/card"
	"Bank-system/pkg/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	cardSvc   *card.Service
	UserSvc   *user.Service
	cinemaSvc *cinema.Service
	CacheSvc  *cache.Service
}

func NewHandler(pool *pgxpool.Pool, mongoDB *mongo.Database, cachePool *redis.Pool) *Handler {
	return &Handler{
		cardSvc:   card.NewService(pool),
		UserSvc:   user.NewService(pool),
		cinemaSvc: cinema.NewService(mongoDB),
		CacheSvc:  cache.NewService(cachePool),
	}
}

func marshalDtos(anyObject interface{}) ([]byte, error) {
	respBody, err := json.Marshal(anyObject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return respBody, nil
}

func internalServerErrorJson(writer http.ResponseWriter) (w http.ResponseWriter) {
	w = writer
	w.WriteHeader(http.StatusInternalServerError)

	intErr := dto.Error{Error: "err.internal_server_error"}
	respBody, err := json.Marshal(intErr)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		return
	}
	return
}

func userAlreadyUsedErrorJson(writer http.ResponseWriter) (w http.ResponseWriter) {
	w = writer
	w.WriteHeader(http.StatusConflict)

	intErr := dto.Error{Error: "err.username_already_used"}
	respBody, err := json.Marshal(intErr)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		return
	}
	return
}

func (h *Handler) ReturnPanic(w http.ResponseWriter, r *http.Request) {
	panic("I will be caught by serviceMiddleware!")
}

func (h *Handler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	log.Println(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := h.cardSvc.FindUserCardsById(r.Context(), int64(id)) //FindUserCardsById(r.Context(), int64(id))
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

func (h *Handler) GetCardTransactions(writer http.ResponseWriter, request *http.Request) {
	cardID, err := strconv.Atoi(request.URL.Query().Get("cardID"))
	if err != nil {
		log.Println(err)
		return
	}

	transactions, err := h.cardSvc.FindCardTransactionsByCardId(request.Context(), int64(cardID))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	respondBody, err := marshalDtos(transactions)
	if err != nil {
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

func (h *Handler) RegisterUser(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	err = request.Body.Close()
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	userDTO := &dto.UserRegAndLognDTO{}
	err = json.Unmarshal(body, userDTO)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	id, err := h.UserSvc.RegisterUser(request.Context(), userDTO)
	if err != nil {
		writer = userAlreadyUsedErrorJson(writer)
		return
	}

	responseBody, err := marshalDtos(id)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBody)

}

func (h *Handler) TokenGenerator(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	err = request.Body.Close()
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	userDTO := &dto.UserRegAndLognDTO{}
	err = json.Unmarshal(body, userDTO)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	token, err := h.UserSvc.AuthorizeUser(request.Context(), userDTO)
	if err != nil {
		if err == user.ErrPasswordWrong {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer = internalServerErrorJson(writer)
		return
	}

	responsBody, err := marshalDtos(dto.Token{Token: token})
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responsBody)

}

func (h *Handler) GetAllCards(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	login := &dto.UserLoginDTO{}
	err = json.Unmarshal(body, login)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	role := request.Context().Value(contextKey.AuthenticationContextKey).(user.Role)

	cards, err := h.cardSvc.All(request.Context(), login, string(role))
	if err != nil {
		return
	}

	responseBody, err := marshalDtos(cards)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBody)
}

func (h *Handler) FindAll(writer http.ResponseWriter, request *http.Request) {
	body, err := h.cinemaSvc.GetAllOrders(request.Context())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) FindByID(writer http.ResponseWriter, request *http.Request) {
	body, err := h.cinemaSvc.GetOrderByID(request.Context(), chi.URLParam(request, "id"))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Search(writer http.ResponseWriter, request *http.Request) {
	rating, err := strconv.ParseFloat(request.URL.Query().Get("min_rating"), 64)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := h.cinemaSvc.SearchByMinRating(request.Context(), rating)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Save(writer http.ResponseWriter, request *http.Request) {
	var order model.Order
	err := json.NewDecoder(request.Body).Decode(&order)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if order.ID == primitive.NilObjectID {
		err := h.cinemaSvc.Save(request.Context(), &order)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := h.cinemaSvc.Update(request.Context(), &order)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	body, err := json.Marshal(order)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) All(writer http.ResponseWriter, request *http.Request) {
	if cached, err := h.CacheSvc.FromCache(request.Context(), "films:all"); err == nil {
		log.Printf("Got from cache: %s", cached)
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(cached)
		if err != nil {
			log.Print(err)
		}
		return
	}

	body, err := h.cinemaSvc.GetAllFilmsInfo(request.Context())
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// После получения данных из основной БД и отправки клиенту, можем сохранить в кэш
	go func() {
		_ = h.CacheSvc.ToCache(context.Background(), "films:all", body)
	}()

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ByID(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	if idParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if cached, err := h.CacheSvc.FromCache(request.Context(), fmt.Sprintf("films:%s", idParam)); err == nil {
		log.Printf("Got from cache: %s", cached)
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(cached)
		if err != nil {
			log.Print(err)
		}
		return
	}

	body, err := h.cinemaSvc.GetFilmInfoById(request.Context(), idParam)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	// После получения данных из основной БД и отправки клиенту, можем сохранить в кэш
	go func() {
		_ = h.CacheSvc.ToCache(context.Background(), fmt.Sprintf("films:%s", idParam), body)
	}()

}

func (h *Handler) Upload(writer http.ResponseWriter, request *http.Request) {
	var film model.FilmInfo

	err := json.NewDecoder(request.Body).Decode(&film)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	filmId, err := h.cinemaSvc.UploadFilm(request.Context(), film)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(filmId)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//Присутствует уязвимость, что юзер может увидеть персональные предложения другого юзера
func (h *Handler) PersonalSuggestion(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	if idParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("users:%s:suggestions", idParam)

	if !(h.UserSvc.IsUserExist(request.Context(), userID)) {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if cached, err := h.CacheSvc.FromCache(request.Context(), key); err == nil {
		log.Printf("Got from cache: %s", cached)
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(cached)
		if err != nil {
			log.Print(err)
		}
		return
	}

	time.Sleep(time.Second * 3) // Симуляция долгого запроса к бд

	body, err := h.cinemaSvc.GetUserSuggestion(request.Context(), idParam)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	// После получения данных из основной БД и отправки клиенту, можем сохранить в кэш
	go func() {
		_ = h.CacheSvc.ToCache(context.Background(), key, body)
	}()
}
