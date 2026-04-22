package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"image"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"

	_ "image/jpeg"
	_ "image/png"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/disintegration/imaging"
)

type ImageProvider interface {
	UploadImage(ctx context.Context, key string, body io.Reader, contentType string) error
	DeleteImage(ctx context.Context, key string) error
	GetPathImage(ctx context.Context) string
	GetPublicUrl(ctx context.Context, key string) string
}

type ImageService struct {
	Log      *slog.Logger
	Provider ImageProvider
}

func (s *ImageService) UploadAvatar(
	ctx context.Context,
	userID []byte,
	fileHeader *multipart.FileHeader,
) (avatarKey models.AvatarPath, err error) {

	const op = "ImageService.UploadImage"

	defer func() { err = utils.ErrWrap(op, err) }()

	log := s.Log.With("op", op)

	file, err := fileHeader.Open()
	if err != nil {
		log.Error("can't open fileHeader", utils.SlogErr(err))
		return models.AvatarPath{}, err
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		log.Error("can't read fileHeader", utils.SlogErr(err))
		return models.AvatarPath{}, err
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = guessExt(contentType)
	}

	noise := make([]byte, 8)
	if _, err := rand.Read(noise); err != nil {
		log.Error("can't generate noise", utils.SlogErr(err))
		return models.AvatarPath{}, err
	}

	origKey := fmt.Sprintf("%x%x/original%s", noise, userID, ext)
	thumbKey := fmt.Sprintf("%x%x/thumb%s", noise, userID, ext)

	thumb, err := generateThumbnail(bytes.NewReader(buf), 256)
	if err != nil {
		log.Error("can't generate thumbnail", utils.SlogErr(err))
		return models.AvatarPath{}, err
	}

	if err := s.Provider.UploadImage(
		ctx, origKey, bytes.NewReader(buf), contentType,
	); err != nil {
		log.Error("can't upload original image", utils.SlogErr(err))
		return models.AvatarPath{}, err
	}

	if err := s.Provider.UploadImage(
		ctx, thumbKey, thumb, contentType,
	); err != nil {
		log.Error("can't upload thumbnail image", utils.SlogErr(err))
		return avatarKey, err
	}

	return models.AvatarPath{
		Original:  s.Provider.GetPublicUrl(ctx, origKey),
		Thumbnail: s.Provider.GetPublicUrl(ctx, thumbKey),
	}, nil
}

func guessExt(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".img"
	}
}

func generateThumbnail(r io.Reader, maxSize int) (io.Reader, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	thumb := imaging.Thumbnail(img, maxSize, maxSize, imaging.Lanczos)

	out := new(bytes.Buffer)
	err = imaging.Encode(out, thumb, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return out, nil
}
