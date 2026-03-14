package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danindudesilva/payments-service/internal/config"
	"github.com/danindudesilva/payments-service/internal/httpserver"
	"github.com/danindudesilva/payments-service/internal/payments/domain"
	memoryrepo "github.com/danindudesilva/payments-service/internal/payments/repository/memory"
	paymentservice "github.com/danindudesilva/payments-service/internal/payments/service"
	paymenthttp "github.com/danindudesilva/payments-service/internal/payments/transport/http"
)

type App struct {
	cfg    config.Config
	server *http.Server
	logger *slog.Logger
}

func New(cfg config.Config) *App {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := memoryrepo.NewRepository()
	gateway := newNoopGateway()
	service := paymentservice.New(
		repo,
		gateway,
		time.Now,
		func() string {
			return fmt.Sprintf("attempt_%d", time.Now().UnixNano())
		},
	)

	router := httpserver.NewRouter(cfg, logger)
	mux := http.NewServeMux()
	mux.Handle("/", router)

	handler := paymenthttp.NewHandler(service, logger)
	handler.Register(mux)

	server := &http.Server{
		Addr:              cfg.HTTPAddress(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: server,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("http server starting",
			slog.String("addr", a.server.Addr),
			slog.String("env", a.cfg.AppEnv),
		)

		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen and serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	a.logger.Info("http server stopped")
	return nil
}

type noopGateway struct{}

func newNoopGateway() *noopGateway {
	return &noopGateway{}
}

func (g *noopGateway) CreatePayment(ctx context.Context, request domain.CreateProviderPaymentRequest) (domain.CreateProviderPaymentResult, error) {
	return domain.CreateProviderPaymentResult{
		ProviderName:      "fake",
		ProviderPaymentID: "fake_payment_id",
		ClientSecret:      "fake_client_secret",
		Status:            domain.PaymentStatusPending,
	}, nil
}

func (g *noopGateway) GetPayment(ctx context.Context, providerPaymentID string) (domain.CreateProviderPaymentResult, error) {
	return domain.CreateProviderPaymentResult{
		ProviderName:      "fake",
		ProviderPaymentID: providerPaymentID,
		Status:            domain.PaymentStatusPending,
	}, nil
}
