package forms

type CategoryForm struct {
	Name           string `form:"name" json:"name" binding:"required,min=3,max=20"`
	ParentCategory int32  `form:"parent" json:"parent"`
	Level          int32  `form:"level" json:"level" binding:"required,oneof=1 2 3"` //限定在123中选择
	IsTab          *bool  `form:"is_tab" json:"is_tab" binding:"required"`
}

type UpdateCategoryForm struct {
	Name  string `form:"name" json:"name" binding:"required,min=3,max=20"`
	IsTab *bool  `form:"is_tab" json:"is_tab"`
}
