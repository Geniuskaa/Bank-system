package cinema

import (
	"Bank-system/cmd/bank/app/dto"
	"Bank-system/pkg/additionalservice/cinema/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

type Service struct {
	mongoDB *mongo.Database
}

func NewService(mongoDB *mongo.Database) *Service {
	return &Service{mongoDB: mongoDB}
}

func (s *Service) GetAllOrders(ctx context.Context) ([]byte, error) {
	cursor, err := s.mongoDB.Collection("orders").Find(ctx, bson.D{})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Print(cerr)
		}
	}()

	orders := make([]model.Order, 0)
	for cursor.Next(ctx) {
		var order model.Order
		err = cursor.Decode(&order)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	body, err := json.Marshal(orders)
	if err != nil {
		e := fmt.Errorf("Error when marshaling orders: %w", err)
		log.Println(e)
		return nil, e
	}

	return body, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = s.mongoDB.Collection("orders").FindOne(
		ctx,
		bson.D{{"_id", id}},
	).Decode(&order)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) SearchByMinRating(ctx context.Context, rating float64) ([]byte, error) {
	cursor, err := s.mongoDB.Collection("orders").Find(
		ctx,
		bson.D{
			{"film.rating", bson.D{
				{"$gt", rating},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Print(cerr)
		}
	}()

	orders := make([]model.Order, 0)
	for cursor.Next(ctx) {
		var order model.Order
		err = cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	body, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) Save(ctx context.Context, order *model.Order) error {
	order.Created = time.Now().Unix()
	result, err := s.mongoDB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *Service) Update(ctx context.Context, order *model.Order) error {
	result, err := s.mongoDB.Collection("orders").UpdateOne(
		ctx,
		bson.D{{"_id", order.ID}},
		bson.D{
			{"$set", bson.D{
				{"start", order.Start},
				{"price", order.Price},
			}},
		})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("Something going wrong and entity wasn`t updated", errors.New("Update failed"))
	}

	return nil
}

func (s *Service) GetAllFilmsInfo(ctx context.Context) ([]byte, error) {

	cursor, err := s.mongoDB.Collection("films").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Print(cerr)
		}
	}()

	films := make([]model.FilmInfo, 0)
	for cursor.Next(ctx) {
		var film model.FilmInfo
		err = cursor.Decode(&film)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(films)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) GetFilmInfoById(ctx context.Context, idOfFilm string) ([]byte, error) {

	id, err := primitive.ObjectIDFromHex(idOfFilm)
	if err != nil {
		return nil, err
	}

	var film model.FilmInfo

	err = s.mongoDB.Collection("films").FindOne(ctx, bson.D{{"_id", id}}).Decode(&film)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}

	body, err := json.Marshal(film)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) UploadFilm(ctx context.Context, film model.FilmInfo) (string, error) {
	result, err := s.mongoDB.Collection("films").InsertOne(ctx, film)
	if err != nil {
		return "", err
	}

	filmID := result.InsertedID.(primitive.ObjectID).String()

	return filmID, nil
}

func (s *Service) GetUserSuggestion(ctx context.Context, userID string) ([]byte, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	var suggestion dto.Suggestion

	err = s.mongoDB.Collection("suggestions").FindOne(ctx, bson.D{{"user_id", id}}).Decode(&suggestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, err
	}

	body, err := json.Marshal(suggestion)
	if err != nil {
		return nil, err
	}

	return body, nil
}
