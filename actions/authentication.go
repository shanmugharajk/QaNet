package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/shanmugharajk/qanet/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginIndex - Handler for login index page.
func LoginIndex(c buffalo.Context) error {
	u := &models.User{}
	url := "/login"

	if c.Value("userId") != nil {
		return c.Redirect(302, "/")
	}

	c.Set("user", u)

	if c.Param("returnUrl") != "" {
		url = url + "?&returnUrl=" + c.Param("returnUrl")
	}

	if c.Param("src") == "question" {
		c.Set("AuthRedirect", "To ask question you need to login.")
	}

	c.Set("FormAction", url)

	return c.Render(200, r.HTML("authentication/login.html"))
}

// LoginNew - Handler for creating new login.
func LoginNew(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, _ := c.Value("tx").(*gorm.DB)

	user, err := loginUser(tx, u)
	if err != nil {
		c.Set("user", u)
		c.Set("error", err.Error())
		return c.Render(200, r.HTML("authentication/login.html"))
	}

	c.Session().Set("currentUserId", user.ID)
	c.Cookies().SetWithExpirationTime("qaid", "t", time.Now().Add(30*24*time.Hour))

	if c.Param("returnUrl") != "" {
		return c.Redirect(302, c.Param("returnUrl"))
	}

	return c.Redirect(302, "/")
}

// Logout clears the session and logs a user out
func Logout(c buffalo.Context) error {
	c.Session().Clear()
	DeleteCookie("qaid", c)
	DeleteCookie("qcursor", c)

	if c.Param("returnUrl") != "" {
		return c.Redirect(302, c.Param("returnUrl"))
	}

	return c.Redirect(302, "/")
}

// SignupIndex - Handler for signup page.
func SignupIndex(c buffalo.Context) error {
	u := &models.User{}
	c.Set("user", u)
	return c.Render(200, r.HTML("authentication/signup.html"))
}

// SignupNew creates a new user and redirects to the home if the
// user gets created successfully.
func SignupNew(c buffalo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx, _ := c.Value("tx").(*gorm.DB)

	verrors, err := models.CreateUser(tx, u)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrors.HasAny() {
		c.Set("user", u)
		c.Set("verrors", verrors)
		return c.Render(200, r.HTML("authentication/signup.html"))
	}

	c.Session().Set("currentUserId", u.ID)

	return c.Redirect(302, "/")
}

// SetCurrentUser attempts to find a user based on the currentUserId
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if id := c.Session().Get("currentUserId"); id != nil {
			u := models.User{}
			tx, _ := c.Value("tx").(*gorm.DB)
			res := tx.Where("id = ?", id).Find(&u)
			if res.Error != nil {
				return errors.WithStack(res.Error)
			}
			c.Set("currentUser", u)
			c.Set("userId", u.ID)
			c.Set("userRole", u.RoleID)
		}
		return next(c)
	}
}

func loginUser(tx *gorm.DB, u *models.User) (*models.User, error) {
	existing := new(models.User)

	if err := models.GetById(tx, u.ID, &existing); err != nil {
		return u, err
	}

	// confirm that the given password matches the hashed password from the db
	if err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(u.Password)); err != nil {
		return u, errors.New("invalid credentials")
	}

	return existing, nil
}
