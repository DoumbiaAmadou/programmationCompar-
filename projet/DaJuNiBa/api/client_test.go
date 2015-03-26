package api

import (
	"fmt"
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeAPIServer() (*string, *httptest.Server) {
	// poor man's state
	users := make(map[string]string)
	var logged string

	wrongParamsHTML := `
        <!DOCTYPE html><!-- Page generated by OCaml with Ocsigen.
        See http://ocsigen.org/ and http://caml.inria.fr/ for information -->
        <html xmlns="http://www.w3.org/1999/xhtml"><head><title>Wrong
        parameters</title></head><body><h1>Wrong parameters</h1></body></html>`

	err404HTML := `
        <!DOCTYPE html><!-- Page generated by OCaml with Ocsigen.
        See http://ocsigen.org/ and http://caml.inria.fr/ for information -->
        <html xmlns="http://www.w3.org/1999/xhtml"><head>
        <title>Error 404</title></head><body><h1>Not Found</h1>
        <p>Error 404</p></body></html>`

	return &logged, httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		r.ParseForm()

		switch route {
		case "GET /0/api":
			if r.Form.Encode() != "" {
				w.WriteHeader(500)
				fmt.Fprint(w, wrongParamsHTML)
				break
			}
			w.WriteHeader(200)
			fmt.Fprint(w, `{
              "status": "completed",
              "response": {
              "doc": {
                "m1": {
                  "method": "post",
                  "input": [ "i1 : string", "i2 : string" ],
                  "output": [],
                  "errors": [
                    { "code": 334269347, "description": "USER_ALREADY_EXISTS" },
                    { "code": 114981602, "description": "INVALID_LOGIN" }
                  ],
                  "description": "Do something"
                },
                "m2": {
                  "method": "get",
                  "input": [],
                  "output": [ "foo : string" ],
                  "errors": [],
                  "description": "Do something else"
                }
              }
            }
          }`)

		case "POST /0/register":
			logins, hasLogin := r.PostForm["login"]
			passwds, hasPassword := r.PostForm["password"]

			if !hasLogin || !hasPassword ||
				len(logins) != 1 || len(passwds) != 1 || len(r.PostForm) != 2 {

				w.WriteHeader(500)
				fmt.Fprint(w, wrongParamsHTML)
				break
			}

			w.WriteHeader(200)

			login := logins[0]

			if _, userExists := users[login]; userExists {
				fmt.Fprint(w, fmt.Sprintf(`{
                  "status": "error",
                  "response": {
                    "error_code": 334269347,
                    "error_msg": "You already exist, %s."
                  }
                }`, login))
				break
			}

			// we register the user locally
			users[login] = passwds[0]

			fmt.Fprint(w, `{ "status": "completed", "response": {} }`)

		case "POST /0/auth":
			logins, hasLogin := r.PostForm["login"]
			passwds, hasPassword := r.PostForm["password"]

			if !hasLogin || !hasPassword ||
				len(logins) != 1 || len(passwds) != 1 || len(r.PostForm) != 2 {

				w.WriteHeader(500)
				fmt.Fprint(w, wrongParamsHTML)
				break
			}

			w.WriteHeader(200)

			login := logins[0]
			passwd := passwds[0]

			if p, userExists := users[login]; !userExists || p != passwd {
				fmt.Fprint(w, fmt.Sprintf(`{
                  "status": "error",
                  "response": {
                    "error_code": 502441794,
                    "error_msg": "I do not know you, %s."
                  }
                }`, login))
				break
			}

			logged = login

			w.Header().Add("Set-Cookie",
				"eliompersistentsession|||ref|=u9; path=/; "+
					"expires=Sun, 09-Mar-2025 15:03:26 GMT")

			fmt.Fprint(w, `{ "status": "completed", "response": {} }`)

		case "GET /0/logout":
			if len(r.Form) != 0 {
				w.WriteHeader(404)
				fmt.Fprintf(w, err404HTML)
				break
			}
			logged = ""
			w.WriteHeader(200)
			fmt.Fprint(w, `{ "status": "completed", "response": {} }`)

		// TODO mock the rest of the API

		default:
			w.WriteHeader(500)
			fmt.Fprint(w, err404HTML)
		}
	}))
}

