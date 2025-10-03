package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ko"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ko_translations "github.com/go-playground/validator/v10/translations/ko"
)

// ============================================================================
// 모델 정의 with Validation Tags
// ============================================================================

// User registration model
type UserRegistration struct {
	Email           string    `json:"email" binding:"required,email" label:"이메일"`
	Username        string    `json:"username" binding:"required,min=3,max=20,alphanum" label:"사용자명"`
	Password        string    `json:"password" binding:"required,min=8,max=50,strong_password" label:"비밀번호"`
	ConfirmPassword string    `json:"confirm_password" binding:"required,eqfield=Password" label:"비밀번호 확인"`
	Age             int       `json:"age" binding:"required,min=18,max=120" label:"나이"`
	Phone           string    `json:"phone" binding:"required,e164|korean_phone" label:"전화번호"`
	Website         string    `json:"website" binding:"omitempty,url" label:"웹사이트"`
	BirthDate       time.Time `json:"birth_date" binding:"required,before_today" time_format:"2006-01-02" label:"생년월일"`
	Gender          string    `json:"gender" binding:"required,oneof=male female other" label:"성별"`
	Country         string    `json:"country" binding:"required,iso3166_1_alpha2" label:"국가"`
	PostalCode      string    `json:"postal_code" binding:"required,postal_code" label:"우편번호"`
	Terms           bool      `json:"terms" binding:"required,eq=true" label:"약관동의"`
}

// Product model
type Product struct {
	Name        string   `json:"name" binding:"required,min=1,max=100" label:"상품명"`
	Description string   `json:"description" binding:"required,min=10,max=1000" label:"상품설명"`
	Price       float64  `json:"price" binding:"required,min=0.01,max=1000000" label:"가격"`
	SKU         string   `json:"sku" binding:"required,alphanum,len=10" label:"SKU"`
	Category    string   `json:"category" binding:"required,category" label:"카테고리"`
	Tags        []string `json:"tags" binding:"required,min=1,max=5,dive,min=2,max=20" label:"태그"`
	Stock       int      `json:"stock" binding:"required,min=0" label:"재고"`
	Images      []string `json:"images" binding:"required,min=1,max=10,dive,url" label:"이미지"`
	Available   bool     `json:"available" binding:"" label:"판매가능"`
}

// Address model with nested validation
type Address struct {
	Street     string `json:"street" binding:"required,min=5,max=100" label:"도로명"`
	City       string `json:"city" binding:"required,min=2,max=50" label:"도시"`
	State      string `json:"state" binding:"omitempty,min=2,max=50" label:"주/도"`
	Country    string `json:"country" binding:"required,iso3166_1_alpha2" label:"국가"`
	PostalCode string `json:"postal_code" binding:"required,postal_code" label:"우편번호"`
}

// Order model with complex validation
type Order struct {
	CustomerID   uint       `json:"customer_id" binding:"required,min=1" label:"고객ID"`
	Items        []OrderItem `json:"items" binding:"required,min=1,dive" label:"주문항목"`
	ShippingAddr Address    `json:"shipping_address" binding:"required" label:"배송주소"`
	BillingAddr  *Address   `json:"billing_address" binding:"omitempty" label:"청구주소"`
	PaymentMethod string    `json:"payment_method" binding:"required,oneof=card cash transfer" label:"결제방법"`
	CouponCode   string     `json:"coupon_code" binding:"omitempty,alphanum,len=8" label:"쿠폰코드"`
	Notes        string     `json:"notes" binding:"omitempty,max=500" label:"메모"`
}

type OrderItem struct {
	ProductID uint `json:"product_id" binding:"required,min=1" label:"상품ID"`
	Quantity  int  `json:"quantity" binding:"required,min=1,max=100" label:"수량"`
	Price     float64 `json:"price" binding:"required,min=0.01" label:"가격"`
}

// Credit card validation
type CreditCard struct {
	Number   string `json:"number" binding:"required,credit_card" label:"카드번호"`
	Name     string `json:"name" binding:"required,min=2,max=50" label:"카드소유자"`
	ExpMonth int    `json:"exp_month" binding:"required,min=1,max=12" label:"만료월"`
	ExpYear  int    `json:"exp_year" binding:"required,min=2024,max=2050" label:"만료년도"`
	CVV      string `json:"cvv" binding:"required,len=3|len=4,numeric" label:"CVV"`
}

