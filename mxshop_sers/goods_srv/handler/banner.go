package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_sers/goods_srv/global"
	"mxshop_sers/goods_srv/model"
	"mxshop_sers/goods_srv/proto"
)

func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (*proto.BannerListResponse, error) {
	var banner []model.Banner
	result := global.DB.Find(&banner)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Banner{}).Count(&total)

	var bannerResponse []*proto.BannerResponse
	for _, banner := range banner {
		bannerResponse = append(bannerResponse, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}
	bannerListResponse := proto.BannerListResponse{
		Total: int32(total),
		Data:  bannerResponse,
	}
	return &bannerListResponse, nil
}
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{}
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	global.DB.Save(&banner)
	return &proto.BannerResponse{Id: banner.ID}, nil
}
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner

	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	global.DB.Save(&banner)
	return &emptypb.Empty{}, nil
}
