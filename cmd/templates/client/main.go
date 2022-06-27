package main

import (
	phoneBalanceV1Pb "Bank-system/pkg/topUpPhoneBalancePattern/v1"
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := phoneBalanceV1Pb.NewPhoneBalancePatternServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	{
		response, err := client.CreatePattern(ctx, &phoneBalanceV1Pb.Pattern{
			Title:       "Мамин телефон",
			PhoneNumber: "+79644272739",
			Created:     &timestamp.Timestamp{Seconds: 1656338658},
			Updated:     nil,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Print(st.Code())
				log.Print(st.Message())
			}
			return err
		}

		log.Print(response)
	}

	{
		response, err := client.GetAllPatterns(ctx, &phoneBalanceV1Pb.EmptyRequest{})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Print(st.Code())
				log.Print(st.Message())
			}
			return err
		}

		log.Print(response)
	}

	{
		response, err := client.GetPatternById(ctx, &phoneBalanceV1Pb.PatternId{Id: 1})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Print(st.Code())
				log.Print(st.Message())
			}
			return err
		}

		log.Print(response)
	}

	{
		response, err := client.EditPatterById(ctx, &phoneBalanceV1Pb.Pattern{
			Id:          1,
			Title:       "Мамин телефон",
			PhoneNumber: "+79644272738",
			Created:     &timestamp.Timestamp{Seconds: 1656338658},
			Updated:     nil,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Print(st.Code())
				log.Print(st.Message())
			}
			return err
		}

		log.Print(response)
	}

	{
		response, err := client.DeleteById(ctx, &phoneBalanceV1Pb.PatternId{Id: 1})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Print(st.Code())
				log.Print(st.Message())
			}
			return err
		}

		log.Print(response)
	}

	return nil
}
