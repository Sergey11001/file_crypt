package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"univer/internal/dto"
	"univer/internal/pgrepository"
	"univer/pkg/lib/errs"
)

type FilesService struct {
	pgClient PgClient
	s3Client S3Client
}

func NewFilesService(
	pgClient PgClient,
	s3Client S3Client,
) *FilesService {
	return &FilesService{
		pgClient: pgClient,
		s3Client: s3Client,
	}
}

func (s *FilesService) Files(ctx context.Context, userUUID uuid.UUID) ([]dto.File, error) {
	rows, err := pgrepository.New(s.pgClient).Files(ctx, pgrepository.FilesParams{
		UserUuid: userUUID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.File, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.File{
			Uuid:         row.File.Uuid,
			Name:         row.File.Name,
			Size:         row.File.Size,
			CreatedAt:    row.File.CreatedAt,
			SymmetricKey: row.SymmetricKey,
		})
	}

	return result, nil
}

func (s *FilesService) File(ctx context.Context, userUUID, fileUUID uuid.UUID) (string, []byte, error) {
	row, err := pgrepository.New(s.pgClient).File(ctx, pgrepository.FileParams{
		UserUuid: userUUID,
		Uuid:     fileUUID,
	})
	if err != nil {
		if pgrepository.IsNoRows(err) {
			return "", nil, errs.NotFound.New("File not found", "file not found")
		}

		return "", nil, err
	}

	fileData, err := s.s3Client.Download(ctx, "filecrypto", row.File.Uuid.String())
	if err != nil {
		return "", nil, err
	}

	return row.File.Name, fileData, nil
}

func (s *FilesService) CommonFile(ctx context.Context, fileUUID uuid.UUID) (string, []byte, error) {
	row, err := pgrepository.New(s.pgClient).CommonFile(ctx, pgrepository.CommonFileParams{
		Uuid: fileUUID,
	})
	if err != nil {
		if pgrepository.IsNoRows(err) {
			return "", nil, errs.NotFound.New("File not found", "file not found")
		}

		return "", nil, err
	}

	fileData, err := s.s3Client.Download(ctx, "filecrypto", row.File.Uuid.String())
	if err != nil {
		return "", nil, err
	}

	return row.File.Name, fileData, nil
}

func (s *FilesService) CreateFile(ctx context.Context, userUUID uuid.UUID, input dto.CreateFileInput) (dto.File, error) {
	var zero dto.File

	fileUUID := uuid.New()

	pg := pgrepository.New(s.pgClient)

	file, err := pg.CreateFile(ctx, pgrepository.CreateFileParams{
		Uuid:     fileUUID,
		UserUuid: userUUID,
		Name:     input.Name,
		Size:     input.File.FileSize(),
		IsCrypt:  input.SymmetricKey != nil,
	})
	if err != nil {
		return zero, err
	}

	if input.SymmetricKey != nil {
		_, err := pg.UpsertFileCryptoKey(ctx, pgrepository.UpsertFileCryptoKeyParams{
			FileUuid:     file.Uuid,
			UserUuid:     userUUID,
			SymmetricKey: *input.SymmetricKey,
		})
		if err != nil {
			return zero, err
		}
	}

	data, err := input.File.Bytes()
	if err != nil {
		return zero, err
	}

	err = s.s3Client.Upload(ctx, "filecrypto", fileUUID.String(), data)
	if err != nil {
		return zero, err
	}

	return dto.File{
		Uuid:         fileUUID,
		Name:         file.Name,
		Size:         file.Size,
		CreatedAt:    file.CreatedAt,
		SymmetricKey: input.SymmetricKey,
	}, err
}

func (s *FilesService) DeleteFile(ctx context.Context, userUUID, fileUUID uuid.UUID) error {
	err := pgrepository.New(s.pgClient).DeleteFile(ctx, pgrepository.DeleteFileParams{
		UserUuid: userUUID,
		Uuid:     fileUUID,
	})
	if err != nil {
		return err
	}

	err = s.s3Client.Delete(ctx, "filecrypto", fileUUID.String())
	if err != nil {
		return err
	}

	return nil
}

func (s *FilesService) ShareFile(ctx context.Context, userUUID, fileUUID uuid.UUID, input dto.ShareFileInput) error {
	pg := pgrepository.New(s.pgClient)

	file, err := pg.FileByUUID(ctx, pgrepository.FileByUUIDParams{
		Uuid: fileUUID,
	})
	if err != nil {
		if pgrepository.IsNoRows(err) {
			return errs.NotFound.New("FileNotFound", "file not found")
		}

		return err
	}

	if file.UserUuid != userUUID {
		return errs.PermissionDenied.New("PermissionDenied", "permission denied")
	}

	_, err = pg.UpsertFileCryptoKey(ctx, pgrepository.UpsertFileCryptoKeyParams{
		FileUuid:     fileUUID,
		UserUuid:     input.RecipientUUID,
		SymmetricKey: input.SymmetricKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *FilesService) DeleteFileAccess(ctx context.Context, input dto.DeleteFileAccessInput) error {
	err := pgrepository.New(s.pgClient).DeleteFileAccess(ctx, pgrepository.DeleteFileAccessParams{
		RecipientUuid: input.RecipientUUID,
		FileUuid:      input.FileUUID,
		OwnerUuid:     input.OwnerUUID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *FilesService) UpdateUserKeys(ctx context.Context, userUUID uuid.UUID) error {
	fileUUIDs, err := pgrepository.New(s.pgClient).DeleteAllUserFiles(ctx, pgrepository.DeleteAllUserFilesParams{
		UserUuid: userUUID,
	})
	if err != nil {
		return err
	}

	err = pgrepository.New(s.pgClient).DeleteAlleFileAccess(ctx, pgrepository.DeleteAlleFileAccessParams{
		UserUuid: userUUID,
	})
	if err != nil {
		return err
	}

	for _, fileUUID := range fileUUIDs {
		err = s.s3Client.Delete(ctx, "filecrypto", fileUUID.String())
		if err != nil {
			log.Errorf("file service: failed to delete file %s: %s", fileUUID.String(), err.Error())
		}
	}

	return nil
}

func (s *FilesService) AvailableFiles(ctx context.Context, userUUID uuid.UUID) ([]dto.File, error) {
	rows, err := pgrepository.New(s.pgClient).AvailableFiles(ctx, pgrepository.AvailableFilesParams{
		UserUuid: userUUID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.File, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.File{
			Uuid:         row.File.Uuid,
			Name:         row.File.Name,
			Size:         row.File.Size,
			CreatedAt:    row.File.CreatedAt,
			SymmetricKey: &row.SymmetricKey,
		})
	}

	return result, nil
}
