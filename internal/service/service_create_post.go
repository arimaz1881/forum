package service

import (
	"context"
	"fmt"
	"forum/internal/domain"
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TODO: add comment for block of code

type CreatePostInput struct {
	Title      string
	Content    string
	Categories []string
	File       *httphelper.File
	UserID     int64
}

type DeletePostInput struct {
	PostID	string
	UserID	int64
}

type EditPostInput struct {
	PostID	string
	Content string
	UserID	int64
}

func (s *service) CreatePost(ctx context.Context, input CreatePostInput) (_ int64, err error) {
	if err := input.validate(); err != nil {
		return 0, err
	}

	// check exist category IDs
	for _, categoreiaID := range input.Categories {
		if _, err := s.categories.GetOne(ctx, categoreiaID); err != nil {
			return 0, domain.ErrInvalidCategoryID
		}
	}

	var fileKey string

	if input.File != nil {
		fileExt := strings.ToLower(filepath.Ext(input.File.FileName))
		fileKey, err = s.saveFile(fileExt, input.File)
		if err != nil {
			return 0, err
		}
	}

	postID, err := s.posts.Create(ctx, domain.CreatePostInput{
		Title:   input.Title,
		Content: input.Content,
		UserID:  input.UserID,
		FileKey: fileKey,
	})
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	for _, categoriaID := range input.Categories {
		if err := s.postCategories.Create(ctx, domain.CreateCategoriaInput{
			PostID:      postID,
			CategoriaID: categoriaID,
		}); err != nil {
			return 0, err
		}
	}
	return postID, nil
}

func (i CreatePostInput) validate() error {
	if i.UserID == 0 {
		return domain.ErrInvalidToken
	}

	if len(i.Categories) == 0 {
		return domain.ErrInvalidCategory
	}

	i.Title = strings.ReplaceAll(i.Title, "ㅤ", "")

	i.Title = strings.TrimSpace(i.Title)
	if i.Title == "" {
		return domain.ErrInvalidContent
	}

	i.Content = strings.ReplaceAll(i.Content, "ㅤ", "")

	i.Content = strings.TrimSpace(i.Content)
	if i.Content == "" {
		return domain.ErrInvalidContent
	}

	if i.File == nil {
		return nil
	}

	if i.File.FileSize == 0 {
		return domain.ErrInvalidFile
	}

	fileExt := strings.ToLower(filepath.Ext(i.File.FileName))
	if fileExt != ".jpeg" && fileExt != ".jpg" && fileExt != ".png" && fileExt != ".gif" {
		return domain.ErrInvalidFile
	}

	return nil
}

func (s *service) saveFile(fileExt string, file *httphelper.File) (string, error) {
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	imagePath := filepath.Join(s.fileStorage, fileName)
	err := os.MkdirAll(s.fileStorage, os.ModePerm)
	if err != nil {
		return "", e3r.Internal("Failed to create dir for uploads")
	}

	dst, err := os.Create(imagePath)
	if err != nil {
		return "", e3r.Internal("Failed to save uploaded file")
	}
	defer dst.Close()

	_, err = io.Copy(dst, file.FileReader)
	if err != nil {
		return "", e3r.Internal("Failed to save uploaded file")
	}

	return fileName, nil
}



func (s *service) DeletePost(ctx context.Context, input DeletePostInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	post, err := s.posts.GetOne(ctx, input.PostID); 
	if err != nil {
		return err
	}

	if post.UserID != strconv.Itoa(int(input.UserID)) {
		return domain.ErrForbidden
	}

	if err := s.posts.Delete(ctx, input.PostID); err != nil {
		return err
	}

	return nil
}


func (i DeletePostInput) validate() error {
	if i.PostID == "" {
		return domain.ErrInvalidPostID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	return nil
}


func (s *service) EditPost(ctx context.Context, input EditPostInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	post, err := s.posts.GetOne(ctx, input.PostID); 
	if err != nil {
		return err
	}

	if post.UserID != strconv.Itoa(int(input.UserID)) {
		return domain.ErrForbidden
	}

	if err := s.posts.Edit(ctx, domain.EditPostInput{
		PostID: input.PostID,
		Content: input.Content}); 
		err != nil {
		return err
	}

	return nil
}


func (i EditPostInput) validate() error {
	if i.PostID == "" {
		return domain.ErrInvalidPostID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	i.Content = strings.ReplaceAll(i.Content, "ㅤ", "")

	i.Content = strings.TrimSpace(i.Content)
	if i.Content == "" {
		return domain.ErrInvalidContent
	}

	return nil
}