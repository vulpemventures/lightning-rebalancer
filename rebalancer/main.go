package main

import (
	"context"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vulpemventures/lightning-rebalancer/rebalancer/util"

	"github.com/btcsuite/btcutil"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	pb "github.com/lightningnetwork/lnd/lnrpc"
)

var config *util.Config
var lnURI string

// Deposit Handler
func Deposit(w http.ResponseWriter, r *http.Request) {
	body := util.GetRequestBody(r.Body)
	defaultNet := util.GetNetworkTypeByString(body["network"])
	log.Println(defaultNet)
	//invoice := body["invoice"]

	conn, err := util.GetClientConnection(lnURI)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewLightningClient(conn)
	defer conn.Close()

	resp, err := client.NewAddress(context.Background(), &pb.NewAddressRequest{Type: 0})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating deposit address \n"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.Address))
}

// Withdraw Handler
func Withdraw(w http.ResponseWriter, r *http.Request) {
	body := util.GetRequestBody(r.Body)
	defaultNet := util.GetNetworkTypeByString(body["network"])
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

	conn, err := util.GetClientConnection(lnURI)
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
	conn, err := util.GetClientConnection(lnURI)
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

func main() {
	config = util.GetConfigFromEnv()
	lnURI = config.LNDHost + ":" + config.LNDPort
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/ping", Ping).Methods("GET")
	r.HandleFunc("/deposit", Deposit).Methods("POST")
	r.HandleFunc("/withdraw", Withdraw).Methods("POST")

	//CORS
	router := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(r)

	// Bind to a port and pass our router in
	log.Println("running rebalancer server on localhost:" + config.HTTPPort)
	log.Println("looking for LND gRPC server on " + lnURI)
	log.Fatal(http.ListenAndServe(":"+config.HTTPPort, router))
}