// Search query with validation
type SearchQuery struct {
	Query    string `form:"q" binding:"required,min=1,max=100" label:"검색어"`
	Category string `form:"category" binding:"omitempty,category" label:"카테고리"`
	MinPrice float64 `form:"min_price" binding:"omitempty,min=0" label:"최소가격"`
	MaxPrice float64 `form:"max_price" binding:"omitempty,gtfield=MinPrice" label:"최대가격"`
	Sort     string `form:"sort" binding:"omitempty,oneof=price name date" label:"정렬"`
	Page     int    `form:"page" binding:"omitempty,min=1" label:"페이지"`
	PerPage  int    `form:"per_page" binding:"omitempty,min=1,max=100" label:"페이지당항목"`
}

// File upload validation
type FileUpload struct {
	File     string `json:"file" binding:"required,base64" label:"파일"`
	Filename string `json:"filename" binding:"required,min=1,max=255" label:"파일명"`
	MimeType string `json:"mime_type" binding:"required,oneof=image/jpeg image/png application/pdf" label:"MIME타입"`
	Size     int64  `json:"size" binding:"required,min=1,max=10485760" label:"파일크기"` // max 10MB
}

// ============================================================================
// Custom Validators
// ============================================================================

var (
	validate *validator.Validate
	trans    ut.Translator
)

// Initialize validators
func initValidators() {
	// Create validator instance
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v

		// Register custom validators
		v.RegisterValidation("strong_password", strongPassword)
		v.RegisterValidation("korean_phone", koreanPhone)
		v.RegisterValidation("postal_code", postalCode)
		v.RegisterValidation("category", categoryValidator)
		v.RegisterValidation("before_today", beforeToday)
		v.RegisterValidation("credit_card", creditCard)

		// Register custom tag name func
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// Strong password validator
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// Korean phone number validator
func koreanPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Korean phone format: 010-1234-5678 or 01012345678
	pattern := `^(010|011|016|017|018|019)[-]?\d{3,4}[-]?\d{4}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// Postal code validator (Korean)
func postalCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// Korean postal code: 5 digits
	pattern := `^\d{5}$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// Category validator
func categoryValidator(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	validCategories := []string{
		"electronics", "clothing", "food", "books", "toys",
		"sports", "home", "beauty", "automotive", "other",
	}

	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

// Before today validator
func beforeToday(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.Before(time.Now())
}

// Credit card validator (Luhn algorithm)
func creditCard(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	// Remove spaces and dashes
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	alternate := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')
		if alternate {
			n *= 2
			if n > 9 {
				n = n%10 + 1
			}
		}
		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

// ============================================================================
// Translator Setup
// ============================================================================

func setupTranslator() ut.Translator {
	en := en.New()
	ko := ko.New()
	uni := ut.New(en, en, ko)

	// Get translator for Korean
	trans, _ = uni.GetTranslator("ko")

	// Register translations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Korean translations
		ko_translations.RegisterDefaultTranslations(v, trans)

		// Custom error messages
		registerCustomTranslations(v, trans)
	}

	return trans
}

func registerCustomTranslations(v *validator.Validate, trans ut.Translator) {
	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "strong_password",
			translation: "{0}은(는) 대문자, 소문자, 숫자, 특수문자를 포함해야 합니다",
		},
		{
			tag:         "korean_phone",
			translation: "{0}은(는) 올바른 한국 전화번호 형식이어야 합니다",
		},
		{
			tag:         "postal_code",
			translation: "{0}은(는) 5자리 우편번호여야 합니다",
		},
		{
			tag:         "category",
			translation: "{0}은(는) 유효한 카테고리여야 합니다",
		},
		{
			tag:         "before_today",
			translation: "{0}은(는) 오늘 이전 날짜여야 합니다",
		},
		{
			tag:         "credit_card",
			translation: "{0}은(는) 유효한 신용카드 번호여야 합니다",
		},
	}

	for _, t := range translations {
		v.RegisterTranslation(t.tag, trans,
			func(ut ut.Translator) error {
				return ut.Add(t.tag, t.translation, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(t.tag, fe.Field())
				return t
			},
		)
	}
}

// ============================================================================
// Error Handling
// ============================================================================

// Custom error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   interface{} `json:"value,omitempty"`
}

// Format validation errors
func formatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: e.Translate(trans),
				Tag:     e.Tag(),
				Value:   e.Value(),
			})
		}
	}

	return errors
}

// ============================================================================
// Handlers
// ============================================================================

// User registration handler
func handleUserRegistration(c *gin.Context) {
	var user UserRegistration

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": formatValidationErrors(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Product creation handler
func handleProductCreation(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": formatValidationErrors(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

// Order creation handler
func handleOrderCreation(c *gin.Context) {
	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": formatValidationErrors(err),
		})
		return
	}

	// Additional business logic validation
	if len(order.Items) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order cannot contain more than 10 items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

