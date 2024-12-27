package helper

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Registrasi format gambar
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rs/zerolog/log"
)

// Daftar ekstensi gambar yang diperbolehkan
var allowedImageExtension = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// ValidateImageFile memvalidasi file gambar berdasarkan ekstensi, MIME type, dan validitas data gambar
func ValidateImageFile(file *multipart.FileHeader) error {
	// Periksa ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtension[ext] {
		errMsg := fmt.Sprintf("ekstensi file tidak diizinkan: %s", ext)
		log.Error().Str("filename", file.Filename).Msg(errMsg)
		return errors.New(errMsg)
	}

	// Buka file untuk validasi lebih lanjut
	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("gagal membuka file")
		return fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	contentType, _ := mimetype.DetectReader(src)

	if !strings.HasPrefix(contentType.String(), "image/") {
		errMsg := fmt.Sprintf("tipe MIME file tidak sesuai: %s", contentType.String())
		log.Error().Str("filename", file.Filename).Str("mime_type", contentType.String()).Msg(errMsg)
		return errors.New(errMsg)
	}

	// Pastikan file benar-benar gambar
	if _, _, err := image.Decode(src); err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("file bukan gambar valid")
		return fmt.Errorf("file bukan gambar valid: %w", err)
	}

	log.Info().Str("filename", file.Filename).Str("mime_type", contentType.String()).Msg("file valid")
	return nil
}
