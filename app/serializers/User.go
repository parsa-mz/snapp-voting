package serializers

type UserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ResetPassword struct {
	Password string `json:"password" binding:"required"`
}
type UserJWT struct {
	Access string `json:"access"`
}