// Search handler
func handleSearch(c *gin.Context) {
	var query SearchQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid search parameters",
			"details": formatValidationErrors(err),
		})
		return
	}

	// Set defaults
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PerPage == 0 {
		query.PerPage = 20
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search results",
		"query":   query,
		"results": []string{"result1", "result2", "result3"},
	})
}

// Credit card validation handler
func handleCreditCardValidation(c *gin.Context) {
	var card CreditCard

	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid credit card information",
			"details": formatValidationErrors(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Credit card is valid",
		"card": gin.H{
			"last_four": card.Number[len(card.Number)-4:],
			"exp":       fmt.Sprintf("%02d/%d", card.ExpMonth, card.ExpYear),
		},
	})
}

// File upload handler
func handleFileUpload(c *gin.Context) {
	var file FileUpload

	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid file upload",
			"details": formatValidationErrors(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
		"type":     file.MimeType,
	})
}

// Dynamic validation handler
func handleDynamicValidation(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dynamic validation rules based on data
	rules := make(map[string]string)

	for key := range data {
		switch key {
		case "email":
			rules[key] = "required,email"
		case "age":
			rules[key] = "required,min=0,max=150"
		case "website":
			rules[key] = "omitempty,url"
		default:
			rules[key] = "required"
		}
	}

	// Validate dynamically
	errors := validateMap(data, rules)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data validated successfully",
		"data":    data,
	})
}

// Validate map with dynamic rules
func validateMap(data map[string]interface{}, rules map[string]string) []ValidationError {
	var errors []ValidationError

	for field, rule := range rules {
		value, exists := data[field]
		if !exists && strings.Contains(rule, "required") {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s is required", field),
				Tag:     "required",
			})
			continue
		}

		// Additional validation logic here
		if err := validate.Var(value, rule); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				for _, e := range errs {
					errors = append(errors, ValidationError{
						Field:   field,
						Message: e.Translate(trans),
						Tag:     e.Tag(),
						Value:   value,
					})
				}
			}
		}
	}

	return errors
}

// ============================================================================
// Router Setup
// ============================================================================

func setupRouter() *gin.Engine {
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/register", handleUserRegistration)
		api.POST("/products", handleProductCreation)
		api.POST("/orders", handleOrderCreation)
		api.GET("/search", handleSearch)
		api.POST("/credit-card/validate", handleCreditCardValidation)
		api.POST("/upload", handleFileUpload)
		api.POST("/validate", handleDynamicValidation)
	}

	// Validation info endpoint
	router.GET("/validation/rules", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"built_in_validators": []string{
				"required", "email", "url", "min", "max", "len",
				"eq", "ne", "gt", "gte", "lt", "lte",
				"alpha", "alphanum", "numeric", "hexadecimal",
				"lowercase", "uppercase", "contains", "containsany",
				"excludes", "startswith", "endswith",
				"isbn", "isbn10", "isbn13", "uuid", "uuid3", "uuid4", "uuid5",
				"ascii", "base64", "ip", "ipv4", "ipv6",
				"datetime", "timezone",
			},
			"custom_validators": []string{
				"strong_password", "korean_phone", "postal_code",
				"category", "before_today", "credit_card",
			},
		})
	})

	// Test data endpoint
	router.GET("/test/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"valid_user": gin.H{
				"email":            "test@example.com",
				"username":         "testuser123",
				"password":         "Test123!@#",
				"confirm_password": "Test123!@#",
				"age":              25,
				"phone":            "010-1234-5678",
				"website":          "https://example.com",
				"birth_date":       "1999-01-01",
				"gender":           "male",
				"country":          "KR",
				"postal_code":      "12345",
				"terms":            true,
			},
			"invalid_user": gin.H{
				"email":            "invalid-email",
				"username":         "a",
				"password":         "weak",
				"confirm_password": "different",
				"age":              15,
				"phone":            "123",
				"website":          "not-a-url",
				"birth_date":       "2030-01-01",
				"gender":           "unknown",
				"country":          "XXX",
				"postal_code":      "abc",
				"terms":            false,
			},
		})
	})

	return router
}

// ============================================================================
// Main
// ============================================================================

func main() {
	// Initialize validators
	initValidators()

	// Setup translator
	setupTranslator()

	// Setup router
	router := setupRouter()

	log.Println("🚀 Validation Server starting on :8080")
	log.Println("✅ Custom validators registered")
	log.Println("🌍 Translator configured (Korean)")
	log.Println("")
	log.Println("Try the test endpoints:")
	log.Println("  GET  /test/data - Get sample valid/invalid data")
	log.Println("  GET  /validation/rules - List all validation rules")
	log.Println("  POST /api/v1/register - Test user registration")
	log.Println("")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}