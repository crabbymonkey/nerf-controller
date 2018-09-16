package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetPort(t *testing.T) {
	// Save value so we can set it back
	orig, isSet := os.LookupEnv("PORT")

	err := os.Unsetenv("PORT")
	if err != nil {
		t.Fatal(err)
	}
	expected := ":8080"

	if getPort() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			getPort(), expected)
	}

	err = os.Setenv("PORT", "3000")
	if err != nil {
		t.Fatal(err)
	}
	expected = ":3000"

	if getPort() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			getPort(), expected)
	}

	if isSet {
		err = os.Setenv("PORT", orig)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		err = os.Unsetenv("PORT")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestHomeHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestFireHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/fire", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("POST", "/fire", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fireHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusSeeOther {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusSeeOther)
		}

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestTokenHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("POST", "/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(activateHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusSeeOther {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusSeeOther)
		}

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestRandomPageHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("GET", "/fire", nil)
	if err != nil {
		t.Fatal(err)
	}

	req3, err := http.NewRequest("GET", "/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	req4, err := http.NewRequest("GET", "/testindex.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	req5, err := http.NewRequest("GET", "/Thisshouldnotbeaurl", nil)
	if err != nil {
		t.Fatal(err)
	}

	req6, err := http.NewRequest("GET", "/This should be & a ) bad / request", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2, req3, req4, req5, req6}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(randomPageHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestFire(t *testing.T) {
	fire()
}

func TestRandomValue(t *testing.T) {
	sum := 0
	//TODO: This should be a map so it can be broken out into a helper function
	hasBeen5 := false
	hasBeen6 := false
	hasBeen7 := false
	hasBeen8 := false
	hasBeen9 := false
	hasBeen10 := false
	hasBeen11 := false
	hasBeen12 := false
	hasBeen13 := false
	hasBeen14 := false
	hasBeen15 := false
	allNumHit := false

	for i := 0; i < 10000; i++ {
		num := randomValue(5, 15)
		switch num {
		case 5:
			if !hasBeen5 {
				fmt.Println("5")
				hasBeen5 = true
			}
		case 6:
			if !hasBeen6 {
				fmt.Println("6")
				hasBeen6 = true
			}
		case 7:
			if !hasBeen7 {
				fmt.Println("7")
				hasBeen7 = true
			}
		case 8:
			if !hasBeen8 {
				fmt.Println("8")
				hasBeen8 = true
			}
		case 9:
			if !hasBeen9 {
				fmt.Println("9")
				hasBeen9 = true
			}
		case 10:
			if !hasBeen10 {
				fmt.Println("10")
				hasBeen10 = true
			}
		case 11:
			if !hasBeen11 {
				fmt.Println("11")
				hasBeen11 = true
			}
		case 12:
			if !hasBeen12 {
				fmt.Println("12")
				hasBeen12 = true
			}
		case 13:
			if !hasBeen13 {
				fmt.Println("13")
				hasBeen13 = true
			}
		case 14:
			if !hasBeen14 {
				fmt.Println("14")
				hasBeen14 = true
			}
		case 15:
			if !hasBeen15 {
				fmt.Println("15")
				hasBeen15 = true
			}
		default:
			break
		}

		if hasBeen5 &&
			hasBeen6 &&
			hasBeen7 &&
			hasBeen8 &&
			hasBeen9 &&
			hasBeen10 &&
			hasBeen11 &&
			hasBeen12 &&
			hasBeen13 &&
			hasBeen14 &&
			hasBeen15 {
			allNumHit = true
			break
		}
		sum += i
	}

	if !allNumHit {
		t.Errorf("not all numbers hit in test")
	}
}
