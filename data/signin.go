package data

import (
	"time"

	"github.com/ritsec/ops-bot-iii/ent"
	"github.com/ritsec/ops-bot-iii/ent/signin"
	"github.com/ritsec/ops-bot-iii/ent/user"
	"github.com/ritsec/ops-bot-iii/structs"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Signin is the interface for interacting with the signin table
type signin_s struct{}

// Create creates a new signin for a user
func (*signin_s) Create(userID string, signinType signin.Type, ctx ddtrace.SpanContext) (*ent.Signin, error) {
	span := tracer.StartSpan(
		"data.signin:Create",
		tracer.ResourceName("Data.Signin.Create"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	entUser, err := User.Get(userID, span.Context())
	if err != nil {
		return nil, err
	}

	return Client.Signin.Create().
		SetUser(entUser).
		SetType(signinType).
		Save(Ctx)
}

// GetSignins gets all signins for a user
func (*signin_s) GetSignins(id string, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.signin:GetSignins",
		tracer.ResourceName("Data.Signin.GetSignins"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Signin.Query().
		Where(signin.HasUserWith(user.IDEQ(id))).
		Count(Ctx)
}

// GetSigninsByType gets all signins for a user of a specific type
func (*signin_s) GetSigninsByType(id string, signinType signin.Type, ctx ddtrace.SpanContext) (int, error) {
	span := tracer.StartSpan(
		"data.signin:GetSigninsByType",
		tracer.ResourceName("Data.Signin.GetSigninsByType"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	return Client.Signin.Query().
		Where(
			signin.HasUserWith(user.IDEQ(id)),
			signin.TypeEQ(signinType),
		).
		Count(Ctx)
}

// RecentSignin checks if a user has signed in recently
func (*signin_s) RecentSignin(userID string, signinType signin.Type, ctx ddtrace.SpanContext) (bool, error) {
	span := tracer.StartSpan(
		"data.signin:RecentSignin",
		tracer.ResourceName("Data.Signin.RecentSignin"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	ok, err := Client.Signin.Query().
		Where(
			signin.HasUserWith(user.IDEQ(userID)),
			signin.TypeEQ(signinType),
			signin.TimestampGTE(time.Now().Add(-12*time.Hour)),
		).
		Exist(Ctx)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (*signin_s) Query(delta time.Duration, signinType signin.Type, ctx ddtrace.SpanContext) (structs.PairList[string], error) {
	span := tracer.StartSpan(
		"data.signin:Query",
		tracer.ResourceName("Data.Signin.Query"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	var (
		entSignins []*ent.Signin
		err        error
	)

	if signinType == "All" {
		entSignins, err = Client.Signin.Query().
			Where(
				signin.TimestampGTE(time.Now().Add(-delta)),
			).
			WithUser().
			All(Ctx)
	} else {
		entSignins, err = Client.Signin.Query().
			Where(
				signin.TypeEQ(signinType),
				signin.TimestampGTE(time.Now().Add(-delta)),
			).
			WithUser().
			All(Ctx)
	}
	if err != nil {
		return nil, err
	}

	userCount := make(map[string]int)

	for _, entSignin := range entSignins {
		userCount[entSignin.Edges.User.ID]++
	}

	pairList := make(structs.PairList[string], len(userCount))

	i := 0
	for userID, count := range userCount {
		pairList[i] = structs.Pair[string]{Key: userID, Value: count}
		i++
	}

	pairList.Sort()
	pairList.Reverse()

	return pairList, nil
}

func (s *signin_s) DQuery(date time.Time, signinType signin.Type, ctx ddtrace.SpanContext) (structs.PairList[string], error) {
	span := tracer.StartSpan(
		"data.signin:Query",
		tracer.ResourceName("Data.Signin.Query"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	var (
		entSignins []*ent.Signin
		err        error
	)

	// Calculate start and end of the specified date
	// Make the time 00:00:00
	startOfDay := date.Truncate(24 * time.Hour)
	// Set end time to end of the day 11:59:59
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)

	if signinType == "All" {
		entSignins, err = Client.Signin.Query().
			Where(
				signin.TimestampGTE(startOfDay),
				signin.TimestampLTE(endOfDay),
			).
			WithUser().
			All(Ctx)
	} else {
		entSignins, err = Client.Signin.Query().
			Where(
				signin.TypeEQ(signinType),
				signin.TimestampGTE(startOfDay),
				signin.TimestampLTE(endOfDay),
			).
			WithUser().
			All(Ctx)
	}
	if err != nil {
		return nil, err
	}

	userCount := make(map[string]int)

	for _, entSignin := range entSignins {
		userCount[entSignin.Edges.User.ID]++
	}

	pairList := make(structs.PairList[string], len(userCount))

	i := 0
	for userID, count := range userCount {
		pairList[i] = structs.Pair[string]{Key: userID, Value: count}
		i++
	}

	pairList.Sort()
	pairList.Reverse()

	return pairList, nil
}
