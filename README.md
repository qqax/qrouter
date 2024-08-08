#qrouter
### Set server
```golang
func setServer(ctx context.Context) *http.Server {
    mux := http.NewServeMux()

    locales := qrouter.InitApi(mux, "locales").Group("add")
    addEn := locales.Group("en-US")
    addEn.Handle(
        qrouter.PathHandler{Path: "translations", POST: addEnTransl},
    )

    api := qrouter.InitApi(mux, "api")

    office := api.Group("office")

    office.Group("currencies").
        Handle(
            qrouter.PathHandler{Path: "/", GET: getCurrencies, PATCH: updateCurrency},
        )

    wrappedMux := qrouter.RequestLogger(mux)

    server := &http.Server{
        Handler:        wrappedMux,
        Addr:           host + ":" + port,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    //// Initializing the server in a goroutine so that
    //// it won't block the graceful shutdown handling below
    go func() {
        // service connections
        if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatal().Timestamp().Err(err).Msg("listen: %s\n" + err.Error())
        }
    }()

    return server
}
```

### Start & stop server
```golang
func Start(cfg *conf.Config) {
	ctx := context.Background()
	
	qrouter.SetRecaptchaKey(cfg.RecaptchaSecretKey)
	
	server := setServer(ctx)

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Timestamp().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	
	defer cancel()
	
	if err = server.Shutdown(ctx); err != nil {
		log.Fatal().Timestamp().Err(err).Msg("Server Shutdown:" + err.Error())
	}
	
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Info().Timestamp().Msg("timeout of 5 seconds.")
	}
	
	log.Info().Timestamp().Msg("Server exiting")
}
```