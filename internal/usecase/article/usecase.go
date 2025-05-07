package article

import (
	"context"
	"io"

	"finance-backend/internal/domain"
	"finance-backend/internal/gateways/file_gateway"
	"finance-backend/internal/repository/article"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"

	"github.com/google/uuid"
)

type ArticleUseCase struct {
	repo              article.IArticleRepository
	file_gw           file_gateway.IFileGateway
	image_bucket_name string
	log               *logger.Logger
}

func NewArticleUseCase(logger *logger.Logger, repo article.IArticleRepository, file_gw file_gateway.IFileGateway, image_bucket_name string) *ArticleUseCase {
	return &ArticleUseCase{
		log:               logger,
		repo:              repo,
		file_gw:           file_gw,
		image_bucket_name: image_bucket_name,
	}
}

func (uc *ArticleUseCase) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ArticleUseCase) SearchArticlesPaginated(
	ctx context.Context,
	limit int,
	offset int,
	search *string,
	categoriesIDs []int,
) (utils.RestfullPaginatedEntities[domain.Article], error) {
	return uc.repo.SearchPaginated(ctx, limit, offset, search, categoriesIDs)
}

func (uc *ArticleUseCase) CreateArticle(
	ctx context.Context,
	header string,
	subHeader string,
	description string,
) (*domain.Article, error) {
	return uc.repo.Create(ctx, header, subHeader, description)
}

func (uc *ArticleUseCase) UpdateArticle(
	ctx context.Context,
	id int64,
	header *string,
	subHeader *string,
	description *string,
) (*domain.Article, error) {
	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, domain.ErrArticleNotFound
	}

	uc.repo.Update(ctx, id, header, nil, subHeader, description)

	return uc.repo.GetByID(ctx, id)
}

func (uc *ArticleUseCase) DeleteArticle(ctx context.Context, id int64) error {
	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return domain.ErrArticleNotFound
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *ArticleUseCase) LinkCategories(ctx context.Context, id int64, categoriesIDs []int) (*domain.Article, error) {
	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, domain.ErrArticleNotFound
	}

	uc.repo.LinkCategories(ctx, id, categoriesIDs)

	return uc.repo.GetByID(ctx, id)
}

func (uc *ArticleUseCase) LinkImage(ctx context.Context, id int64, body io.ReadSeeker, size int64, contentType string) (*domain.Article, error) {
	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, domain.ErrArticleNotFound
	}

	if article.Image != nil {
		err = uc.file_gw.DeleteObject(ctx, uc.image_bucket_name, *article.Image)
		if err != nil {
			return nil, domain.ErrS3Connection
		}
	}

	name := uuid.New().String()

	_, err = uc.file_gw.UploadObject(ctx, uc.image_bucket_name, name, body, size, contentType)

	if err != nil {
		return nil, domain.ErrS3Connection
	}

	err = uc.repo.Update(ctx, id, nil, &name, nil, nil)

	if err != nil {
		return nil, domain.ErrDBConnection
	}

	article.Image = &name

	return article, nil
}

var _ IArticleUseCase = (*ArticleUseCase)(nil)
