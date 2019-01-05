package main

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/vulpemventures/lightning-rebalancer/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	macaroon "gopkg.in/macaroon.v2"
)

//Client type
type Client struct {
	uri string
	//client pb.LightningClient
	conn *grpc.ClientConn
}

func getClientConnection(uri string) (*grpc.ClientConn, error) {
	if !fileExists(TLSCertPath) || !fileExists(MACPath) {
		return nil, errors.New("Missing either tls certficate or macaroon")
	}

	macBytes, err := ioutil.ReadFile(MACPath)
	if err != nil {
		return nil, errors.New("Missing macaroon")
	}

	cCreds, err := credentials.NewClientTLSFromFile(TLSCertPath, "")
	if err != nil {
		return nil, errors.New("Missing tls certficate")
	}
	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cCreds),
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		return nil, err
	}

	// Now we append the macaroon credentials to the dial options.
	macCreds := macaroons.NewMacaroonCredential(mac)
	opts = append(opts, grpc.WithPerRPCCredentials(macCreds))

	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		log.Println(err)
	}
	//defer conn.Close()

	//client := pb.NewLightningClient(conn)
	return conn, nil
}
