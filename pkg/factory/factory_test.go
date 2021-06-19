package factory_test

import (
	"errors"
	"github.com/tjpnz/structbot/internal/transformer"
	"github.com/tjpnz/structbot/pkg/config"
	"github.com/tjpnz/structbot/pkg/factory"
	"reflect"
	"testing"
	"time"
)

type Book struct {
	Title     string
	ISBN      string
	Price     int
	Author    *Author
	Published time.Time
}

type Author struct {
	FamilyName  string
	FirstName   string
	Birthday    time.Time
	Nationality Country
}

type Country struct {
	Code string
	Name string
}

type TestingT struct {
	testing.TB
	fatalfCalled bool
}

func (tt *TestingT) Fatalf(_ string, _ ...interface{}) {
	tt.fatalfCalled = true
}

var bookWithNoAuthor, bookWithAuthorOfUnknownNationality, bookWithAllFields, bookWithMissingPages factory.Factory

func init() {
	c := &config.Config{Transformers: transformer.TimeTransformer{}}

	bookWithNoAuthor = factory.New(c, func() (interface{}, error) {
		return &Book{
			Title:     "Diary of an Oxygen Thief",
			ISBN:      "978-0615275062",
			Price:     2999,
			Published: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil
	})

	bookWithAuthorOfUnknownNationality = factory.New(c, func() (interface{}, error) {
		return &Book{
			Title: "The Terminal Man",
			ISBN:  "978-0552152747",
			Price: 9838,
			Author: &Author{
				FamilyName: "Mehran",
				FirstName:  "Alfred",
				Birthday:   time.Date(1946, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Published: time.Date(2004, 8, 31, 0, 0, 0, 0, time.UTC),
		}, nil
	})

	bookWithAllFields = factory.New(c, func() (interface{}, error) {
		return &Book{
			Title: "Harry Potter and the Philosopher's Stone",
			ISBN:  "0-7475-3269-9",
			Price: 1500,
			Author: &Author{
				FamilyName: "Rowling",
				FirstName:  "J. K.",
				Birthday:   time.Date(1965, 7, 31, 0, 0, 0, 0, time.UTC),
				Nationality: Country{
					Code: "uk",
					Name: "United Kingdom",
				},
			},
			Published: time.Date(1997, 6, 26, 0, 0, 0, 0, time.UTC),
		}, nil
	})

	bookWithMissingPages = factory.New(c, func() (interface{}, error) {
		return nil, errors.New("error")
	})
}

func TestFactory_Create(t *testing.T) {
	for name, tc := range map[string]struct {
		sut    factory.Factory
		outVal *Book
		outErr error
	}{
		"BookWithNoAuthor": {
			sut: bookWithNoAuthor,
			outVal: &Book{
				Title:     "Diary of an Oxygen Thief",
				ISBN:      "978-0615275062",
				Price:     2999,
				Published: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"BookWithAuthorOfUnknownNationality": {
			sut: bookWithAuthorOfUnknownNationality,
			outVal: &Book{
				Title: "The Terminal Man",
				ISBN:  "978-0552152747",
				Price: 9838,
				Author: &Author{
					FamilyName: "Mehran",
					FirstName:  "Alfred",
					Birthday:   time.Date(1946, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Published: time.Date(2004, 8, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		"BookWithAllFields": {
			sut: bookWithAllFields,
			outVal: &Book{
				Title: "Harry Potter and the Philosopher's Stone",
				ISBN:  "0-7475-3269-9",
				Price: 1500,
				Author: &Author{
					FamilyName: "Rowling",
					FirstName:  "J. K.",
					Birthday:   time.Date(1965, 7, 31, 0, 0, 0, 0, time.UTC),
					Nationality: Country{
						Code: "uk",
						Name: "United Kingdom",
					},
				},
				Published: time.Date(1997, 6, 26, 0, 0, 0, 0, time.UTC),
			},
		},
		"BookWithMissingPages": {
			sut:    bookWithMissingPages,
			outErr: errors.New("error"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotVal, gotErr := tc.sut.Create()
			if got, want := gotVal, tc.outVal; tc.outVal != nil && !reflect.DeepEqual(got, want) {
				t.Errorf("got: %v, want: %v", got, want)
			}
			if got, want := gotErr, tc.outErr; got != nil && got.Error() != tc.outErr.Error() {
				t.Errorf("got: %v, want: %v", got, want)
			}
		})
	}
}

func TestFactory_Patch(t *testing.T) {
	for name, tc := range map[string]struct {
		sut    factory.Factory
		in     *Book
		outVal *Book
		outErr error
	}{
		"PatchNothing": {
			sut: bookWithAllFields,
			in:  &Book{},
			outVal: &Book{
				Title: "Harry Potter and the Philosopher's Stone",
				ISBN:  "0-7475-3269-9",
				Price: 1500,
				Author: &Author{
					FamilyName: "Rowling",
					FirstName:  "J. K.",
					Birthday:   time.Date(1965, 7, 31, 0, 0, 0, 0, time.UTC),
					Nationality: Country{
						Code: "uk",
						Name: "United Kingdom",
					},
				},
				Published: time.Date(1997, 6, 26, 0, 0, 0, 0, time.UTC),
			},
		},
		"PatchBookPrice": {
			sut: bookWithAllFields,
			in: &Book{
				Price: 2000,
			},
			outVal: &Book{
				Title: "Harry Potter and the Philosopher's Stone",
				ISBN:  "0-7475-3269-9",
				Price: 2000,
				Author: &Author{
					FamilyName: "Rowling",
					FirstName:  "J. K.",
					Birthday:   time.Date(1965, 7, 31, 0, 0, 0, 0, time.UTC),
					Nationality: Country{
						Code: "uk",
						Name: "United Kingdom",
					},
				},
				Published: time.Date(1997, 6, 26, 0, 0, 0, 0, time.UTC),
			},
		},
		"PatchBookAuthorNationality": {
			sut: bookWithAllFields,
			in: &Book{
				Author: &Author{
					Nationality: Country{
						Code: "gb",
						Name: "Great Britain",
					},
				},
			},
			outVal: &Book{
				Title: "Harry Potter and the Philosopher's Stone",
				ISBN:  "0-7475-3269-9",
				Price: 1500,
				Author: &Author{
					FamilyName: "Rowling",
					FirstName:  "J. K.",
					Birthday:   time.Date(1965, 7, 31, 0, 0, 0, 0, time.UTC),
					Nationality: Country{
						Code: "gb",
						Name: "Great Britain",
					},
				},
				Published: time.Date(1997, 6, 26, 0, 0, 0, 0, time.UTC),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			gotVal, gotErr := tc.sut.Patch(tc.in)
			if got, want := gotVal, tc.outVal; tc.outVal != nil && !reflect.DeepEqual(got, want) {
				t.Errorf("got: %v, want: %v", got, want)
			}
			if got, want := gotErr, tc.outErr; got != nil && got.Error() != tc.outErr.Error() {
				t.Errorf("got: %v, want: %v", got, want)
			}
		})
	}
}

func TestFactory_MustCreate(t *testing.T) {
	sut := bookWithMissingPages
	tt := &TestingT{}
	sut.MustCreate(tt)
	if !tt.fatalfCalled {
		t.Errorf("expected call to fatalf")
	}
}

func TestFactory_MustPatch(t *testing.T) {
	sut := bookWithMissingPages
	tt := &TestingT{}
	sut.MustPatch(tt, &Book{})
	if !tt.fatalfCalled {
		t.Errorf("expected call to fatalf")
	}
}
