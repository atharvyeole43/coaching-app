package middleware

import (
	"coaching-app-backend/constant"
	"coaching-app-backend/pkg/response"
	"coaching-app-backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type MiddlewareResponseBody struct {
	StatusCode int
	Message    string
	Body       struct{}
}
type RequestInfo struct {
	URL     string
	Agent   string
	Headers map[string][]string
}

// There are a few issues in the given code:
//
// 1. Using logrus.Info with multiple arguments doesn't produce a formatted string; it's preferred to use logrus.Infof or use structured logging with fields.
// 2. The logging code logs potentially sensitive headers and user-agent to logs, which could be a security/log-overflow concern.
// 3. The code reuses requestInfo for logging three times, some of which is redundant.
// 4. The middleware does not call c.Next(), so handlers downstream won't execute if required in some middlewares (though this may be intentional if TokenAuthenticationHandler fully handles it).
// 5. There's no error recovery or panic handling in the middleware.
//
// Improved version below, with suggestions to mitigate issues:
//
// - Use structured logging for request info.
// - Reduce redundant logs.
// - Consider adding c.Next() call if you want further handlers to run.
//
// (If TokenAuthenticationHandler fully handles unauthorized, there's usually no need for c.Next(). Adjust as per actual intent.)

// func TokenAuthentication() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// logrus.WithFields(logrus.Fields{
// 		// 	"url":     c.Request.RequestURI,
// 		// 	"headers": c.Request.Header,
// 		// 	"agent":   c.Request.Header.Get("User-Agent"),
// 		// }).Info("Incoming request for token authentication")
// 		TokenAuthenticationHandler(c)
// 	}
// }

type GetUserDetails struct {
	UserId          string `gorm:"column:user_id"`
	EmployeeCode    string `gorm:"column:employee_code"`
	CbaId           string `gorm:"column:cba_id"`
	IsRmUser        string `gorm:"column:is_rm_user"`
	AccountType     string `gorm:"column:account_type"`
	DataAccessLevel string `gorm:"column:data_access_level"`
	CreatedAt       string `gorm:"column:created_at"`
	UpdatedAt       string `gorm:"column:updated_at"`
	IsApproved      int    `gorm:"column:is_approved"`
}

// func TokenAuthenticationHandler(c *gin.Context) {
// 	handlePanicRecovery(c)

// 	authorization := c.Request.Header.Get("Authorization")
// 	if authorization == "" {
// 		RespondWithError(c, http.StatusUnauthorized, constant.ErrorConstants.UnauthorizedAccess)
// 		return
// 	}

// 	tokenValidClaims, tokenErr := getTokenClaims(authorization)
// 	logrus.Info(tokenValidClaims)
// 	generateDailyLogFile()
// 	if tokenErr != nil {
// 		RespondWithError(c, http.StatusUnauthorized, "Invalid Token.")
// 		return
// 	}

// 	user, err := getUserDetails(c, tokenValidClaims)
// 	if err != nil {
// 		RespondWithError(c, http.StatusUnauthorized, constant.ErrorConstants.UnauthorizedAccess)
// 		return
// 	}
// 	logrus.Info("TokenAuthenticationHandler User date Response ", user, " User Subject ID ", tokenValidClaims["sub"])

// 	clientId := os.Getenv("SSO_CLIENT_ID")
// 	if clientId == "" {
// 		clientId = tokenValidClaims["client_id"].(string)
// 		logrus.Info("@TokenAuthenticationHandler  clientId from token: ", clientId)
// 	}

// 	setHeaders(c, user, tokenValidClaims, clientId)

// 	chkReq := parseRequestBody(c)

// 	if val, exists := chkReq["cba_id"]; !exists || val == "" {
// 		c.Request.Header.Set("cba_id", user.UserId)
// 	}
// 	if val, exists := chkReq["cba_code"]; !exists || val == "" {
// 		if dataAccessType, daExists := chkReq["data_access_type"]; !daExists || dataAccessType != "admin" {
// 			c.Request.Header.Set("cba_code", user.CbaId)
// 		}
// 	}
// 	if val, exists := chkReq["cba_hierarchy_code"]; !exists || val == "" {
// 		c.Request.Header.Set("cba_hierarchy_code", user.CbaId)
// 	}

// 	logTracingId := uuid.NewString()
// 	c.Request.Header.Set("log-tracing-id", logTracingId)
// 	c.Writer.Header().Set("log-tracing-id", logTracingId)
// 	logrus.Info("user : ", user)

// 	logrus.Info("Before calling next handler", c.Request.RequestURI, tokenValidClaims["sub"])
// 	c.Next()
// 	logrus.Info("After calling next handler", c.Request.RequestURI, tokenValidClaims["sub"])
// }

func parseRequestBody(c *gin.Context) map[string]interface{} {
	reqBody := response.RequestBodyLogger(c)
	logrus.Info("Request Body:", reqBody)
	var chkReq map[string]interface{}
	// Handle empty request body safely
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		// Only try to parse JSON if content-type is application/json
		if reqBody == "" {
			logrus.Warn("Empty request body received")
			chkReq = make(map[string]interface{})
		} else {
			if err := json.Unmarshal([]byte(reqBody), &chkReq); err != nil {
				logrus.Error("JSON Parsing Error:", err)
				chkReq = make(map[string]interface{})
			}
		}
	} else {
		logrus.Info("Non-JSON request, skipping body parse: Content-Type =", contentType)
		chkReq = make(map[string]interface{}) // fallback
	}
	return chkReq
}

