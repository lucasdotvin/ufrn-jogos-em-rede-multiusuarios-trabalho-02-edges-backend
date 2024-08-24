package main

import (
	"net/http"
	"trabalho-02-edges/config"
	roomCase "trabalho-02-edges/lib/application/case/room"
	userCase "trabalho-02-edges/lib/application/case/user"
	"trabalho-02-edges/lib/domain/service"
	"trabalho-02-edges/lib/infrastructure/repository"
	"trabalho-02-edges/lib/infrastructure/repository/database/sqlite"
	"trabalho-02-edges/lib/presentation/auth/token/jwt"
	"trabalho-02-edges/lib/presentation/broadcast"
	"trabalho-02-edges/lib/presentation/broadcast/drivers/websocket"
	"trabalho-02-edges/lib/presentation/http/controller"
	"trabalho-02-edges/lib/presentation/middleware"
	"trabalho-02-edges/lib/presentation/middleware/auth"
)

func main() {
	cfg := config.GetConfig()

	sqliteDb, err := sqlite.NewDatabase(cfg)
	defer func() {
		_ = sqliteDb.Close()
	}()

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	bcryptHashRepository := repository.NewBcryptHashRepository(cfg)

	jwtTokenService := jwt.NewJwtService(cfg)

	websocketBroadcaster := websocket.NewWebSocketDriver(cfg)

	roomBroadcast := broadcast.NewRoomBroadcast(websocketBroadcaster)

	sqliteRoomRepository := sqlite.NewSqliteRoomRepository(sqliteDb)
	sqliteRoomUserRepository := sqlite.NewSqliteRoomUserRepository(sqliteDb)
	sqliteUserRepository := sqlite.NewSqliteUserRepository(sqliteDb)

	roomService := service.NewRoomService(roomBroadcast, sqliteRoomRepository, sqliteRoomUserRepository, sqliteUserRepository)
	userService := service.NewUserService(bcryptHashRepository, sqliteUserRepository)

	signUpUseCase := userCase.NewSignUpUseCase(userService)
	signInUseCase := userCase.NewSignInUseCase(userService)
	findUserUseCase := userCase.NewFindUserUseCase(userService)
	validateUserIsPresentUseCase := roomCase.NewValidateUserIsPresent(roomService)
	handleUserDisconnectUseCase := roomCase.NewHandleUserDisconnectUseCase(roomService)
	listRoomUseCase := roomCase.NewListRoomUseCase(roomService, userService)
	findUserActiveRoomUseCase := roomCase.NewFindUserActiveRoomUseCase(roomService, userService)
	createRoomUseCase := roomCase.NewCreateRoomUseCase(cfg, roomService)
	ingressUserUseCase := roomCase.NewIngressUserUseCase(roomService, userService)

	authController := controller.NewAuthController(signUpUseCase, signInUseCase, findUserUseCase, jwtTokenService)
	broadcastController := controller.NewBroadcastController(websocketBroadcaster, validateUserIsPresentUseCase, handleUserDisconnectUseCase)
	roomController := controller.NewRoomController(listRoomUseCase, findUserActiveRoomUseCase, createRoomUseCase, ingressUserUseCase)

	cookieTokenAuthMiddleware := auth.NewCookieTokenAuthMiddleware(jwtTokenService)
	webSocketQueryTokenAuthMiddleware := auth.NewWebSocketQueryTokenAuthMiddleware(jwtTokenService)

	router := http.NewServeMux()

	apiRouter := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api", apiRouter))

	v1Router := http.NewServeMux()
	apiRouter.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	authRouter := http.NewServeMux()
	v1Router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	authRouter.HandleFunc("POST /sign-up", authController.SignUp)
	authRouter.HandleFunc("POST /sign-in", authController.SignIn)
	authRouter.HandleFunc("POST /refresh", authController.Refresh)
	authRouter.Handle("GET /me", middleware.ApplyToFunc(authController.Me, cookieTokenAuthMiddleware.Handle))

	roomRouter := http.NewServeMux()
	v1Router.Handle("/rooms/", http.StripPrefix("/rooms", middleware.Apply(roomRouter, cookieTokenAuthMiddleware.Handle)))
	roomRouter.HandleFunc("GET /", roomController.Index)
	roomRouter.HandleFunc("GET /my-active-room", roomController.FindMyActiveRoom)
	roomRouter.HandleFunc("POST /", roomController.Store)
	roomRouter.HandleFunc("POST /{room}/ingress/", roomController.IngressUser)

	wsRouter := http.NewServeMux()
	v1Router.Handle("/ws/", http.StripPrefix("/ws", middleware.Apply(wsRouter, webSocketQueryTokenAuthMiddleware.Handle)))
	wsRouter.HandleFunc("/rooms/events", broadcastController.SubscribeForGlobalRoomEvents)
	wsRouter.HandleFunc("/rooms/{room}/events", broadcastController.SubscribeForRoomEvents)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	err = server.ListenAndServe()

	if err != nil {
		panic("failed to start server: " + err.Error())
	}
}
