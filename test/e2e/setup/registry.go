package setup

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi"
	"github.com/tunaiku/mobilebanking/internal/app/authentication"
	"github.com/tunaiku/mobilebanking/internal/app/savings"
	"github.com/tunaiku/mobilebanking/internal/app/transaction"
	"github.com/tunaiku/mobilebanking/internal/app/user"
	"github.com/tunaiku/mobilebanking/internal/pkg/pg"
	"go.uber.org/dig"
)

var (
	Container = dig.New()
)

func init() {
	log.Println("register ...")
	transaction.Register(Container)
	pg.Register(Container)
	authentication.Register(Container)
	savings.Register(Container)
	user.Register(Container)
	Container.Provide(func() chi.Router {
		return chi.NewRouter()
	})

}

func InvokeHttpTest(t *testing.T, testFunc func(expect *httpexpect.Expect)) {
	transaction.Invoke(Container)
	authentication.Invoke(Container)
	savings.Invoke(Container)
	user.Invoke(Container)
	Container.Invoke(func(router chi.Router) {
		server := httptest.NewServer(router)
		defer server.Close()
		e := httpexpect.New(t, server.URL)
		testFunc(e)
	})
}
