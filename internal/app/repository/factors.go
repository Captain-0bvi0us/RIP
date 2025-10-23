package repository

import (
	"RIP/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GET /api/factors - список факторов с фильтрацией
func (r *Repository) FactorsList(title string) ([]ds.Factors, int64, error) {
	var factors []ds.Factors
	var total int64

	query := r.db.Model(&ds.Factors{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	factorsQuery := query.Order("id asc")
	if err := factorsQuery.Find(&factors).Error; err != nil {
		return nil, 0, err
	}

	if factors == nil {
		factors = []ds.Factors{}
	}

	return factors, total, nil
}

// GET /api/factors/:id - один фактор
func (r *Repository) GetFactorByID(id int) (*ds.Factors, error) {
	var factor ds.Factors
	err := r.db.First(&factor, id).Error
	if err != nil {
		return nil, err
	}
	return &factor, nil
}

// POST /api/factors - создание фактора
func (r *Repository) CreateFactor(factor *ds.Factors) error {
	return r.db.Create(factor).Error
}

// PUT /api/factors/:id - обновление фактора
func (r *Repository) UpdateFactor(id uint, req ds.FactorUpdateRequest) (*ds.Factors, error) {
	var factor ds.Factors
	if err := r.db.First(&factor, id).Error; err != nil {
		return nil, err
	}

	if req.Title != nil {
		factor.Title = *req.Title
	}
	if req.Text != nil {
		factor.Text = *req.Text
	}
	if req.Argument != nil {
		factor.Argument = req.Argument
	}

	if err := r.db.Save(&factor).Error; err != nil {
		return nil, err
	}

	return &factor, nil
}

// DELETE /api/factors/:id - удаление фактора
func (r *Repository) DeleteFactor(id uint) error {
	var factor ds.Factors
	var imageURLToDelete string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&factor, id).Error; err != nil {
			return err
		}
		if factor.Image != nil {
			imageURLToDelete = *factor.Image
		}
		if err := tx.Delete(&ds.Factors{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	if imageURLToDelete != "" {
		parsedURL, err := url.Parse(imageURLToDelete)
		if err != nil {
			log.Printf("ERROR: could not parse image URL for deletion: %v", err)
			return nil
		}

		objectName := strings.TrimPrefix(parsedURL.Path, fmt.Sprintf("/%s/", r.bucketName))

		err = r.minioClient.RemoveObject(context.Background(), r.bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("ERROR: failed to delete object '%s' from MinIO: %v", objectName, err)
		}
	}

	return nil
}

// POST /api/frax/draft/factors/:factor_id - добавление фактора в черновик
func (r *Repository) AddFactorToDraft(userID, factorID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var frax ds.FraxSearching
		err := tx.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&frax).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newFrax := ds.FraxSearching{
					CreatorID:    userID,
					Status:       ds.StatusDraft,
					CreationDate: time.Now(),
				}
				if err := tx.Create(&newFrax).Error; err != nil {
					return fmt.Errorf("failed to create draft frax: %w", err)
				}
				frax = newFrax
			} else {
				return err
			}
		}

		var count int64
		tx.Model(&ds.FactorToFrax{}).Where("frax_id = ? AND factor_id = ?", frax.ID, factorID).Count(&count)
		if count > 0 {
			return errors.New("factor already in frax")
		}

		link := ds.FactorToFrax{
			FraxID:   frax.ID,
			FactorID: factorID,
		}

		if err := tx.Create(&link).Error; err != nil {
			return fmt.Errorf("failed to add factor to frax: %w", err)
		}

		if err := tx.Model(&ds.Factors{}).Where("id = ?", factorID).Update("status", true).Error; err != nil {
			return fmt.Errorf("failed to update factor status: %w", err)
		}
		return nil
	})
}

// POST /api/factors/:id/image - загрузка изображения фактора
func (r *Repository) UploadFactorImage(factorID uint, fileHeader *multipart.FileHeader) (string, error) {
	var finalImageURL string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var factor ds.Factors
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&factor, factorID).Error; err != nil {
			return fmt.Errorf("factor with id %d not found: %w", factorID, err)
		}

		const imagePathPrefix = "Images/"

		if factor.Image != nil && *factor.Image != "" {
			oldImageURL, err := url.Parse(*factor.Image)
			if err == nil {
				oldObjectName := strings.TrimPrefix(oldImageURL.Path, fmt.Sprintf("/%s/", r.bucketName))
				r.minioClient.RemoveObject(context.Background(), r.bucketName, oldObjectName, minio.RemoveObjectOptions{})
			}
		}

		fileName := filepath.Base(fileHeader.Filename)
		objectName := imagePathPrefix + fileName

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = r.minioClient.PutObject(context.Background(), r.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})

		if err != nil {
			return fmt.Errorf("failed to upload to minio: %w", err)
		}

		imageURL := fmt.Sprintf("http://%s/%s/%s", r.minioEndpoint, r.bucketName, objectName)

		if err := tx.Model(&factor).Update("image", imageURL).Error; err != nil {
			return fmt.Errorf("failed to update factor image url in db: %w", err)
		}

		finalImageURL = imageURL
		return nil
	})
	if err != nil {
		return "", err
	}
	return finalImageURL, nil
}
