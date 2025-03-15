package configs

const JwtKey = "my_secret_key"

const (
    // Объявляем привилегии нашей системы
    ReadPermission  = "read"
	AdminPermission  = "admin"

    // Объявляем роли нашей системы
    UserRole  = "user"
	AdminRole  = "admin"
)

var (
    // Связка роль — привилегии
    RolePermissions = map[string][]string{
        AdminRole:  {AdminPermission},
    }
)

var (
    // Связка пользователь — роль
    UserRoles = map[string][]string{
        "arthurCool": {AdminRole},
    }
)

var (
    // Связка путь — роль
    PathsRoles = map[string][]string{
        "/smokers": {AdminRole},
    }
)