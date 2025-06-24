package service

import (
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CategoryService 栏目服务
type CategoryService struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	categoryModel *model.CategoryModel
}

// NewCategoryService 创建栏目服务
func NewCategoryService(db *database.DB, cache cache.Cache, config *config.Config) *CategoryService {
	return &CategoryService{
		db:            db,
		cache:         cache,
		config:        config,
		categoryModel: model.NewCategoryModel(db),
	}
}

// GetChildCategories 获取子栏目
func (s *CategoryService) GetChildCategories(parentID int64) ([]*model.Category, error) {
	return s.GetSubCategories(parentID)
}

// GetCategoryPath 获取栏目路径
func (s *CategoryService) GetCategoryPath(categoryID int64) ([]*model.Category, error) {
	// 缓存键
	cacheKey := "category_path_" + string(categoryID)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if categories, ok := cached.([]*model.Category); ok {
			return categories, nil
		}
	}

	// 获取栏目
	category, err := s.categoryModel.GetByID(categoryID)
	if err != nil {
		logger.Error("获取栏目失败", "id", categoryID, "error", err)
		return nil, err
	}

	// 构建路径
	path := []*model.Category{category}

	// 如果有父栏目，递归获取
	if category.ParentID > 0 {
		parentPath, err := s.GetCategoryPath(category.ParentID)
		if err != nil {
			logger.Error("获取父栏目路径失败", "parentID", category.ParentID, "error", err)
		} else {
			// 将父路径添加到前面
			path = append(parentPath, path...)
		}
	}

	// 缓存路径
	cache.SafeSet(s.cache, cacheKey, path, time.Duration(3600)*time.Second)

	return path, nil
}

// GetTopCategories 获取顶级栏目
func (s *CategoryService) GetTopCategories() ([]*model.Category, error) {
	// 缓存键
	cacheKey := "top_categories"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if categories, ok := cached.([]*model.Category); ok {
			return categories, nil
		}
	}

	// 获取顶级栏目
	categories, err := s.categoryModel.GetTopCategories()
	if err != nil {
		logger.Error("获取顶级栏目失败", "error", err)
		return nil, err
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, categories, time.Duration(3600)*time.Second)

	return categories, nil
}

// GetCategoryByID 根据ID获取栏目
func (s *CategoryService) GetCategoryByID(id int64) (*model.Category, error) {
	// 缓存键
	cacheKey := "category_" + string(id)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if category, ok := cached.(*model.Category); ok {
			return category, nil
		}
	}

	// 获取栏目
	category, err := s.categoryModel.GetByID(id)
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		return nil, err
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, category, time.Duration(3600)*time.Second)

	return category, nil
}

// GetCategoryByDir 根据目录获取栏目
func (s *CategoryService) GetCategoryByDir(dir string) (*model.Category, error) {
	// 缓存键
	cacheKey := "category_dir_" + dir

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if category, ok := cached.(*model.Category); ok {
			return category, nil
		}
	}

	// 获取栏目
	category, err := s.categoryModel.GetByDir(dir)
	if err != nil {
		logger.Error("获取栏目失败", "dir", dir, "error", err)
		return nil, err
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, category, time.Duration(3600)*time.Second)

	return category, nil
}

// GetSubCategories 获取子栏目
func (s *CategoryService) GetSubCategories(parentID int64) ([]*model.Category, error) {
	// 缓存键
	cacheKey := "sub_categories_" + string(parentID)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if categories, ok := cached.([]*model.Category); ok {
			return categories, nil
		}
	}

	// 获取子栏目
	categories, err := s.categoryModel.GetSubCategories(parentID)
	if err != nil {
		logger.Error("获取子栏目失败", "parentID", parentID, "error", err)
		return nil, err
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, categories, time.Duration(3600)*time.Second)

	return categories, nil
}

// GetAllCategories 获取所有栏目
func (s *CategoryService) GetAllCategories() ([]*model.Category, error) {
	// 缓存键
	cacheKey := "all_categories"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if categories, ok := cached.([]*model.Category); ok {
			return categories, nil
		}
	}

	// 获取所有栏目
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		return nil, err
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, categories, time.Duration(3600)*time.Second)

	return categories, nil
}

// GetCategoryTree 获取栏目树
func (s *CategoryService) GetCategoryTree() ([]*model.Category, error) {
	// 缓存键
	cacheKey := "category_tree"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if categories, ok := cached.([]*model.Category); ok {
			return categories, nil
		}
	}

	// 获取所有栏目
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		return nil, err
	}

	// 构建栏目树
	categoryMap := make(map[int64]*model.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	// 缓存栏目
	cache.SafeSet(s.cache, cacheKey, categories, time.Duration(3600)*time.Second)

	return categories, nil
}
