package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kajjagtenberg/eventflowdb/cluster"
	"github.com/kajjagtenberg/eventflowdb/env"
	"github.com/kajjagtenberg/eventflowdb/graph/generated"
	"github.com/kajjagtenberg/eventflowdb/graph/resolvers"
	"github.com/kajjagtenberg/eventflowdb/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

func main() {
	/////////////
	//  Hello  //
	/////////////

	log.Println("Hello, world!")

	//////////////
	//  Config  //
	//////////////

	godotenv.Load()

	grpcAddr := env.GetEnv("GRPC_LISTENER", ":6543")
	httpAddr := env.GetEnv("HTTP_LISTENER", ":16543")
	eventsFile := env.GetEnv("EVENTS_FILE", "events.db")

	///////////////
	//  Storage  //
	///////////////

	log.Printf("Initializing Storage service at: %v", eventsFile)

	db, err := bbolt.Open(eventsFile, 0777, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}
	defer db.Close()

	storage, err := store.NewStorage(db)
	if err != nil {
		log.Fatalf("Failed to initialize Storage service: %v", err)
	}

	///////////////
	//  Cluster  //
	///////////////

	log.Println("Setting up a cluster")

	cl, err := cluster.NewCluster()
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}
	defer cl.Leave()

	if err := cl.Join(); err != nil {
		log.Fatalf("Failed to join cluster: %v", err)
	}

	////////////
	//  gRPC  //
	////////////

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s", grpcAddr)
	}
	defer lis.Close()

	grpcSrv := grpc.NewServer()

	log.Println("Initializing gRPC services")

	store.RegisterStreamsServer(grpcSrv, store.NewEventStoreService(storage))
	cluster.RegisterClusterServiceServer(grpcSrv, cluster.NewClusterService(cl))

	go func() {
		log.Printf("Starting gRPC server on %s", grpcAddr)

		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	////////////
	//  HTTP  //
	////////////

	log.Println("Initializing Prometheus metrics")

	httpSrv := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ReadTimeout:           time.Second * 10,
		WriteTimeout:          time.Second * 10,
		IdleTimeout:           time.Second * 10,
	})
	httpSrv.Use(cors.New())
	httpSrv.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	log.Println("Initializing GraphQL")

	httpSrv.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/")))
	httpSrv.Post("/", adaptor.HTTPHandler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		Cluster: cl,
		Storage: storage,
		Start:   time.Now(),
	}}))))
	httpSrv.Get("/backup", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/octet-stream")
		c.Set("Content-Disposition", "attachment;filename=backup.db")

		return storage.Backup(c.Response().BodyWriter())
	})

	go func() {
		log.Printf("Starting HTTP server on %s", httpAddr)

		httpSrv.Listen(httpAddr)
	}()

	/*go func() {
		types := []string{"ProductAdded", "UserRegistered", "CashDeposited", "CashWithdrawn", "UserLocked", "UserDeleted", "ProductDiscontinued"}

		faker := faker.New(rand.Int63())

		for {
			stream := uuid.New()

			event := &store.AddRequest_Event{
				Type: types[rand.Intn(len(types))],
			}

			var data interface{}

			switch event.Type {
			case "ProductAdded":
				data = struct {
					Name  string `fake:"{carModel}" json:"name"`
					Brand string `fake:"{company}" json:"brand"`
					Price int    `fake:"{price:100,10000}" json:"price"`
				}{}
			case "UserRegistered":
				data = struct {
					Firstname string `fake:"{firstname}" json:"firstname"`
					Lastname  string `fake:"{lastname}" json:"lastname"`
					Email     string `fake:"{email}" json:"email"`
					Phone     string `fake:"{phone}" json:"phone"`
				}{}
			case "CashDeposited":
				data = struct {
					Amount int `fake:"{price:100,10000}" json:"amount"`
				}{}
			case "CashWithdrawn":
				data = struct {
					Amount int `fake:"{price:100,10000}" json:"amount"`
				}{}
			case "UserLocked":
				data = struct{}{}
			case "UserDeleted":
				data = struct{}{}
			case "ProductDiscontinued":
				data = struct{}{}
			default:
				data = struct{}{}
			}

			faker.Struct(&data)
			raw, err := json.Marshal(&data)
			if err != nil {
				log.Fatalf("Failed to marshal: %v", err)
			}
			event.Data = raw

			go func() {
				if _, err := storage.Add(&store.AddRequest{
					Stream:  stream[:],
					Version: 0,
					Events: []*store.AddRequest_Event{
						event,
					},
				}); err != nil {
					log.Fatalf("Failed to add event: %v", err)
				}
			}()

			time.Sleep(time.Millisecond * 1)
		}
	}()*/

	////////////////
	//  Shutdown  //
	////////////////

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-sig

	log.Println("Stopping all services")

	grpcSrv.GracefulStop()
	// httpSrv.Shutdown()
	db.Close()

	log.Println("Stopped Storage")

	log.Println("Stopped all services")
}
