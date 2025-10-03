package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ========================================
	// 1. 기본 정적 파일 서빙
	// ========================================

	// /static 경로로 static 폴더의 파일들 서빙
	// URL: http://localhost:8080/static/...
	// 파일: ./static/...
	r.Static("/static", "./07/static")

	// ========================================
	// 2. 단일 파일 서빙
	// ========================================

	// favicon.ico 서빙
	r.StaticFile("/favicon.ico", "./07/static/favicon.ico")

	// robots.txt 서빙
	r.StaticFile("/robots.txt", "./07/static/robots.txt")

	// ========================================
	// 3. 파일 시스템 서빙 (고급)
	// ========================================

	// http.FileSystem 인터페이스 사용
	r.StaticFS("/assets", http.Dir("./07/static"))

	// ========================================
	// 4. 다운로드 전용 엔드포인트
	// ========================================

	// 파일 다운로드 (Content-Disposition 헤더 포함)
	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filepath := fmt.Sprintf("./07/uploads/%s", filename)

		// 파일 존재 확인
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
			return
		}

		// Content-Disposition 헤더로 다운로드 강제
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Type", "application/octet-stream")
		c.File(filepath)
	})

	// ========================================
	// 5. 업로드와 정적 서빙 조합
	// ========================================

	// 파일 업로드
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file uploaded",
			})
			return
		}

		// 업로드 폴더 확인 및 생성
		uploadPath := "./07/uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.MkdirAll(uploadPath, 0755)
		}

		// 파일 저장
		filename := filepath.Base(file.Filename)
		dst := filepath.Join(uploadPath, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"filename": filename,
			"size":     file.Size,
			"url":      fmt.Sprintf("/uploads/%s", filename),
		})
	})

	// 업로드된 파일 서빙
	r.Static("/uploads", "./07/uploads")

	// ========================================
	// 6. SPA (Single Page Application) 지원
	// ========================================

	// SPA를 위한 모든 경로 처리
	r.NoRoute(func(c *gin.Context) {
		// API 경로는 제외
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:5] == "/api/" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "API endpoint not found",
			})
			return
		}

		// 정적 파일이 존재하면 서빙
		requestedPath := c.Request.URL.Path
		fullPath := filepath.Join("./07/static", requestedPath)

		if _, err := os.Stat(fullPath); err == nil {
			c.File(fullPath)
			return
		}

		// 파일이 없으면 index.html 반환 (SPA)
		c.File("./07/static/index.html")
	})

	// ========================================
	// 7. 캐시 제어
	// ========================================

	r.GET("/cached/*filepath", func(c *gin.Context) {
		// 캐시 헤더 설정
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("ETag", "W/\"123456\"")

		filepath := c.Param("filepath")
		c.File("./07/static" + filepath)
	})

	// ========================================
	// API 엔드포인트 (정적 파일과 함께 사용 예제)
	// ========================================

	api := r.Group("/api")
	{
		// 파일 목록 API
		api.GET("/files", func(c *gin.Context) {
			files := []gin.H{}

			// uploads 폴더의 파일 목록
			uploadPath := "./07/uploads"
			if entries, err := os.ReadDir(uploadPath); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						info, _ := entry.Info()
						files = append(files, gin.H{
							"name": entry.Name(),
							"size": info.Size(),
							"url":  fmt.Sprintf("/uploads/%s", entry.Name()),
						})
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"files": files,
				"total": len(files),
			})
		})

		// 파일 삭제 API
		api.DELETE("/files/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			filepath := fmt.Sprintf("./07/uploads/%s", filename)

			if err := os.Remove(filepath); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "File not found or cannot be deleted",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "File deleted successfully",
				"filename": filename,
			})
		})

		// 파일 정보 API
		api.GET("/files/:filename/info", func(c *gin.Context) {
			filename := c.Param("filename")
			filepath := fmt.Sprintf("./07/uploads/%s", filename)

			info, err := os.Stat(filepath)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "File not found",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"name":         info.Name(),
				"size":         info.Size(),
				"modified":     info.ModTime(),
				"is_directory": info.IsDir(),
			})
		})
	}

	// 루트 경로 - 파일 관리 UI
	r.GET("/", func(c *gin.Context) {
		c.File("./07/static/index.html")
	})

	// 서버 시작
	fmt.Println("Server is running on :8080")
	fmt.Println("Static files: http://localhost:8080/static/")
	fmt.Println("File manager: http://localhost:8080/")
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}