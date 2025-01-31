package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *economyServiceServer) GetShop(ctx context.Context, req *v1.GetShopRequest) (*v1.GetShopResponse, error) {
	fmt.Println("GetShop")

	shop, err := s.shopRepository.Get(ctx, req.GetShopId())

	if err != nil {
		return nil, status.Error(codes.NotFound, "shop not found")
	}

	return &v1.GetShopResponse{
		Shop: shop,
	}, nil
}

func (s *economyServiceServer) CreateShop(ctx context.Context, req *v1.CreateShopRequest) (*v1.CreateShopResponse, error) {
	fmt.Println("CreateShop")

	shop, err := s.shopRepository.Create(
		ctx,
		req.GetName(),
		req.GetMetadata(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create shop")
	}

	return &v1.CreateShopResponse{
		Shop: shop,
	}, nil
}

func (s *economyServiceServer) UpdateShop(ctx context.Context, req *v1.UpdateShopRequest) (*v1.UpdateShopResponse, error) {
	fmt.Println("UpdateShop")

	shop, err := s.shopRepository.Update(
		ctx,
		req.GetShopId(),
		req.GetName(),
		req.GetMetadata(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update shop")
	}

	return &v1.UpdateShopResponse{
		Shop: shop,
	}, nil
}

func (s *economyServiceServer) ListShop(ctx context.Context, req *v1.ListShopRequest) (*v1.ListShopResponse, error) {
	fmt.Println("ListShop")

	// Parse the page token
	var parsedToken int64
	parsedToken, _ = strconv.ParseInt(req.GetPageToken(), 10, 32)

	// Get the limit
	limit := req.GetPageSize()
	if limit == 0 {
		limit = 100
	}

	// Get the offset
	offset := int32(0)
	if len(req.GetPageToken()) > 0 {
		offset = int32(parsedToken) * limit
	}

	// Get the shops
	shops, totalSize, err := s.shopRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve shop list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListShopResponse{
		Shops:         shops,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *economyServiceServer) AttachProduct(ctx context.Context, req *v1.AttachProductRequest) (*v1.AttachProductResponse, error) {
	fmt.Println("AttachProduct")

	// Add product to the database return the generated UUID
	shop, err := s.shopRepository.AttachProduct(
		ctx,
		req.GetShopId(),
		req.GetProductId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to attach product to shop")
	}

	return &v1.AttachProductResponse{
		Shop: shop,
	}, nil
}

func (s *economyServiceServer) DetachProduct(ctx context.Context, req *v1.DetachProductRequest) (*v1.DetachProductResponse, error) {
	fmt.Println("DetachProduct")

	// Add product to the databased return the generated UUID
	shop, err := s.shopRepository.DetachProduct(
		ctx,
		req.GetShopProductId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to detach product from shop")
	}

	return &v1.DetachProductResponse{
		Shop: shop,
	}, nil
}
