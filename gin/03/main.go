package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// User 구조체 - JSON 바인딩을 위한 태그 포함
type User struct {
	ID    string `json:"id" form:"id" binding:"required"`
	Name  string `json:"name" form:"name" binding:"required"`
	Email string `json:"email" form:"email" binding:"required,email"`
	Age   int    `json:"age" form:"age" binding:"min=1,max=120"`
}

// SearchRequest - Query 파라미터 바인딩용
type SearchRequest struct {
	Query  string `form:"q" binding:"required"`
	Page   int    `form:"page,default=1" binding:"min=1"`
	Limit  int    `form:"limit,default=10" binding:"min=1,max=100"`
	Sort   string `form:"sort,default=desc"`
}

// UpdateRequest - PATCH 요청용 (부분 업데이트)
type UpdateRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}

func main() {
	r := gin.Default()

	// 1. Path 파라미터 바인딩
	// /users/:id - 단일 Path 파라미터
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "Path parameter example",
			"user_id": id,
		})
	})

	// /users/:id/posts/:postId - 다중 Path 파라미터
	r.GET("/users/:id/posts/:postId", func(c *gin.Context) {
		userID := c.Param("id")
		postID := c.Param("postId")
		c.JSON(http.StatusOK, gin.H{
			"message": "Multiple path parameters",
			"user_id": userID,
			"post_id": postID,
		})
	})

	// 2. Query 파라미터 바인딩
	// /search?q=golang&page=1&limit=10&sort=desc
	r.GET("/search", func(c *gin.Context) {
		// 방법 1: 개별적으로 가져오기
		query := c.Query("q")                    // 없으면 빈 문자열
		page := c.DefaultQuery("page", "1")      // 기본값 설정

		// 방법 2: 구조체에 바인딩
		var req SearchRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid query parameters",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Search results",
			"query": query,
			"page": page,
			"structured": req,
		})
	})

	// 3. JSON Body 파라미터 바인딩
	// POST /users - 새 사용자 생성
	r.POST("/users", func(c *gin.Context) {
		var user User

		// JSON 바인딩 (필수 필드 검증 포함)
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data",
				"details": err.Error(),
			})
			return
		}

		// 성공적으로 바인딩된 데이터 반환
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user": user,
		})
	})

	// 4. Form 데이터 바인딩
	// POST /users/form - HTML form 데이터 처리
	r.POST("/users/form", func(c *gin.Context) {
		var user User

		// Form 데이터 바인딩
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Form data received",
			"user": user,
		})
	})

	// 5. 복합 예제: Path + Query + Body
	// PUT /users/:id?notify=true
	r.PUT("/users/:id", func(c *gin.Context) {
		// Path 파라미터
		id := c.Param("id")

		// Query 파라미터
		notify := c.DefaultQuery("notify", "false")

		// Body 파라미터 (JSON)
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON data",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated",
			"id": id,
			"notify": notify,
			"updated_data": user,
		})
	})

	// 6. PATCH - 부분 업데이트 (옵셔널 필드)
	r.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var updateReq UpdateRequest
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid update data",
				"details": err.Error(),
			})
			return
		}

		// 업데이트할 필드만 처리
		updates := gin.H{"id": id}
		if updateReq.Name != nil {
			updates["name"] = *updateReq.Name
		}
		if updateReq.Email != nil {
			updates["email"] = *updateReq.Email
		}
		if updateReq.Age != nil {
			updates["age"] = *updateReq.Age
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User partially updated",
			"updates": updates,
		})
	})

	// 7. 파일 업로드
	r.POST("/upload", func(c *gin.Context) {
		// 단일 파일
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No file uploaded",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"filename": file.Filename,
			"size": file.Size,
		})
	})

	// 8. 다중 파일 업로드
	r.POST("/upload/multiple", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse multipart form",
				"details": err.Error(),
			})
			return
		}

		files := form.File["files"]
		var uploadedFiles []gin.H

		for _, file := range files {
			uploadedFiles = append(uploadedFiles, gin.H{
				"filename": file.Filename,
				"size": file.Size,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Multiple files uploaded",
			"count": len(files),
			"files": uploadedFiles,
		})
	})

	// 헬스 체크 엔드포인트
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"message": "Parameter binding examples server is running",
		})
	})

	// 서버 시작
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}