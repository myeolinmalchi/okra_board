package services

import (
    "okra_board/config"
    "okra_board/common/encryption"
    "okra_board/models"
    "regexp"
    "github.com/golang-jwt/jwt"
    "time"
    "errors"
)

type ValidationResult struct {
    msg     string
}

func (r *ValidationResult) getOrNil() *string {
    if r == nil {
        return nil
    } else {
        return &r.msg
    }
}

func checkAdminExists(column string, value string) (exists bool) {
    config.DB.Table("admin").
        Select("count(*) > 0").
        Where(column+" = ?", value).
        Find(&exists)
    return 
}

func AdminValidationHandler(isRegist bool) func(admin *models.Admin) *models.AdminValidationResult {
    checkID := func (id string) *string {
        var msg string
        if match, _ := regexp.MatchString("^[a-z]+[a-z0-9]{5,19}$", id); !match {
            msg = "유효하지 않은 ID입니다."
        } else if checkAdminExists("id", id) {
            msg = "이미 존재하는 ID입니다."
        } else {
            return nil
        }
        return &msg
    }
    checkPW := func (pw string) *string {
        var msg string
        if match, _ := regexp.MatchString("^(?=.*\\d)(?=.*[a-zA-Z])[0-9a-zA-Z]{8,16}$", pw); !match {
            msg = "유효하지 않은 PW입니다."
        } else {
            return nil
        }
        return &msg
    }
    checkName := func (name string) *string {
        var msg string
        if match, _ := regexp.MatchString("^[ㄱ-힣]+$", name); !match {
            msg = "유효하지 않은 이름입니다."
        } else {
            return nil
        }
        return &msg
    }
    checkEmail := func (email string) *string {
        var msg string
        if match, _ := regexp.
            MatchString("^[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*@[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*.[a-zA-Z]{2,3}$", email); !match {
            msg = "유효하지 않은 이메일입니다."
        } else if checkAdminExists("email", email) {
            msg = "이미 사용중인 이메일입니다."
        } else {
            return nil
        }
        return &msg
    }
    checkPhone := func (phone string) *string {
        var msg string
        if match, _ := regexp.MatchString("^\\d{3}-\\d{3,4}-\\d{4}$", phone); !match {
            msg = "유효하지 않은 전화번호입니다."
        } else if checkAdminExists("phone", phone) {
            msg = "이미 사용중인 전화번호입니다."
        } else {
            return nil
        }
        return &msg
    }

    return func(admin *models.Admin) *models.AdminValidationResult {
        result := &models.AdminValidationResult{}
        if isRegist {
            result = &models.AdminValidationResult{
                ID: checkID(admin.ID),
                Password: checkPW(admin.Password),
                Name: checkName(admin.Name),
                Email: checkEmail(admin.Email),
                Phone: checkPhone(admin.Phone),
            }
        } else {
            result = &models.AdminValidationResult{
                Password: checkPW(admin.Password),
                Name: checkName(admin.Name),
                Email: checkEmail(admin.Email),
                Phone: checkPhone(admin.Phone),
            }
        }
        if  result.ID == nil && 
            result.Password == nil && 
            result.Name == nil &&
            result.Email == nil &&
            result.Phone == nil {
            return nil
        } else {
            return result
        }
    }
}

func AdminInfoHandler(isRegist bool) func(admin *models.Admin) (bool, *models.AdminValidationResult) {
    return func(admin*models.Admin) (bool, *models.AdminValidationResult) {
        result := AdminValidationHandler(true)(admin)
        if result != nil {
            return false, result
        } else {
            admin.Password = encryption.EncryptSHA256(admin.Password)
            query := config.DB.Table("admin")
            if isRegist {
                query = query.Create(admin)
            } else {
                query = query.UpdateColumns(admin)
            }
            aff := query.RowsAffected
            return aff == 1, nil
        }
    }
}

func Regist(admin *models.Admin) (bool, *models.AdminValidationResult) {
    return AdminInfoHandler(true)(admin)
}

func Modify(admin *models.Admin) (bool, *models.AdminValidationResult) {
    return AdminInfoHandler(true)(admin)
}

func Login(admin *models.Admin) bool {
    insertedPassword := admin.Password
    row := config.DB.Table("admin").Find(admin, "id", admin.ID)
    if row.RowsAffected == 0 {
        return false
    } else {
        return encryption.EncryptSHA256(insertedPassword) == admin.Password 
    }
}

func CreateToken(id string) (string, error) {
    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["id"] = id
    atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    token, err := at.SignedString([]byte("SECRET_KEY"))
    if err != nil {
        return "", err
    }
    return token, nil
}

func VerifyToken(token string) error {
    claims := jwt.MapClaims{}
    verifying := func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, errors.New("Unexpected Signing Method")
        }
        return []byte("SECRET_KEY"), nil
    }
    _, err := jwt.ParseWithClaims(token, &claims, verifying)
    return err
}
