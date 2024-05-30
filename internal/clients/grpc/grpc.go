package grpc

import (
	"context"
	ssov1 "diplom/gen"
	"fyne.io/fyne/v2"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/errgo.v2/fmt/errors"
	"strings"
	"time"
)

type Client struct {
	authApi       ssov1.AuthClient
	userApi       ssov1.UserClient
	machineDepApi ssov1.MachineDepartmentClient
	radarDepApi   ssov1.RadarClient
	App           fyne.App
	UserName      string
	Token         string
	Password      string
}

func New(ctx context.Context, addr string, timeout time.Duration, retriesCount int, app fyne.App) (*Client, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpc_retry.WithMax(uint(retriesCount)),
		grpc_retry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, errors.Newf("failed to create new gRPC client: %w", err)
	}

	return &Client{
		authApi:       ssov1.NewAuthClient(cc),
		userApi:       ssov1.NewUserClient(cc),
		machineDepApi: ssov1.NewMachineDepartmentClient(cc),
		radarDepApi:   ssov1.NewRadarClient(cc),
		App:           app,
	}, nil
}

func (c *Client) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	resp, err := c.authApi.Register(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to register user: %w", err)
	}

	return resp, nil
}

func (c *Client) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	resp, err := c.authApi.Login(ctx, request)
	if err != nil {
		if isSpecificError(err) {
			return nil, errors.New("Не удалось подключится к серверу")
		}

		return nil, errors.New("Невернные данные")
	}

	return resp, nil
}

func (c *Client) ChangeEmail(ctx context.Context, request *ssov1.Email) (*ssov1.GetUserResponse, error) {
	resp, err := c.userApi.ChangeEmail(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ChangePassword(ctx context.Context, request *ssov1.NewPassword) error {
	_, err := c.userApi.ChangePassword(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func isSpecificError(err error) bool {
	specificError := "connection error: desc = \"transport: Error while dialing: dial tcp [::1]:44044: connectex: No connection could be made because the target machine actively refused it.\""
	return strings.Contains(err.Error(), specificError)
}

func (c *Client) GetUserData(ctx context.Context, request *ssov1.GetUserRequest) (*ssov1.GetUserResponse, error) {
	resp, err := c.userApi.GetUser(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to get user data: %w", err)
	}

	return resp, nil
}

func (c *Client) GetAllUsers(ctx context.Context, request *ssov1.UserName) (*ssov1.GetUsersDataResponse, error) {
	resp, err := c.userApi.GetUsersData(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to get users data: %w", err)
	}

	return resp, nil
}

func (c *Client) ChangeProfilePhoto(ctx context.Context, request *ssov1.ChangeProfilePhotoRequest) (*ssov1.ChangeProfilePhotoResponse, error) {
	resp, err := c.userApi.ChangeProfilePhoto(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to change profile photo: %w", err)
	}

	return resp, nil
}

func (c *Client) GetMachineInfo(ctx context.Context, request *ssov1.UserName) (ssov1.MachineDepartment_GetInfoMachineDepClient, error) {
	resp, err := c.machineDepApi.GetInfoMachineDep(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to get machine info: %w", err)
	}

	return resp, nil
}

func (c *Client) GetRadarObjects(ctx context.Context, request *ssov1.UserName) (ssov1.Radar_GetRadarInfoClient, error) {
	resp, err := c.radarDepApi.GetRadarInfo(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to get radar info: %w", err)
	}

	return resp, nil
}

func (c *Client) ChangeShipParameters(ctx context.Context, request *ssov1.UpdateShipParameters) (*ssov1.Empty, error) {
	_, err := c.radarDepApi.ChangeShipParameters(ctx, request)
	if err != nil {
		return nil, errors.Newf("failed to update ship parameters: %w", err)
	}

	return nil, nil
}
