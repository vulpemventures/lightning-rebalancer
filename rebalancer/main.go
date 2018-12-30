package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vulpemventures/lightning-rebalancer/macaroons"
	macaroon "gopkg.in/macaroon.v2"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/gorilla/mux"
	pb "github.com/lightningnetwork/lnd/lnrpc"
	//"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var HTTPPort = "3000"
var LNDPort = "10009"
var LNDHost = "localhost"
var TLSCertPath = "/home/bitcoin/.lnd/tls.cert"
var MACPath = "/home/bitcoin/.lnd/data/chain/bitcoin/testnet/admin.macaroon"

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

//Deposit Handler
func Deposit(w http.ResponseWriter, r *http.Request) {
	//body := getRequestBody(r.Body)
	//defaultNet := getNetworkTypeByString(body["network"])
	//Doing my things
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Not ready for this network \n"))
}

//Withdraw Handler
func Withdraw(w http.ResponseWriter, r *http.Request) {
	body := getRequestBody(r.Body)
	defaultNet := getNetworkTypeByString(body["network"])
	addrString := body["address"]
	amtSatoshi, err := strconv.ParseInt(body["amount"], 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading the amount\n"))
		return
	}

	_, err = btcutil.DecodeAddress(addrString, defaultNet)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding address\n"))
		return

	}

	conn, err := getClientConnection(LNDHost + ":" + LNDPort)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewLightningClient(conn)
	defer conn.Close()

	resp, err := client.AddInvoice(context.Background(), &pb.Invoice{Value: amtSatoshi, Memo: addrString})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating LN invoice\n"))
		return
	}

	//Doing my things
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.String()))
}

//Ping func
func Ping(w http.ResponseWriter, r *http.Request) {
	conn, err := getClientConnection(LNDHost + ":" + LNDPort)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewLightningClient(conn)
	defer conn.Close()

	resp, err := client.GetInfo(context.Background(), &pb.GetInfoRequest{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error GetInfo\n"))
		return
	}

	//Doing my things
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.String()))

}

func getNetworkTypeByString(key string) *chaincfg.Params {

	if key == "mainnet" {
		return &chaincfg.MainNetParams
	}

	return &chaincfg.TestNet3Params
}

func getRequestBody(body io.ReadCloser) map[string]string {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]string
	decoder.Decode(&decodedBody)

	return decodedBody
}

// fileExists reports whether the named file or directory exists.
// This function is taken from https://github.com/btcsuite/btcd
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {

	if _, present := os.LookupEnv("HTTP_PORT"); present {
		HTTPPort = os.Getenv("HTTP_PORT")
	}

	if _, present := os.LookupEnv("LND_PORT"); present {
		LNDPort = os.Getenv("LND_PORT")
	}

	if _, present := os.LookupEnv("LND_HOST"); present {
		LNDHost = os.Getenv("LND_HOST")
	}

	if _, present := os.LookupEnv("TLS_PATH"); present {
		TLSCertPath = os.Getenv("TLS_PATH")
	}

	if _, present := os.LookupEnv("MAC_PATH"); present {
		MACPath = os.Getenv("MAC_PATH")
	}

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/ping", Ping).Methods("GET")
	r.HandleFunc("/deposit", Deposit).Methods("POST")
	r.HandleFunc("/withdraw", Withdraw).Methods("POST")

	// Bind to a port and pass our router in
	log.Println("running rebalancer server on localhost:" + HTTPPort)
	log.Println("looking for LND gRPC server on " + LNDHost + ":" + LNDPort)
	log.Fatal(http.ListenAndServe(":"+HTTPPort, r))
}
