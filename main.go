package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jawher/mow.cli"
)

const appDescription = "UPP - event reader, reads from DynamoDB"

func main() {
	app := cli.App("transaction_reader", appDescription)

	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "transaction_reader",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})

	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "transaction_reader",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})

	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})

	log.SetLevel(log.InfoLevel)
	log.Infof("[Startup] transaction_reader is starting ")

	app.Action = func() {
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		go func() {
			serveEndpoints(*appSystemCode, *appName, *port)
		}()

		waitForSignal()
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func serveEndpoints(appSystemCode string, appName string, port string) {
	router := http.NewServeMux()

	//admin endpoints
	healthService := newHealthService(&healthConfig{appSystemCode: appSystemCode, appName: appName, port: port})

	hc := health.HealthCheck{SystemCode: appSystemCode, Name: appName, Description: appDescription, Checks: healthService.checks}
	router.HandleFunc(healthPath, health.Handler(hc))
	router.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(healthService.gtgCheck))
	router.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)

	//functional endpoints
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000")})
	if err != nil {
		log.Fatalf("Could not create AWS session %v", err)
	}
	hh := httpHandlers{transactionReaderService{dynamodb.New(sess)}}

	router.HandleFunc("/transactions/inprogress", hh.getInProgressTransactionsForType)
	router.HandleFunc("/transactions/byuuid", hh.getTransactionsForUUID)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Unable to start: %v", err)
	}
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
