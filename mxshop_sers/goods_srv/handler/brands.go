package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_sers/goods_srv/global"
	"mxshop_sers/goods_srv/model"
	pb "mxshop_sers/goods_srv/proto"
)

// 品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, req *pb.BrandFilterRequest) (*pb.BrandListResponse, error) {
	brandListResponse := pb.BrandListResponse{}

	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandResponse []*pb.BrandInfoResponse
	for _, brand := range brands {
		brandResponse = append(brandResponse, &pb.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponse
	return &brandListResponse, nil
}

// 新建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.BrandInfoResponse, error) {
	if result := global.DB.Where("name = ?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)

	return &pb.BrandInfoResponse{Id: brand.ID}, nil
}

// 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新品牌
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *pb.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.Where("name = ?", req.Name).First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	global.DB.Save(&brands)
	return &emptypb.Empty{}, nil
}
