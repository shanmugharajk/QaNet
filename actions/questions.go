package actions

import (
	"html/template"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/shanmugharajk/qanet/models"
	"github.com/shanmugharajk/qanet/services"
)

// AskQuestionIndex returns the form for creating new post.
func AskQuestionIndex(c buffalo.Context) error {
	tx, _ := c.Value("tx").(*gorm.DB)

	tags, err := services.FetchAllTags(tx)
	if err != nil {
		return err
	}

	c.Set("tags", tags)

	return c.Render(200, r.HTML("questions/ask.html"))
}

// AskQuestion accepts the posted data and creates a new question.
func AskQuestion(c buffalo.Context) error {
	q := &models.Question{}
	if err := c.Bind(q); err != nil {
		return errors.WithStack(err)
	}

	q.CreatedBy = c.Value("userId").(string)
	q.UpdatedBy = c.Value("userId").(string)
	q.Author = c.Value("userId").(string)

	tx, _ := c.Value("tx").(*gorm.DB)

	verrors, err := services.CreateQuestion(tx, q)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrors.HasAny() {
		c.Set("question", q)
		c.Set("verrors", verrors)
		return c.Render(200, r.HTML("questions/ask.html"))
	}

	return c.Redirect(302, "/questions/"+strconv.FormatInt(q.ID, 10))
}

// QuestionDetail returns the question with all its details.
// 1st 5 comments for questions + answers and 1st 5 Answers.
func QuestionDetail(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, _ := c.Value("tx").(*gorm.DB)
	questionID := c.Param("questionID")

	qid, err := strconv.ParseInt(questionID, 10, 64)
	if err != nil {
		return errors.WithStack(err)
	}

	question, err := services.GetQuestionDetails(tx, c.Value("currentUser"), qid)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("Question", question)
	c.Set("getContent", func(content string) template.HTML {
		return template.HTML(template.JSEscapeString(content))
	})
	return c.Render(200, r.HTML("questions/detail.html"))
}