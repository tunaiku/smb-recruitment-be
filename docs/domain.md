# Package domain

`import "github.com/tunaiku/mobilebanking/internal/app/domain"`

## Variables {#pkg-variables}

    var (
        ErrTransactionDetailNotFound = errors.BadRequest("com.tunaiku.service.cbs", "transaction detail not found")
        ErrAccountNotFound           = errors.BadRequest("com.tunaiku.service.cbs", "account not found")
    )

    var (
        ErrCredentialNotMatch   error = errors.BadRequest("com.tunaiku.service.mbanking", "invalid credential")
        ErrOtpNotConfigured     error = errors.BadRequest("com.tunaiku.service.mbanking", "otp credential not configured on the user")
        ErrOtpAlreadyConfigured error = errors.BadRequest("com.tunaiku.service.mbanking", "otp credential not configured on the user")
        ErrPinNotConfigured     error = errors.BadRequest("com.tunaiku.service.mbanking", "pin credential not nonfigured on the user")
        ErrUserNotFound         error = errors.BadRequest("com.tunaiku.service.mbanking", "user not found")
    )

    var (
        ErrUnauthorized = errors.Unauthorized("com.tunaiku.service.mbanking", "invalid credential")
    )

## type [AccountInformationService](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=649:816#L23) [¶](#AccountInformationService) {#AccountInformationService}

    type AccountInformationService interface {
        IsAccountExists(accountNumber string) bool
        GetTransactionPrivileges(accountNumber string) (TransactionPrivileges, error)
    }

## type [AuthenticationResult](/src/github.com/tunaiku/mobilebanking/internal/app/domain/authentication.go?s=181:259#L3) [¶](#AuthenticationResult) {#AuthenticationResult}

    type AuthenticationResult struct {
        AccessToken string `json:"access_token"`
    }

## type [AuthenticationService](/src/github.com/tunaiku/mobilebanking/internal/app/domain/authentication.go?s=261:379#L7) [¶](#AuthenticationService) {#AuthenticationService}

    type AuthenticationService interface {
        Authenticate(username string, password string) (AuthenticationResult, error)
    }

## type [AuthorizationMethod](/src/github.com/tunaiku/mobilebanking/internal/app/domain/transaction.go?s=172:200#L7) [¶](#AuthorizationMethod) {#AuthorizationMethod}

    type AuthorizationMethod int

    const (
        UnknownAuthorizationMethod AuthorizationMethod = iota
        OtpAuthorization
        PinAuthorization
    )

## type [ConfiguredCredential](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=914:990#L18) [¶](#ConfiguredCredential) {#ConfiguredCredential}

ConfiguredCredential Represent user credential

    type ConfiguredCredential struct {
        Pin *PinCredential
        Otp *OtpCredential
    }

### func (\*ConfiguredCredential) [IsOtpConfigured](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1194:1247#L29) [¶](#ConfiguredCredential.IsOtpConfigured) {#ConfiguredCredential.IsOtpConfigured}

    func (c *ConfiguredCredential) IsOtpConfigured() bool

IsOtpConfigured it would return true if user configure otp

### func (\*ConfiguredCredential) [IsPinConfigured](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1053:1106#L24) [¶](#ConfiguredCredential.IsPinConfigured) {#ConfiguredCredential.IsPinConfigured}

    func (c *ConfiguredCredential) IsPinConfigured() bool

IsPinConfigured it would return true if user configure pin

## type [FindUserResult](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1594:1727#L43) [¶](#FindUserResult) {#FindUserResult}

    type FindUserResult struct {
        ID               string
        Name             string
        AccountReference string
        JoinAt           time.Time
    }

## type [OtpCredential](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=723:772#L8) [¶](#OtpCredential) {#OtpCredential}

OtpCredential Represent user's otp credential

    type OtpCredential struct {
        PhoneNumber string
    }

## type [OtpCredentialManager](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=2089:2189#L67) [¶](#OtpCredentialManager) {#OtpCredentialManager}

    type OtpCredentialManager interface {
        UserCredentialValidator
        RequestNewOtp(userId string) error
    }

