package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_sers/goods_srv/global"
	"mxshop_sers/goods_srv/model"
	pb "mxshop_sers/goods_srv/proto"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*pb.CategoryListResponse, error) {
	var categorys []model.Category
	// SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	// 配置指明了外键后，可以使用Preload预加载，来把品牌的子分类也取出来
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//var tot int64
	//global.DB.Model(&model.Category{}).Where(&model.Category{Level: 1}).Count(&tot)
	tot := len(categorys)
	b, _ := json.Marshal(&categorys)
	data := []*pb.CategoryInfoResponse{}
	for _, category := range categorys {
		data = append(data, &pb.CategoryInfoResponse{
			Name:           category.Name,
			Id:             category.ID,
			ParentCategory: category.ParentCategoryID,
			Level:          category.Level,
			IsTab:          category.IsTab,
		})
	}
	return &pb.CategoryListResponse{
		JsonData: string(b),
		Total:    int32(tot),
		Data:     data,
	}, nil
}

// 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *pb.CategoryListRequest) (*pb.SubCategoryListResponse, error) {
	categoryListResponse := &pb.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &pb.CategoryInfoResponse{
		Name:           category.Name,
		Id:             category.ID,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*pb.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &pb.CategoryInfoResponse{
			Name:           subCategory.Name,
			Id:             subCategory.ID,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})
	}
	categoryListResponse.Total = int32(len(subCategorys))
	categoryListResponse.SubCategorys = subCategoryResponse
	return categoryListResponse, nil
}

// 新增商品分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.CategoryInfoResponse, error) {
	category := &model.Category{
		Level: req.Level,
		Name:  req.Name,
		IsTab: req.IsTab,
	}
	if req.Level != 1 {
		//可以添加一个父类是否存在
		category.ParentCategoryID = req.ParentCategory
	}
	global.DB.Create(category)
	return &pb.CategoryInfoResponse{
		Id: category.ID,
	}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	global.DB.Save(&category)
	return &emptypb.Empty{}, nil
}