func setHeaders(c *gin.Context, user GetUserDetails, claims jwt.MapClaims, clientId string) {
	c.Request.Header.Set("common_prefix_path", os.Getenv("COMMON_PORTAL_PREFIX")+"/"+os.Getenv("PORTAL_ENVIRONMENT"))
	c.Request.Header.Set("loan_prefix_path", os.Getenv("LOAN_PORTAL_PREFIX")+"/"+os.Getenv("PORTAL_ENVIRONMENT"))
	c.Request.Header.Set("insurance_prefix_path", os.Getenv("INSURANCE_PORTAL_PREFIX")+"/"+os.Getenv("PORTAL_ENVIRONMENT"))
	c.Request.Header.Set("subject_id", claims["sub"].(string))
	c.Request.Header.Set("logged_in_user_id", user.UserId)
	c.Request.Header.Set("logged_in_user_cba_id", user.CbaId)
	c.Request.Header.Set("is_rm_user", user.IsRmUser)
	c.Request.Header.Set("name", claims["name"].(string))
	c.Request.Header.Set("client_id", clientId)
	c.Request.Header.Set("agent_no", claims["agent_no"].(string))
	c.Request.Header.Set("unique_id", claims["id"].(string))
	c.Request.Header.Set("user_type", "web")
	c.Request.Header.Set("data_access_level", user.DataAccessLevel)
	c.Request.Header.Set("employee_code", user.EmployeeCode)
	c.Request.Header.Set("is_approved", strconv.Itoa(user.IsApproved))

	if user.UserId != "" {
		c.Request.Header.Set("logged_in_user_id", user.UserId)
		c.Request.Header.Set("user_id", user.UserId)
	}
	if user.AccountType != "" {
		c.Request.Header.Set("logged_in_user_account_type", user.AccountType)
	}
	if user.IsRmUser != "" {
		c.Request.Header.Set("is_rm_user", user.IsRmUser)
	}

	logrus.Info("user : ", user)
}

// func getUserDetails(c *gin.Context, claims jwt.MapClaims) (GetUserDetails, error) {
// 	var user GetUserDetails
// 	if err := dbstore.GetDatabaseConfig().ConnectDB.Raw("SELECT user_details.user_id,user_details.employee_code,user_details.cba_id,user_details.account_type,user_details.is_rm_user,user_details.data_access_level, user_details.is_approved FROM user_details WHERE user_details.sso_subject_id = ?", claims["sub"]).Scan(&user).Error; err != nil {
// 		logrus.Error("TokenAuthentication Error: ", err)
// 		return user, err
// 	}
// 	return user, nil
// }

func handlePanicRecovery(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("TokenAuthenticationHandler Panic Error:", c.Request.RequestURI, panicInfo)
			errorResponse := map[string]interface{}{
				"StatusCode": 500,
				"Message":    "Internal Server Error",
				"DevMessage": "",
				"Body":       map[string]interface{}{},
			}
			c.JSON(http.StatusInternalServerError, errorResponse)
			c.Abort()
			return
		}
	}()
}

func inArray(val string, array []string) (ok bool) {
	var i int
	for i = range array {
		if ok = array[i] == val; ok {
			return true
		}
	}
	return false
}

func TokenValid(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//if claims, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return claims, nil
	}
	return nil, err
}

func RespondWithError(c *gin.Context, code int, message string) {
	response := MiddlewareResponseBody{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
	c.AbortWithStatusJSON(code, response)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {

	pubKeyFileExt := "pem"
	if os.Getenv("APP_ENV") == "production" {
		pubKeyFileExt = "pub"
	}
	keyData, err := ioutil.ReadFile("storage/ssl/" + os.Getenv("APP_ENV") + "/jwtRSA256-public-" + os.Getenv("APP_ENV") + "." + pubKeyFileExt)
	if err != nil {
		logrus.Error("Failed to read public key file:", err)
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		logrus.Error("Failed to parse RSA public key from PEM:", err)
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

func generateDailyLogFile() {
	now := time.Now() //or time.Now().UTC()
	logFileName := now.Format(constant.DateFormat) + ".log"

	file, err := os.OpenFile(path.Join("./storage/logs", logFileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error("Failed to open log file:", err)
		return // Exit if file cannot be opened
	}
	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableHTMLEscape: true,
		PrettyPrint:       true,
		TimestampFormat:   constant.DateFormat + " 15:04:05",
	})
	gin.DefaultWriter = io.MultiWriter(file)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}

func APIKeyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := c.Request.RequestURI
		logrus.Info("requestURL:", requestURL)
		APIKeyAuthenticator(c)
	}
}

func DecryptAndVerifyAPIKey(auth_api_key string) error {
	token, err := jwt.Parse(auth_api_key, func(token *jwt.Token) (interface{}, error) {
		secret := strings.TrimSpace(utils.MustGetenv("IR_X_API_ENCRYPTION_KEY"))
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return err
	}

	if claims["api_key"].(string) != os.Getenv("IR_X_API_KEY") {
		return errors.New("can not match provided x-api-key or expired")
	}
	return nil
}

func APIKeyAuthenticator(c *gin.Context) {
	logrus.Info("app_env : ---", os.Getenv("APP_ENV"))
	authorization := c.Request.Header.Get("X-Api-Key")
	if authorization != "" {
		err := DecryptAndVerifyAPIKey(authorization)
		if err != nil {
			logrus.Errorf("can not Match Provided API key for the Authorization : %v", err)
			response.UnAuthorizedResponse(c, err, "Provided API key is not matching  or expired")
			return
		}
		c.Next()
	} else {
		response.UnAuthorizedResponse(c, errors.New("Authorization can not be Empty"), constant.ErrorConstants.UnauthorizedAccess)
		return
	}
}

func getTokenClaims(authorization string) (jwt.MapClaims, error) {
	tokenValidClaims, tokenErr := TokenValid(authorization)
	if tokenErr != nil {
		return nil, tokenErr
	}
	return tokenValidClaims, nil
}