func TestClient(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewClient", func() {
		g.It("Should not return nil", func() {
			o.Expect(NewClient()).NotTo(o.BeNil())
		})
	})

	g.Describe("Client", func() {
		var ts *httptest.Server
		var logged *string
		var c *Client

		g.BeforeEach(func() {
			logged, ts = NewFakeAPIServer()
			c = NewClient()
			c.http.baseURL = ts.URL
		})

		g.AfterEach(func() { ts.Close() })

		g.Describe(".getUserCredentialsParams()", func() {
			g.It("Should initially return an empty UserCredentialsParams", func() {
				p := UserCredentialsParams{}
				o.Expect(c.getUserCredentialsParams()).To(o.Equal(p))
			})

			g.It("Should return an UserCredentialsParams w/ username/password", func() {
				c.username = "foo"
				c.password = "secret"
				p := UserCredentialsParams{Login: "foo", Password: "secret"}
				o.Expect(c.getUserCredentialsParams()).To(o.Equal(p))
			})
		})

		g.Describe(".Authenticated()", func() {
			g.It("Should be initially false", func() {
				o.Expect(c.Authenticated()).To(o.BeFalse())
			})
		})

		g.Describe(".APIInfo()", func() {
			g.It("Should return API infos", func() {
				info, err := c.APIInfo()

				o.Expect(err).To(o.BeNil())
				o.Expect(info).NotTo(o.BeNil())
				o.Expect(info.Doc).NotTo(o.BeNil())
			})

			g.It("Should populate .Doc with all methods", func() {
				info, err := c.APIInfo()

				o.Expect(err).To(o.BeNil())
				o.Expect(info).NotTo(o.BeNil())
				o.Expect(info.Doc).NotTo(o.BeNil())
				o.Expect(len(info.Doc)).To(o.Equal(2))

				m1, m2 := info.Doc["m1"], info.Doc["m2"]

				o.Expect(m1).To(o.Equal(APIMethod{
					Verb:  "post",
					Input: []string{"i1 : string", "i2 : string"},
					//Output: []string{},
					Errors: []APIError{
						APIError{Code: 334269347,
							Description: "USER_ALREADY_EXISTS"},
						APIError{Code: 114981602,
							Description: "INVALID_LOGIN"},
					},
					Description: "Do something",
				}))

				o.Expect(m2).To(o.Equal(APIMethod{
					Verb:  "get",
					Input: []string{},
					//Output:      []string{"foo : string"},
					Errors:      []APIError{},
					Description: "Do something else",
				}))
			})
		})

		g.Describe(".RegisterWithCredentials(u,p)", func() {
			g.It("Should return nil if the username doesn't exist yet", func() {
				err := c.RegisterWithCredentials("foo", "bar")

				o.Expect(err).To(o.BeNil())
			})

			g.It("Should return an ErrUserAlreadyExists if the username already exists", func() {
				err1 := c.RegisterWithCredentials("foo", "bar")
				err2 := c.RegisterWithCredentials("foo", "qux")

				o.Expect(err1).To(o.BeNil())
				o.Expect(err2).To(o.Equal(ErrUserAlreadyExists))
			})
		})

		g.Describe(".LoginWithCredentials(u,p)", func() {
			g.It("Should return ErrUnknownUser if the user doesn't exist", func() {
				err := c.LoginWithCredentials("foo", "bar")
				o.Expect(err).To(o.Equal(ErrUnknownUser))
			})

			g.It("Should return ErrUnknownUser if the password doesn't match", func() {
				c.RegisterWithCredentials("foo", "barx")
				err := c.LoginWithCredentials("foo", "bar")
				o.Expect(err).To(o.Equal(ErrUnknownUser))
			})

			g.It("Should return nil if the credentials match", func() {
				c.RegisterWithCredentials("foo", "barx")
				err := c.LoginWithCredentials("foo", "barx")
				o.Expect(err).To(o.BeNil())
			})

			g.It("Should return nil if the user is already connected", func() {
				c.RegisterWithCredentials("foo", "barx")
				c.LoginWithCredentials("foo", "barx")
				err := c.LoginWithCredentials("foo", "barx")
				o.Expect(err).To(o.BeNil())
			})

			g.It("Should set .authenticated", func() {
				c.RegisterWithCredentials("foo", "barx")
				c.LoginWithCredentials("foo", "barx")
				o.Expect(c.authenticated).To(o.BeTrue())
			})

			g.It("Should logout if already logged with other credentials", func() {
				c.RegisterWithCredentials("foo1", "barx")
				c.RegisterWithCredentials("foo2", "barx")

				c.LoginWithCredentials("foo1", "barx")
				err := c.LoginWithCredentials("foo2", "barx")

				o.Expect(err).To(o.BeNil())
				o.Expect(*logged).To(o.Equal("foo2"))
			})
		})

		g.Describe(".Login()", func() {
			g.It("Should return ErrUserAlreadyExists if the credentials are not set", func() {
				o.Expect(c.Login()).To(o.Equal(ErrUnknownUser))
			})

			g.It("Should return nil if the credentials match", func() {
				c.RegisterWithCredentials("foo", "barx")
				err := c.Login()
				o.Expect(err).To(o.BeNil())
			})

			g.It("Should set .authenticated", func() {
				c.RegisterWithCredentials("foo", "barx")
				c.Login()
				o.Expect(c.authenticated).To(o.BeTrue())
			})
		})

		g.Describe(".Logout()", func() {
			g.It("Should return nil without calling the server "+
				"if the user is not logged", func() {

				ts.Close()
				o.Expect(c.Logout()).To(o.BeNil())
			})

			g.It("Should return nil if the user is logged", func() {
				c.RegisterWithCredentials("foo", "barx")
				c.Login()
				o.Expect(c.Logout()).To(o.BeNil())
			})

			g.It("Should unset .authenticated", func() {
				c.RegisterWithCredentials("foo", "barx")
				c.Login()
				o.Expect(c.Logout()).To(o.BeNil())
				o.Expect(c.authenticated).To(o.BeFalse())
			})
		})

		g.Describe(".ShutdownIdentifier(id)", func() {
			g.It("Should return ErrNotImplemented", func() {
				o.Expect(c.ShutdownIdentifier("foo")).
					To(o.Equal(ErrNotImplemented))
			})
		})
	})
}
