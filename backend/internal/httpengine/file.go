package httpengine

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"univer/internal/dto"
	"univer/internal/httpengine/openapi"
	"univer/pkg/lib/errs"
	"univer/pkg/lib/httpapi"
)

func (c *controller) UpdateUserKeys(w http.ResponseWriter, r *http.Request) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		err := c.filesService.UpdateUserKeys(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		return true, nil
	}))
}

func (c *controller) Files(w http.ResponseWriter, r *http.Request) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		files, err := c.filesService.Files(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		return openapi.FilesResult{
			Files: files,
		}, nil
	}))
}

func (c *controller) DownloadFile(w http.ResponseWriter, r *http.Request, fileUUID uuid.UUID) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		fileName, fileData, err := c.filesService.File(ctx, userUUID, fileUUID)
		if err != nil {
			return nil, err
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))

		_, err = w.Write(fileData)
		if err != nil {
			return nil, errs.Internal.New("Failed to send file", "failed to send file")
		}

		return true, nil
	}))
}

func (c *controller) DownloadCommonFile(w http.ResponseWriter, r *http.Request, fileUUID types.UUID) {
	c.base(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		fileName, fileData, err := c.filesService.CommonFile(ctx, fileUUID)
		if err != nil {
			return nil, err
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))

		_, err = w.Write(fileData)
		if err != nil {
			return nil, errs.Internal.New("Failed to send file", "failed to send file")
		}

		return true, nil
	}))
}

func (c *controller) CreateFile(w http.ResponseWriter, r *http.Request) {
	c.auth(w, r, httpapi.HandlerWithForm(func(ctx context.Context, input openapi.CreateFileMultipartBody) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		var key *string
		if input.SymmetricKey != "" {
			key = &input.SymmetricKey
		}

		file, err := c.filesService.CreateFile(ctx, userUUID, dto.CreateFileInput{
			Name:         input.Name,
			SymmetricKey: key,
			File:         input.File,
		})
		if err != nil {
			return nil, err
		}

		return openapi.CreateFileResult{
			File: file,
		}, nil
	}))
}

func (c *controller) DeleteFile(w http.ResponseWriter, r *http.Request, fileUUID types.UUID) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		err := c.filesService.DeleteFile(ctx, userUUID, fileUUID)
		if err != nil {
			return nil, err
		}

		return true, nil
	}))
}

func (c *controller) ShareFile(w http.ResponseWriter, r *http.Request, fileUUID uuid.UUID) {
	c.auth(w, r, httpapi.HandlerWithInput(func(ctx context.Context, input openapi.ShareFileInput) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		err := c.filesService.ShareFile(ctx, userUUID, fileUUID, dto.ShareFileInput{
			RecipientUUID: input.RecipientUuid,
			SymmetricKey:  input.SymmetricKey,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}))
}

func (c *controller) DeleteFileAccess(w http.ResponseWriter, r *http.Request, fileUUID types.UUID) {
	c.auth(w, r, httpapi.HandlerWithInput(func(ctx context.Context, input openapi.DeleteFileAccessInput) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		err := c.filesService.DeleteFileAccess(ctx, dto.DeleteFileAccessInput{
			OwnerUUID:     userUUID,
			FileUUID:      fileUUID,
			RecipientUUID: input.RecipientUuid,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}))
}
func (c *controller) AvailableFiles(w http.ResponseWriter, r *http.Request) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		files, err := c.filesService.AvailableFiles(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		return openapi.AvailableFilesResult{
			Files: files,
		}, nil
	}))
}