## type [PinCredential](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=822:863#L13) [¶](#PinCredential) {#PinCredential}

PinCredential Represent user's pin credential

    type PinCredential struct {
        Pin string
    }

## type [PinCredentialManager](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=2023:2087#L63) [¶](#PinCredentialManager) {#PinCredentialManager}

    type PinCredentialManager interface {
        UserCredentialValidator
    }

## type [Transaction](/src/github.com/tunaiku/mobilebanking/internal/app/domain/transaction.go?s=304:613#L15) [¶](#Transaction) {#Transaction}

    type Transaction struct {
        ID                  string
        UserID              string
        State               TransactionState
        AuthorizationMethod AuthorizationMethod
        TransactionCode     string
        Amount              *big.Float
        SourceAccount       string
        DestinationAccount  string
        CreatedAt           time.Time
    }

## type [TransactionCreation](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=442:647#L14) [¶](#TransactionCreation) {#TransactionCreation}

    type TransactionCreation struct {
        SourceAccount      string
        DestinationAccount string
        TransactionCode    string
        Amount             *big.Float
        Currency           string
        TransactionDate    *time.Time
    }

## type [TransactionDetail](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=304:385#L5) [¶](#TransactionDetail) {#TransactionDetail}

    type TransactionDetail struct {
        Code          string
        MinimumAmount *big.Float
    }

## type [TransactionInformationService](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=818:935#L28) [¶](#TransactionInformationService) {#TransactionInformationService}

    type TransactionInformationService interface {
        FindTransactionDetailByCode(code string) (TransactionDetail, error)
    }

## type [TransactionPrivileges](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=387:440#L10) [¶](#TransactionPrivileges) {#TransactionPrivileges}

    type TransactionPrivileges struct {
        Codes []string
    }

## type [TransactionService](/src/github.com/tunaiku/mobilebanking/internal/app/domain/savings.go?s=937:1040#L32) [¶](#TransactionService) {#TransactionService}

    type TransactionService interface {
        CreateTransaction(transactionCreation TransactionCreation) error
    }

## type [TransactionState](/src/github.com/tunaiku/mobilebanking/internal/app/domain/transaction.go?s=48:73#L1) [¶](#TransactionState) {#TransactionState}

    type TransactionState int

    const (
        UnknownTransactionStatus TransactionState = iota
        WaitAuthorization
        Failed
        Success
    )

## type [User](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1274:1592#L33) [¶](#User) {#User}

    type User struct {
        ID                              string
        Name                            string
        AccountReference                string
        JoinDate                        time.Time
        Username                        string
        Password                        string
        ConfiguredTransactionCredential *ConfiguredCredential
    }

## type [UserCredentialValidator](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1929:2021#L59) [¶](#UserCredentialValidator) {#UserCredentialValidator}

    type UserCredentialValidator interface {
        Validate(userId string, credential string) error
    }

## type [UserRepository](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1729:1846#L50) [¶](#UserRepository) {#UserRepository}

    type UserRepository interface {
        LoadUser(id string) (*User, error)
        LoadByUsername(username string) (*User, error)
    }

## type [UserService](/src/github.com/tunaiku/mobilebanking/internal/app/domain/user.go?s=1848:1927#L55) [¶](#UserService) {#UserService}

    type UserService interface {
        FindUser(userId string) (FindUserResult, error)
    }

## type [UserSession](/src/github.com/tunaiku/mobilebanking/internal/app/domain/authentication.go?s=381:415#L11) [¶](#UserSession) {#UserSession}

    type UserSession struct {
        *User
    }

## type [UserSessionHelper](/src/github.com/tunaiku/mobilebanking/internal/app/domain/authentication.go?s=417:523#L15) [¶](#UserSessionHelper) {#UserSessionHelper}

    type UserSessionHelper interface {
        GetFromContext(ctx context.Context) (session UserSession, err error)
    }
